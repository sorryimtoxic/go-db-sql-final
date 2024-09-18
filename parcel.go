package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:n, :c, :s, :a, :c_a)",
		sql.Named("c", p.Client),
		sql.Named("s", p.Status),
		sql.Named("a", p.Address),
		sql.Named("c_a", p.CreatedAt))
	// верните идентификатор последней добавленной записи
	if err != nil {
		return 0, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row, err := s.db.Query("SELECT client, status, adress, created_at FROM parcel WHERE number = :num", sql.Named("num", number))
	if err != nil {
		return Parcel{}, err
	}
	defer row.Close()
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err = row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	rows, err := s.db.Query("SELECT number, status, adress, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :num",
		sql.Named("status", status),
		sql.Named("num", number))

	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :addr WHERE number = :num AND status = 'registered'",
		sql.Named("addr", address),
		sql.Named("num", number))

	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :num AND status = 'registered'",
		sql.Named("num", number))

	if err != nil {
		return err
	}
	return nil
}
