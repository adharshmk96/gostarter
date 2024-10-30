package common

import (
	"encoding/json"
	"io"
)

func UnmarshalData[T any](data io.ReadCloser) (T, error) {
	var obj T
	err := json.NewDecoder(data).Decode(&obj)
	return obj, err
}
