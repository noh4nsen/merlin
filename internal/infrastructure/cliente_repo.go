package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
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
		return nil, fmt.Errorf("erro ao buscar usuario pelo id %s: %v", id, err)
	}

	return &cliente, nil
}

func (c *ClienteRepo) Create(cliente *domain.Cliente) error {
	query := "INSERT INTO clientes (id, nome, telefone, email) VALUES (?, ?, ?, ?)"

	cliente.Id = uuid.New().String()

	_, err := c.db.Exec(query, cliente.Id, cliente.Nome, cliente.Telefone, cliente.Email)
	return fmt.Errorf("erro ao criar usuário de nome %s: %v", cliente.Nome, err)
}

func (c *ClienteRepo) Update(cliente *domain.Cliente) error {
	query := "UPDATE clientes set nome = ?, telefone = ?, email = ? WHERE id = ?"

	_, err := c.db.Exec(query, cliente.Nome, cliente.Telefone, cliente.Email, cliente.Id)
	return fmt.Errorf("erro ao atualizar usuario de id %s: %v", cliente.Id, err)
}

func (c *ClienteRepo) Delete(id string) error {
	query := "DELETE FROM clientes WHERE id = ?"

	_, err := c.db.Exec(query, id)
	return fmt.Errorf("erro ao deletar usuario de id %s: %v", id, err)
}

func (c *ClienteRepo) List() ([]*domain.Cliente, error) {
	query := "SELECT id, nome, telefone, email FROM clientes"

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %v", err)
	}
	defer rows.Close()

	var clientes []*domain.Cliente
	for rows.Next() {
		var cliente domain.Cliente
		err := rows.Scan(&cliente.Id, &cliente.Nome, &cliente.Telefone, &cliente.Email)
		if err != nil {
			return nil, fmt.Errorf("erro ao varrer lista de usuários retornada pelo banco: %v", err)
		}
		clientes = append(clientes, &cliente)
	}

	return clientes, nil
}
