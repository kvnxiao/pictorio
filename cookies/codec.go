package cookies

import (
	"encoding/base64"
)

func encode(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

func decode(input string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(input)
}
