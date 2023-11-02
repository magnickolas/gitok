package parser_test

import (
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/magnickolas/gitok/config/parser"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    []parser.KeyValue
		wantErr bool
	}{
		{
			input: `[section]`,
			want:  nil,
		},
		{
			input: `[SECtion "XyZ"]
				    key1=value1`,
			want: []parser.KeyValue{
				{Key: "section.XyZ.key1", Value: "value1"},
			},
		},
		{
			input: `[SECtion1.subSECtion2]
				    key1=value1`,
			want: []parser.KeyValue{
				{Key: "section1.subsection2.key1", Value: "value1"},
			},
		},
		{
			input:   `[section`,
			wantErr: true,
		},
		{
			input: `[section]
		            key1 = part1 part2  ; comment`,
			want: []parser.KeyValue{
				{Key: "section.key1", Value: "part1 part2"},
			},
		},
		{
			input: stripMargin(`[section]
				          	   |key1 = part1 \
							   |part2`),
			want: []parser.KeyValue{
				{Key: "section.key1", Value: "part1 part2"},
			},
		},
		{
			input: `[section]
				    key1 = "value1  " # comment
					key2 = part1\b\n\t\"\\part2`,
			want: []parser.KeyValue{
				{Key: "section.key1", Value: "value1  "},
				{Key: "section.key2", Value: "part1\b\n\t\"\\part2"},
			},
		},
		{
			input: `[section1]
				    key1 =value1
					[section2]
					key2= value2`,
			want: []parser.KeyValue{
				{Key: "section1.key1", Value: "value1"},
				{Key: "section2.key2", Value: "value2"},
			},
		},
		{
			input:   "[section \"subsec\x00tion\"]",
			wantErr: true,
		},
		{
			input:   "[section \"subsec\ntion\"]",
			wantErr: true,
		},
		{
			input: `[section "subsec\\\"tion"]
			          key1 = value1`,
			want: []parser.KeyValue{
				{Key: `section.subsec\"tion.key1`, Value: "value1"},
			},
		},
	}
	for _, test := range tests {
		r := strings.NewReader(test.input)
		p, err := parser.NewParser(r)
		if err != nil {
			t.Error("failed to initialize parser")
		}
		got, err := p.Parse()
		if test.wantErr {
			if err == nil {
				t.Errorf("wanted error for %#v", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("failed to parse %#v", test.input)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("incorrect result for %#v: wanted %#v, got %#v", test.input, test.want, got)
			}
		}
	}
}

func stripMargin(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		strippedLine := strings.TrimLeftFunc(line, unicode.IsSpace)
		if len(strippedLine) > 0 && strippedLine[0] == '|' {
			strippedLine = strippedLine[1:]
		}
		lines[i] = strippedLine
	}
	return strings.Join(lines, "\n")
}
