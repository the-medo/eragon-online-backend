package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/the-medo/talebound-backend/db/migration/migrator"
	_ "github.com/the-medo/talebound-backend/doc/statik"
	"github.com/the-medo/talebound-backend/util"
	"os"
	"strings"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config:")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	runMigrator(config.MigrationURL, config.MigrationObjectsURL, config.DBSource)

}

func runMigrator(migrationURL string, migrationObjectsURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new migrate instance: ")
		return
	}

	path := strings.TrimPrefix(migrationObjectsURL, "file://")
	log.Info().Msgf("Path: %s", path)

	mg, err := migrator.New(&migrator.Config{
		DbObjectPath: path,
	}, migration)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new mg instance! ")
		return
	}

	fmt.Println(mg.GetObjectList())

}
