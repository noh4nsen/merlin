package infrastructure

import "database/sql"

func NewDBConnection(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS clientes (
			id TEXT PRIMARY KEY,
			nome TEXT,
			telefone TEXT,
			email TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS veiculos (
			id TEXT PRIMARY KEY,
			cliente_id TEXT,
			marca TEXT,
			modelo TEXT,
			ano INTEGER,
			placa TEXT,
			FOREIGN KEY (cliente_id) references clientes (id)
		);`,
		`CREATE TABLE IF NOT EXISTS notas (
			id TEXT PRIMARY KEY,
			cliente_id TEXT,
			veiculo_id TEXT,
			data TEXT,
			custo_total REAL,
			FOREIGN KEY (cliente_id) REFERENCES clientes (id)
			FOREIGN KEY (veiculo_id) REFERENCES veiculos (id)
		);`,
		`CREATE TABKE IF NOT EXISTS servicos (
			nota_id TEXT,
			descricao TEXT,
			custo REAL,
			FOREIGN KEY (nota_id) REFERENCES notas (id)
		);`,
		`CREATE TABLE IF NOT EXISTS partes (
			nota_id TEXT,
			nome TEXT,
			custo REAL,
			quantidade INTEGER,
			FOREIGN KEY (nota_id) REFERENCES notas (id)
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
