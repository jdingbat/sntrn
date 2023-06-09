package sntrn_test

import (
	"errors"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/misterbianco/sntrn/sntrn"
)

func TestEarlyClose(t *testing.T) {
	client := sntrn.New(&sntrn.Options{LogLevel: log.DebugLevel})

	err := client.Close()
	if errors.Is(err, sntrn.ErrFailedLogin) {
		t.Errorf("Errors matched")
	}
}
