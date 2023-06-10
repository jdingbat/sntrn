package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jdingbat/sntrn"
	"github.com/joho/godotenv"
)

var (
	username string
	password string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
}

func main() {
	client := sntrn.New(&sntrn.Options{LogLevel: log.FatalLevel})
	defer client.Close()

	err := client.Login(context.TODO(), username, password)
	if err != nil {
		log.Fatal(err)
	}

	sr, err := client.Search(context.TODO(), "family") // IE rick and morty
	if err != nil {
		log.Fatal(err)
	}

	var dl *sntrn.SearchResponse
	for _, s := range sr {
		if s.Id == "1459" {
			dl = &s
			break
		}
	}
	if dl == nil {
		log.Fatal("Failed to find media with given id")
	}

	fmt.Println(dl)

	lr, err := client.Links(context.TODO(), sntrn.SearchResponse{Id: "1459"})
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range lr {
		if l.DlHd != "" {
			fmt.Printf("Download hd media: %s\n", l.DlHd)
		} else {
			fmt.Printf("Download media: %s\n", l.Dl)
		}
	}
}
