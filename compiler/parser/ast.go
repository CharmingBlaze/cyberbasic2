package parser

import "fmt"

// NodeType represents different types of AST nodes
type NodeType int

const (
	NodeProgram NodeType = iota
	NodeStatement
	NodeExpression
	NodeIfStatement
	NodeForStatement
	NodeWhileStatement
	NodeFunctionDecl
	NodeSubDecl
	NodeAssignment
	NodeBinaryOp
	NodeUnaryOp
	NodeCall
	NodeNumber
	NodeString
	NodeBoolean
	NodeNil
	NodeIdentifier
	NodeBlock
	NodeReturnStatement
	NodeDimStatement
	NodeGameCommand
	NodeSelectCaseStatement
	NodeRepeatStatement
	NodeMemberAccess
	NodeConstStatement
	NodeEnumStatement
	NodeTypeDecl
	NodeEntityDecl
	NodeCompoundAssign
	NodeExitLoop
	NodeContinueLoop
	NodeAssertStatement
	NodeModuleStatement
	NodeOnEventStatement
	NodeStartCoroutineStatement
	NodeYieldStatement
	NodeWaitSecondsStatement
	NodeJSONIndexAccess
	NodeDictLiteral
)

// Node represents a node in the Abstract Syntax Tree
type Node interface {
	Type() NodeType
	String() string
}

// HasSourceLoc is implemented by nodes that carry source location for error reporting.
type HasSourceLoc interface {
	GetLine() int
	GetCol() int
}

// Program represents the entire program
type Program struct {
	Statements []Node
}

func (p *Program) Type() NodeType { return NodeProgram }
func (p *Program) String() string {
	var result string
	for _, stmt := range p.Statements {
		result += stmt.String() + "\n"
	}
	return result
}

// Statement represents a generic statement
type Statement struct {
	Value Node
}

func (s *Statement) Type() NodeType { return NodeStatement }
func (s *Statement) String() string { return s.Value.String() }

// ElseIfBranch represents one ELSEIF condition THEN block in an IF statement.
type ElseIfBranch struct {
	Condition Node
	Block     *Block
}

// IfStatement represents an IF...THEN...[ELSEIF...THEN...]*[ELSE...]ENDIF block
type IfStatement struct {
	Condition Node
	ThenBlock *Block
	ElseIfs   []ElseIfBranch
	ElseBlock *Block
}

func (i *IfStatement) Type() NodeType { return NodeIfStatement }
func (i *IfStatement) String() string {
	result := "IF " + i.Condition.String() + " THEN\n"
	result += i.ThenBlock.String()
	for _, b := range i.ElseIfs {
		result += "ELSEIF " + b.Condition.String() + " THEN\n"
		result += b.Block.String()
	}
	if i.ElseBlock != nil {
		result += "ELSE\n" + i.ElseBlock.String()
	}
	result += "ENDIF"
	return result
}

// ForStatement represents a FOR...TO...STEP...NEXT loop
type ForStatement struct {
	Variable string
	Start    Node
	End      Node
	Step     Node
	Body     *Block
}

func (f *ForStatement) Type() NodeType { return NodeForStatement }
func (f *ForStatement) String() string {
	result := "FOR " + f.Variable + " = " + f.Start.String() + " TO " + f.End.String()
	if f.Step != nil {
		result += " STEP " + f.Step.String()
	}
	result += "\n" + f.Body.String() + "NEXT"
	return result
}

// WhileStatement represents a WHILE...WEND loop
type WhileStatement struct {
	Condition Node
	Body      *Block
}

func (w *WhileStatement) Type() NodeType { return NodeWhileStatement }
func (w *WhileStatement) String() string {
	return "WHILE " + w.Condition.String() + "\n" + w.Body.String() + "WEND"
}

// FunctionDecl represents a function declaration (optionally inside a Module)
type FunctionDecl struct {
	Name       string
	ModuleName string // set when inside Module X ... End Module
	Parameters []string
	ReturnType string
	Body       *Block
}

func (f *FunctionDecl) Type() NodeType { return NodeFunctionDecl }
func (f *FunctionDecl) String() string {
	result := "FUNCTION " + f.Name + "("
	for i, param := range f.Parameters {
		if i > 0 {
			result += ", "
		}
		result += param
	}
	result += ")"
	if f.ReturnType != "" {
		result += " AS " + f.ReturnType
	}
	result += "\n" + f.Body.String() + "END FUNCTION"
	return result
}

// SubDecl represents a sub procedure declaration (optionally inside a Module)
type SubDecl struct {
	Name       string
	ModuleName string // set when inside Module X ... End Module
	Parameters []string
	Body       *Block
}

