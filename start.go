package main

import (
	"auth-api-go/internal"
	"auth-api-go/internal/data"
	"auth-api-go/internal/service/config"
	"auth-api-go/internal/service/server"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

func main() {
	err := start()
	if err != nil {
		fmt.Println("error has occurred: ", err.Error())
	}
}

func start() error {
	log.Info().Msg("starting auth-api-go...")
	err, cfg := config.LoadConfig()
	if err != nil {
		log.Error().Msgf("error loading configuration: %s", err.Error())

		return err
	}

	log.Info().Msg("configuration loaded correctly")
	logConfig(cfg)

	// get the data connection pool
	log.Info().Msg("initializing connection pool")
	db, err := getConnectionPool(cfg)
	log.Info().Msg("connection pool properly initialized")
	if err != nil {
		log.Info().Msgf("error initializing the connection pool: %s", err.Error())

		return err
	}

	// now start the routing gorilla mux tech
	var router = internal.Start(cfg, db)

	log.Info().Msg("starting server...")
	err = server.StartServer(cfg, router)
	if err != nil {
		return err
	}

	return nil
}

func logConfig(cfg data.Config) {
	log.Info().Msgf("server port: %s", cfg.Port)
	log.Info().Msgf("token duration [ms]: %d", cfg.TokenDurationMs)
	log.Info().Msgf("cors: %t", cfg.Cors)
	log.Info().Msgf("db user: %s", cfg.Db.User)
	log.Info().Msgf("db machine: %s", cfg.Db.Machine)
	log.Info().Msgf("db port: %d", cfg.Db.Port)
	log.Info().Msgf("db poolsize: %d", cfg.Db.PoolSize)
}

func getConnectionPool(config data.Config) (*sql.DB, error) {
	dbConfiguration := config.Db
	user := dbConfiguration.User
	password := dbConfiguration.Password
	machine := dbConfiguration.Machine
	port := dbConfiguration.Port
	database := dbConfiguration.Database

	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, machine, port, database)
	db, err := sql.Open("mysql", url)

	if err != nil {
		return nil, err
	}

	return db, nil
}
