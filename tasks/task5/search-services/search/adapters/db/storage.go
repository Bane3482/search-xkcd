package db

import (
	"context"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
	"yadro.com/course/search/core"
)

type DB struct {
	conn *sqlx.DB
	log  *slog.Logger
}

func New(log *slog.Logger, address string) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	return &DB{
		log:  log,
		conn: db,
	}, nil
}

func (db *DB) Search(ctx context.Context, keyword string) ([]int, error) {
	var IDs []int
	err := db.conn.SelectContext(
		ctx,
		&IDs,
		"SELECT comics_id FROM comics WHERE $1=ANY(words)",
		keyword,
	)

	if err != nil {
		db.log.Error("db client search", "error", err)
		return nil, err
	}

	return IDs, nil
}

func (db *DB) Get(ctx context.Context, id int) (core.ComicsInfo, error) {
	var comics core.ComicsInfo

	err := db.conn.GetContext(
		ctx,
		&comics,
		"SELECT comics_id, comics_url FROM comics WHERE comics_id=$1",
		id,
	)

	if err != nil {
		db.log.Error("db client get", "error", err)
		return core.ComicsInfo{}, err
	}

	return comics, nil
}
