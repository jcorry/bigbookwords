package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/jcorry/bigbookwords/dictionary-service/proto/dictionary"

	micro "github.com/micro/go-micro"
)

const (
	sourceFile  = "./dictionary.json"
	defaultHost = "datastore:27017"
)

func main() {
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)

	if err != nil {
		log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	} else {
		log.Println("Connected to datastore...")
	}

	service := service{session}

	defer service.GetRepo().Remove()
	defer session.Close()

	repoWords, err := service.GetRepo().GetAll()
	check(err)

	if len(repoWords) == 0 {
		log.Println("Loading words from dictionary file.")
		data, err := ioutil.ReadFile(sourceFile)
		check(err)
		// load words from data into repo
		words := make([]*pb.Word, 0)
		json.Unmarshal(data, &words)
		for _, word := range words {
			err = service.GetRepo().Insert(word)
			check(err)
		}
	} else {
		log.Println(fmt.Sprintf("Found %d words in dictionary", len(repoWords)))
	}

	srv := micro.NewService(
		micro.Name("go.micro.srv.dictionary"),
		micro.Version("latest"),
	)
	srv.Init()

	pb.RegisterDictionaryServiceHandler(srv.Server(), &service)

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}

func check(e error) {
	if e != nil {
		log.Fatalf("Error: %v", e)
		panic(e)
	}
}
