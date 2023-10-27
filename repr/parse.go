package repr

import (
	"bytes"
	"strconv"
)

func ParseObject(compressed []byte) (Object, error) {
	b, err := uncompress(compressed)
	if err != nil {
		return nil, err
	}

	parts := bytes.SplitN(b, []byte{0}, 2)
	if len(parts) != 2 {
		return nil, ErrorCorruptedObject
	}
	headerBytes, content := parts[0], parts[1]

	headerParts := bytes.SplitN(headerBytes, []byte{' '}, 2)
	if len(headerParts) != 2 {
		return nil, ErrorCorruptedObjectHeader
	}

	objType, sizeStr := string(headerParts[0]), string(headerParts[1])

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, err
	}
	if size != len(content) {
		return nil, ErrorSizeNotMatch
	}

	r := bytes.NewReader(content)

	switch objType {
	case "blob":
		return NewBlob(r)
	case "tree":
		return NewTree(r)
	}
	return nil, formatErrorUnknownObjectType(objType)
}
