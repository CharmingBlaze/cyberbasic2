package lexer

// TokenType represents different types of tokens in BASIC
type TokenType int

// String returns the string representation of TokenType
func (t TokenType) String() string {
	switch t {
	case TokenNumber:
		return "TokenNumber"
	case TokenString:
		return "TokenString"
	case TokenIdentifier:
		return "TokenIdentifier"
	case TokenIf:
		return "TokenIf"
	case TokenThen:
		return "TokenThen"
	case TokenElse:
		return "TokenElse"
	case TokenEndIf:
		return "TokenEndIf"
	case TokenFor:
		return "TokenFor"
	case TokenTo:
		return "TokenTo"
	case TokenStep:
		return "TokenStep"
	case TokenNext:
		return "TokenNext"
	case TokenWhile:
		return "TokenWhile"
	case TokenWend:
		return "TokenWend"
	case TokenFunction:
		return "TokenFunction"
	case TokenSub:
		return "TokenSub"
	case TokenModule:
		return "TokenModule"
	case TokenEnd:
		return "TokenEnd"
	case TokenReturn:
		return "TokenReturn"
	case TokenDim:
		return "TokenDim"
	case TokenAs:
		return "TokenAs"
	case TokenInteger:
		return "TokenInteger"
	case TokenStringType:
		return "TokenStringType"
	case TokenFloat:
		return "TokenFloat"
	case TokenBoolean:
		return "TokenBoolean"
	case TokenPrint:
		return "TokenPrint"
	case TokenStr:
		return "TokenStr"
	case TokenTrue:
		return "TokenTrue"
	case TokenFalse:
		return "TokenFalse"
	case TokenNil:
		return "TokenNil"
	case TokenInitGraphics3D:
		return "TokenInitGraphics3D"
	case TokenBegin3DMode:
		return "TokenBegin3DMode"
	case TokenEnd3DMode:
		return "TokenEnd3DMode"
	case TokenDrawModel3D:
		return "TokenDrawModel3D"
	case TokenDrawGrid3D:
		return "TokenDrawGrid3D"
	case TokenDrawAxes3D:
		return "TokenDrawAxes3D"
	case TokenCreatePhysicsWorld2D:
		return "TokenCreatePhysicsWorld2D"
	case TokenDestroyPhysicsWorld2D:
		return "TokenDestroyPhysicsWorld2D"
	case TokenStepPhysics2D:
		return "TokenStepPhysics2D"
	case TokenCreatePhysicsBody2D:
		return "TokenCreatePhysicsBody2D"
	case TokenDestroyPhysicsBody2D:
		return "TokenDestroyPhysicsBody2D"
	case TokenSetPhysicsPosition2D:
		return "TokenSetPhysicsPosition2D"
	case TokenGetPhysicsPosition2D:
		return "TokenGetPhysicsPosition2D"
	case TokenSetPhysicsAngle2D:
		return "TokenSetPhysicsAngle2D"
	case TokenGetPhysicsAngle2D:
		return "TokenGetPhysicsAngle2D"
	case TokenSetPhysicsVelocity2D:
		return "TokenSetPhysicsVelocity2D"
	case TokenGetPhysicsVelocity2D:
		return "TokenGetPhysicsVelocity2D"
	case TokenApplyPhysicsForce2D:
		return "TokenApplyPhysicsForce2D"
	case TokenApplyPhysicsImpulse2D:
		return "TokenApplyPhysicsImpulse2D"
	case TokenSetPhysicsDensity2D:
		return "TokenSetPhysicsDensity2D"
	case TokenSetPhysicsFriction2D:
		return "TokenSetPhysicsFriction2D"
	case TokenSetPhysicsRestitution2D:
		return "TokenSetPhysicsRestitution2D"
	case TokenRayCast2D:
		return "TokenRayCast2D"
	case TokenCheckCollision2D:
		return "TokenCheckCollision2D"
	case TokenQueryAABB2D:
		return "TokenQueryAABB2D"
	case TokenCreatePhysicsWorld3D:
		return "TokenCreatePhysicsWorld3D"
	case TokenDestroyPhysicsWorld3D:
		return "TokenDestroyPhysicsWorld3D"
	case TokenStepPhysics3D:
		return "TokenStepPhysics3D"
	case TokenCreatePhysicsBody3D:
		return "TokenCreatePhysicsBody3D"
	case TokenDestroyPhysicsBody3D:
		return "TokenDestroyPhysicsBody3D"
	case TokenSetPhysicsPosition3D:
		return "TokenSetPhysicsPosition3D"
	case TokenGetPhysicsPosition3D:
		return "TokenGetPhysicsPosition3D"
	case TokenSetPhysicsRotation3D:
		return "TokenSetPhysicsRotation3D"
	case TokenGetPhysicsRotation3D:
		return "TokenGetPhysicsRotation3D"
	case TokenSetPhysicsVelocity3D:
		return "TokenSetPhysicsVelocity3D"
	case TokenGetPhysicsVelocity3D:
		return "TokenGetPhysicsVelocity3D"
	case TokenApplyPhysicsForce3D:
		return "TokenApplyPhysicsForce3D"
	case TokenApplyPhysicsImpulse3D:
		return "TokenApplyPhysicsImpulse3D"
	case TokenSetPhysicsMass3D:
		return "TokenSetPhysicsMass3D"
	case TokenCheckCollision3D:
		return "TokenCheckCollision3D"
	case TokenQueryAABB3D:
		return "TokenQueryAABB3D"
	case TokenLoadImage:
		return "TokenLoadImage"
	case TokenCreateSprite:
		return "TokenCreateSprite"
	case TokenSetSpritePosition:
		return "TokenSetSpritePosition"
	case TokenDrawSprite:
		return "TokenDrawSprite"
	case TokenLoadModel:
		return "TokenLoadModel"
	case TokenCreateCamera:
		return "TokenCreateCamera"
	case TokenSetCameraPosition:
		return "TokenSetCameraPosition"
	case TokenDrawModel:
		return "TokenDrawModel"
	case TokenPlayMusic:
		return "TokenPlayMusic"
	case TokenPlaySound:
		return "TokenPlaySound"
	case TokenLoadSound:
		return "TokenLoadSound"
	case TokenCreatePhysicsBody:
		return "TokenCreatePhysicsBody"
	case TokenSetVelocity:
		return "TokenSetVelocity"
	case TokenApplyForce:
		return "TokenApplyForce"
	case TokenRayCast3D:
		return "TokenRayCast3D"
	case TokenSync:
		return "TokenSync"
	case TokenShouldClose:
		return "TokenShouldClose"
	case TokenSleep:
		return "TokenSleep"
	case TokenWait:
		return "TokenWait"
	case TokenSelect:
		return "TokenSelect"
	case TokenCase:
		return "TokenCase"
	case TokenEndSelect:
		return "TokenEndSelect"
	case TokenQuit:
		return "TokenQuit"
	case TokenLet:
		return "TokenLet"
	case TokenVar:
		return "TokenVar"
	case TokenRepeat:
		return "TokenRepeat"
	case TokenUntil:
		return "TokenUntil"
	case TokenConst:
		return "TokenConst"
	case TokenExit:
		return "TokenExit"
	case TokenEnum:
		return "TokenEnum"
	case TokenTypeKw:
		return "TokenTypeKw"
	case TokenEndType:
		return "TokenEndType"
	case TokenOn:
		return "TokenOn"
	case TokenEndOn:
		return "TokenEndOn"
	case TokenKeyDown:
		return "TokenKeyDown"
	case TokenKeyPressed:
		return "TokenKeyPressed"
	case TokenStartCoroutine:
		return "TokenStartCoroutine"
	case TokenYield:
		return "TokenYield"
	case TokenWaitSeconds:
		return "TokenWaitSeconds"
	case TokenMain:
		return "TokenMain"
	case TokenEndMain:
		return "TokenEndMain"
	case TokenEqual:
		return "TokenEqual"
	case TokenNotEqual:
		return "TokenNotEqual"
	case TokenLess:
		return "TokenLess"
	case TokenLessEqual:
		return "TokenLessEqual"
	case TokenGreater:
		return "TokenGreater"
	case TokenGreaterEqual:
		return "TokenGreaterEqual"
	case TokenPlus:
		return "TokenPlus"
	case TokenMinus:
		return "TokenMinus"
	case TokenMultiply:
		return "TokenMultiply"
	case TokenDivide:
		return "TokenDivide"
	case TokenMod:
		return "TokenMod"
	case TokenAnd:
		return "TokenAnd"
	case TokenOr:
		return "TokenOr"
	case TokenNot:
		return "TokenNot"
	case TokenAssign:
		return "TokenAssign"
	case TokenPlusAssign:
		return "TokenPlusAssign"
	case TokenMinusAssign:
		return "TokenMinusAssign"
	case TokenStarAssign:
		return "TokenStarAssign"
	case TokenSlashAssign:
		return "TokenSlashAssign"
	case TokenLeftParen:
		return "TokenLeftParen"
	case TokenRightParen:
		return "TokenRightParen"
	case TokenLeftBracket:
		return "TokenLeftBracket"
	case TokenRightBracket:
		return "TokenRightBracket"
	case TokenComma:
		return "TokenComma"
	case TokenColon:
		return "TokenColon"
	case TokenSemicolon:
		return "TokenSemicolon"
	case TokenDot:
		return "TokenDot"
	case TokenNewLine:
		return "TokenNewLine"
	case TokenEOF:
		return "TokenEOF"
	case TokenUnknown:
		return "TokenUnknown"
	default:
		return "UnknownToken"
	}
}

