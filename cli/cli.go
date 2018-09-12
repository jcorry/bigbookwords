package main

import (
	"log"
	"time"

	pb "github.com/jcorry/bigbookwords/dictionary-service/proto/dictionary"

	"github.com/micro/go-micro/cmd"

	"fmt"

	microclient "github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

const (
	address = "localhost:50052"
)

func main() {
	cmd.Init()

	client := pb.NewDictionaryServiceClient("go.micro.srv.dictionary", microclient.DefaultClient)

	// Get all the words
	start := time.Now()
	getAll, err := client.GetWords(context.Background(), &pb.GetRequest{})
	log.Printf("Word count: %d", len(getAll.GetWords()))
	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Took %s to get all words", elapsed))

	// Get a single word
	start = time.Now()
	result, err := client.GetWord(context.Background(), &pb.GetRequest{Query: `loath`})
	if err != nil {
		log.Fatalf("Could not get single word: %v", err)
	}
	log.Printf("The word %s appears %d times", result.Word.GetWord(), result.Word.GetAppearances())
	fmt.Println("")
	elapsed = time.Since(start)
	fmt.Println(fmt.Sprintf("Took %s to get word '%s'", elapsed, result.Word.GetWord()))

	fmt.Println("Search for a word:")
	start = time.Now()
	searchTerm := "will"
	log.Printf("Searching for term: %s", searchTerm)
	results, err := client.Search(context.Background(), &pb.GetRequest{Query: searchTerm})
	elapsed = time.Since(start)
	fmt.Println(fmt.Sprintf("Took %s to search words matching '%s'", elapsed, searchTerm))
	if err != nil {
		log.Fatalf("Count not get words for search term: %s", searchTerm)
		log.Fatalf("Error: %v", err)
	}

	for _, word := range results.GetWords() {
		fmt.Println(fmt.Sprintf("match: %s", word.GetWord()))
	}
}
