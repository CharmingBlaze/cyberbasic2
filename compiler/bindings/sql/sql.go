// Package sql provides SQLite bindings for CyberBasic: OpenDatabase, Exec, Query, result access (GetCell, GetRowCount, etc.), transactions, and parameterized statements.
package sql

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"

	"cyberbasic/compiler/vm"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case int64:
		return int(x)
	case float64:
		return int(x)
	default:
		return 0
	}
}

var (
	dbs      = make(map[string]*sql.DB)
	dbCounter int
	sqlMu    sync.Mutex

	// current result set (last Query/QueryParams)
	resultRows   [][]interface{}
	resultCols   []string
	resultMu     sync.Mutex
	lastErr      string
	lastErrMu    sync.Mutex
)

func setLastErr(err error) {
	lastErrMu.Lock()
	defer lastErrMu.Unlock()
	if err != nil {
		lastErr = err.Error()
	} else {
		lastErr = ""
	}
}

// cellValue normalizes a scanned value for BASIC: nil for NULL, int64/float64/string otherwise.
func cellValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	// sqlite may return int64, float64, []byte (text), *[]byte, string
	switch x := v.(type) {
	case []byte:
		return string(x)
	case *[]byte:
		if x != nil {
			return string(*x)
		}
		return nil
	case int64, float64, string:
		return x
	default:
		return fmt.Sprint(v)
	}
}

