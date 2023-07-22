package sntrn_test

import (
	"context"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/jdingbat/sntrn"
	"github.com/joho/godotenv"
)

var username string
var password string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
}

func TestSearch(t *testing.T) {
	client := sntrn.New(&sntrn.Options{LogLevel: log.DebugLevel})

	err := client.Login(context.Background(), username, password)
	if err != nil {
		t.Errorf("Failed to login")
	}

	_, err = client.Search(context.Background(), "rick")
	if err != nil {
		t.Error(err)
	}

	defer client.Close()
}
