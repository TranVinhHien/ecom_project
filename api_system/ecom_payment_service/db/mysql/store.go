package db

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/TranVinhHien/ecom_payment_service/db/sqlc"
	entity "github.com/TranVinhHien/ecom_payment_service/services/entity"
)

type SQLStore struct {
	*db.Queries
	connPool *sql.DB
}

type Store interface {
	db.Querier
	ExecTS(ctx context.Context, fn func(tx db.Querier) error) error
}

// create new store

func NewStore(connectPool *sql.DB) Store {
	return &SQLStore{
		Queries:  db.New(connectPool),
		connPool: connectPool,
	}
}

// write a function transaction using package github.com/ja ckc/pgx/v5/pgxpool
func (s *SQLStore) ExecTS(ctx context.Context, fn func(tx db.Querier) error) error {
	tx, err := s.connPool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		fmt.Printf("lôi err")
		if errTran := tx.Rollback(); errTran != nil {
			return fmt.Errorf("transaction error %v ,rollback trancsaction error : %v", err, errTran)
		}
		return err
	}

	return tx.Commit()
}
func listData(ctx context.Context, connPool *sql.DB, table string, query entity.QueryFilter) (*sql.Rows, int, error) {
	querySQL := fmt.Sprintf("SELECT *  FROM %s WHERE 1=1", table)
	querySQLCount := fmt.Sprintf("SELECT COUNT(*) as totalElements  FROM %s WHERE 1=1", table)
	args := []interface{}{}
	argsCount := []interface{}{}

	// Xây dựng SQL từ QueryFilter
	for _, condition := range query.Conditions {
		querySQL += fmt.Sprintf(" AND %s %s ?", condition.Field, condition.Operator)
		querySQLCount += fmt.Sprintf(" AND %s %s ?", condition.Field, condition.Operator)
		args = append(args, condition.Value)
		argsCount = append(argsCount, condition.Value)
	}

	// Thêm sắp xếp
	if query.OrderBy != nil {
		querySQL += fmt.Sprintf(" ORDER BY %s %s", query.OrderBy.Field, query.OrderBy.Value)
	}

	// Thêm phân trang
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		querySQL += " LIMIT ? OFFSET ?"
		args = append(args, query.PageSize, offset)
	}

	// Thực thi truy vấn

	rows, err := connPool.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error connPool.QueryContext: %s", err.Error())
	}
	row := connPool.QueryRowContext(ctx, querySQLCount, argsCount...)
	var sc int64
	err = row.Scan(&sc)
	if err != nil {
		return nil, 0, fmt.Errorf("error row.Scan totalElements: %s", err.Error())
	}
	totalElements := int(sc)
	return rows, totalElements, nil
}
