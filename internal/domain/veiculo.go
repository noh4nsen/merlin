package domain

type Veiculo struct {
	Id      string
	Cliente Cliente
	Marca   string
	Modelo  string
	Ano     int
	Placa   string
}
