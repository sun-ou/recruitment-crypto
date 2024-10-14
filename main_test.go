package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleCommandHelp(t *testing.T) {
	assert.Contains(t, GetIoBuffer("?"), "usage", "unexpected result")
}

func TestHandleCommandUnknown(t *testing.T) {
	assert.Contains(t, GetIoBuffer("><"), "Unknown command", "unexpected result")
	assert.Contains(t, GetIoBuffer(">< <>"), "Unknown command", "unexpected result")
}

func GetIoBuffer(command string) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	handleCommand(command)

	w.Close()

	out, _ := io.ReadAll(r)

	// restore the stdout
	os.Stdout = old
	return string(out)
}
