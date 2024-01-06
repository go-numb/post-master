package database

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Client struct {
	ctx context.Context
	db  *bun.DB
}

func New(dbsourcepath string, models ...any) *Client {
	ctx := context.Background()
	sqlite, err := sql.Open(sqliteshim.ShimName, dbsourcepath)
	if err != nil {
		log.Fatal().Msgf("does not open database, %v", err)
	}
	sqlite.SetMaxOpenConns(1)

	db := bun.NewDB(sqlite, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Create models table.
	for _, model := range models {
		_, err = db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			log.Info().Msgf("already exists, %v", err)
			continue
		}
	}

	db.RegisterModel(models...)

	return &Client{
		ctx: ctx,
		db:  db,
	}
}

func (p *Client) Close() error {
	return p.db.Close()
}

func (p *Client) Create(data any) error {
	if _, err := p.db.NewInsert().Model(data).Exec(p.ctx); err != nil {
		return err
	}

	return nil
}

func (p *Client) Read(query string, data any) error {
	if _, err := p.db.NewSelect().
		Model(data).
		Where(query).
		Exec(p.ctx); err != nil {
		return err
	}

	return nil
}

func (p *Client) Update(query string, data any) error {
	if _, err := p.db.NewUpdate().
		Model(data).
		Where(query).
		Exec(p.ctx); err != nil {
		return err
	}

	return nil
}

func (p *Client) Delete(privateKey string, data any) error {
	if _, err := p.db.NewDelete().
		Model(data).
		WherePK(privateKey).
		Exec(p.ctx); err != nil {
		return err
	}

	return nil
}
