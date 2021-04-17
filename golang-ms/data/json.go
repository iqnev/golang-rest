package data

import (
	"encoding/json"

	"io"
)

func ToJson(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}

func FromJson(i interface{}, r io.Reader) error {
	e := json.NewDecoder(r)

	return e.Decode(i)
}
