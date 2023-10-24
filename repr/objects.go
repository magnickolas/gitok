package repr

import (
	"bytes"
	"fmt"
	"io"
)

type Object interface {
	// Uncompressed representation of the object
	Raw() []byte
	// Compressed representation of the object that is stored
	Compressed() ([]byte, error)
	Hash() string
	Representation() []byte
}

// verify interface compliance
var _ Object = (*Blob)(nil)

type Blob struct {
	content        []byte
	lazyCompressed []byte
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
	return blob, nil
}

func (blob *Blob) Raw() (res []byte) {
	res = append(res, fmt.Sprintf("blob %d", len(blob.content))...)
	res = append(res, 0)
	res = append(res, blob.content...)
	return
}

func (blob *Blob) Compressed() ([]byte, error) {
	if blob.lazyCompressed == nil {
		var err error
		blob.lazyCompressed, err = compress(blob.Raw())
		if err != nil {
			return nil, err
		}
	}
	return blob.lazyCompressed, nil
}

func (blob *Blob) Hash() string {
	return hasherHex(blob.Raw())
}

func (blob *Blob) Representation() []byte {
	return blob.content
}
