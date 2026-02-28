package parser

import (
	"cyberbasic/compiler/lexer"
	"strings"
)

func (p *Parser) ifStatement() (Node, error) {
	p.advance() // Skip IF

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(lexer.TokenThen) {
		return nil, &Error{Message: "expected THEN after IF condition", Line: p.line(), Col: p.col()}
	}

	thenBlock, err := p.block(true) // true = allow single-line IF (stop at newline+IF)
	if err != nil {
		return nil, err
	}

	var elseIfs []ElseIfBranch
	for p.match(lexer.TokenElseIf) {
		cond, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(lexer.TokenThen) {
			return nil, &Error{Message: "expected THEN after ELSEIF condition", Line: p.line(), Col: p.col()}
		}
		blk, err := p.block(true)
		if err != nil {
			return nil, err
		}
		elseIfs = append(elseIfs, ElseIfBranch{Condition: cond, Block: blk})
	}

	var elseBlock *Block
	if p.match(lexer.TokenElse) {
		elseBlock, err = p.block(false)
		if err != nil {
			return nil, err
		}
	}

	// Single-line IF: one statement and we're not at ENDIF/END IF/ELSE/ELSEIF => no ENDIF required
	if len(thenBlock.Statements) == 1 && len(elseIfs) == 0 && elseBlock == nil && !p.checkEndIf() && !p.check(lexer.TokenElse) && !p.check(lexer.TokenElseIf) {
		// do not consume ENDIF
	} else if !p.matchEndIf() {
		return nil, &Error{Message: "expected ENDIF or END IF", Line: p.line(), Col: p.col()}
	}

	return &IfStatement{
		Condition: condition,
		ThenBlock: thenBlock,
		ElseIfs:   elseIfs,
		ElseBlock: elseBlock,
	}, nil
}

