package sql

var sqlV2 = map[string]string{
	"opendatabase":    "OpenDatabase",
	"closedatabase":   "CloseDatabase",
	"exec":            "Exec",
	"query":           "Query",
	"getrowcount":     "GetRowCount",
	"getcolumncount":  "GetColumnCount",
	"getcolumnname":   "GetColumnName",
	"getcell":         "GetCell",
	"lasterror":       "LastError",
	"begin":           "Begin",
	"commit":          "Commit",
	"rollback":        "Rollback",
	"execparams":      "ExecParams",
	"queryparams":     "QueryParams",
}