const (
	// Literals
	TokenNumber TokenType = iota
	TokenString
	TokenIdentifier

	// Keywords
	TokenIf
	TokenThen
	TokenElse
	TokenEndIf
	TokenFor
	TokenTo
	TokenStep
	TokenNext
	TokenWhile
	TokenWend
	TokenFunction
	TokenSub
	TokenModule
	TokenEnd
	TokenReturn
	TokenDim
	TokenAs
	TokenInteger
	TokenStringType
	TokenFloat
	TokenBoolean
	TokenPrint
	TokenStr
	TokenTrue
	TokenFalse
	TokenNil
	TokenInitGraphics3D
	TokenBegin3DMode
	TokenEnd3DMode
	TokenDrawModel3D
	TokenDrawGrid3D
	TokenDrawAxes3D

	// 2D Physics keywords
	TokenCreatePhysicsWorld2D
	TokenDestroyPhysicsWorld2D
	TokenStepPhysics2D
	TokenCreatePhysicsBody2D
	TokenDestroyPhysicsBody2D
	TokenSetPhysicsPosition2D
	TokenGetPhysicsPosition2D
	TokenSetPhysicsAngle2D
	TokenGetPhysicsAngle2D
	TokenSetPhysicsVelocity2D
	TokenGetPhysicsVelocity2D
	TokenApplyPhysicsForce2D
	TokenApplyPhysicsImpulse2D
	TokenSetPhysicsDensity2D
	TokenSetPhysicsFriction2D
	TokenSetPhysicsRestitution2D
	TokenRayCast2D
	TokenCheckCollision2D
	TokenQueryAABB2D

	// 3D Physics keywords
	TokenCreatePhysicsWorld3D
	TokenDestroyPhysicsWorld3D
	TokenStepPhysics3D
	TokenCreatePhysicsBody3D
	TokenDestroyPhysicsBody3D
	TokenSetPhysicsPosition3D
	TokenGetPhysicsPosition3D
	TokenSetPhysicsRotation3D
	TokenGetPhysicsRotation3D
	TokenSetPhysicsVelocity3D
	TokenGetPhysicsVelocity3D
	TokenApplyPhysicsForce3D
	TokenApplyPhysicsImpulse3D
	TokenSetPhysicsMass3D
	TokenCheckCollision3D
	TokenQueryAABB3D

	// Game-specific keywords
	TokenLoadImage
	TokenCreateSprite
	TokenSetSpritePosition
	TokenDrawSprite
	TokenLoadModel
	TokenCreateCamera
	TokenSetCameraPosition
	TokenDrawModel
	TokenPlayMusic
	TokenPlaySound
	TokenLoadSound
	TokenCreatePhysicsBody
	TokenSetVelocity
	TokenApplyForce
	TokenRayCast3D
	TokenSync
	TokenShouldClose
	TokenSleep
	TokenWait
	TokenSelect
	TokenCase
	TokenEndSelect
	TokenQuit
	TokenLet
	TokenVar
	TokenRepeat
	TokenUntil
	TokenConst
	TokenExit
	TokenEnum
	TokenTypeKw   // TYPE keyword (UDT)
	TokenEndType  // ENDTYPE
	TokenOn       // ON (event)
	TokenEndOn    // ENDON
	TokenKeyDown  // KeyDown
	TokenKeyPressed
	TokenStartCoroutine
	TokenYield
	TokenWaitSeconds
	TokenMain
	TokenEndMain

	// Operators
	TokenEqual
	TokenNotEqual
	TokenLess
	TokenLessEqual
	TokenGreater
	TokenGreaterEqual
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenMod
	TokenAnd
	TokenOr
	TokenNot
	TokenAssign
	TokenPlusAssign  // +=
	TokenMinusAssign // -=
	TokenStarAssign  // *=
	TokenSlashAssign // /=

	// Delimiters
	TokenLeftParen
	TokenRightParen
	TokenLeftBracket
	TokenRightBracket
	TokenComma
	TokenColon
	TokenSemicolon
	TokenNewLine
	TokenDot // for qualified names: RL.InitWindow

	// Special
	TokenEOF
	TokenUnknown
)