// forStatement parses FOR...TO...STEP...NEXT
func (p *Parser) forStatement() (Node, error) {
	p.advance() // Skip FOR

	if !p.match(lexer.TokenIdentifier) {
		return nil, &Error{Message: "expected variable name after FOR", Line: p.line(), Col: p.col()}
	}
	variable := p.previous().Value

	if !p.match(lexer.TokenAssign) {
		return nil, &Error{Message: "expected '=' after variable name", Line: p.line(), Col: p.col()}
	}

	start, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(lexer.TokenTo) {
		return nil, &Error{Message: "expected TO", Line: p.line(), Col: p.col()}
	}

	end, err := p.expression()
	if err != nil {
		return nil, err
	}

	var step Node
	if p.match(lexer.TokenStep) {
		step, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	body, err := p.block(false)
	if err != nil {
		return nil, err
	}

	if !p.match(lexer.TokenNext) {
		return nil, &Error{Message: "expected NEXT", Line: p.line(), Col: p.col()}
	}

	return &ForStatement{
		Variable: variable,
		Start:    start,
		End:      end,
		Step:     step,
		Body:     body,
	}, nil
}

// whileStatement parses WHILE...WEND
func (p *Parser) whileStatement() (Node, error) {
	p.advance() // Skip WHILE

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	body, err := p.block(false)
	if err != nil {
		return nil, err
	}

	if !p.match(lexer.TokenWend) {
		return nil, &Error{Message: "expected WEND", Line: p.line(), Col: p.col()}
	}

	return &WhileStatement{
		Condition: condition,
		Body:      body,
	}, nil
}

// functionDecl parses FUNCTION...END FUNCTION
func (p *Parser) functionDecl() (Node, error) {
	p.advance() // Skip FUNCTION

	if !p.match(lexer.TokenIdentifier) {
		return nil, &Error{Message: "expected function name", Line: p.line(), Col: p.col()}
	}
	name := p.previous().Value

	var parameters []string
	if p.match(lexer.TokenLeftParen) {
		for !p.check(lexer.TokenRightParen) {
			if !p.match(lexer.TokenIdentifier) {
				return nil, &Error{Message: "expected parameter name", Line: p.line(), Col: p.col()}
			}
			parameters = append(parameters, p.previous().Value)

			if !p.match(lexer.TokenComma) {
				break
			}
		}
		if !p.match(lexer.TokenRightParen) {
			return nil, &Error{Message: "expected ')'", Line: p.line(), Col: p.col()}
		}
	}

	var returnType string
	if p.match(lexer.TokenAs) {
		if !p.match(lexer.TokenIdentifier) {
			return nil, &Error{Message: "expected return type", Line: p.line(), Col: p.col()}
		}
		returnType = p.previous().Value
	}

	body, err := p.block(false)
	if err != nil {
		return nil, err
	}

	// ENDFUNCTION (single-word) or END FUNCTION (two words)
	if p.match(lexer.TokenEndFunction) {
		// consumed
	} else if p.match(lexer.TokenEnd) {
		if !p.match(lexer.TokenFunction) {
			return nil, &Error{Message: "expected FUNCTION after END", Line: p.line(), Col: p.col()}
		}
	} else {
		return nil, &Error{Message: "expected ENDFUNCTION or END FUNCTION", Line: p.line(), Col: p.col()}
	}

	return &FunctionDecl{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}, nil
}

// subDecl parses SUB...END SUB
func (p *Parser) subDecl() (Node, error) {
	p.advance() // Skip SUB

	if !p.match(lexer.TokenIdentifier) {
		return nil, &Error{Message: "expected sub name", Line: p.line(), Col: p.col()}
	}
	name := p.previous().Value

	var parameters []string
	if p.match(lexer.TokenLeftParen) {
		for !p.check(lexer.TokenRightParen) {
			if !p.match(lexer.TokenIdentifier) {
				return nil, &Error{Message: "expected parameter name", Line: p.line(), Col: p.col()}
			}
			parameters = append(parameters, p.previous().Value)

			if !p.match(lexer.TokenComma) {
				break
			}
		}
		if !p.match(lexer.TokenRightParen) {
			return nil, &Error{Message: "expected ')'", Line: p.line(), Col: p.col()}
		}
	}

	body, err := p.block(false)
	if err != nil {
		return nil, err
	}

	// ENDSUB (single-word) or END SUB (two words)
	if p.match(lexer.TokenEndSub) {
		// consumed
	} else if p.match(lexer.TokenEnd) {
		if !p.match(lexer.TokenSub) {
			return nil, &Error{Message: "expected SUB after END", Line: p.line(), Col: p.col()}
		}
	} else {
		return nil, &Error{Message: "expected ENDSUB or END SUB", Line: p.line(), Col: p.col()}
	}

	return &SubDecl{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}, nil
}

// moduleDecl parses MODULE name ... (Function|Sub)* ... END MODULE
func (p *Parser) moduleDecl() (Node, error) {
	p.advance() // Skip MODULE
	if !p.match(lexer.TokenIdentifier) {
		return nil, &Error{Message: "expected module name", Line: p.line(), Col: p.col()}
	}
	moduleName := p.previous().Value
	var body []Node
	for !p.isAtEnd() {
		for p.match(lexer.TokenNewLine) {
			continue
		}
		// ENDMODULE (single-word) or END MODULE (two words); only consume END when next is MODULE
		if p.match(lexer.TokenEndModule) {
			break
		}
		if p.check(lexer.TokenEnd) && p.current+1 < len(p.tokens) && p.tokens[p.current+1].Type == lexer.TokenModule {
			p.advance()
			p.advance()
			break
		}
		if p.check(lexer.TokenFunction) {
			fn, err := p.functionDecl()
			if err != nil {
				return nil, err
			}
			if fd, ok := fn.(*FunctionDecl); ok {
				fd.ModuleName = moduleName
			}
			body = append(body, fn)
			continue
		}
		if p.check(lexer.TokenSub) {
			sub, err := p.subDecl()
			if err != nil {
				return nil, err
			}
			if sd, ok := sub.(*SubDecl); ok {
				sd.ModuleName = moduleName
			}
			body = append(body, sub)
			continue
		}
		if p.isAtEnd() {
			return nil, &Error{Message: "expected END MODULE", Line: p.line(), Col: p.col()}
		}
		return nil, &Error{Message: "module body may only contain FUNCTION or SUB", Line: p.line(), Col: p.col()}
	}
	return &ModuleStatement{Name: moduleName, Body: body}, nil
}

// onEventStatement parses On KeyDown("ESCAPE") ... End On
func (p *Parser) onEventStatement() (Node, error) {
	p.advance() // Skip ON
	eventType := ""
	switch {
	case p.match(lexer.TokenKeyDown):
		eventType = "keydown"
	case p.match(lexer.TokenKeyPressed):
		eventType = "keypressed"
	default:
		return nil, &Error{Message: "expected KeyDown or KeyPressed after On", Line: p.line(), Col: p.col()}
	}
	key := ""
	if p.match(lexer.TokenLeftParen) {
		if p.match(lexer.TokenString) {
			key = p.previous().Value
		} else if p.match(lexer.TokenNumber) {
			key = p.previous().Value
		}
		if !p.match(lexer.TokenRightParen) {
			return nil, &Error{Message: "expected ')' after key", Line: p.line(), Col: p.col()}
		}
	}
	body, err := p.block(false)
	if err != nil {
		return nil, err
	}
	if p.match(lexer.TokenEndOn) {
		// single token ENDON
	} else if p.match(lexer.TokenEnd) {
		if !p.match(lexer.TokenOn) {
			return nil, &Error{Message: "expected On after END", Line: p.line(), Col: p.col()}
		}
	} else {
		return nil, &Error{Message: "expected End On", Line: p.line(), Col: p.col()}
	}
	return &OnEventStatement{EventType: eventType, Key: key, Body: body}, nil
}

// startCoroutineStatement parses StartCoroutine SubName()
func (p *Parser) startCoroutineStatement() (Node, error) {
	p.advance() // Skip StartCoroutine
	if !p.match(lexer.TokenIdentifier) {
		return nil, &Error{Message: "expected sub name after StartCoroutine", Line: p.line(), Col: p.col()}
	}
	subName := p.previous().Value
	if p.match(lexer.TokenLeftParen) {
		if !p.match(lexer.TokenRightParen) {
			return nil, &Error{Message: "expected ')'", Line: p.line(), Col: p.col()}
		}
	}
	return &StartCoroutineStatement{SubName: subName}, nil
}

// yieldStatement parses Yield
func (p *Parser) yieldStatement() (Node, error) {
	p.advance() // Skip Yield
	return &YieldStatement{}, nil
}

// waitSecondsStatement parses WaitSeconds(expr)
func (p *Parser) waitSecondsStatement() (Node, error) {
	p.advance() // Skip WaitSeconds
	if !p.match(lexer.TokenLeftParen) {
		return nil, &Error{Message: "expected '(' after WaitSeconds", Line: p.line(), Col: p.col()}
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.match(lexer.TokenRightParen) {
		return nil, &Error{Message: "expected ')'", Line: p.line(), Col: p.col()}
	}
	return &WaitSecondsStatement{Seconds: expr}, nil
}

// constStatement parses CONST name = expr (, name = expr)*
func (p *Parser) constStatement() (Node, error) {
	p.advance() // Skip CONST
	var decls []ConstDecl
	for {
		if !p.match(lexer.TokenIdentifier) {
			return nil, &Error{Message: "expected constant name after CONST", Line: p.line(), Col: p.col()}
		}
		name := p.previous().Value
		if !p.match(lexer.TokenAssign) {
			return nil, &Error{Message: "expected '=' after constant name", Line: p.line(), Col: p.col()}
		}
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		decls = append(decls, ConstDecl{Name: name, Value: value})
		if !p.match(lexer.TokenComma) {
			break
		}
	}
	return &ConstStatement{Decls: decls}, nil
}

// checkEndEnum returns true if current position is ENDENUM or END ENUM.
func (p *Parser) checkEndEnum() bool {
	if p.isAtEnd() {
		return false
	}
	if p.check(lexer.TokenEndEnum) {
		return true
	}
	return p.check(lexer.TokenEnd) && p.current+1 < len(p.tokens) && p.tokens[p.current+1].Type == lexer.TokenEnum
}

// matchEndEnum consumes ENDENUM or END ENUM; returns true if matched.
func (p *Parser) matchEndEnum() bool {
	if p.match(lexer.TokenEndEnum) {
		return true
	}
	if p.check(lexer.TokenEnd) && p.current+1 < len(p.tokens) && p.tokens[p.current+1].Type == lexer.TokenEnum {
		p.advance()
		p.advance()
		return true
	}
	return false
}

// enumStatement parses ENUM [Name] : member1, ... or ENUM [Name] newline member1 ... END ENUM (multi-line/unnamed).
func (p *Parser) enumStatement() (Node, error) {
	p.advance() // Skip ENUM
	enumName := ""
	if p.match(lexer.TokenIdentifier) {
		enumName = p.previous().Value
	}
	// Single-line form: ENUM Name : a, b = 2, c
	if enumName != "" && p.match(lexer.TokenColon) {
		var members []EnumMember
		for {
			if !p.match(lexer.TokenIdentifier) {
				return nil, &Error{Message: "expected enum member name", Line: p.line(), Col: p.col()}
			}
			memName := p.previous().Value
			var value Node
			if p.match(lexer.TokenAssign) {
				var err error
				value, err = p.expression()
				if err != nil {
					return nil, err
				}
			}
			members = append(members, EnumMember{Name: memName, Value: value})
			if !p.match(lexer.TokenComma) {
				break
			}
		}
		for p.match(lexer.TokenNewLine) {
			continue
		}
		if p.matchEndEnum() {
			// consumed
		} else if p.check(lexer.TokenEnd) {
			p.advance()
			if !p.match(lexer.TokenEnum) {
				return nil, &Error{Message: "expected ENUM after END", Line: p.line(), Col: p.col()}
			}
		}
		return &EnumStatement{Name: enumName, Members: members}, nil
	}
	// Multi-line form: ENUM [Name] newline members... END ENUM (unnamed if no name)
	for p.match(lexer.TokenNewLine) {
		continue
	}
	var members []EnumMember
	for !p.isAtEnd() && !p.checkEndEnum() {
		for p.match(lexer.TokenNewLine) {
			continue
		}
		if p.checkEndEnum() {
			break
		}
		if !p.match(lexer.TokenIdentifier) {
			return nil, &Error{Message: "expected enum member name or END ENUM", Line: p.line(), Col: p.col()}
		}
		memName := p.previous().Value
		var value Node
		if p.match(lexer.TokenAssign) {
			var err error
			value, err = p.expression()
			if err != nil {
				return nil, err
			}
		}
		members = append(members, EnumMember{Name: memName, Value: value})
		for p.match(lexer.TokenNewLine) {
			continue
		}
		if p.match(lexer.TokenComma) {
			continue
		}
		if p.checkEndEnum() {
			break
		}
	}
	if !p.matchEndEnum() {
		return nil, &Error{Message: "expected END ENUM or ENDENUM", Line: p.line(), Col: p.col()}
	}
	return &EnumStatement{Name: enumName, Members: members}, nil
}

// typeDeclStatement parses TYPE Name ... ENDTYPE (or END TYPE).
func (p *Parser) typeDeclStatement() (Node, error) {
	p.advance() // Skip TYPE
	if !p.match(lexer.TokenIdentifier) {
		return nil, &Error{Message: "expected type name after TYPE", Line: p.line(), Col: p.col()}
	}
	typeName := p.previous().Value
	var fields []TypeField
	for {
		for p.match(lexer.TokenNewLine) {
		}
		if p.match(lexer.TokenEndType) {
			break
		}
		if p.check(lexer.TokenEnd) {
			p.advance()
			if p.match(lexer.TokenTypeKw) {
				break
			}
			return nil, &Error{Message: "expected TYPE after END", Line: p.line(), Col: p.col()}
		}
		if p.isAtEnd() {
			return nil, &Error{Message: "expected ENDTYPE or END TYPE", Line: p.line(), Col: p.col()}
		}
		if !p.match(lexer.TokenIdentifier) {
			return nil, &Error{Message: "expected field name or ENDTYPE", Line: p.line(), Col: p.col()}
		}
		fieldName := p.previous().Value
		fieldType := ""
		var constVal Node
		if p.match(lexer.TokenAs) {
			if !p.match(lexer.TokenIdentifier, lexer.TokenInteger, lexer.TokenStringType, lexer.TokenFloat, lexer.TokenBoolean) {
				return nil, &Error{Message: "expected type after AS", Line: p.line(), Col: p.col()}
			}
			fieldType = p.previous().Value
		}
		if p.match(lexer.TokenAssign) {
			var err error
			constVal, err = p.expression()
			if err != nil {
				return nil, err
			}
		}
		fields = append(fields, TypeField{Name: fieldName, FieldType: fieldType, ConstValue: constVal})
	}
	return &TypeDecl{Name: typeName, Fields: fields}, nil
}

// checkEndEntity returns true if current position is ENDENTITY or END ENTITY.
func (p *Parser) checkEndEntity() bool {
	if p.isAtEnd() {
		return false
	}
	if p.check(lexer.TokenEndEntity) {
		return true
	}
	return p.check(lexer.TokenEnd) && p.current+1 < len(p.tokens) && p.tokens[p.current+1].Type == lexer.TokenEntity
}

// entityDeclStatement parses ENTITY Name ... END ENTITY (or ENDENTITY).
func (p *Parser) entityDeclStatement() (Node, error) {
	p.advance() // Skip ENTITY
	if !p.match(lexer.TokenIdentifier) {
		return nil, &Error{Message: "expected entity name after ENTITY", Line: p.line(), Col: p.col()}
	}
	entityName := p.previous().Value
	var props []EntityProperty
	for {
		for p.match(lexer.TokenNewLine) {
		}
		if p.match(lexer.TokenEndEntity) {
			break
		}
		if p.check(lexer.TokenEnd) {
			p.advance()
			if p.match(lexer.TokenEntity) {
				break
			}
			return nil, &Error{Message: "expected ENTITY after END", Line: p.line(), Col: p.col()}
		}
		if p.isAtEnd() {
			return nil, &Error{Message: "expected END ENTITY or ENDENTITY", Line: p.line(), Col: p.col()}
		}
		if !p.match(lexer.TokenIdentifier) {
			return nil, &Error{Message: "expected property name or END ENTITY", Line: p.line(), Col: p.col()}
		}
		propName := p.previous().Value
		if !p.match(lexer.TokenAssign) {
			return nil, &Error{Message: "expected '=' after property name", Line: p.line(), Col: p.col()}
		}
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		props = append(props, EntityProperty{Name: propName, Value: val})
	}
	return &EntityDecl{Name: entityName, Properties: props}, nil
}

// dimStatement parses DIM statement (scalar or array: DIM x AS Integer, DIM a(10,20) AS Integer)
func (p *Parser) dimStatement() (Node, error) {
	p.advance() // Skip DIM

	var variables []VariableDecl

	for {
		if !p.match(lexer.TokenIdentifier) {
			return nil, &Error{Message: "expected variable name", Line: p.line(), Col: p.col()}
		}
		name := p.previous().Value

		var dimensions []Node
		if p.match(lexer.TokenLeftParen) {
			for !p.check(lexer.TokenRightParen) {
				dim, err := p.expression()
				if err != nil {
					return nil, err
				}
				dimensions = append(dimensions, dim)
				if !p.match(lexer.TokenComma) {
					break
				}
			}
			if !p.match(lexer.TokenRightParen) {
				return nil, &Error{Message: "expected ')' after array dimensions", Line: p.line(), Col: p.col()}
			}
		}

		varType := ""
		if p.match(lexer.TokenAs) {
			if !p.match(lexer.TokenIdentifier, lexer.TokenInteger, lexer.TokenStringType, lexer.TokenFloat, lexer.TokenBoolean) {
				return nil, &Error{Message: "expected type after AS", Line: p.line(), Col: p.col()}
			}
			varType = p.previous().Value
		}

		variables = append(variables, VariableDecl{Name: name, Type: varType, Dimensions: dimensions})

		if !p.match(lexer.TokenComma) {
			break
		}
	}

	return &DimStatement{Variables: variables}, nil
}

// printStatement parses PRINT statement
func (p *Parser) printStatement() (Node, error) {
	p.advance() // Skip PRINT

	var arguments []Node

	// Parse arguments until newline or EOF
	for !p.isAtEnd() && !p.check(lexer.TokenNewLine) && !p.check(lexer.TokenEOF) {
		arg, err := p.expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)

		// Break if next token is not a comma (for multiple arguments)
		if !p.match(lexer.TokenComma) {
			break
		}
	}

	return &Call{
		Name:      "print",
		Arguments: arguments,
		Line:      p.line(),
		Col:       p.col(),
	}, nil
}

// strStatement parses STR statement
func (p *Parser) strStatement() (Node, error) {
	p.advance() // Skip STR

	var arguments []Node

	// Parse arguments until newline or EOF
	for !p.isAtEnd() && !p.check(lexer.TokenNewLine) && !p.check(lexer.TokenEOF) {
		arg, err := p.expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)

		// Break if next token is not a comma (for multiple arguments)
		if !p.match(lexer.TokenComma) {
			break
		}
	}

	return &Call{
		Name:      "str",
		Arguments: arguments,
		Line:      p.line(),
		Col:       p.col(),
	}, nil
}

