package db

import (
	"context"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"yadro.com/course/update/core"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
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

func (db *DB) Add(ctx context.Context, comics core.Comics) error {
	_, err := db.conn.Exec(`INSERT INTO update_schema.comics (comics_id, comics_url, words) VALUES ($1, $2, $3);`, comics.ID, comics.URL, comics.Words)

	if err != nil {
		db.log.Error("add db exec", "error", err, "comics", comics)
	}

	return err
}

func (db *DB) Stats(ctx context.Context) (core.DBStats, error) {
	rows, err := db.conn.QueryContext(ctx, `SELECT words FROM update_schema.comics;`)

	var stats core.DBStats

	err = db.conn.GetContext(ctx, &stats.ComicsFetched, `SELECT COUNT(*) FROM update_schema.comics`)

	if err != nil {
		db.log.Error("stats db get", "error", err)
		return core.DBStats{}, err
	}

	defer rows.Close()

	wordsMap := make(map[string]struct{})

	for rows.Next() {
		var words []string

		if err := rows.Scan(pq.Array(&words)); err != nil {
			db.log.Error("stats db scan", "error", err)
			return core.DBStats{}, err
		}

		stats.WordsTotal += len(words)

		for _, word := range words {
			wordsMap[word] = struct{}{}
		}
	}

	stats.WordsUnique = len(wordsMap)

	return stats, nil
}

func (db *DB) IDs(ctx context.Context) ([]int, error) {
	rows, err := db.conn.QueryContext(ctx, `SELECT comics_id FROM update_schema.comics;`)

	if err != nil {
		db.log.Error("ids db query", "error", err)
		return nil, err
	}

	defer rows.Close()

	idList := make([]int, 0)

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			db.log.Error("ids db scan", "error", err)
			return nil, err
		}

		idList = append(idList, id)
	}

	return idList, nil
}

func (db *DB) Drop(ctx context.Context) error {
	_, err := db.conn.ExecContext(ctx, `DELETE FROM update_schema.comics;`)

	if err != nil {
		db.log.Error("drop db exec", "error", err)
	}

	return err
}