// Token represents a single token in the source code
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

// KeywordMap maps keyword strings to token types
var KeywordMap = map[string]TokenType{
	"IF":                      TokenIf,
	"THEN":                    TokenThen,
	"ELSE":                    TokenElse,
	"ENDIF":                   TokenEndIf,
	"FOR":                     TokenFor,
	"TO":                      TokenTo,
	"STEP":                    TokenStep,
	"NEXT":                    TokenNext,
	"WHILE":                   TokenWhile,
	"WEND":                    TokenWend,
	"FUNCTION":                TokenFunction,
	"SUB":                     TokenSub,
	"MODULE":                  TokenModule,
	"END":                     TokenEnd,
	"RETURN":                  TokenReturn,
	"DIM":                     TokenDim,
	"AS":                      TokenAs,
	"INTEGER":                 TokenInteger,
	"STRING":                  TokenStringType,
	"FLOAT":                   TokenFloat,
	"BOOLEAN":                 TokenBoolean,
	"PRINT":                   TokenPrint,
	"STR":                     TokenStr,
	"TRUE":                    TokenTrue,
	"FALSE":                   TokenFalse,
	"NIL":                     TokenNil,
	"NULL":                    TokenNil,
	"LOADIMAGE":               TokenLoadImage,
	"CREATESPRITE":            TokenCreateSprite,
	"SETSPRITEPOSITION":       TokenSetSpritePosition,
	"DRAWSPRITE":              TokenDrawSprite,
	"LOADMODEL":               TokenLoadModel,
	"CREATECAMERA":            TokenCreateCamera,
	"SETCAMERAPOSITION":       TokenSetCameraPosition,
	"DRAWMODEL":               TokenDrawModel,
	"PLAYMUSIC":               TokenPlayMusic,
	"PLAYSOUND":               TokenPlaySound,
	"LOADSOUND":               TokenLoadSound,
	"CREATEPHYSICSBODY":       TokenCreatePhysicsBody,
	"SETVELOCITY":             TokenSetVelocity,
	"APPLYFORCE":              TokenApplyForce,
	"RAYCAST3D":               TokenRayCast3D,
	"SYNC":                    TokenSync,
	"SHOULDCLOSE":             TokenShouldClose,
	"SLEEP":                   TokenSleep,
	"WAIT":                    TokenWait,
	"AND":                     TokenAnd,
	"OR":                      TokenOr,
	"NOT":                     TokenNot,
	"SELECT":                  TokenSelect,
	"CASE":                    TokenCase,
	"ENDSELECT":               TokenEndSelect,
	"QUIT":                    TokenQuit,
	"LET":                     TokenLet,
	"VAR":                     TokenVar,
	"REPEAT":                  TokenRepeat,
	"UNTIL":                   TokenUntil,
	"CONST":                   TokenConst,
	"EXIT":                    TokenExit,
	"ENUM":                    TokenEnum,
	"TYPE":                    TokenTypeKw,
	"ENDTYPE":                 TokenEndType,
	"ON":                      TokenOn,
	"ENDON":                   TokenEndOn,
	"KEYDOWN":                 TokenKeyDown,
	"KEYPRESSED":              TokenKeyPressed,
	"STARTCOROUTINE":           TokenStartCoroutine,
	"YIELD":                    TokenYield,
	"WAITSECONDS":              TokenWaitSeconds,
	"MAIN":                     TokenMain,
	"ENDMAIN":                  TokenEndMain,
}

// OperatorMap maps operator strings to token types
var OperatorMap = map[string]TokenType{
	"=":   TokenEqual,
	"<>":  TokenNotEqual,
	"<":   TokenLess,
	"<=":  TokenLessEqual,
	">":   TokenGreater,
	">=":  TokenGreaterEqual,
	"+":   TokenPlus,
	"-":   TokenMinus,
	"*":   TokenMultiply,
	"/":   TokenDivide,
	"%":   TokenMod,
	"AND": TokenAnd,
	"OR":  TokenOr,
	"NOT": TokenNot,
}