// sleepStatement parses SLEEP ms or SLEEP(ms)
func (p *Parser) sleepStatement() (Node, error) {
	p.advance() // Skip SLEEP
	arg, err := p.expression()
	if err != nil {
		return nil, err
	}
	return &Call{Name: "sleep", Arguments: []Node{arg}, Line: p.line(), Col: p.col()}, nil
}

// waitStatement parses WAIT ms or WAIT(ms)
func (p *Parser) waitStatement() (Node, error) {
	p.advance() // Skip WAIT
	arg, err := p.expression()
	if err != nil {
		return nil, err
	}
	return &Call{Name: "wait", Arguments: []Node{arg}, Line: p.line(), Col: p.col()}, nil
}

// returnStatement parses RETURN statement
func (p *Parser) returnStatement() (Node, error) {
	p.advance() // Skip RETURN

	var value Node
	if !p.check(lexer.TokenNewLine) && !p.isAtEnd() {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	return &ReturnStatement{Value: value}, nil
}

// gameCommand parses game-specific commands
func (p *Parser) gameCommand() (Node, error) {
	command := p.advance().Value

	var arguments []Node

	// Check for arguments in parentheses
	if p.match(lexer.TokenLeftParen) {
		for !p.check(lexer.TokenRightParen) {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			if !p.match(lexer.TokenComma) {
				break
			}
		}
		if !p.match(lexer.TokenRightParen) {
			return nil, &Error{Message: "expected ')'", Line: p.line(), Col: p.col()}
		}
	} else {
		// SYNC takes no arguments
		if strings.EqualFold(command, "sync") {
			return &GameCommand{Command: "sync", Arguments: nil}, nil
		}
		// Handle arguments without parentheses (like LOADIMAGE "filename")
		// Keep consuming tokens until we hit a newline, EOF, or another command
		for !p.isAtEnd() && !p.check(lexer.TokenNewLine) && !p.check(lexer.TokenEOF) {
			// Stop if we encounter another game command keyword
			nextToken := p.peek()
			switch nextToken.Type {
			case lexer.TokenLoadImage, lexer.TokenCreateSprite, lexer.TokenSetSpritePosition,
				lexer.TokenDrawSprite, lexer.TokenLoadModel, lexer.TokenCreateCamera,
				lexer.TokenSetCameraPosition, lexer.TokenDrawModel, lexer.TokenPlayMusic,
				lexer.TokenPlaySound, lexer.TokenLoadSound, lexer.TokenCreatePhysicsBody,
				lexer.TokenSetVelocity, lexer.TokenApplyForce, lexer.TokenRayCast3D,
				lexer.TokenSync, lexer.TokenShouldClose,
				lexer.TokenIf, lexer.TokenFor, lexer.TokenWhile, lexer.TokenFunction,
				lexer.TokenSub, lexer.TokenDim, lexer.TokenReturn:
				break
			}

			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			// Break if next token is a comma (for multiple arguments)
			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}

	return &GameCommand{
		Command:   command,
		Arguments: arguments,
	}, nil
}

// readQualifiedName reads identifier and optional . identifier/keyword (e.g. RL.InitWindow, BULLET.Step)
func (p *Parser) readQualifiedName() string {
	name := p.previous().Value
	for p.match(lexer.TokenDot) {
		if p.isAtEnd() {
			return name
		}
		p.advance() // accept identifier or keyword (e.g. Step) as part of qualified name
		name += "." + p.previous().Value
	}
	return name
}

// assignmentOrCall parses assignment or function call
func (p *Parser) assignmentOrCall() (Node, error) {
	if !p.match(lexer.TokenIdentifier) {
		return nil, &Error{Message: "expected identifier", Line: p.line(), Col: p.col()}
	}
	name := p.readQualifiedName()

	if p.match(lexer.TokenPlusAssign) {
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		return &CompoundAssign{Variable: name, Op: "+=", Value: value, Line: p.line(), Col: p.col()}, nil
	}
	if p.match(lexer.TokenMinusAssign) {
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		return &CompoundAssign{Variable: name, Op: "-=", Value: value, Line: p.line(), Col: p.col()}, nil
	}
	if p.match(lexer.TokenStarAssign) {
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		return &CompoundAssign{Variable: name, Op: "*=", Value: value, Line: p.line(), Col: p.col()}, nil
	}
	if p.match(lexer.TokenSlashAssign) {
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		return &CompoundAssign{Variable: name, Op: "/=", Value: value, Line: p.line(), Col: p.col()}, nil
	}
	if p.match(lexer.TokenAssign) {
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		return &Assignment{Variable: name, Value: value, Line: p.line(), Col: p.col()}, nil
	}
	if p.match(lexer.TokenLeftParen) {
		var arguments []Node
		for !p.check(lexer.TokenRightParen) {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			if !p.match(lexer.TokenComma) {
				break
			}
		}
		if !p.match(lexer.TokenRightParen) {
			return nil, &Error{Message: "expected ')'", Line: p.line(), Col: p.col()}
		}
		// Array assignment: a(i,j) = value
		if p.match(lexer.TokenAssign) {
			value, err := p.expression()
			if err != nil {
				return nil, err
			}
			return &Assignment{Variable: name, Indices: arguments, Value: value, Line: p.line(), Col: p.col()}, nil
		}
		return &Call{Name: name, Arguments: arguments, Line: p.line(), Col: p.col()}, nil
	} else {
		return &Identifier{Name: name, Line: p.line(), Col: p.col()}, nil
	}
}

// block parses a block of statements. When stopAtNewlineIf is true (IF then-block only),
// a single statement followed by newline and IF causes an early return so the next IF is not consumed.
func (p *Parser) block(stopAtNewlineIf bool) (*Block, error) {
	block := &Block{}

	for !p.isAtEnd() && !p.check(lexer.TokenEndIf) && !p.check(lexer.TokenNext) &&
		!p.check(lexer.TokenWend) && !p.check(lexer.TokenEnd) &&
		!p.check(lexer.TokenElseIf) && !p.check(lexer.TokenElse) &&
		!p.check(lexer.TokenEndFunction) && !p.check(lexer.TokenEndSub) && !p.check(lexer.TokenEndModule) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
			if stopAtNewlineIf && len(block.Statements) == 1 {
				for p.match(lexer.TokenNewLine) {
				}
				if p.check(lexer.TokenIf) {
					return block, nil
				}
			}
		}
		if p.match(lexer.TokenNewLine) {
			continue
		}
	}

	return block, nil
}

