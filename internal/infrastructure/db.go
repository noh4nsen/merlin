package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func NewDBConnection(appCtx context.Context) (*sql.DB, error) {
	dataDir, err := getAppDataDirectory()
	if err != nil {
		return nil, fmt.Errorf("erro durante abertura de conexão com o banco: %v", err)
	}
	dbFile := filepath.Join(dataDir, "merlin.db")

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %v", err)
	}

	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("erro durante abertura de conexão com o banco: %v", err)
	}

	return db, nil
}

func getAppDataDirectory() (string, error) {
	var dataDir string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao buscar diretório do usuário: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		dataDir = filepath.Join(os.Getenv("APPDATA"), "Merlin")
	default:
		dataDir = filepath.Join(homeDir, ".config", "Merlin")
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			return "", fmt.Errorf("erro ao criar diretorio de configuração: %v", err)
		}
	}

	return dataDir, nil
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
			return fmt.Errorf("erro ao criar tabelas: %v", err)
		}
	}

	return nil
}
