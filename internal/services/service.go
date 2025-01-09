package services

import "merlin/internal/ports"

type Service struct {
	clienteRepo ports.ClienteRepo
	veiculoRepo ports.VeiculoRepo
	notaRepo    ports.NotaRepo
}

func NewService(clienteRepo ports.ClienteRepo, veiculoRepo ports.VeiculoRepo, notaRepo ports.NotaRepo) ports.Service {
	return &Service{
		clienteRepo: clienteRepo,
		veiculoRepo: veiculoRepo,
		notaRepo:    notaRepo,
	}
}
