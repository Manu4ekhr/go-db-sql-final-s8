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
	// добавление строки в таблицу parcel, заполняя данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client,address,status,created_at) VALUES (:client, :address, :status, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("address", p.Address),
		sql.Named("status", p.Status),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	row := s.db.QueryRow("SELECT number,client,address,status,created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Address, &p.Status, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	rows, err := s.db.Query("SELECT number,client,address,status,created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Address, &p.Status, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address), sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}
