package restic

import (
	"fmt"

	"restic/errors"
)

// Blob is one part of a file or a tree.
type Blob struct {
	Type   BlobType
	Length uint
	ID     ID
	Offset uint
}

// PackedBlob is a blob stored within a file.
type PackedBlob struct {
	Blob
	PackID ID
}

// BlobHandle identifies a blob of a given type.
type BlobHandle struct {
	ID   ID
	Type BlobType
}

func (h BlobHandle) String() string {
	return fmt.Sprintf("<%s/%s>", h.Type, h.ID.Str())
}

// BlobType specifies what a blob stored in a pack is.
type BlobType uint8

// These are the blob types that can be stored in a pack.
const (
	InvalidBlob BlobType = iota
	DataBlob
	TreeBlob
)

func (t BlobType) String() string {
	switch t {
	case DataBlob:
		return "data"
	case TreeBlob:
		return "tree"
	}

	return fmt.Sprintf("<BlobType %d>", t)
}

// MarshalJSON encodes the BlobType into JSON.
func (t BlobType) MarshalJSON() ([]byte, error) {
	switch t {
	case DataBlob:
		return []byte(`"data"`), nil
	case TreeBlob:
		return []byte(`"tree"`), nil
	}

	return nil, errors.New("unknown blob type")
}

// UnmarshalJSON decodes the BlobType from JSON.
func (t *BlobType) UnmarshalJSON(buf []byte) error {
	switch string(buf) {
	case `"data"`:
		*t = DataBlob
	case `"tree"`:
		*t = TreeBlob
	default:
		return errors.New("unknown blob type")
	}

	return nil
}

// BlobHandles is an ordered list of BlobHandles that implements sort.Interface.
type BlobHandles []BlobHandle

func (h BlobHandles) Len() int {
	return len(h)
}

func (h BlobHandles) Less(i, j int) bool {
	for k, b := range h[i].ID {
		if b == h[j].ID[k] {
			continue
		}

		if b < h[j].ID[k] {
			return true
		}

		return false
	}

	return h[i].Type < h[j].Type
}

func (h BlobHandles) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h BlobHandles) String() string {
	elements := make([]string, 0, len(h))
	for _, e := range h {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("%v", elements)
}