// RegisterSQL registers SQLite functions with the VM.
func RegisterSQL(v *vm.VM) {
	// --- Connection ---
	v.RegisterForeign("OpenDatabase", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("OpenDatabase(path) requires 1 argument")
		}
		path := toString(args[0])
		db, err := sql.Open("sqlite", path)
		if err != nil {
			setLastErr(err)
			return nil, nil
		}
		if err := db.Ping(); err != nil {
			_ = db.Close()
			setLastErr(err)
			return nil, nil
		}
		sqlMu.Lock()
		dbCounter++
		id := fmt.Sprintf("db_%d", dbCounter)
		dbs[id] = db
		sqlMu.Unlock()
		setLastErr(nil)
		return id, nil
	})

	v.RegisterForeign("CloseDatabase", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CloseDatabase(dbId) requires 1 argument")
		}
		id := toString(args[0])
		sqlMu.Lock()
		db, ok := dbs[id]
		if ok {
			delete(dbs, id)
		}
		sqlMu.Unlock()
		if ok {
			_ = db.Close()
		}
		return nil, nil
	})

	// --- Exec / Query (no params) ---
	v.RegisterForeign("Exec", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Exec(dbId, sql) requires 2 arguments")
		}
		id := toString(args[0])
		stmt := toString(args[1])
		sqlMu.Lock()
		db, ok := dbs[id]
		sqlMu.Unlock()
		if !ok {
			setLastErr(fmt.Errorf("unknown database: %s", id))
			return -1, nil
		}
		res, err := db.Exec(stmt)
		if err != nil {
			setLastErr(err)
			return -1, nil
		}
		affected, _ := res.RowsAffected()
		setLastErr(nil)
		return int(affected), nil
	})

	v.RegisterForeign("Query", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Query(dbId, sql) requires 2 arguments")
		}
		id := toString(args[0])
		stmt := toString(args[1])
		sqlMu.Lock()
		db, ok := dbs[id]
		sqlMu.Unlock()
		if !ok {
			setLastErr(fmt.Errorf("unknown database: %s", id))
			return -1, nil
		}
		rows, err := db.Query(stmt)
		if err != nil {
			setLastErr(err)
			return -1, nil
		}
		defer rows.Close()
		cols, err := rows.Columns()
		if err != nil {
			setLastErr(err)
			return -1, nil
		}
		var result [][]interface{}
		dest := make([]interface{}, len(cols))
		for i := range dest {
			var v interface{}
			dest[i] = &v
		}
		for rows.Next() {
			if err := rows.Scan(dest...); err != nil {
				setLastErr(err)
				return -1, nil
			}
			row := make([]interface{}, len(cols))
			for i := range cols {
				row[i] = cellValue(*(dest[i].(*interface{})))
			}
			result = append(result, row)
		}
		if err := rows.Err(); err != nil {
			setLastErr(err)
			return -1, nil
		}
		resultMu.Lock()
		resultRows = result
		resultCols = cols
		resultMu.Unlock()
		setLastErr(nil)
		return len(result), nil
	})

	// --- Result access ---
	v.RegisterForeign("GetRowCount", func(args []interface{}) (interface{}, error) {
		resultMu.Lock()
		n := len(resultRows)
		resultMu.Unlock()
		return n, nil
	})

	v.RegisterForeign("GetColumnCount", func(args []interface{}) (interface{}, error) {
		resultMu.Lock()
		n := len(resultCols)
		resultMu.Unlock()
		return n, nil
	})

	v.RegisterForeign("GetColumnName", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetColumnName(colIndex) requires 1 argument")
		}
		col := toInt(args[0])
		resultMu.Lock()
		defer resultMu.Unlock()
		if col < 0 || col >= len(resultCols) {
			return "", nil
		}
		return resultCols[col], nil
	})

	v.RegisterForeign("GetCell", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetCell(row, col) requires 2 arguments")
		}
		row, col := toInt(args[0]), toInt(args[1])
		resultMu.Lock()
		defer resultMu.Unlock()
		if row < 0 || row >= len(resultRows) {
			return nil, nil
		}
		r := resultRows[row]
		if col < 0 || col >= len(r) {
			return nil, nil
		}
		return r[col], nil
	})

	v.RegisterForeign("LastError", func(args []interface{}) (interface{}, error) {
		lastErrMu.Lock()
		s := lastErr
		lastErrMu.Unlock()
		return s, nil
	})

	// --- Transactions ---
	v.RegisterForeign("Begin", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Begin(dbId) requires 1 argument")
		}
		id := toString(args[0])
		sqlMu.Lock()
		db, ok := dbs[id]
		sqlMu.Unlock()
		if !ok {
			setLastErr(fmt.Errorf("unknown database: %s", id))
			return 0, nil
		}
		_, err := db.Exec("BEGIN")
		if err != nil {
			setLastErr(err)
			return 0, nil
		}
		setLastErr(nil)
		return 1, nil
	})

	v.RegisterForeign("Commit", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Commit(dbId) requires 1 argument")
		}
		id := toString(args[0])
		sqlMu.Lock()
		db, ok := dbs[id]
		sqlMu.Unlock()
		if !ok {
			setLastErr(fmt.Errorf("unknown database: %s", id))
			return 0, nil
		}
		_, err := db.Exec("COMMIT")
		if err != nil {
			setLastErr(err)
			return 0, nil
		}
		setLastErr(nil)
		return 1, nil
	})

	v.RegisterForeign("Rollback", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Rollback(dbId) requires 1 argument")
		}
		id := toString(args[0])
		sqlMu.Lock()
		db, ok := dbs[id]
		sqlMu.Unlock()
		if !ok {
			setLastErr(fmt.Errorf("unknown database: %s", id))
			return 0, nil
		}
		_, err := db.Exec("ROLLBACK")
		if err != nil {
			setLastErr(err)
			return 0, nil
		}
		setLastErr(nil)
		return 1, nil
	})

	// --- Parameterized Exec / Query ---
	v.RegisterForeign("ExecParams", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ExecParams(dbId, sql, ...args) requires at least 2 arguments")
		}
		id := toString(args[0])
		stmt := toString(args[1])
		sqlMu.Lock()
		db, ok := dbs[id]
		sqlMu.Unlock()
		if !ok {
			setLastErr(fmt.Errorf("unknown database: %s", id))
			return -1, nil
		}
		params := make([]interface{}, 0, len(args)-2)
		for i := 2; i < len(args); i++ {
			params = append(params, args[i])
		}
		res, err := db.Exec(stmt, params...)
		if err != nil {
			setLastErr(err)
			return -1, nil
		}
		affected, _ := res.RowsAffected()
		setLastErr(nil)
		return int(affected), nil
	})

	v.RegisterForeign("QueryParams", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("QueryParams(dbId, sql, ...args) requires at least 2 arguments")
		}
		id := toString(args[0])
		stmt := toString(args[1])
		sqlMu.Lock()
		db, ok := dbs[id]
		sqlMu.Unlock()
		if !ok {
			setLastErr(fmt.Errorf("unknown database: %s", id))
			return -1, nil
		}
		params := make([]interface{}, 0, len(args)-2)
		for i := 2; i < len(args); i++ {
			params = append(params, args[i])
		}
		rows, err := db.Query(stmt, params...)
		if err != nil {
			setLastErr(err)
			return -1, nil
		}
		defer rows.Close()
		cols, err := rows.Columns()
		if err != nil {
			setLastErr(err)
			return -1, nil
		}
		var result [][]interface{}
		dest := make([]interface{}, len(cols))
		for i := range dest {
			var v interface{}
			dest[i] = &v
		}
		for rows.Next() {
			if err := rows.Scan(dest...); err != nil {
				setLastErr(err)
				return -1, nil
			}
			row := make([]interface{}, len(cols))
			for i := range cols {
				row[i] = cellValue(*(dest[i].(*interface{})))
			}
			result = append(result, row)
		}
		if err := rows.Err(); err != nil {
			setLastErr(err)
			return -1, nil
		}
		resultMu.Lock()
		resultRows = result
		resultCols = cols
		resultMu.Unlock()
		setLastErr(nil)
		return len(result), nil
	})
}
