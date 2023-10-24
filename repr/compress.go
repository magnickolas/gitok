package repr

import (
	"bytes"
	"compress/zlib"
)

func compress(raw []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	zw := zlib.NewWriter(buf)
	_, err := zw.Write(raw)
	if err != nil {
		return nil, err
	}
	// close so that the checksum is written to the buffer
	err = zw.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func uncompress(b []byte) ([]byte, error) {
	r := bytes.NewReader(b)
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