func (s *SubDecl) Type() NodeType { return NodeSubDecl }
func (s *SubDecl) String() string {
	result := "SUB " + s.Name + "("
	for i, param := range s.Parameters {
		if i > 0 {
			result += ", "
		}
		result += param
	}
	result += ")\n" + s.Body.String() + "END SUB"
	return result
}

// ModuleStatement represents Module Name ... End Module (body is FunctionDecl/SubDecl only)
type ModuleStatement struct {
	Name string
	Body []Node
}

func (m *ModuleStatement) Type() NodeType { return NodeModuleStatement }
func (m *ModuleStatement) String() string {
	result := "MODULE " + m.Name + "\n"
	for _, n := range m.Body {
		result += n.String()
	}
	return result + "END MODULE"
}

// OnEventStatement represents On KeyDown("ESCAPE") ... End On
type OnEventStatement struct {
	EventType string // "keydown", "keypressed"
	Key       string // key name e.g. "ESCAPE" or empty for any
	Body      *Block
}

func (o *OnEventStatement) Type() NodeType { return NodeOnEventStatement }
func (o *OnEventStatement) String() string {
	return "On " + o.EventType + "(\"" + o.Key + "\") ... End On"
}

// StartCoroutineStatement represents StartCoroutine SubName()
type StartCoroutineStatement struct {
	SubName string
}

func (s *StartCoroutineStatement) Type() NodeType { return NodeStartCoroutineStatement }
func (s *StartCoroutineStatement) String() string { return "StartCoroutine " + s.SubName + "()" }

// YieldStatement represents Yield
type YieldStatement struct{}

func (y *YieldStatement) Type() NodeType   { return NodeYieldStatement }
func (y *YieldStatement) String() string   { return "Yield" }

// WaitSecondsStatement represents WaitSeconds(seconds)
type WaitSecondsStatement struct {
	Seconds Node
}

func (w *WaitSecondsStatement) Type() NodeType { return NodeWaitSecondsStatement }
func (w *WaitSecondsStatement) String() string  { return "WaitSeconds(...)" }

// Assignment represents variable assignment (scalar or array element)
type Assignment struct {
	Variable string
	Indices  []Node // nil = scalar assignment; non-nil = array element
	Value    Node
	Line     int
	Col      int
}

func (a *Assignment) Type() NodeType   { return NodeAssignment }
func (a *Assignment) GetLine() int    { return a.Line }
func (a *Assignment) GetCol() int     { return a.Col }
func (a *Assignment) String() string {
	if len(a.Indices) > 0 {
		s := a.Variable + "("
		for i, idx := range a.Indices {
			if i > 0 {
				s += ", "
			}
			s += idx.String()
		}
		s += ") = " + a.Value.String()
		return s
	}
	return a.Variable + " = " + a.Value.String()
}

// CompoundAssign represents +=, -=, *=, /= (scalar only).
type CompoundAssign struct {
	Variable string
	Op       string // "+=", "-=", "*=", "/="
	Value    Node
	Line     int
	Col      int
}

func (a *CompoundAssign) Type() NodeType { return NodeCompoundAssign }
func (a *CompoundAssign) GetLine() int   { return a.Line }
func (a *CompoundAssign) GetCol() int   { return a.Col }
func (a *CompoundAssign) String() string { return a.Variable + " " + a.Op + " " + a.Value.String() }

// ExitLoopStatement represents EXIT FOR or EXIT WHILE (or BREAK FOR / BREAK WHILE).
type ExitLoopStatement struct {
	Kind string // "FOR" or "WHILE"
}

func (e *ExitLoopStatement) Type() NodeType { return NodeExitLoop }
func (e *ExitLoopStatement) String() string  { return "EXIT " + e.Kind }

// ContinueLoopStatement represents CONTINUE FOR or CONTINUE WHILE.
type ContinueLoopStatement struct {
	Kind string // "FOR" or "WHILE"
}

func (c *ContinueLoopStatement) Type() NodeType { return NodeContinueLoop }
func (c *ContinueLoopStatement) String() string { return "CONTINUE " + c.Kind }

// AssertStatement represents ASSERT condition [, message].
type AssertStatement struct {
	Condition Node
	Message   Node // optional; nil = use default "assertion failed"
}

func (a *AssertStatement) Type() NodeType { return NodeAssertStatement }
func (a *AssertStatement) String() string {
	if a.Message != nil {
		return "ASSERT " + a.Condition.String() + ", " + a.Message.String()
	}
	return "ASSERT " + a.Condition.String()
}

// MemberAccess represents dot notation: expr.member (e.g. pos.x, GetMousePosition().y)
type MemberAccess struct {
	Object Node
	Member string
}

func (m *MemberAccess) Type() NodeType { return NodeMemberAccess }
func (m *MemberAccess) String() string  { return m.Object.String() + "." + m.Member }

