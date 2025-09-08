package db

import (
	"log/slog"

	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/huandu/go-sqlbuilder"
)

type Database struct {
	db *sqlx.DB
}

type ReqQueryTableRow struct {
	Tags    []string
	Table   string
	Builder func(sb *sqlbuilder.SelectBuilder) *sqlbuilder.SelectBuilder
}

func New(source string) (*Database, error) {

	db, err := sqlx.Connect("mysql", source)
	if err != nil {
		return nil, errors.Wrap(err, "cant sqlx.Connect")
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "cant ping")
	}

	return &Database{
		db: db,
	}, nil
}

func (db *Database) Sqlx() *sqlx.DB {
	return db.db
}

func ItemCount(sqlxDB *sqlx.DB, selectBuilder *sqlbuilder.SelectBuilder, logger *slog.Logger) int64 {
	selectBuilder.Limit(-1)
	selectBuilder.Offset(-1)
	queryCount, args := selectBuilder.Select("count(*) as count").Build()

	row := sqlxDB.QueryRow(queryCount, args...)

	var count int64
	err := row.Scan(&count)
	if err != nil && logger != nil {
		logger.Error("cant row.Scan", slog.Any("err", err))
	}

	return count
}

type QueryTableRowRes[T any] struct {
	Item T
	Sb   *sqlbuilder.SelectBuilder
}

func QueryTableRow[T any](sqlxDB *sqlx.DB, req ReqQueryTableRow) (*QueryTableRowRes[T], error) {
	var item T
	itemStruct := sqlbuilder.NewStruct(&item)
	if len(req.Tags) > 0 {
		itemStruct = itemStruct.WithTag(req.Tags...)
	}
	bSb := itemStruct.SelectFrom(req.Table)
	if req.Builder != nil {
		bSb = req.Builder(bSb)
	}
	itemQuery, args := bSb.Build()
	row := sqlxDB.QueryRowx(itemQuery, args...)
	err := row.StructScan(&item)
	if err != nil {
		return nil, errors.Wrap(err, "cant scan item")
	}

	return &QueryTableRowRes[T]{
		Item: item,
		Sb:   bSb,
	}, nil
}

type QueryTableRowsRes[T any] struct {
	Item []T
	Sb   *sqlbuilder.SelectBuilder
}

func QueryTableRows[T any](sqlxDB *sqlx.DB, req ReqQueryTableRow) (*QueryTableRowsRes[T], error) {
	var item T
	itemStruct := sqlbuilder.NewStruct(&item)
	if len(req.Tags) > 0 {
		itemStruct = itemStruct.WithTag(req.Tags...)
	}
	bSb := itemStruct.SelectFrom(req.Table)
	if req.Builder != nil {
		bSb = req.Builder(bSb)
	}
	itemQuery, args := bSb.Build()
	rows, err := sqlxDB.Queryx(itemQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "cant query rows")
	}
	defer rows.Close()

	items := []T{}
	for rows.Next() {
		err = rows.StructScan(&item)
		if err != nil {
			return nil, errors.Wrap(err, "cant scan item")
		}
		items = append(items, item)
	}

	return &QueryTableRowsRes[T]{
		Item: items,
		Sb:   bSb,
	}, nil
}
