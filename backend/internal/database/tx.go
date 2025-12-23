package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager управляет транзакциями и создаёт репозитории с tx.
type TxManager struct {
	pool *pgxpool.Pool
}

// NewTxManager создаёт новый TxManager.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// RunInTx выполняет функцию в рамках транзакции.
// Querier передаётся в функцию для создания репозиториев.
func (m *TxManager) RunInTx(ctx context.Context, fn func(ctx context.Context, tx Querier) error) error {
	return WithTx(ctx, m.pool, fn)
}

// Pool возвращает пул соединений с БД.
func (m *TxManager) Pool() *pgxpool.Pool {
	return m.pool
}
