package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	pb "github.com/jcorry/bigbookwords/dictionary-service/proto/dictionary"

	micro "github.com/micro/go-micro"
)

const (
	sourceFile = "./dictionary.json"
)

type Repository interface {
	Get(string) (*pb.Word, error)
	GetAll() []*pb.Word
	Search(string) ([]*pb.Word, error)
}

type DictionaryRepository struct {
	words []*pb.Word
}

func (repo *DictionaryRepository) Get(word string) (*pb.Word, error) {
	allWords := repo.GetAll()
	for i := range allWords {
		if allWords[i].Word == word {
			return allWords[i], nil
		}
	}
	return nil, fmt.Errorf("No word found in dictionary: %s", word)
}

func (repo *DictionaryRepository) GetAll() []*pb.Word {
	return repo.words
}

func (repo *DictionaryRepository) Search(searchTerm string) ([]*pb.Word, error) {
	allWords := repo.GetAll()
	var returnWords []*pb.Word
	for i := range allWords {
		if strings.Contains(allWords[i].Word, searchTerm) {
			returnWords = append(returnWords, allWords[i])
		}
	}
	if len(returnWords) == 0 {
		return nil, fmt.Errorf("No word found in dictionary matching search string: %s", searchTerm)
	}
	return returnWords, nil
}

type service struct {
	repo DictionaryRepository
}

func (s *service) GetWord(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	word, err := s.repo.Get(req.Query)
	if err != nil {
		return err
	}
	res.Word = word
	return nil
}

func (s *service) GetWords(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	words := s.repo.GetAll()
	res.Words = words
	return nil
}

func (s *service) Search(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	words, err := s.repo.Search(req.Query)
	if err != nil {
		return err
	}
	res.Words = words
	return nil
}

func main() {
	repo := &DictionaryRepository{}

	if len(repo.words) == 0 {
		log.Println("Loading words from dictionary file.")
		data, err := ioutil.ReadFile(sourceFile)
		check(err)
		err = json.Unmarshal(data, &repo.words)
		check(err)

		log.Printf("Loaded %d words", len(repo.words))
	}

	srv := micro.NewService(
		micro.Name("go.micro.srv.dictionary"),
		micro.Version("latest"),
	)
	srv.Init()

	pb.RegisterDictionaryServiceHandler(srv.Server(), &service{repo})

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
