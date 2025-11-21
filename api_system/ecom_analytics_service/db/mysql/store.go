package db

import (
	"database/sql"

	db_agent_ai_db "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/agent_ai_db"
	db_interact "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/interact"
	db_order "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/order"
	db_transaction "github.com/TranVinhHien/ecom_analytics_service/db/sqlc/transaction"
)

type SQLStoreOrder struct {
	*db_order.Queries
	connPool *sql.DB
}

type StoreOrder interface {
	db_order.Querier
}

// create new store

func NewStoreOrder(connectPool *sql.DB) StoreOrder {
	return &SQLStoreOrder{
		Queries:  db_order.New(connectPool),
		connPool: connectPool,
	}
}

type SQLStoreTransaction struct {
	*db_transaction.Queries
	connPool *sql.DB
}

type StoreTransaction interface {
	db_transaction.Querier
}

// create new store

func NewStoreTransaction(connectPool *sql.DB) StoreTransaction {
	return &SQLStoreTransaction{
		Queries:  db_transaction.New(connectPool),
		connPool: connectPool,
	}
}

type SQLStoreInteract struct {
	*db_interact.Queries
	connPool *sql.DB
}

type StoreInteract interface {
	db_interact.Querier
}

// create new store

func NewStoreInteract(connectPool *sql.DB) StoreInteract {
	return &SQLStoreInteract{
		Queries:  db_interact.New(connectPool),
		connPool: connectPool,
	}
}

type SQLStoreAgentAIDB struct {
	*db_agent_ai_db.Queries
	connPool *sql.DB
}

type StoreAgentAIDB interface {
	db_agent_ai_db.Querier
}

// create new store

func NewStoreAgentAIDB(connectPool *sql.DB) StoreAgentAIDB {
	return &SQLStoreAgentAIDB{
		Queries:  db_agent_ai_db.New(connectPool),
		connPool: connectPool,
	}
}
