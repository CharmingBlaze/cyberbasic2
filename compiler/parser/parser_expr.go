package parser

import (
	"cyberbasic/compiler/lexer"
	"strings"
)

// expression parses an expression
func (p *Parser) expression() (Node, error) {
	return p.logicalOr()
}

// logicalOr parses OR operations
func (p *Parser) logicalOr() (Node, error) {
	left, err := p.logicalXor()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.TokenOr) {
		op := p.previous().Value
		right, err := p.logicalXor()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Operator: op, Left: left, Right: right, Line: p.line(), Col: p.col()}
	}

	return left, nil
}

// logicalXor parses XOR operations
func (p *Parser) logicalXor() (Node, error) {
	left, err := p.logicalAnd()
	if err != nil {
		return nil, err
	}
	for p.match(lexer.TokenXor) {
		right, err := p.logicalAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Operator: "XOR", Left: left, Right: right, Line: p.line(), Col: p.col()}
	}
	return left, nil
}

// logicalAnd parses AND operations
func (p *Parser) logicalAnd() (Node, error) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.TokenAnd) {
		op := p.previous().Value
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Operator: op, Left: left, Right: right, Line: p.line(), Col: p.col()}
	}

	return left, nil
}

// equality parses equality operations
func (p *Parser) equality() (Node, error) {
	left, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.TokenEqual, lexer.TokenNotEqual, lexer.TokenAssign) {
		op := p.previous().Value
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Operator: op, Left: left, Right: right, Line: p.line(), Col: p.col()}
	}

	return left, nil
}

// comparison parses comparison operations
func (p *Parser) comparison() (Node, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.TokenGreater, lexer.TokenGreaterEqual, lexer.TokenLess, lexer.TokenLessEqual) {
		op := p.previous().Value
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Operator: op, Left: left, Right: right, Line: p.line(), Col: p.col()}
	}

	return left, nil
}

// term parses addition and subtraction
func (p *Parser) term() (Node, error) {
	left, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.TokenPlus, lexer.TokenMinus) {
		op := p.previous().Value
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Operator: op, Left: left, Right: right, Line: p.line(), Col: p.col()}
	}

	return left, nil
}

// factor parses multiplication, division, modulo, and integer division
func (p *Parser) factor() (Node, error) {
	left, err := p.power()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.TokenMultiply, lexer.TokenDivide, lexer.TokenMod, lexer.TokenIntDiv) {
		op := p.previous().Value
		right, err := p.power()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Operator: op, Left: left, Right: right, Line: p.line(), Col: p.col()}
	}

	return left, nil
}

// power parses exponentiation (right-associative): unary (^ power)*
func (p *Parser) power() (Node, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(lexer.TokenPower) {
		right, err := p.power()
		if err != nil {
			return nil, err
		}
		left = &BinaryOp{Operator: "^", Left: left, Right: right, Line: p.line(), Col: p.col()}
	}
	return left, nil
}

// unary parses unary operations
func (p *Parser) unary() (Node, error) {
	if p.match(lexer.TokenNot, lexer.TokenMinus) {
		op := p.previous().Value
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryOp{Operator: op, Operand: right, Line: p.line(), Col: p.col()}, nil
	}
	// Treat identifier "NOT" as unary operator when lexer returns TokenIdentifier
	if !p.isAtEnd() && p.peek().Type == lexer.TokenIdentifier && strings.EqualFold(p.peek().Value, "NOT") {
		p.advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryOp{Operator: "NOT", Operand: right, Line: p.line(), Col: p.col()}, nil
	}

	return p.primary()
}

// dictLiteral parses { "key": value } or { key = value } (JSON-style or BASIC-style).
func (p *Parser) dictLiteral() (Node, error) {
	p.advance() // {
	var pairs []DictPair
	for !p.check(lexer.TokenRightBrace) && !p.isAtEnd() {
		for p.match(lexer.TokenNewLine) {
			continue
		}
		var key string
		if p.match(lexer.TokenString) {
			key = p.previous().Value
			if !p.match(lexer.TokenColon) {
				return nil, &Error{Message: "expected ':' after string key in dict", Line: p.line(), Col: p.col()}
			}
		} else if p.match(lexer.TokenIdentifier) {
			key = p.previous().Value
			if !p.match(lexer.TokenAssign) {
				return nil, &Error{Message: "expected '=' after identifier key in dict", Line: p.line(), Col: p.col()}
			}
		} else if p.match(lexer.TokenNumber) {
			key = p.previous().Value
			if !p.match(lexer.TokenColon) {
				return nil, &Error{Message: "expected ':' after number key in dict", Line: p.line(), Col: p.col()}
			}
		} else {
			return nil, &Error{Message: "expected string, identifier, or number key in dict", Line: p.line(), Col: p.col()}
		}
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, DictPair{Key: key, Value: value})
		for p.match(lexer.TokenNewLine) {
			continue
		}
		if !p.match(lexer.TokenComma) {
			break
		}
	}
	if !p.match(lexer.TokenRightBrace) {
		return nil, &Error{Message: "expected '}' to close dict literal", Line: p.line(), Col: p.col()}
	}
	return &DictLiteral{Pairs: pairs}, nil
}