// JSONIndexAccess represents obj["key"] sugar, compiled to GetJSONKey(obj, "key")
type JSONIndexAccess struct {
	Object Node   // expression yielding JSON handle
	Key    string // string key
}

func (j *JSONIndexAccess) Type() NodeType { return NodeJSONIndexAccess }
func (j *JSONIndexAccess) String() string { return j.Object.String() + "[\"" + j.Key + "\"]" }

// DictPair is one key-value pair in a dict literal (key is string; value is expression).
type DictPair struct {
	Key   string
	Value Node
}

// DictLiteral represents { "key": value } or { key = value } literal.
type DictLiteral struct {
	Pairs []DictPair
}

func (d *DictLiteral) Type() NodeType { return NodeDictLiteral }
func (d *DictLiteral) String() string {
	s := "{"
	for i, p := range d.Pairs {
		if i > 0 {
			s += ", "
		}
		s += "\"" + p.Key + "\": " + p.Value.String()
	}
	return s + "}"
}

// BinaryOp represents a binary operation
type BinaryOp struct {
	Operator string
	Left     Node
	Right    Node
	Line     int
	Col      int
}

func (b *BinaryOp) Type() NodeType { return NodeBinaryOp }
func (b *BinaryOp) GetLine() int   { return b.Line }
func (b *BinaryOp) GetCol() int    { return b.Col }
func (b *BinaryOp) String() string {
	return "(" + b.Left.String() + " " + b.Operator + " " + b.Right.String() + ")"
}

// UnaryOp represents a unary operation
type UnaryOp struct {
	Operator string
	Operand  Node
	Line     int
	Col      int
}

func (u *UnaryOp) Type() NodeType { return NodeUnaryOp }
func (u *UnaryOp) GetLine() int   { return u.Line }
func (u *UnaryOp) GetCol() int    { return u.Col }
func (u *UnaryOp) String() string { return u.Operator + u.Operand.String() }

// Call represents a function or procedure call
type Call struct {
	Name      string
	Arguments []Node
	Line      int
	Col       int
}

func (c *Call) Type() NodeType { return NodeCall }
func (c *Call) GetLine() int   { return c.Line }
func (c *Call) GetCol() int    { return c.Col }
func (c *Call) String() string {
	result := c.Name + "("
	for i, arg := range c.Arguments {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	return result
}

// Number represents a numeric literal
type Number struct {
	Value string
}

func (n *Number) Type() NodeType { return NodeNumber }
func (n *Number) String() string { return n.Value }

// StringLiteral represents a string literal
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Type() NodeType { return NodeString }
func (s *StringLiteral) String() string { return "\"" + s.Value + "\"" }

// Boolean represents a boolean literal
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() NodeType { return NodeBoolean }
func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }

// NilLiteral represents the null/nil literal
type NilLiteral struct{}

func (n *NilLiteral) Type() NodeType { return NodeNil }
func (n *NilLiteral) String() string { return "nil" }

// Identifier represents a variable or function name
type Identifier struct {
	Name string
	Line int
	Col  int
}

func (i *Identifier) Type() NodeType { return NodeIdentifier }
func (i *Identifier) GetLine() int   { return i.Line }
func (i *Identifier) GetCol() int    { return i.Col }
func (i *Identifier) String() string { return i.Name }

// Block represents a block of statements
type Block struct {
	Statements []Node
}

func (b *Block) Type() NodeType { return NodeBlock }
func (b *Block) String() string {
	var result string
	for _, stmt := range b.Statements {
		result += "  " + stmt.String() + "\n"
	}
	return result
}

// ReturnStatement represents a RETURN statement
type ReturnStatement struct {
	Value Node
}

func (r *ReturnStatement) Type() NodeType { return NodeReturnStatement }
func (r *ReturnStatement) String() string {
	if r.Value != nil {
		return "RETURN " + r.Value.String()
	}
	return "RETURN"
}

// ConstStatement represents CONST name = value (one or more, comma-separated).
type ConstStatement struct {
	Decls []ConstDecl
}

// ConstDecl is one name = value in a CONST statement.
type ConstDecl struct {
	Name  string
	Value Node
}

func (c *ConstStatement) Type() NodeType { return NodeConstStatement }
func (c *ConstStatement) String() string {
	result := "CONST "
	for i, d := range c.Decls {
		if i > 0 {
			result += ", "
		}
		result += d.Name + " = " + d.Value.String()
	}
	return result
}

// EnumStatement represents ENUM Name : member1, member2 = expr, ...
type EnumStatement struct {
	Name    string       // enum type name (for reference only)
	Members []EnumMember // member names and optional explicit values
}

// EnumMember is one name with optional = value (nil = auto-increment).
type EnumMember struct {
	Name  string
	Value Node // nil = use auto-increment from previous
}

