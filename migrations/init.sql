CREATE TABLE trades (
   id SERIAL PRIMARY KEY,
   data_referencia DATE NOT NULL,
   codigo_instrumento VARCHAR(20) NOT NULL,
   acao_atualizacao VARCHAR(10) NULL,
   preco_negocio NUMERIC(15,2) NOT NULL,
   quantidade_negociada BIGINT NOT NULL,
   hora_fechamento VARCHAR(10) NULL,
   codigo_identificador_negocio VARCHAR(4) NULL,
   tipo_sessao_pregao VARCHAR(10) NULL,
   data_negocio DATE NULL,
   codigo_participante_comprador VARCHAR(50) NULL,
   codigo_participante_vendedor VARCHAR(50) NULL
);

CREATE INDEX idx_trade_csv_data_negocio ON trades(data_negocio);
CREATE INDEX idx_trade_csv_codigo_instrumento ON trades(codigo_instrumento);
CREATE INDEX idx_trades_composite ON trades(codigo_instrumento, data_negocio);
