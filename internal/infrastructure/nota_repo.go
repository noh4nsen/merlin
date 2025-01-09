package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"merlin/internal/domain"
	"merlin/internal/ports"
	"time"

	"github.com/google/uuid"
)

type NotaRepo struct {
	db *sql.DB
}

func NewNotaRepo(db *sql.DB) ports.NotaRepo {
	return &NotaRepo{db: db}
}

func (n *NotaRepo) GetById(id string) (*domain.Nota, error) {
	query := "SELECT id, cliente_id, veiculo_id, data, custo_total FROM notas WHERE id = ?"

	var nota domain.Nota
	var dataStr string

	row := n.db.QueryRow(query, id)
	err := row.Scan(&nota.Id, &nota.Cliente.Id, &nota.Veiculo.Id, &dataStr, &nota.CustoTotal)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar nota de id %s: %v", id, err)
	}

	nota.Data, err = time.Parse(time.RFC3339, dataStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao transfomar data de nota: %v", err)
	}

	nota.Servicos, err = n.getServicosByNotaId(nota.Id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar serviços atrelados a nota de id %s: %v", nota.Id, err)
	}
	nota.Partes, err = n.getPartesByNotaId(nota.Id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar partes atreladas a nota de id %s: %v", nota.Id, err)
	}

	return &nota, nil
}

func (n *NotaRepo) Create(nota *domain.Nota) error {
	query := "INSERT INTO notas (id, cliente_id, veiculo_id, data, custo_total) VALUES (?, ?, ?, ?, ?)"

	nota.Id = uuid.New().String()

	_, err := n.db.Exec(query, nota.Id, nota.Cliente.Id, nota.Veiculo.Id, nota.Data.Format(time.RFC3339), nota.CustoTotal)
	if err != nil {
		return fmt.Errorf("erro ao criar nova nota: %v", err)
	}

	err = n.insertServicos(nota.Id, nota.Servicos)
	if err != nil {
		return fmt.Errorf("erro ao criar serviços atrelados a nova nota: %v", err)
	}
	err = n.insertPartes(nota.Id, nota.Partes)
	if err != nil {
		return fmt.Errorf("erro ao criar partes atreladas a nova nota: %v", err)
	}

	return nil
}

func (n *NotaRepo) Update(nota *domain.Nota) error {
	query := "UPDATE notas SET cliente_id = ?, veiculo_id = ?, data = ?, custo_total = ? WHERE id = ?"

	_, err := n.db.Exec(query, nota.Cliente.Id, nota.Veiculo.Id, nota.Data.Format(time.RFC3339), nota.CustoTotal, nota.Id)
	if err != nil {
		return err
	}

	err = n.deleteServicosByNotaId(nota.Id)
	if err != nil {
		return err
	}
	err = n.insertServicos(nota.Id, nota.Servicos)
	if err != nil {
		return err
	}

	err = n.deletePartesByNotaId(nota.Id)
	if err != nil {
		return err
	}
	err = n.insertPartes(nota.Id, nota.Partes)
	if err != nil {
		return err
	}

	return nil
}

func (n *NotaRepo) Delete(id string) error {
	query := "DELETE FROM notas WHERE id = ?"

	err := n.deleteServicosByNotaId(id)
	if err != nil {
		return err
	}

	err = n.deletePartesByNotaId(id)
	if err != nil {
		return err
	}

	_, err = n.db.Exec(query, id)
	return err
}

func (n *NotaRepo) List() ([]*domain.Nota, error) {
	query := "SELECT id, cliente_id, veiculo_id, data, custo_total FROM notas"

	rows, err := n.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notas []*domain.Nota
	for rows.Next() {
		var nota domain.Nota
		var dataStr string

		err := rows.Scan(&nota.Id, &nota.Cliente.Id, &nota.Veiculo.Id, &dataStr, &nota.CustoTotal)
		if err != nil {
			return nil, err
		}

		nota.Data, err = time.Parse(time.RFC3339, dataStr)
		if err != nil {
			return nil, err
		}

		nota.Servicos, err = n.getServicosByNotaId(nota.Id)
		if err != nil {
			return nil, err
		}

		nota.Partes, err = n.getPartesByNotaId(nota.Id)
		if err != nil {
			return nil, err
		}

		notas = append(notas, &nota)
	}

	return notas, nil
}

func (n *NotaRepo) getServicosByNotaId(notaId string) ([]domain.Servico, error) {
	query := "SELECT descricao, custo FROM servicos WHERE nota_id = ?"

	rows, err := n.db.Query(query, notaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servicos []domain.Servico
	for rows.Next() {
		var servico domain.Servico
		err := rows.Scan(&servico.Descricao, &servico.Custo)
		if err != nil {
			return nil, err
		}

		servicos = append(servicos, servico)
	}

	return servicos, nil
}

func (n *NotaRepo) getPartesByNotaId(notaId string) ([]domain.Parte, error) {
	query := "SELECT nome, custo, quantidade FROM partes WHERE nota_id = ?"

	rows, err := n.db.Query(query, notaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partes []domain.Parte
	for rows.Next() {
		var parte domain.Parte
		err := rows.Scan(&parte.Nome, &parte.Custo, &parte.Quantidade)
		if err != nil {
			return nil, err
		}
		partes = append(partes, parte)
	}

	return partes, nil
}

func (n *NotaRepo) insertServicos(notaId string, servicos []domain.Servico) error {
	query := "INSERT INTO servicos (nota_id, descricao, custo) VALUES (?, ?, ?)"

	for _, servico := range servicos {
		_, err := n.db.Exec(query, notaId, servico.Descricao, servico.Custo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NotaRepo) insertPartes(notaId string, partes []domain.Parte) error {
	query := "INSERT INTO partes (nota_id, nome, custo, quantidade) VALUES (?, ?, ?, ?)"

	for _, parte := range partes {
		_, err := n.db.Exec(query, notaId, parte.Nome, parte.Custo, parte.Quantidade)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NotaRepo) deleteServicosByNotaId(notaId string) error {
	query := "DELETE FROM servicos WHERE nota_id = ?"

	_, err := n.db.Exec(query, notaId)
	return err
}

func (n *NotaRepo) deletePartesByNotaId(notaId string) error {
	query := "DELETE FROM partes WHERE nota_id = ?"

	_, err := n.db.Exec(query, notaId)
	return err
}
