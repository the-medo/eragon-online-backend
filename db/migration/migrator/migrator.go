package migrator

//
//import (
//	"fmt"
//	"github.com/golang-migrate/migrate/v4"
//	database "github.com/golang-migrate/migrate/v4/database"
//	_ "github.com/golang-migrate/migrate/v4/database/postgres"
//	_ "github.com/golang-migrate/migrate/v4/source/file"
//)
//
//func runBeforeMigrationScripts(db database.Driver, version uint) {
//	// run your SQL or other things here
//	fmt.Printf("Before migration - Version: %d\n", version)
//}
//
//func runAfterMigrationScripts(db database.Driver, version uint) {
//	// run your SQL or other things here
//	fmt.Printf("After migration - Version: %d\n", version)
//}
//
//func runDBMigration(migrationURL string, dbSource string) {
//	m, err := migrate.New(migrationURL, dbSource)
//	if err != nil {
//		log.Fatal().Err(err).Msg("Cannot create new migrate instance.")
//		return
//	}
//
//	steps := 1 // Number of steps you'd like to migrate at a time
//	currentVersion, dirty, _ := m.Version()
//
//	for {
//		if dirty {
//			log.Fatal().Err(err).Msg("Migration failed, database in dirty state.")
//		}
//
//		// Execute script before migration
//		runBeforeMigrationScripts(m.DB, currentVersion)
//
//		// Migrate steps forward
//		err := m.Steps(steps)
//
//		if err == migrate.ErrNoChange {
//			log.Info().Msg("No migrations left to apply.")
//			break
//		} else if err != nil {
//			log.Fatal().Err(err).Msg("Migration failed.")
//			return
//		}
//
//		// Execute script after migration
//		currentVersion, dirty, _ = m.Version()
//		runAfterMigrationScripts(m.DB, currentVersion)
//	}
//}
