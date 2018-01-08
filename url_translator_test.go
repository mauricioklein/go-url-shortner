package urlshortner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	testCases := map[int]string{
		0:  "a",
		61: "9",
		62: "ba",
		63: "bb",
	}

	for id, expectedCode := range testCases {
		assert.Equal(t, expectedCode, encode(id))
	}
}

func TestDecode(t *testing.T) {
	testCases := map[string]int{
		"a":  0,
		"9":  61,
		"ba": 62,
		"bb": 63,
	}

	for code, expectedID := range testCases {
		assert.Equal(t, expectedID, decode(code))
	}
}
