package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const DRIVER_NAME = "postgres"

var getTimeNow = time.Now().Truncate(time.Second)

type DBParam struct {
	Name    string
	Age     int
	Address string
}

type DBInterface interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int) (User, error)
	CreateUser(ctx context.Context, dbParam DBParam) error
	UpdateUser(ctx context.Context, id int, dbParam DBParam) error
	DeleteUser(ctx context.Context, id int) error
}

type PostgreStore struct {
	db       *sqlx.DB
	user     string
	password string
	dbName   string
	host     string
	sslMode  string
}

func NewPostgreStore(user, password, dbName, host, sslMode string) *PostgreStore {
	return &PostgreStore{
		user:     user,
		password: password,
		dbName:   dbName,
		host:     host,
		sslMode:  sslMode,
	}
}

func (store *PostgreStore) connect(ctx context.Context) error {
	src := fmt.Sprintf("user=%s dbname=%s sslmode=%v password=%s host=%s", store.user, store.dbName, store.sslMode, store.password, store.host)
	db, err := sqlx.ConnectContext(ctx, DRIVER_NAME, src)
	if err != nil {
		return err
	}

	store.db = db
	return nil
}

func (store *PostgreStore) close() error {
	return store.db.Close()
}

func (store *PostgreStore) GetUsers(ctx context.Context) ([]User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := store.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer store.close()

	var users []User

	query := "SELECT * FROM users"

	if err := store.db.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (store *PostgreStore) GetUserByID(ctx context.Context, id string) (*User, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := store.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer store.close()

	query := "SELECT * FROM users WHERE id = $1"
	user := User{}

	if err := store.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}

	return &user, nil
}

func (store *PostgreStore) CreateUser(ctx context.Context, dbParam *DBParam) error {
	if ctx == nil {
		ctx = context.Background()
	}
	err := store.connect(ctx)
	if err != nil {
		return err
	}
	defer store.close()

	user := User{
		Name:      dbParam.Name,
		Age:       dbParam.Age,
		Address:   dbParam.Address,
		CreatedAt: getTimeNow,
		UpdatedAt: &getTimeNow,
	}

	query := "INSERT INTO users (name, age, address, createdat, updatedat) VALUES (:name, :age, :address, :createdat, :updatedat)"

	if _, err := store.db.NamedExecContext(ctx, query, user); err != nil {
		// if strings.Contains(err.Error(), "SQLSTATE 23505") {
		// 	return &Dup
		// }
		return err
	}

	return nil
}

func (store *PostgreStore) UpdateUser(ctx context.Context, id string, dbParam *DBParam) error {
	if ctx == nil {
		ctx = context.Background()
	}
	err := store.connect(ctx)
	if err != nil {
		return err
	}
	defer store.close()

	user := User{
		Name:      dbParam.Name,
		Age:       dbParam.Age,
		Address:   dbParam.Address,
		UpdatedAt: &getTimeNow,
	}

	query := fmt.Sprintf("UPDATE users SET name = :name, age = :age, address = :address, updatedat = :updatedat WHERE id = %s", id)

	if _, err := store.db.NamedExecContext(ctx, query, user); err != nil {
		return err
	}

	return nil
}

func (store *PostgreStore) DeleteUser(ctx context.Context, id string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	err := store.connect(ctx)
	if err != nil {
		return err
	}
	defer store.close()

	query := "DELETE FROM users WHERE id = $1"

	if _, err := store.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
