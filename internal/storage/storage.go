package storage

import (
	"database/sql"

	"wildberries_L0/internal/model"
)

type Storage struct {
	db *sql.DB
}

func (st *Storage) SaveBD(orderUid string, b []byte) error {

	_, err := st.db.Exec(`INSERT INTO "Order" (order_uid, order_json) VALUES($1, $2)`, orderUid, b)
	if err != nil {
		return err
	}

	return nil
}

func (st *Storage) RecoveryCacheBD() ([]model.OrderDB, error) {
	orders := []model.OrderDB{}
	s := &model.OrderDB{}

	rows, err := st.db.Query(`SELECT * From "Order"`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		err := rows.Scan(&s.OrderUid, &s.OrderJson)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *s)
	}

	return orders, err
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}
