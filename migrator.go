package main

import (
	"database/sql"
	"errors"
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

	mg := getMigratorInstance(config.MigrationURL, config.MigrationObjectsURL, config.DBSource, config.MigrationCreateObjectsFilename, config.MigrationDropObjectsFilename)

	createObjectFile := flag.Bool("sumfile", false, "merge newest versions into single file ")
	migrateUp := flag.Bool("up", false, "migrate files up ")
	migrateDown := flag.Bool("down", false, "migrate files down - \"step\" parameter must be provided.")
	step := flag.Int("step", 1, "how much steps should be migrated? If not provided, all ")
	flag.Parse()

	if *createObjectFile {
		err = mg.CreateObjectsFile()
		if err != nil {
			log.Fatal().Err(err).Msg("Creating sumfile failed! ")
		}
	}

	if *migrateUp && *migrateDown {
		log.Fatal().Err(err).Msg("Can not do migrate UP and DOWN at the same time! ")
	} else if *migrateDown && step == nil {
		log.Fatal().Err(err).Msg("Step argument must be provided when doing down migration! ")
	}

	cv, dirty, _ := mg.Migrate.Version()
	currentVersion := int(cv)
	if dirty {
		log.Fatal().Err(errors.New("starting from dirty database")).Msg("Fix errors manually")
	}
	finalVersion := currentVersion + *step
	highestVersion, err := mg.GetHighestAvailableVersion()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to get highest available version")
	}

	if finalVersion > highestVersion {
		finalVersion = highestVersion
	} else if finalVersion < 0 {
		finalVersion = 0
	}

	for currentVersion != finalVersion {
		fmt.Println("Current version: ", currentVersion)
		if *migrateUp {
			currentVersion, err = mg.RunStep(currentVersion, migrator.DirectionUp)
		} else if *migrateDown {
			currentVersion, err = mg.RunStep(currentVersion, migrator.DirectionDown)
		}
		if err != nil {
			break
		}
	}
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
