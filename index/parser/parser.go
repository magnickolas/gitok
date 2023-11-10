package parser

import (
	"bytes"
	"fmt"
	"io"

	"github.com/magnickolas/gitok/repr"
)

type Index struct {
	Magic      int32
	Version    int32
	NumEntries int32
	Entries    []Entry
}

type Entry struct {
	CTime   int32
	CTimeNS int32
	MTime   int32
	MTimeNS int32
	Dev     int32
	Ino     int32
	Mode    int32
	Uid     int32
	Gid     int32
	Size    int32
	Name    string
}

type Parser struct {
	b []byte
}

func (p *Parser) Init(r io.Reader) (*Parser, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	p.b = buf.Bytes()
	return p, nil
}

func NewParser(r io.Reader) (*Parser, error) {
	return new(Parser).Init(r)
}

func (p *Parser) Parse() (*Index, error) {
	const CE_NAMEMASK uint32 = 0x0FFF
	const CE_EXTENDED uint32 = 0x4000
	const CE_INTENT_TO_ADD uint32 = (1 << 29)
	const CE_SKIP_WORKTREE uint32 = (1 << 30)
	const CE_EXTENDED_FLAGS uint32 = CE_INTENT_TO_ADD | CE_SKIP_WORKTREE

	index := Index{}
	index.Magic = p.parseInt32()
	if index.Magic != int32('D')<<24|int32('I')<<16|int32('R')<<8|int32('C') {
		return nil, fmt.Errorf("wrong magic")
	}
	index.Version = p.parseInt32()
	index.NumEntries = p.parseInt32()
	for i := int32(0); i < index.NumEntries; i += 1 {
		entry := Entry{}
		entry.CTime = p.parseInt32()
		entry.CTimeNS = p.parseInt32()
		entry.MTime = p.parseInt32()
		entry.MTimeNS = p.parseInt32()
		entry.Dev = p.parseInt32()
		entry.Ino = p.parseInt32()
		entry.Mode = p.parseInt32()
		entry.Uid = p.parseInt32()
		entry.Gid = p.parseInt32()
		entry.Size = p.parseInt32()
		p.shift(repr.HashSize())
		flags := uint32(p.parseInt16())
		length := uint32(flags & CE_NAMEMASK)
		if flags&CE_EXTENDED != 0 {
			extended_flags := uint32(p.parseInt16()) << 16
			if extended_flags&(^CE_EXTENDED_FLAGS) != 0 {
				return nil, fmt.Errorf("unknown index entry format 0x%08x", extended_flags)
			}
			flags |= extended_flags
		}
		if length == CE_NAMEMASK {
			length = uint32(bytes.Index(p.b, []byte{0}))
		}
		entry.Name = string(p.b[:length])
		p.shift(int(length) + 1)
		index.Entries = append(index.Entries, entry)
	}
	return &index, nil
}

func (p *Parser) parseInt16() (res int16) {
	res = int16(p.b[0])<<8 | int16(p.b[1])
	p.shift(2)
	return
}

func (p *Parser) parseInt32() (res int32) {
	res = int32(p.b[0])<<24 | int32(p.b[1])<<16 | int32(p.b[2])<<8 | int32(p.b[3])
	p.shift(4)
	return
}

func (p *Parser) shift(n int) {
	p.b = p.b[n:]
}
