package main

import (
	"fmt"
	"strings"

	pb "github.com/jcorry/bigbookwords/dictionary-service/proto/dictionary"
	"gopkg.in/mgo.v2"
)

const (
	dbName               = "bigbookwords"
	dictionaryCollection = "dictionary"
)

type Repository interface {
	Get(string) (*pb.Word, error)
	GetAll() ([]*pb.Word, error)
	Search(string) ([]*pb.Word, error)
	Insert(*pb.Word) error
	Close()
	Remove() (*mgo.ChangeInfo, error)
}

type DictionaryRepository struct {
	session *mgo.Session
}

func (repo *DictionaryRepository) Insert(word *pb.Word) error {
	return repo.collection().Insert(word)
}

func (repo *DictionaryRepository) Remove() (*mgo.ChangeInfo, error) {
	return repo.collection().RemoveAll(nil)
}

func (repo *DictionaryRepository) Get(word string) (*pb.Word, error) {
	var allWords []*pb.Word
	err := repo.collection().Find(nil).All(&allWords)
	if err != nil {
		return nil, err
	}

	for i := range allWords {
		if allWords[i].Word == word {
			return allWords[i], nil
		}
	}

	return nil, fmt.Errorf("No word found in dictionary: %s", word)
}

func (repo *DictionaryRepository) GetAll() ([]*pb.Word, error) {
	var words []*pb.Word

	err := repo.collection().Find(nil).All(&words)
	return words, err
}

func (repo *DictionaryRepository) Search(searchTerm string) ([]*pb.Word, error) {
	allWords, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

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

func (repo *DictionaryRepository) Close() {
	repo.session.Close()
}

func (repo *DictionaryRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(dictionaryCollection)
}
