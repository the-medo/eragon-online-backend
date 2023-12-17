package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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

	createObjectFile := flag.Bool("sumfile", false, "Merge newest versions into single file ")
	flag.Parse()

	mg := getMigratorInstance(config.MigrationURL, config.MigrationObjectsURL, config.DBSource, config.MigrationCreateObjectsFilename, config.MigrationDropObjectsFilename)

	fmt.Println(mg)

	if *createObjectFile {
		err = mg.CreateObjectsFile()
		if err != nil {
			log.Fatal().Err(err).Msg("Creating sumfile failed! ")
		}
	}

	fmt.Println(*createObjectFile)

}

func getMigratorInstance(migrationURL string, migrationObjectsURL string, dbSource string, createObjectsFilename string, dropObjectsFilename string) *migrator.Migrator {
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect! ")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	migration, err := migrate.NewWithDatabaseInstance(
		migrationURL,
		"talebound", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new migrate instance! ")
	}

	path := strings.TrimPrefix(migrationObjectsURL, "file://")
	log.Info().Msgf("Path: %s", path)

	mg, err := migrator.New(&migrator.Config{
		DB:                    db,
		DbObjectPath:          path,
		MigrationFilesPath:    strings.TrimPrefix(migrationURL, "file://"),
		CreateObjectsFilename: createObjectsFilename,
		DropObjectsFilename:   dropObjectsFilename,
	}, migration)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new mg instance! ")
		return nil
	}

	return mg
}

func runMigrator() {

	//fmt.Println(mg.GetHighestAvailableVersion())
	//
	//err = mg.CreateObjectsFile()
	//if err != nil {
	//	fmt.Println("OOOPS")
	//}

}