func (e *EnumStatement) Type() NodeType { return NodeEnumStatement }
func (e *EnumStatement) String() string {
	s := "ENUM " + e.Name + " : "
	for i, m := range e.Members {
		if i > 0 {
			s += ", "
		}
		s += m.Name
		if m.Value != nil {
			s += " = " + m.Value.String()
		}
	}
	return s
}

// TypeDecl represents TYPE Name ... ENDTYPE (UDT definition).
type TypeDecl struct {
	Name   string
	Fields []TypeField
}

// TypeField is one field in a TYPE: Name, optional AS FieldType, optional = ConstValue (for constant groups).
type TypeField struct {
	Name       string
	FieldType  string // e.g. "FLOAT", "STRING", "Vector3", or ""
	ConstValue Node   // if set, this field is a constant-group member; nil = data field
}

func (t *TypeDecl) Type() NodeType { return NodeTypeDecl }

// EntityDecl represents ENTITY Name ... END ENTITY (single instance with properties).
type EntityDecl struct {
	Name       string
	Properties []EntityProperty
}

// EntityProperty is one property in an ENTITY: Name = initial Value.
type EntityProperty struct {
	Name  string
	Value Node
}

func (e *EntityDecl) Type() NodeType { return NodeEntityDecl }
func (e *EntityDecl) String() string {
	s := "ENTITY " + e.Name + "\n"
	for _, p := range e.Properties {
		s += "  " + p.Name + " = " + p.Value.String() + "\n"
	}
	return s + "END ENTITY"
}
func (t *TypeDecl) String() string {
	s := "TYPE " + t.Name + "\n"
	for _, f := range t.Fields {
		s += "  " + f.Name
		if f.FieldType != "" {
			s += " AS " + f.FieldType
		}
		if f.ConstValue != nil {
			s += " = " + f.ConstValue.String()
		}
		s += "\n"
	}
	return s + "ENDTYPE"
}

// DimStatement represents a DIM statement for variable declaration
type DimStatement struct {
	Variables []VariableDecl
}

type VariableDecl struct {
	Name       string
	Type       string
	Dimensions []Node // nil = scalar; e.g. [10], [10,20] for DIM a(10) or DIM a(10,20)
}

func (d *DimStatement) Type() NodeType { return NodeDimStatement }
func (d *DimStatement) String() string {
	result := "DIM "
	for i, v := range d.Variables {
		if i > 0 {
			result += ", "
		}
		result += v.Name
		if len(v.Dimensions) > 0 {
			result += "("
			for j, dim := range v.Dimensions {
				if j > 0 {
					result += ", "
				}
				result += dim.String()
			}
			result += ")"
		}
		if v.Type != "" {
			result += " AS " + v.Type
		}
	}
	return result
}

// SelectCaseStatement represents SELECT CASE expr ... CASE value: block ... CASE ELSE: block ... END SELECT
type SelectCaseStatement struct {
	Expr       Node        // value to match
	Cases      []CaseClause // CASE value: block (Value nil for CASE ELSE)
	ElseBlock  *Block      // optional CASE ELSE block
}

// CaseClause is one CASE value: block (Value nil for ELSE)
type CaseClause struct {
	Value Node  // match value; nil for CASE ELSE
	Block *Block
}

// RepeatStatement represents REPEAT ... UNTIL condition
type RepeatStatement struct {
	Body      *Block
	Condition Node
}

func (r *RepeatStatement) Type() NodeType { return NodeRepeatStatement }
func (r *RepeatStatement) String() string {
	return "REPEAT\n" + r.Body.String() + "UNTIL " + r.Condition.String()
}

func (s *SelectCaseStatement) Type() NodeType { return NodeSelectCaseStatement }
func (s *SelectCaseStatement) String() string {
	out := "SELECT CASE " + s.Expr.String() + "\n"
	for _, c := range s.Cases {
		if c.Value != nil {
			out += "  CASE " + c.Value.String() + ":\n" + c.Block.String()
		} else {
			out += "  CASE ELSE:\n" + c.Block.String()
		}
	}
	if s.ElseBlock != nil {
		out += "  CASE ELSE:\n" + s.ElseBlock.String()
	}
	return out + "END SELECT"
}

// GameCommand represents game-specific commands
type GameCommand struct {
	Command   string
	Arguments []Node
}

func (g *GameCommand) Type() NodeType { return NodeGameCommand }
func (g *GameCommand) String() string {
	result := g.Command
	if len(g.Arguments) > 0 {
		result += "("
		for i, arg := range g.Arguments {
			if i > 0 {
				result += ", "
			}
			result += arg.String()
		}
		result += ")"
	}
	return result
}

// Error represents a parsing error
type Error struct {
	Message string
	Line    int
	Col     int
}

func (e *Error) Error() string {
	return fmt.Sprintf("Parse error at line %d, col %d: %s", e.Line, e.Col, e.Message)
}
