package utils

import "encoding/json"

func MustMarshal(data any) []byte {
	res, err := json.Marshal(data)
	Must(err)
	return res
}