// primary parses primary expressions
func (p *Parser) primary() (Node, error) {
	switch p.peek().Type {
	case lexer.TokenNumber:
		return &Number{Value: p.advance().Value}, nil
	case lexer.TokenString:
		return &StringLiteral{Value: p.advance().Value}, nil
	case lexer.TokenTrue:
		p.advance()
		return &Boolean{Value: true}, nil
	case lexer.TokenFalse:
		p.advance()
		return &Boolean{Value: false}, nil
	case lexer.TokenNil:
		p.advance()
		return &NilLiteral{}, nil
	case lexer.TokenLeftBrace:
		return p.dictLiteral()
	case lexer.TokenShouldClose:
		p.advance()
		if p.match(lexer.TokenLeftParen) {
			if !p.match(lexer.TokenRightParen) {
				return nil, &Error{Message: "expected ')'", Line: p.line(), Col: p.col()}
			}
		}
		return &Call{Name: "shouldclose", Arguments: nil, Line: p.line(), Col: p.col()}, nil
	case lexer.TokenStr:
		p.advance()
		if !p.match(lexer.TokenLeftParen) {
			return nil, &Error{Message: "STR in expression requires (", Line: p.line(), Col: p.col()}
		}
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
		return &Call{Name: "str", Arguments: arguments, Line: p.line(), Col: p.col()}, nil
	}

	if p.match(lexer.TokenIdentifier) {
		prev := p.previous()
		left, err := p.parseMemberAccessChain(&Identifier{Name: prev.Value, Line: prev.Line, Col: prev.Col})
		if err != nil {
			return nil, err
		}
		// Optional JSON index sugar: id["key"] or id["key"]["key2"] ...
		for p.match(lexer.TokenLeftBracket) {
			if !p.match(lexer.TokenString) {
				return nil, &Error{Message: "expected string key in [ ] for JSON index", Line: p.line(), Col: p.col()}
			}
			key := p.previous().Value
			if !p.match(lexer.TokenRightBracket) {
				return nil, &Error{Message: "expected ']' after JSON key", Line: p.line(), Col: p.col()}
			}
			left = &JSONIndexAccess{Object: left, Key: key}
		}
		// MemberAccess followed by ( is a call (e.g. RL.InitWindow(...))
		if ma, ok := left.(*MemberAccess); ok && p.match(lexer.TokenLeftParen) {
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
			callName := memberAccessToQualifiedName(ma)
			return p.parseMemberAccessChain(&Call{Name: callName, Arguments: arguments})
		}
		// Bare identifier followed by ( is a call (e.g. Sin(x), Sqrt(y))
		if id, ok := left.(*Identifier); ok && p.match(lexer.TokenLeftParen) {
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
			return p.parseMemberAccessChain(&Call{Name: id.Name, Arguments: arguments})
		}
		return left, nil
	}

	if p.match(lexer.TokenLeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(lexer.TokenRightParen) {
			return nil, &Error{Message: "expected ')'", Line: p.line(), Col: p.col()}
		}
		return p.parseMemberAccessChain(expr)
	}

	return nil, &Error{Message: "unexpected token", Line: p.line(), Col: p.col()}
}

// parseMemberAccessChain parses optional . member . member ... after a primary expression.
func (p *Parser) parseMemberAccessChain(left Node) (Node, error) {
	for p.match(lexer.TokenDot) {
		if !p.match(lexer.TokenIdentifier) {
			return nil, &Error{Message: "expected identifier after '.'", Line: p.line(), Col: p.col()}
		}
		member := p.previous().Value
		left = &MemberAccess{Object: left, Member: member}
	}
	return left, nil
}

// memberAccessToQualifiedName returns "obj.member" or "a.b.c" for Call name (lowercase).
func memberAccessToQualifiedName(ma *MemberAccess) string {
	var prefix string
	switch obj := ma.Object.(type) {
	case *Identifier:
		prefix = obj.Name
	case *MemberAccess:
		prefix = memberAccessToQualifiedName(obj)
	default:
		prefix = "unknown"
	}
	return strings.ToLower(prefix) + "." + strings.ToLower(ma.Member)
}
