package vtt

type Token int

const (

	// Special Tokens
	NIL       Token = iota // Indicates EOF or error.
	WS                     // One or multiple spaces. Does not include LF.
	LF                     // LF/Newline
	DOUBLE_LF              // \n\n

	// "Primitives"
	NUMBER // [0-9]*
	LETTER // [A-Z]|[a-z]

	// Literal Identifiers
	WEBVTT
	STYLE
	NOTE

	// Styling Literal Identifiers
	VERTICAL
	LINE
	POSITION
	SIZE
	ALIGN

	// Symbols
	RARROW     // -->
	LT         // <
	GT         // >
	COLON      // :
	DOT        // .
	LEG_SYMBOL // All remaining legal symbols in text that has no special functionality. //Newline, Space, -->, & , <, >, are reserved as separate tokens.

	// Escape sequences
	ESC_AMP  // &amp
	ESC_LT   // &lt;
	ESC_GT   //	&gt;
	ESC_LRM  // &lrm;
	ESC_RLM  // &rlm;
	ESC_NBSP // &nbsp;

)

func StringFromToken(tok Token) string {
	switch tok {
	case NIL:
		return "NIL"
	case WS:
		return "WS"
	case LF:
		return "LF"
	case DOUBLE_LF:
		return "DOUBLE_LF"
	case NUMBER:
		return "NUMBER"
	case LETTER:
		return "LETTER"
	case WEBVTT:
		return "WEBVTT"
	case STYLE:
		return "STYLE"
	case NOTE:
		return "NOTE"
	case VERTICAL:
		return "VERTICAL"
	case LINE:
		return "LINE"
	case POSITION:
		return "POSITION"
	case SIZE:
		return "SIZE"
	case ALIGN:
		return "ALIGN"
	case RARROW:
		return "RARROW"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case COLON:
		return "COLON"
	case DOT:
		return "DOT"
	case LEG_SYMBOL:
		return "LEG_SYMBOL"
	case ESC_AMP:
		return "ESC_AMP"
	case ESC_LT:
		return "ESC_LT"
	case ESC_GT:
		return "ESC_GT"
	case ESC_LRM:
		return "ESC_LRM"
	case ESC_RLM:
		return "ESC_RLM"
	case ESC_NBSP:
		return "ESC_NBSP"
	default:
		panic("StringToToken failed. Unknown token passed.")
	}
}