// blockUntilSelectEnd parses statements until CASE or END SELECT (does not consume them)
func (p *Parser) blockUntilSelectEnd() (*Block, error) {
	block := &Block{}
	for !p.isAtEnd() && !p.check(lexer.TokenCase) && !p.check(lexer.TokenEndSelect) && !p.check(lexer.TokenEnd) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		if p.match(lexer.TokenNewLine) {
			continue
		}
	}
	return block, nil
}

// selectCaseStatement parses SELECT CASE expr ... CASE value: block ... CASE ELSE: block ... END SELECT
func (p *Parser) selectCaseStatement() (Node, error) {
	p.advance() // SELECT
	if !p.match(lexer.TokenCase) {
		return nil, &Error{Message: "expected CASE after SELECT", Line: p.line(), Col: p.col()}
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	var cases []CaseClause
	var elseBlock *Block
	for !p.check(lexer.TokenEnd) && !p.check(lexer.TokenEndSelect) {
		for p.match(lexer.TokenNewLine) {
		}
		if p.check(lexer.TokenEnd) || p.check(lexer.TokenEndSelect) {
			break
		}
		if !p.match(lexer.TokenCase) {
			return nil, &Error{Message: "expected CASE or END SELECT", Line: p.line(), Col: p.col()}
		}
		if p.match(lexer.TokenElse) {
			elseBlock, err = p.blockUntilSelectEnd()
			if err != nil {
				return nil, err
			}
			break
		}
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		blk, err := p.blockUntilSelectEnd()
		if err != nil {
			return nil, err
		}
		cases = append(cases, CaseClause{Value: val, Block: blk})
	}
	// END SELECT (one token ENDSELECT or two tokens END SELECT)
	if p.match(lexer.TokenEndSelect) {
		// single token
	} else if p.match(lexer.TokenEnd) {
		if !p.match(lexer.TokenSelect) {
			return nil, &Error{Message: "expected SELECT after END", Line: p.line(), Col: p.col()}
		}
	} else {
		return nil, &Error{Message: "expected END SELECT", Line: p.line(), Col: p.col()}
	}
	return &SelectCaseStatement{Expr: expr, Cases: cases, ElseBlock: elseBlock}, nil
}

// quitStatement parses QUIT or END (exit program)
func (p *Parser) quitStatement() (Node, error) {
	p.advance() // QUIT
	return &Call{Name: "quit", Arguments: nil, Line: p.line(), Col: p.col()}, nil
}

// blockUntilUntil parses statements until UNTIL (does not consume UNTIL)
func (p *Parser) blockUntilUntil() (*Block, error) {
	block := &Block{}
	for !p.isAtEnd() && !p.check(lexer.TokenUntil) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		if p.match(lexer.TokenNewLine) {
			continue
		}
	}
	return block, nil
}

