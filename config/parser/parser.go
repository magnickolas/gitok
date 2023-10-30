package parser

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Entry interface{}

var _ Entry = Section{}
var _ Entry = KeyValue{}

type Section struct {
	section string
}

type SectionSubsection struct {
	section    string
	subsection string
}

type KeyValue struct {
	Key   string
	Value string
}

type Parser struct {
	lines      []string
	lineNumber int
}

func (p *Parser) Init(r io.Reader) (*Parser, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	for _, b := range bytes.Split(buf.Bytes(), []byte{'\n'}) {
		p.lines = append(p.lines, string(b))
	}
	p.lineNumber = 1
	return p, nil
}

func NewParser(r io.Reader) (*Parser, error) {
	return new(Parser).Init(r)
}

func (p *Parser) Parse() ([]KeyValue, error) {
	entries, err := p.parseEntries()
	if err != nil {
		return nil, err
	}
	var kvs []KeyValue
	var curPrefix string
	for _, entry := range entries {
		switch v := entry.(type) {
		case Section:
			curPrefix = v.section
		case SectionSubsection:
			curPrefix = fmt.Sprintf("%s.%s", v.section, v.subsection)
		case KeyValue:
			key := fmt.Sprintf("%s.%s", curPrefix, v.Key)
			kvs = append(kvs, KeyValue{
				Key:   key,
				Value: v.Value,
			})
		}
	}
	return kvs, nil
}

func (p *Parser) parseEntries() ([]Entry, error) {
	var entries []Entry
	for len(p.lines) > 0 {
		entry, err := p.nextEntry()
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (p *Parser) fmtError(s string) error {
	return fmt.Errorf("line %v: %v", p.lineNumber, s)
}

func (p *Parser) fmtErrorf(format string, a ...any) error {
	return p.fmtError(fmt.Sprintf(format, a...))
}

func (p *Parser) shift() {
	p.lines = p.lines[1:]
	p.lineNumber += 1
}

func (p *Parser) nextSection() (Entry, error) {
	line := strings.TrimSpace(trimComment(p.lines[0]))
	if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
		return nil, p.fmtError("not a section")
	}
	line = line[1 : len(line)-1]
	parts := strings.SplitN(line, " ", 2)
	if len(parts) == 1 {
		section := parts[0]
		if len(section) == 0 || !isValidSection(section) {
			return nil, p.fmtError("invalid section name")
		}
		return Section{section: section}, nil
	}
	var section string
	section, line = parts[0], strings.TrimSpace(parts[1])
	if !isValidSection(section) {
		return nil, p.fmtError("invalid section name")
	}
	if !strings.HasPrefix(line, "\"") && !strings.HasSuffix(line, "\"") {
		return nil, p.fmtError("subsection not in quotes")
	}
	line = line[1 : len(line)-1]
	if !isValidSubsection(line) {
		return nil, p.fmtError("invalid subsection name")
	}
	subsection := escapeSubsection(line)
	return SectionSubsection{
		section:    section,
		subsection: subsection,
	}, nil
}

func (p *Parser) nextKeyValue() (Entry, error) {
	line := strings.TrimSpace(p.lines[0])
	parts := strings.SplitN(line, "=", 2)
	if len(parts) == 1 {
		key := parts[0]
		if !isValidKey(key) {
			return nil, p.fmtError("invalid key name")
		}
		return KeyValue{
			Key:   key,
			Value: "",
		}, nil
	}
	if len(parts) != 2 {
		return nil, p.fmtError("missing =")
	}
	var key string
	key, line = strings.TrimSpace(parts[0]), strings.TrimLeft(parts[1], " ")
	if !isValidKey(key) {
		return nil, p.fmtError("invalid key name")
	}
	value, err := p.nextValue(line)
	if err != nil {
		return nil, err
	}
	return KeyValue{
		Key:   key,
		Value: value,
	}, nil
}

func (p *Parser) nextValue(cur string) (string, error) {
	res := ""
	insideQuoted := false
	for {
		curRunes := []rune(cur)
		breakLine := false
		trailingSpaces := 0
		for i := 0; i < len(curRunes); i += 1 {
			if curRunes[i] == '\\' {
				if !insideQuoted {
					trailingSpaces = 0
				}
				i += 1
				if i == len(curRunes) {
					breakLine = true
					break
				}
				if curRunes[i] == '\\' || curRunes[i] == '"' {
					res += string(curRunes[i])
				} else if curRunes[i] == 'n' {
					res += "\n"
				} else if curRunes[i] == 't' {
					res += "\t"
				} else if curRunes[i] == 'b' {
					res += "\b"
				} else {
					return "", p.fmtErrorf("invalid value string: incorrect sequence \\%c", curRunes[i])
				}
			} else if curRunes[i] == '"' {
				if !insideQuoted {
					trailingSpaces = 0
				}
				insideQuoted = !insideQuoted
			} else if curRunes[i] == ';' || curRunes[i] == '#' {
				if insideQuoted {
					res += string(curRunes[i])
				} else {
					break
				}
			} else if curRunes[i] == ' ' {
				if !insideQuoted {
					trailingSpaces += 1
				}
				res += string(curRunes[i])
			} else {
				trailingSpaces = 0
				res += string(curRunes[i])
			}
		}
		res = res[:len(res)-trailingSpaces]
		if breakLine {
			p.shift()
			if len(p.lines) == 0 {
				break
			}
			cur = p.lines[0]
			continue
		}
		break
	}
	if insideQuoted {
		return "", p.fmtError("invalid value string: unclosed quote string")
	}
	return res, nil
}

func (p *Parser) nextEntry() (Entry, error) {
	defer p.shift()
	line := strings.TrimSpace(p.lines[0])
	if len(line) == 0 || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
		return nil, nil
	}
	if strings.HasPrefix(line, "[") {
		return p.nextSection()
	}
	return p.nextKeyValue()
}

func trimComment(s string) string {
	i := strings.IndexAny(s, ";#")
	if i != -1 {
		return s[:i]
	}
	return s
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlpha(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func isAlphanum(r rune) bool {
	return isDigit(r) || isAlpha(r)
}

func isValidSectionRune(r rune) bool {
	return isAlphanum(r) || r == '.' || r == '-'
}

func isValidSubsectionRune(r rune) bool {
	return r != 0 && r != '\n'
}

func escapeSubsection(s string) (res string) {
	runes := []rune(s)
	for i := 0; i < len(runes); i += 1 {
		if runes[i] != '\\' {
			res += string(runes[i])
		} else {
			i += 1
			if i < len(runes) && (runes[i] == '"' || runes[i] == '\\') {
				res += string(runes[i])
			}
		}
	}
	return
}

func isValidKeyRune(r rune) bool {
	return isAlphanum(r) || r == '-'
}

func isValidSection(s string) bool {
	for _, r := range s {
		if !isValidSectionRune(r) {
			return false
		}
	}
	return true
}

func isValidSubsection(s string) bool {
	for _, r := range s {
		if !isValidSubsectionRune(r) {
			return false
		}
	}
	return true
}

func isValidKey(s string) bool {
	if len(s) == 0 {
		return false
	}
	if !isAlpha(rune(s[0])) {
		return false
	}
	for _, r := range s {
		if !isValidKeyRune(r) {
			return false
		}
	}
	return true
}
