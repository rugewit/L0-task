package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"subscriber/additional"
)

type PostgresDb struct {
	Post     string
	User     string
	Password string
	DbName   string
	DbUrl    string
	Pool     *pgxpool.Pool
	ctx      context.Context
}

func NewPostgresDb(ctx context.Context) (*PostgresDb, error) {
	err := additional.LoadViper("../env/.env")
	if err != nil {
		log.Fatalln("cannot load viper")
		return nil, nil
	}

	dbPort := viper.Get("DB_PORT").(string)
	dbUser := viper.Get("DB_USER").(string)
	dbPassword := viper.Get("DB_PASSWORD").(string)
	dbName := viper.Get("DB_NAME").(string)

	dbUrl := fmt.Sprintf("host=0.0.0.0 port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbPort, dbUser, dbPassword, dbName)

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatalln("cannot connect to postgres sql", err)
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalln("cannot ping postgres sql", err)
		return nil, err
	}
	return &PostgresDb{dbPort, dbUser, dbPassword, dbName,
		dbUrl, pool, ctx}, nil
}

func (db *PostgresDb) Close(ctx context.Context) {
	db.Pool.Close()
}
