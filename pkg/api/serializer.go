package api

import (
	"encoding/json"
)

type ResponseSerializer interface {
	Serialize() (b []byte, err error) 
}

func (r *Response[T]) Serialize() (b []byte, err error) {
	return json.Marshal(r)
}
