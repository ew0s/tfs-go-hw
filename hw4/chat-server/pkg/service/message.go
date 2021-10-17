package service

import (
	"chat/entities"
	"chat/pkg/repository"
)

type MessageService struct {
	repo repository.Message
}

func NewMessageService(repo repository.Message) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) Create(usrID int, message entities.Message) error {
	return s.repo.Create(usrID, message)
}

func (s *MessageService) GetGlobalMessages() []entities.Message {
	return s.repo.GetGlobalMessages()
}
