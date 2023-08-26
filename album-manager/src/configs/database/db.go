package database

import (
	"album-manager/src/configs"
	"context"
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	DB    *gorm.DB
	Ctx   context.Context
	sqlDB *sql.DB
}

func (pgC *PostgresConfig) Close() {
	pgC.sqlDB.Close()
}

func InitDatabase(ctx context.Context) (*PostgresConfig, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		configs.Env.Postgres.Username,
		configs.Env.Postgres.Password,
		configs.Env.Postgres.Host,
		configs.Env.Postgres.Port,
		configs.Env.Postgres.Database,
	)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &PostgresConfig{db, ctx, sqlDB}, nil
}

func (pgC *PostgresConfig) InitializeFunction() {
	pgC.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	pgC.DB.Exec(`
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status_enum') THEN
			CREATE TYPE user_status_enum AS ENUM ('active', 'inactive');
		END IF;
	END $$;
	`)
}
