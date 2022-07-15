package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

func TestRunCmd(t *testing.T) {
	t.Run("exit codes", func(t *testing.T) {
		env := Environment{}
		require.Equal(t, 0, RunCmd([]string{"ls", "testdata"}, env))
		require.Equal(t, 111, RunCmd([]string{"command-not-exist"}, env))
	})

	t.Run("apply env vars", func(t *testing.T) {
		env := Environment{"HELLO": "WORLD"}
		os.Setenv("PRE", "SOME")
		output := capturer.CaptureStdout(func() {
			RunCmd([]string{"env"}, env)
		})
		require.Contains(t, output, "PRE=SOME")
		require.Contains(t, output, "HELLO=WORLD")
	})

	t.Run("update env vars", func(t *testing.T) {
		env := Environment{"HELLO": "WORLD"}
		os.Setenv("HELLO", "NEWWORLD")
		output := capturer.CaptureStdout(func() {
			RunCmd([]string{"env"}, env)
		})
		require.Contains(t, output, "HELLO=WORLD")
	})

	t.Run("unset env vars", func(t *testing.T) {
		env := Environment{"HELLO": ""}
		os.Setenv("HELLO", "VAL")
		output := capturer.CaptureStdout(func() {
			RunCmd([]string{"env"}, env)
		})
		require.NotContains(t, output, "HELLO")
	})

	t.Run("unset env var, skip setting value, if =", func(t *testing.T) {
		env := Environment{"HELLO": "SOME=VAL"}
		os.Setenv("HELLO", "VAL")
		output := capturer.CaptureStdout(func() {
			RunCmd([]string{"env"}, env)
		})
		require.NotContains(t, output, "HELLO")
	})
}
