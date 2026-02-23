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
	NodeMainStatement
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
	NodeCompoundAssign
	NodeExitLoop
	NodeModuleStatement
	NodeOnEventStatement
	NodeStartCoroutineStatement
	NodeYieldStatement
	NodeWaitSecondsStatement
	NodeJSONIndexAccess
)

// Node represents a node in the Abstract Syntax Tree
type Node interface {
	Type() NodeType
	String() string
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

// IfStatement represents an IF...THEN...ELSE...ENDIF block
type IfStatement struct {
	Condition Node
	ThenBlock *Block
	ElseBlock *Block
}

func (i *IfStatement) Type() NodeType { return NodeIfStatement }
func (i *IfStatement) String() string {
	result := "IF " + i.Condition.String() + " THEN\n"
	result += i.ThenBlock.String()
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

// MainStatement represents Main ... EndMain (main game loop)
type MainStatement struct {
	Body *Block
}

func (m *MainStatement) Type() NodeType { return NodeMainStatement }
func (m *MainStatement) String() string { return "Main ... EndMain" }

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
}

func (a *Assignment) Type() NodeType { return NodeAssignment }
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
}

func (a *CompoundAssign) Type() NodeType { return NodeCompoundAssign }
func (a *CompoundAssign) String() string { return a.Variable + " " + a.Op + " " + a.Value.String() }

// ExitLoopStatement represents EXIT FOR or EXIT WHILE.
type ExitLoopStatement struct {
	Kind string // "FOR" or "WHILE"
}

func (e *ExitLoopStatement) Type() NodeType { return NodeExitLoop }
func (e *ExitLoopStatement) String() string  { return "EXIT " + e.Kind }

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

// BinaryOp represents a binary operation
type BinaryOp struct {
	Operator string
	Left     Node
	Right    Node
}

func (b *BinaryOp) Type() NodeType { return NodeBinaryOp }
func (b *BinaryOp) String() string {
	return "(" + b.Left.String() + " " + b.Operator + " " + b.Right.String() + ")"
}

// UnaryOp represents a unary operation
type UnaryOp struct {
	Operator string
	Operand  Node
}

func (u *UnaryOp) Type() NodeType { return NodeUnaryOp }
func (u *UnaryOp) String() string { return u.Operator + u.Operand.String() }

// Call represents a function or procedure call
type Call struct {
	Name      string
	Arguments []Node
}

func (c *Call) Type() NodeType { return NodeCall }
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
}

func (i *Identifier) Type() NodeType { return NodeIdentifier }
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
