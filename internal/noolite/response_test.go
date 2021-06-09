package noolite

import "testing"

type testPair struct {
	data      []byte
	isCorrect bool
}

var tests = []testPair{
	{[]byte{}, false},
	{make([]byte, 16), false},
	{make([]byte, 17), false},
	{[]byte{173, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 175}, false},
	{[]byte{172, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 174}, false},
	{[]byte{173, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 173, 174}, true},
}

func TestResponse_IsCorrect(t *testing.T) {
	for _, pair := range tests {
		got := isCorrect(pair.data)
		if got != pair.isCorrect {
			t.Error(
				"For", pair.data,
				"expected", pair.isCorrect,
				"got", got,
			)
		}
	}
}
