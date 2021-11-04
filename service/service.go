package service

import "context"

type Service interface {
	Run(ctx context.Context, shutdownCh <-chan struct{}) error
}

type Manager struct {
}

func NewManager() (*Manager, error) {
	return nil, nil
}

func AddService(svc Service) error {
	return nil
}

func Run(ctx context.Context) error {
	return nil
}

func Shutdown(ctx context.Context) error {
	return nil
}
