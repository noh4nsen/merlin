package infrastructure

import (
	"database/sql"
	"errors"
	"merlin/internal/domain"
	"merlin/internal/ports"

	"github.com/google/uuid"
)

type ClienteRepo struct {
	db *sql.DB
}

func NewClienteRepo(db *sql.DB) ports.ClienteRepo {
	return &ClienteRepo{db: db}
}

func (c *ClienteRepo) GetById(id string) (*domain.Cliente, error) {
	query := "SELECT id, nome, telefone, email FROM clientes WHERE id = ?"
	row := c.db.QueryRow(query, id)

	var cliente domain.Cliente

	err := row.Scan(&cliente.Id, &cliente.Nome, &cliente.Telefone, &cliente.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cliente, nil
}

func (c *ClienteRepo) Create(cliente *domain.Cliente) error {
	query := "INSERT INTO clientes (id, nome, telefone, email) VALUES (?, ?, ?, ?)"

	cliente.Id = uuid.New().String()

	_, err := c.db.Exec(query, cliente.Id, cliente.Nome, cliente.Telefone, cliente.Email)
	return err
}

func (c *ClienteRepo) Update(cliente *domain.Cliente) error {
	query := "UPDATE clientes set nome = ?, telefone = ?, email = ? WHERE id = ?"

	_, err := c.db.Exec(query, cliente.Nome, cliente.Telefone, cliente.Email, cliente.Id)
	return err
}

func (c *ClienteRepo) Delete(id string) error {
	query := "DELETE FROM clientes WHERE id = ?"

	_, err := c.db.Exec(query, id)
	return err
}

func (c *ClienteRepo) List() ([]*domain.Cliente, error) {
	query := "SELECT id, nome, telefone, email FROM clientes"

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []*domain.Cliente
	for rows.Next() {
		var cliente domain.Cliente
		err := rows.Scan(&cliente.Id, &cliente.Nome, &cliente.Telefone, &cliente.Email)
		if err != nil {
			return nil, err
		}
		clientes = append(clientes, &cliente)
	}

	return clientes, nil
}
