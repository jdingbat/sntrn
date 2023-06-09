package sntrn_test

import (
	"context"
	"errors"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/misterbianco/sntrn/sntrn"
)

// Could totally mock this but who cares rn?
func TestBadLogin(t *testing.T) {
	client := sntrn.New(&sntrn.Options{LogLevel: log.DebugLevel})

	err := client.Login(context.Background(), "", "")
	if !errors.Is(err, sntrn.ErrFailedLogin) {
		t.Errorf("Expected to fail but did not")
	}
}

func TestBadLoginWithClose(t *testing.T) {
	client := sntrn.New(&sntrn.Options{LogLevel: log.DebugLevel})

	err := client.Login(context.Background(), "", "")
	if !errors.Is(err, sntrn.ErrFailedLogin) {
		t.Errorf("Expected to fail but did not")
	}

	err = client.Close()
	if err != nil {
		t.Errorf("Unexpected failure closing client: %v", err)
	}
}
