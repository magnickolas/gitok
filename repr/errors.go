package repr

import (
	"errors"
	"fmt"
)

var (
	ErrorCorruptedObject         = errors.New("corrupted object")
	ErrorCorruptedObjectHeader   = fmt.Errorf("%w: corrupted header", ErrorCorruptedObject)
	ErrorUnknownObjectType       = errors.New("unknown object type")
	formatErrorUnknownObjectType = func(t string) error {
		return fmt.Errorf("%w: %v", ErrorUnknownObjectType, t)
	}
	ErrorSizeNotMatch          = errors.New("object size does not match")
	ErrorUnknownFileMode       = errors.New("unknown file mode")
	formatErrorUnknownFileMode = func(mode string) error {
		return fmt.Errorf("%w: %v", ErrorUnknownFileMode, mode)
	}
)
