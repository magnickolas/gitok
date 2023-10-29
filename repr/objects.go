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

func (blob *Blob) Init(r io.Reader) (*Blob, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	blob.content = buf.Bytes()
	blob.raw = blob.Raw()
	return blob, nil
}

func (blob *Blob) Raw() []byte {
	if blob.raw == nil {
		blob.raw = append(blob.raw, fmt.Sprintf("blob %d", len(blob.content))...)
		blob.raw = append(blob.raw, 0)
		blob.raw = append(blob.raw, blob.content...)
	}
	return blob.raw
}

func (blob *Blob) String() string {
	return string(blob.content)
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

func (node *nodeType) getType() string {
	if node.mode == ModeTree {
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

func (tree *Tree) Init(r io.Reader) (*Tree, error) {
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
		tree.children = append(tree.children, nodeType{
			name:   name,
			mode:   mode,
			digest: digest,
		})
	}
	tree.raw = tree.Raw()
	return tree, nil
}

func (tree *Tree) Raw() []byte {
	if tree.raw == nil {
		var entries []byte
		for _, child := range tree.children {
			var rawDigest []byte
			rawDigest, _ = hex.DecodeString(child.digest)
			entries = append(entries, child.mode...)
			entries = append(entries, ' ')
			entries = append(entries, child.name...)
			entries = append(entries, 0)
			entries = append(entries, rawDigest...)
		}
		tree.raw = append(tree.raw, fmt.Sprintf("tree %d", len(entries))...)
		tree.raw = append(tree.raw, 0)
		tree.raw = append(tree.raw, entries...)
	}
	return tree.raw
}

func (blob *Tree) String() (res string) {
	for _, child := range blob.children {
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
