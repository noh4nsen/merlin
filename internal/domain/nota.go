package domain

import "time"

type Servico struct {
	Descricao string
	Custo     float64
}

type Parte struct {
	Nome       string
	Custo      float64
	Quantidade int
}

type Nota struct {
	Id         string
	Cliente    Cliente
	Veiculo    Veiculo
	Data       time.Time
	Servicos   []Servico
	Partes     []Parte
	CustoTotal float64
}
