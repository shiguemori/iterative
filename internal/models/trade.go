package models

import (
	"time"
)

type Trade struct {
	ID                          uint      `gorm:"primaryKey" json:"id"`
	DataReferencia              time.Time `gorm:"column:data_referencia" json:"data_referencia"`
	CodigoInstrumento           string    `gorm:"column:codigo_instrumento" json:"codigo_instrumento"`
	AcaoAtualizacao             string    `gorm:"column:acao_atualizacao" json:"acao_atualizacao"`
	PrecoNegocio                float64   `gorm:"column:preco_negocio;type:decimal(15,2)" json:"preco_negocio"`
	QuantidadeNegociada         int64     `gorm:"column:quantidade_negociada" json:"quantidade_negociada"`
	HoraFechamento              string    `gorm:"column:hora_fechamento" json:"hora_fechamento"`
	CodigoIdentificadorNegocio  string    `gorm:"column:codigo_identificador_negocio" json:"codigo_identificador_negocio"`
	TipoSessaoPregao            string    `gorm:"column:tipo_sessao_pregao" json:"tipo_sessao_pregao"`
	DataNegocio                 time.Time `gorm:"column:data_negocio;index" json:"data_negocio"`
	CodigoParticipanteComprador string    `gorm:"column:codigo_participante_comprador" json:"codigo_participante_comprador"`
	CodigoParticipanteVendedor  string    `gorm:"column:codigo_participante_vendedor" json:"codigo_participante_vendedor"`
}

// TableName specifies the table name for the Trade model
func (Trade) TableName() string {
	return "trades"
}

// TradeAggregation represents the aggregated data for API response
type TradeAggregation struct {
	Ticker         string  `json:"ticker"`
	MaxRangeValue  float64 `json:"max_range_value"`
	MaxDailyVolume int64   `json:"max_daily_volume"`
}
