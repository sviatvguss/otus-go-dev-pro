package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const dir = "temp"

func setup() {
	if err := os.Mkdir(dir, 0o755); err != nil {
		panic("can not create test directory")
	}
}

func clear() {
	if err := os.RemoveAll(dir); err != nil {
		panic("can not remove test directory")
	}
}

func prepare(file string, content string) {
	f, err := os.Create(path.Join(dir, file))
	if err != nil {
		panic("can not create test file")
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		panic("can not write in test file")
	}
}

func TestReadDir(t *testing.T) {
	t.Run("correct read", func(t *testing.T) {
		setup()
		defer clear()

		prepare("HELLO", "WORLD")
		prepare("WITH_SPACES", "  VAL 	 ")
		prepare("EMPTY", "")
		prepare("WITH_TERM_NULLS", "V\x00A\000L\u0000")
		prepare("WITH_MULTLINE", "FIRST LINE\nSECOND LINE\nTHIRD LINE")

		result, err := ReadDir(dir)

		require.Nil(t, err)
		require.Equal(t, Environment{
			"HELLO":           "WORLD",
			"WITH_SPACES":     "  VAL",
			"EMPTY":           "",
			"WITH_TERM_NULLS": "V\nA\nL\n",
			"WITH_MULTLINE":   "FIRST LINE",
		}, result)
	})
}
