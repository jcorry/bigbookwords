package main

import (
	"context"

	"gopkg.in/mgo.v2"

	pb "github.com/jcorry/bigbookwords/dictionary-service/proto/dictionary"
)

type service struct {
	session *mgo.Session
}

func (s *service) GetRepo() Repository {
	return &DictionaryRepository{s.session.Clone()}
}

func (s *service) GetWord(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	repo := s.GetRepo()
	word, err := repo.Get(req.Query)
	if err != nil {
		return err
	}
	res.Word = word
	return nil
}

func (s *service) GetWords(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	repo := s.GetRepo()
	words, err := repo.GetAll()
	if err != nil {
		return err
	}
	res.Words = words
	return nil
}

func (s *service) Search(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	repo := s.GetRepo()
	words, err := repo.Search(req.Query)
	if err != nil {
		return err
	}
	res.Words = words
	return nil
}
