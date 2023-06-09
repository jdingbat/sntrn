package sntrn_test

import (
	"context"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/jdingbat/sntrn"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
}

func TestLinks(t *testing.T) {
	client := sntrn.New(&sntrn.Options{LogLevel: log.DebugLevel})

	err := client.Login(context.Background(), username, password)
	if err != nil {
		t.Errorf("Failed to login")
	}

	searchResponse := sntrn.SearchResponse{
		Id:         "16755",
		Title:      "Top Gun: Maverick",
		Year:       "2022",
		Genre:      "Action, Drama",
		ImdbRating: "8.3",
		Poster:     "https://image.tmdb.org/t/p/w342/62HCnUTziyWcpDaBO2i1DX17ljH.jpg",
		Link:       "/movies/16755-Top-Gun:-Maverick",
	}

	_, err = client.Links(context.Background(), searchResponse)
	if err != nil {
		t.Error(err)
	}

	defer client.Close()
}
