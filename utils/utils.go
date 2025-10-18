package utils

import (
	"encoding/json"
	"strconv"
)

func MustMarshal(data any) []byte {
	res, err := json.Marshal(data)
	Must(err)
	return res
}

func MustStrToInt(str string) int {
	res, err := strconv.Atoi(str)
	Must(err)
	return res
}
