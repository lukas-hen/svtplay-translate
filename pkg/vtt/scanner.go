package vtt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Scanner struct {
	rd *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{rd: bufio.NewReader(r)}
}

func (s *Scanner) Scan() (Token, string, error) {

	tok, err := s.PeekTokenType()
	if err != nil {
		return NIL, "", err
	}

	switch tok {

	case WEBVTT:
		scanned, err := s.scanExact("WEBVTT")
		return WEBVTT, scanned, err

	case STYLE:
		scanned, err := s.scanExact("STYLE")
		return STYLE, scanned, err

	case NOTE:
		scanned, err := s.scanExact("NOTE")
		return NOTE, scanned, err

	case RARROW:
		scanned, err := s.scanExact("-->")
		return RARROW, scanned, err

	case DOUBLE_LF:
		scanned, err := s.scanExact("\n\n")
		return DOUBLE_LF, scanned, err

	case ESC_AMP:
		scanned, err := s.scanExact("&amp;")
		return ESC_AMP, scanned, err

	case ESC_LT:
		scanned, err := s.scanExact("&lt;")
		return ESC_LT, scanned, err

	case ESC_GT:
		scanned, err := s.scanExact("&gt;")
		return ESC_GT, scanned, err

	case ESC_LRM:
		scanned, err := s.scanExact("&lrm;")
		return ESC_LRM, scanned, err

	case ESC_RLM:
		scanned, err := s.scanExact("&rlm;")
		return ESC_RLM, scanned, err

	case ESC_NBSP:
		scanned, err := s.scanExact("&nbsp;")
		return ESC_NBSP, scanned, err

	case VERTICAL:
		scanned, err := s.scanExact("vertical:")
		return VERTICAL, scanned, err

	case LINE:
		scanned, err := s.scanExact("line:")
		return LINE, scanned, err

	case POSITION:
		scanned, err := s.scanExact("position:")
		return POSITION, scanned, err

	case SIZE:
		scanned, err := s.scanExact("size:")
		return SIZE, scanned, err

	case ALIGN:
		scanned, err := s.scanExact("align:")
		return ALIGN, scanned, err

	case LT:
		scanned, err := s.scanExact("<")
		return LT, scanned, err

	case GT:
		scanned, err := s.scanExact(">")
		return GT, scanned, err

	case COLON:
		scanned, err := s.scanExact(":")
		return COLON, scanned, err

	case LF:
		scanned, err := s.scanExact("\n")
		return LF, scanned, err

	case DOT:
		scanned, err := s.scanExact(".")
		return DOT, scanned, err

	case WS:
		return s.scanWhitespace()

	case NUMBER:
		return s.scanNumber()

	case LETTER:
		r, err := s.next()
		if err != nil {
			return LETTER, "", err
		}

		return LETTER, string(r), nil

	case LEG_SYMBOL:
		// This means the symbol is neither a letter
		// or a number, but also not a reserved/illegal symbol.
		r, err := s.next()
		if err != nil {
			return LEG_SYMBOL, "", err
		}

		return LEG_SYMBOL, string(r), nil

	default:
		r, err := s.next()
		if err != nil {
			return NIL, "", err
		}

		return NIL, string(r), fmt.Errorf("found illegal rune: %s", string(r))
	}

}

func (s *Scanner) PeekTokenType() (Token, error) {
	// Scan multiline tokens first as they
	// have precedence over the single rune ones.
	PEEK_N := 8 // Longest literal keyword is 'vertical', 8 runes.
	// Ignoring error here. We don't want to quit if peeking too far.
	// Just let the peek buffer be wrong so scanning continues until the last byte.
	bytes, _ := s.rd.Peek(PEEK_N)
	if len(bytes) == 0 {
		return NIL, io.EOF
	}
	peek_str := string(bytes)

	switch {

	case strings.HasPrefix(peek_str, "WEBVTT"):
		return WEBVTT, nil

	case strings.HasPrefix(peek_str, "STYLE"):
		return STYLE, nil

	case strings.HasPrefix(peek_str, "NOTE"):
		return NOTE, nil

	case strings.HasPrefix(peek_str, "-->"):
		return RARROW, nil

	case strings.HasPrefix(peek_str, "\n\n"):
		return DOUBLE_LF, nil

	case strings.HasPrefix(peek_str, "&amp;"):
		return ESC_AMP, nil

	case strings.HasPrefix(peek_str, "&lt;"):
		return ESC_LT, nil

	case strings.HasPrefix(peek_str, "&gt;"):
		return ESC_GT, nil

	case strings.HasPrefix(peek_str, "&lrm;"):
		return ESC_LRM, nil

	case strings.HasPrefix(peek_str, "&rlm;"):
		return ESC_RLM, nil

	case strings.HasPrefix(peek_str, "&nbsp;"):
		return ESC_NBSP, nil

	// Below are scanned by including the :
	// A better solution would be that the colon should be a separate token
	// And the tokens should be contextually parsed.
	// But who has time for that?

	case strings.HasPrefix(peek_str, "vertical:"):
		return VERTICAL, nil

	case strings.HasPrefix(peek_str, "line:"):
		return LINE, nil

	case strings.HasPrefix(peek_str, "position:"):
		return POSITION, nil

	case strings.HasPrefix(peek_str, "size:"):
		return SIZE, nil

	case strings.HasPrefix(peek_str, "align:"):
		return ALIGN, nil
	}

	// Tokens where token type can be determined via one single rune.
	ch := []rune(peek_str)[0]

	switch {
	case ch == '<':
		return LT, nil
	case ch == '>':
		return GT, nil
	case ch == ':':
		return COLON, nil
	case ch == '\n':
		return LF, nil
	case ch == '.':
		return DOT, nil
	case isWhitespace(ch):
		return WS, nil
	case unicode.IsDigit(ch):
		return NUMBER, nil
	case unicode.IsLetter(ch):
		return LETTER, nil
	case !isReserved(ch):
		// This means the symbol is neither a letter
		// or a number, but also not a reserved/illegal symbol.
		return LEG_SYMBOL, nil
	default:
		return NIL, fmt.Errorf("found illegar rune: %s", string(ch))
	}
}

func (s *Scanner) scanWhitespace() (Token, string, error) {
	var buf bytes.Buffer

	for {
		ch, err := s.next()
		if err == io.EOF {
			break
		} else if err != nil {
			return NIL, "", err
		} else if !isWhitespace(ch) {
			s.prev()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String(), nil
}

func (s *Scanner) scanExact(ident string) (string, error) {
	var buf bytes.Buffer

	for _, ident_ch := range ident {
		ch, err := s.next()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		} else if ident_ch != ch {
			return "", fmt.Errorf("couldn't scan exact identifier: %s. Found \"%c\", expected \"%c\"", ident, ch, ident_ch)
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String(), nil
}

func (s *Scanner) scanNumber() (Token, string, error) {
	var buf bytes.Buffer

	for {
		ch, err := s.next()
		if err == io.EOF {
			break
		} else if err != nil {
			return NIL, "", err
		} else if !unicode.IsDigit(ch) {
			s.prev()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return NUMBER, buf.String(), nil
}

func (s *Scanner) next() (rune, error) {
	r, _, err := s.rd.ReadRune()
	if err != nil {
		return rune(NIL), err
	}

	return r, nil
}

func (s *Scanner) prev() error {
	err := s.rd.UnreadRune()
	if err != nil {
		return err
	}

	return nil
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isReserved(ch rune) bool {
	return ch == '&' || ch == '<' || ch == '>' || isWhitespace(ch)
}
