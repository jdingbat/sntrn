package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jdingbat/sntrn"
)

var (
	username string
	password string
)

func init() {
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

	sr, err := client.Search(context.TODO(), "rick") // IE rick and morty
	if err != nil {
		log.Fatal(err)
	}

	var dl *sntrn.SearchResponse
	for _, s := range sr {
		if s.Id == "705" {
			dl = &s
			break
		}
	}
	if dl == nil {
		log.Fatal("Failed to find media with given id")
	}

	lr, err := client.Links(context.TODO(), *dl)
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