// repeatStatement parses REPEAT ... UNTIL condition
func (p *Parser) repeatStatement() (Node, error) {
	p.advance() // REPEAT
	body, err := p.blockUntilUntil()
	if err != nil {
		return nil, err
	}
	if !p.match(lexer.TokenUntil) {
		return nil, &Error{Message: "expected UNTIL after REPEAT block", Line: p.line(), Col: p.col()}
	}
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}
	return &RepeatStatement{Body: body, Condition: cond}, nil
}

// exitLoopStatement parses EXIT FOR, EXIT WHILE, BREAK FOR, or BREAK WHILE.
func (p *Parser) exitLoopStatement() (Node, error) {
	p.advance() // EXIT or BREAK
	if p.match(lexer.TokenFor) {
		return &ExitLoopStatement{Kind: "FOR"}, nil
	}
	if p.match(lexer.TokenWhile) {
		return &ExitLoopStatement{Kind: "WHILE"}, nil
	}
	return nil, &Error{Message: "expected FOR or WHILE after EXIT/BREAK", Line: p.line(), Col: p.col()}
}

// continueLoopStatement parses CONTINUE FOR or CONTINUE WHILE.
func (p *Parser) continueLoopStatement() (Node, error) {
	p.advance() // CONTINUE
	if p.match(lexer.TokenFor) {
		return &ContinueLoopStatement{Kind: "FOR"}, nil
	}
	if p.match(lexer.TokenWhile) {
		return &ContinueLoopStatement{Kind: "WHILE"}, nil
	}
	return nil, &Error{Message: "expected FOR or WHILE after CONTINUE", Line: p.line(), Col: p.col()}
}

// assertStatement parses ASSERT condition [, message].
func (p *Parser) assertStatement() (Node, error) {
	p.advance() // ASSERT
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}
	var msg Node
	if p.match(lexer.TokenComma) {
		msg, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	return &AssertStatement{Condition: cond, Message: msg}, nil
}
