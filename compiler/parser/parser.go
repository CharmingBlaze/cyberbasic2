package parser

import (
	"cyberbasic/compiler/lexer"
)

// Parser builds an AST from tokens
type Parser struct {
	tokens  []lexer.Token
	current int
	linePos int
	colPos  int
}

// New creates a new parser instance
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

// Parse builds the AST
func (p *Parser) Parse() (*Program, error) {
	program := &Program{}

	for !p.isAtEnd() {
		// Skip newlines between statements
		for p.match(lexer.TokenNewLine) {
			continue
		}

		if p.isAtEnd() {
			break
		}

		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program, nil
}

// isBlockTerminator returns true if the current token ends a block (do not parse as statement).
func (p *Parser) isBlockTerminator() bool {
	if p.isAtEnd() {
		return true
	}
	switch p.peek().Type {
	case lexer.TokenEndIf, lexer.TokenWend, lexer.TokenNext, lexer.TokenUntil,
		lexer.TokenEndSelect, lexer.TokenEnd, lexer.TokenCase,
		lexer.TokenElseIf, lexer.TokenElse,
		lexer.TokenEndFunction, lexer.TokenEndSub, lexer.TokenEndModule:
		return true
	default:
		return false
	}
}

// statement parses a statement
func (p *Parser) statement() (Node, error) {
	if p.match(lexer.TokenNewLine) {
		return nil, nil
	}
	if p.isBlockTerminator() {
		return nil, nil
	}

	switch p.peek().Type {
	case lexer.TokenIf:
		return p.ifStatement()
	case lexer.TokenFor:
		return p.forStatement()
	case lexer.TokenWhile:
		return p.whileStatement()
	case lexer.TokenFunction:
		return p.functionDecl()
	case lexer.TokenSub:
		return p.subDecl()
	case lexer.TokenModule:
		return p.moduleDecl()
	case lexer.TokenOn:
		return p.onEventStatement()
	case lexer.TokenStartCoroutine:
		return p.startCoroutineStatement()
	case lexer.TokenYield:
		return p.yieldStatement()
	case lexer.TokenWaitSeconds:
		return p.waitSecondsStatement()
	case lexer.TokenDim:
		return p.dimStatement()
	case lexer.TokenConst:
		return p.constStatement()
	case lexer.TokenEnum:
		return p.enumStatement()
	case lexer.TokenTypeKw:
		return p.typeDeclStatement()
	case lexer.TokenEntity:
		return p.entityDeclStatement()
	case lexer.TokenPrint:
		return p.printStatement()
	case lexer.TokenStr:
		return p.strStatement()
	case lexer.TokenSleep:
		return p.sleepStatement()
	case lexer.TokenWait:
		return p.waitStatement()
	case lexer.TokenSelect:
		return p.selectCaseStatement()
	case lexer.TokenQuit:
		return p.quitStatement()
	case lexer.TokenRepeat:
		return p.repeatStatement()
	case lexer.TokenExit, lexer.TokenBreak:
		return p.exitLoopStatement()
	case lexer.TokenContinue:
		return p.continueLoopStatement()
	case lexer.TokenAssert:
		return p.assertStatement()
	case lexer.TokenLet:
		p.advance() // skip LET
		return p.assignmentOrCall()
	case lexer.TokenVar:
		p.advance() // skip VAR (same as LET)
		return p.assignmentOrCall()
	case lexer.TokenReturn:
		return p.returnStatement()
	case lexer.TokenLoadImage, lexer.TokenCreateSprite, lexer.TokenSetSpritePosition,
		lexer.TokenDrawSprite, lexer.TokenLoadModel, lexer.TokenCreateCamera,
		lexer.TokenSetCameraPosition, lexer.TokenDrawModel, lexer.TokenPlayMusic,
		lexer.TokenPlaySound, lexer.TokenLoadSound, lexer.TokenCreatePhysicsBody,
		lexer.TokenSetVelocity, lexer.TokenApplyForce, lexer.TokenRayCast3D,
		lexer.TokenSync, lexer.TokenShouldClose:
		return p.gameCommand()
	default:
		return p.assignmentOrCall()
	}
}

// Helper methods
func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

// checkEndIf returns true if the current position is ENDIF or END IF (two words).
func (p *Parser) checkEndIf() bool {
	if p.isAtEnd() {
		return false
	}
	if p.check(lexer.TokenEndIf) {
		return true
	}
	if p.check(lexer.TokenEnd) && p.current+1 < len(p.tokens) && p.tokens[p.current+1].Type == lexer.TokenIf {
		return true
	}
	return false
}

// matchEndIf consumes ENDIF or END IF (two words); returns true if matched.
func (p *Parser) matchEndIf() bool {
	if p.match(lexer.TokenEndIf) {
		return true
	}
	if p.check(lexer.TokenEnd) && p.current+1 < len(p.tokens) && p.tokens[p.current+1].Type == lexer.TokenIf {
		p.advance() // END
		p.advance() // IF
		return true
	}
	return false
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.TokenEOF
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) line() int {
	if p.current > 0 {
		return p.tokens[p.current-1].Line
	}
	return 0
}

func (p *Parser) col() int {
	if p.current > 0 {
		return p.tokens[p.current-1].Col
	}
	return 0
}
