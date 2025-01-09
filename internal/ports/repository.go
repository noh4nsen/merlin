package ports

import "merlin/internal/domain"

type ClienteRepo interface {
	GetById(id string) (*domain.Cliente, error)
	Create(cliente *domain.Cliente) error
	Update(cliente *domain.Cliente) error
	Delete(id string) error
	List() ([]*domain.Cliente, error)
}

type VeiculoRepo interface {
	GetById(id string) (*domain.Veiculo, error)
	Create(veiculo *domain.Veiculo) error
	Update(veiculo *domain.Veiculo) error
	Delete(id string) error
	ListByClienteId(id string) ([]*domain.Veiculo, error)
}

type NotaRepo interface {
	GetById(id string) (*domain.Nota, error)
	Create(nota *domain.Nota) error
	Update(nota *domain.Nota) error
	Delete(id string) error
	List() ([]*domain.Nota, error)
}
