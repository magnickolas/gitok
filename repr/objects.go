package repr

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"slices"
)

type Object interface {
	// Object's bytes
	Raw() []byte
	// Compressed object's bytes
	Compressed() ([]byte, error)
	// Object's digest (hex string)
	Digest() string
	// Representation for pretty-printing an object
	String() string
	Type() string
}

func StripObjectHeader(o Object) []byte {
	return bytes.SplitN(o.Raw(), []byte{0}, 2)[1]
}

type LazyObject struct {
	raw            []byte
	lazyCompressed []byte
	lazyDigest     string
}

func (lo *LazyObject) Compressed() ([]byte, error) {
	if lo.lazyCompressed == nil {
		var err error
		lo.lazyCompressed, err = compress(lo.raw)
		if err != nil {
			return nil, err
		}
	}
	return lo.lazyCompressed, nil
}

func (lo *LazyObject) Digest() string {
	if lo.lazyDigest == "" {
		lo.lazyDigest = hasherHex(lo.raw)
	}
	return lo.lazyDigest
}

// verify interface compliance
var _ Object = (*Blob)(nil)
var _ Object = (*Tree)(nil)

type Blob struct {
	LazyObject
	content []byte
}

func NewBlob(r io.Reader) (*Blob, error) {
	return new(Blob).Init(r)
}

func (b *Blob) Init(r io.Reader) (*Blob, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	b.content = buf.Bytes()
	b.raw = b.Raw()
	return b, nil
}

func (b *Blob) Raw() []byte {
	if b.raw == nil {
		b.raw = append(b.raw, fmt.Sprintf("b %d", len(b.content))...)
		b.raw = append(b.raw, 0)
		b.raw = append(b.raw, b.content...)
	}
	return b.raw
}

func (b *Blob) String() string {
	return string(b.content)
}

func (b *Blob) Type() string {
	return "blob"
}

type Tree struct {
	LazyObject
	children []nodeType
}

type nodeType struct {
	name   string
	mode   ObjectModeType
	digest string
}

func (n *nodeType) getType() string {
	if n.mode == ModeTree {
		return "tree"
	}
	return "blob"
}

type ObjectModeType string

const (
	ModeNormal       ObjectModeType = "100644"
	ModeExecutable   ObjectModeType = "100755"
	ModeSymbolicLink ObjectModeType = "120000"
	ModeTree         ObjectModeType = "40000"
)

var modes = []ObjectModeType{ModeNormal, ModeExecutable, ModeSymbolicLink, ModeTree}

func NewTree(r io.Reader) (*Tree, error) {
	return new(Tree).Init(r)
}

func (t *Tree) Init(r io.Reader) (*Tree, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()
	for len(b) > 0 {
		parts := bytes.SplitN(b, []byte{' '}, 2)
		if len(parts) != 2 {
			return nil, ErrorCorruptedObject
		}
		var mode ObjectModeType
		mode, b = ObjectModeType(parts[0]), parts[1]
		if !slices.Contains[[]ObjectModeType](modes, mode) {
			return nil, formatErrorUnknownFileMode(string(mode))
		}
		parts = bytes.SplitN(b, []byte{0}, 2)
		if len(parts) != 2 {
			return nil, ErrorCorruptedObject
		}
		var name string
		name, b = string(parts[0]), parts[1]
		var digest string
		digest, b = hex.EncodeToString(b[:20]), b[20:]
		t.children = append(t.children, nodeType{
			name:   name,
			mode:   mode,
			digest: digest,
		})
	}
	t.raw = t.Raw()
	return t, nil
}

func (t *Tree) Raw() []byte {
	if t.raw == nil {
		var entries []byte
		for _, child := range t.children {
			var rawDigest []byte
			rawDigest, _ = hex.DecodeString(child.digest)
			entries = append(entries, child.mode...)
			entries = append(entries, ' ')
			entries = append(entries, child.name...)
			entries = append(entries, 0)
			entries = append(entries, rawDigest...)
		}
		t.raw = append(t.raw, fmt.Sprintf("t %d", len(entries))...)
		t.raw = append(t.raw, 0)
		t.raw = append(t.raw, entries...)
	}
	return t.raw
}

func (t *Tree) String() (res string) {
	for _, child := range t.children {
		res += fmt.Sprintf(
			"%06s %s %s\t%s\n",
			child.mode,
			child.getType(),
			child.digest,
			child.name,
		)
	}
	return
}

func (t *Tree) Type() string {
	return "tree"
}
