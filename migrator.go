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

	sumfile := flag.Bool("sumfile", false, "merge newest versions into single file ")
	migrateUp := flag.Bool("up", false, "migrate files up ")
	migrateDown := flag.Bool("down", false, "migrate files down")
	step := flag.Int("step", 0, "number of migration steps... if 0 or not provided, runs all migrations ")
	flag.Parse()

	// ============== Create sumfile
	if *sumfile {
		err = mg.CreateObjectsFile()
		if err != nil {
			log.Fatal().Err(err).Msg("Creating sumfile failed! ")
		}
		return
	}

	// =============== Run migrations
	//fmt.Println("Parsed:", *sumfile, *migrateUp, *migrateDown, *step)
	if *migrateUp && *migrateDown {
		log.Fatal().Err(nil).Msg("Can not do migrate UP and DOWN at the same time! ")
	} else if *migrateDown && step == nil {
		log.Fatal().Err(nil).Msg("Step argument must be provided when doing down migration! ")
	} else if *step < 0 {
		log.Fatal().Err(nil).Msg("Step cannot be negative!")
	}

	// === Down migration check - running all down migrations is often not wanted
	if *migrateDown {
		*step *= -1
		if *step == 0 {
			fmt.Println("You set the 'step' parameter to 0, this resets all migrations. Do you want to continue? (y/n)")

			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				fmt.Println("Please enter 'y' or 'n'")
			}

			if strings.ToLower(response) != "y" {
				fmt.Println("Operation cancelled")
				return
			}
		}
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

	if finalVersion > highestVersion || *step == 0 && *migrateUp {
		finalVersion = highestVersion
	} else if finalVersion < 0 {
		finalVersion = 0
	}

	drop := true
	for currentVersion != finalVersion {
		//fmt.Println("Current version: ", currentVersion, "Final version: ", finalVersion)
		if *migrateUp {
			currentVersion, err = mg.RunStep(migrator.DirectionUp, drop, currentVersion+1 == finalVersion)

		} else if *migrateDown {
			currentVersion, err = mg.RunStep(migrator.DirectionDown, drop, currentVersion-1 == finalVersion)
		}

		drop = false

		if err != nil {
			if err != nil {
				log.Fatal().Err(err).Msg("step failed")
			}
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
