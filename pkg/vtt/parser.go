package vtt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Parser struct {
	s *Scanner
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func ParseFile(path string) *WebVTT {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := NewParser(f)

	return p.Parse()
}

// Panics
func (p *Parser) Parse() *WebVTT {

	header, err := p.parseHeader()
	if err != nil {
		panic(err)
	}

	var allNotes []string
	var allCues []*Cue

	for {
		tok, err := p.s.PeekTokenType()

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		switch tok {
		case STYLE: // Don't parse for now.
			err := p.skipUntil(DOUBLE_LF)
			if err != nil {
				panic(err)
			}
		case NOTE:
			str, err := p.stringUntil(DOUBLE_LF, isValidNoteToken)
			if err != nil {
				panic(err)
			}
			allNotes = append(allNotes, str)
		default:
			cue, err := p.parseCue()
			if err != nil {
				break // Temp fix
			}
			allCues = append(allCues, cue)
		}
	}

	return &WebVTT{
		Header: header,
		Notes:  allNotes,
		Cues:   allCues,
	}
}

func (p *Parser) parseHeader() (string, error) {
	tok, str, err := p.s.Scan()
	if err != nil {
		return "", err
	} else if tok != WEBVTT {
		return "", fmt.Errorf("expected token \"WEBVTT\", found [%s: %s]", StringFromToken(tok), str)
	}

	tmp := str

	str, err = p.stringUntil(DOUBLE_LF, isValidHeaderToken)
	if err != nil {
		return "", err
	}
	tmp += str

	return tmp, nil
}

func (p *Parser) parseCue() (*Cue, error) {

	timings := Timings{
		From: "",
		To:   "",
	}

	// Note now cues can start with the timings and skip id's OR have a first line with an id.
	// First line of Cue
	// Attempt parsing timestamp
	str, success, _ := p.parseTimestamp()
	if success { // Timestamp parsed successfully, this is an cue without id
		// Parse rest of timing
		timings.From = str

		_, wsexperr := p.expect(WS)
		if wsexperr != nil {
			return nil, wsexperr
		}

		_, arrexperr := p.expect(RARROW)
		if arrexperr != nil {
			return nil, arrexperr
		}

		_, wsexperr = p.expect(WS)
		if wsexperr != nil {
			return nil, wsexperr
		}

		str, success, err := p.parseTimestamp()
		if err != nil {
			return nil, err
		} else if !success {
			return nil, fmt.Errorf("error parsing timestamp, invalid token in: %s", str)
		} else {
			// success & no error
			timings.To = str
		}

		// Timings parsed, get STYLES & Text.

		// Skip scanning styles for now.
		p.skipUntil(LF)

		text, err := p.stringUntil(DOUBLE_LF, isValidCueTextToken)
		if err != nil {
			return nil, err
		}

		return &Cue{
			Id:      "",
			Timings: timings,
			Text:    text,
		}, nil

	} else { // Timestamp failed to parse, but not due to EOF. This has to be an id.
		// Parse rest of id
		id := str

		if strings.HasSuffix(id, "\n") {
			// edge case if newline is what broke the timestamp parsing.
			// better solition would be to unscan last scanned before running stringUntil
			id = id[:len(id)-1]
		} else {
			str, err := p.stringUntil(LF, isValidIdToken)
			if err != nil {
				return nil, err
			}

			id += str
		}

		// Parse timing

		str, success, err := p.parseTimestamp()
		if err != nil {
			return nil, err
		} else if !success {
			return nil, fmt.Errorf("error parsing timestamp, invalid token in: %s", str)
		} else {
			// success & no error
			timings.From = str
		}

		_, wsexperr := p.expect(WS)
		if wsexperr != nil {
			return nil, wsexperr
		}

		_, arrexperr := p.expect(RARROW)
		if arrexperr != nil {
			return nil, arrexperr
		}

		_, wsexperr = p.expect(WS)
		if wsexperr != nil {
			return nil, wsexperr
		}

		str, success, err = p.parseTimestamp()
		if err != nil {
			return nil, err
		} else if !success {
			return nil, fmt.Errorf("error parsing timestamp, invalid token in: %s", str)
		} else {
			// success & no error
			timings.To = str
		}

		// Skip scanning styles for now.
		p.skipUntil(LF)

		text, err := p.stringUntil(DOUBLE_LF, isValidCueTextToken)
		if err != nil {
			return nil, err
		}

		return &Cue{
			Id:      id,
			Timings: timings,
			Text:    text,
		}, nil
	}
}

func (p *Parser) parseTimestamp() (string, bool, error) {
	// mm:ss.ttt
	// ||
	// hh:mm:ss.ttt

	var tmp string

	n1, err := p.expect(NUMBER)
	tmp += n1
	if err != nil {
		return tmp, false, err
	}

	c1, err := p.expect(COLON)
	tmp += c1
	if err != nil {
		return tmp, false, err
	}

	n2, err := p.expect(NUMBER)
	tmp += n2
	if err != nil {
		return tmp, false, err
	}

	// Either : or .
	tok, str, _ := p.s.Scan()

	if tok == COLON {
		tmp += str
		n3, err := p.expect(NUMBER)
		tmp += n3
		if err != nil {
			return tmp, false, err
		}
		d, err := p.expect(DOT)
		tmp += d
		if err != nil {
			return tmp, false, err
		}
		n4, err := p.expect(NUMBER)
		tmp += n4
		if err != nil {
			return tmp, false, err
		}
	} else if tok == DOT {
		tmp += str
		n3, err := p.expect(NUMBER)
		tmp += n3
		if err != nil {
			return tmp, false, err
		}
	} else {
		// Non technical parsing error.
		return tmp, false, fmt.Errorf("could not parse timestamp")
	}
	return tmp, true, nil
}

// Skips parsing up to and including delim token.
func (p *Parser) skipUntil(delim Token) error {
	for {
		tok, _, err := p.s.Scan()
		if err != nil {
			// Since we don't expect EOF here, we also return it.
			return err
		} else if tok == delim {
			break
		}
	}

	return nil
}

func (p *Parser) stringUntil(delim Token, validToken func(Token) bool) (string, error) {
	var tmp string
	for {
		tok, str, err := p.s.Scan()
		if err != nil {
			return "", err
		} else if !validToken(tok) {
			return "", fmt.Errorf("found illegal token [%s], expected [%s]", StringFromToken(tok), str)
		} else if tok == delim {
			break
		} else {
			tmp += str
		}
	}

	return tmp, nil
}

func (p *Parser) expect(expected Token) (string, error) {
	// only scan if found expected.
	tok, str, err := p.s.Scan()
	if err != nil {
		return "", err
	} else if tok != expected {
		return str, fmt.Errorf("found illegal token [%s], expected [%s]", StringFromToken(tok), StringFromToken(expected))
	}
	return str, err
}

func isValidHeaderToken(tok Token) bool {
	return tok != RARROW
}

func isValidNoteToken(tok Token) bool {
	return tok != RARROW && tok != LT && tok != GT
}

func isValidCueTextToken(tok Token) bool {
	return true // All tokens valid
}

func isValidIdToken(tok Token) bool {
	return tok != RARROW // All tokens valid
}

func (webvtt *WebVTT) WriteSRT(w io.Writer) error {
	all_subs := webvtt.Cues
	for i := 0; i < len(all_subs); i++ {
		//n := strconv.Itoa(i + 1)
		//dur := all_subs[i].Timings.String()
		sub := all_subs[i].TextWithoutTags()
		//full_str := n + "\n" + dur + "\n" + sub + "\n\n"
		full_str := sub + "\n\n"

		_, err := w.Write([]byte(full_str))
		if err != nil {
			return err
		}
	}

	return nil
}

func (webvtt *WebVTT) WriteSrtFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = webvtt.WriteSRT(f)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cue) TextWithoutTags() string {
	s := bufio.NewScanner(strings.NewReader(c.Text))
	s.Split(bufio.ScanRunes)

	var buf bytes.Buffer

	for s.Scan() {
		if s.Text() == "<" {
			for s.Scan() {
				if s.Text() == ">" {
					break
				}
			}
		} else {
			buf.WriteString(s.Text())
		}
	}

	return buf.String()
}

func (c *Cue) ToSRT() string {
	return fmt.Sprintf("%s\n%s\n%s\n\n", c.Id, c.Timings.String(), c.Text)
}

func (ti *Timings) String() string {
	return ti.From + " --> " + ti.To
}

type WebVTT struct {
	Header string
	//Style  VTTStyle // Don't parse this for now.
	Notes []string
	Cues  []*Cue
}

type Cue struct {
	Id      string
	Timings Timings
	Styling Stylings
	Text    string
}

type Timings struct {
	From string
	To   string
}

type Stylings struct {
	Vertical string
	Line     string
	Position string
	Size     string
	Align    string
}
