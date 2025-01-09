package infrastructure

import (
	"database/sql"
	"errors"
	"merlin/internal/domain"
	"merlin/internal/ports"

	"github.com/google/uuid"
)

type VeiculoRepo struct {
	db *sql.DB
}

func NewVeiculoRepo(db *sql.DB) ports.VeiculoRepo {
	return &VeiculoRepo{db: db}
}

func (v *VeiculoRepo) GetById(id string) (*domain.Veiculo, error) {
	query := "SELECT id, cliente_id, marca, modelo, ano, placa FROM veiculos WHERE id = ?"

	row := v.db.QueryRow(query, id)

	var veiculo domain.Veiculo
	err := row.Scan(&veiculo.Id, &veiculo.Cliente.Id, &veiculo.Marca, &veiculo.Modelo, &veiculo.Ano, &veiculo.Placa)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &veiculo, nil
}

func (v *VeiculoRepo) Create(veiculo *domain.Veiculo) error {
	query := "INSERT INTO veiculos (id, cliente_id, marca, modelo, ano, placa) VALUES (?, ?, ?, ?, ?, ?)"

	veiculo.Id = uuid.New().String()

	_, err := v.db.Exec(query, veiculo.Id, veiculo.Cliente.Id, veiculo.Marca, veiculo.Modelo, veiculo.Ano, veiculo.Placa)
	return err
}

func (v *VeiculoRepo) Update(veiculo *domain.Veiculo) error {
	query := "UPDATE veiculos SET cliente_id = ?, marca = ?, modelo = ?, ano = ?, placa = ? WHERE id = ?"

	_, err := v.db.Exec(query, veiculo.Cliente.Id, veiculo.Marca, veiculo.Modelo, veiculo.Ano, veiculo.Placa, veiculo.Id)
	return err
}

func (v *VeiculoRepo) Delete(id string) error {
	query := "DELETE FROM veiculos WHERE id = ?"

	_, err := v.db.Exec(query, id)
	return err
}

func (v *VeiculoRepo) ListByClienteId(clienteId string) ([]*domain.Veiculo, error) {
	query := "SELECT id, cliente_id, marca, modelo, ano, placa FROM veiculos WHERE cliente_id = ?"

	rows, err := v.db.Query(query, clienteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var veiculos []*domain.Veiculo
	for rows.Next() {
		var veiculo domain.Veiculo
		err := rows.Scan(&veiculo.Id, &veiculo.Cliente.Id, &veiculo.Marca, &veiculo.Modelo, &veiculo.Ano, &veiculo.Placa)
		if err != nil {
			return nil, err
		}
		veiculos = append(veiculos, &veiculo)
	}

	return veiculos, nil
}
