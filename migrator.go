package main

import (
	"database/sql"
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

	runMigrator(config.MigrationURL, config.MigrationObjectsURL, config.DBSource)

}

func runMigrator(migrationURL string, migrationObjectsURL string, dbSource string) {
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		//log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	migration, err := migrate.NewWithDatabaseInstance(
		migrationURL,
		"talebound", driver)
	if err != nil {
		//log.Fatal(err)
	}

	//migration, err := migrate.New(migrationURL, dbSource)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Cannot create new migrate instance: ")
	//	return
	//}

	path := strings.TrimPrefix(migrationObjectsURL, "file://")
	log.Info().Msgf("Path: %s", path)

	mg, err := migrator.New(&migrator.Config{
		DB:           db,
		DbObjectPath: path,
	}, migration)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new mg instance! ")
		return
	}

	fmt.Println(mg.RunFile("db/migration/objects/migrations_drop_objects.sql"))
	fmt.Println(mg.GetObjectList())
	fmt.Println("==============")
	objectVersions, err := mg.GetObjectsForStep(25, migrator.DirectionDown)
	for _, ov := range objectVersions {
		path := mg.GetDbObjectVersionPath(ov)
		fmt.Println(ov.DbObject.Name, " => ", ov.Version, "; ", path)
		fmt.Println(mg.RunFile(path))
		fmt.Println("==============")
	}

}
