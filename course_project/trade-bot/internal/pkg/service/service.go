package service

import "trade-bot/internal/pkg/repository"

type Authorization interface {
}

type Service struct {
	Authorization
}

func NewService(r *repository.Repository) *Service {
	return &Service{}
}
