package migrator

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrDifferentPriorityLength = errors.New("priorities need to have the same length")
	ErrDifferentVersionLength  = errors.New("versions need to have the same length")
	ErrInvalidPriority         = errors.New("priority should be a number")
)

type Direction string

const (
	DirectionUp   Direction = "up"
	DirectionDown Direction = "down"
)

type Migrator struct {
	migrate *migrate.Migrate
	config  *Config
}

func (m *Migrator) GetObjectList() ([]*DbObject, error) {
	entries, err := os.ReadDir(m.config.DbObjectPath)
	if err != nil {
		return nil, err
	}

	orderDirEntries(entries)

	dbObjects := make([]*DbObject, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			dbObject, err := parseDir(entry, m.config)
			if err != nil {
				return nil, err
			}

			dbObjects = append(dbObjects, dbObject)
		}
	}

	return dbObjects, nil
}

func (m *Migrator) GetObjectsForStep(step int, upOrDown Direction) (result []*DbObjectVersion, err error) {
	entries, err := m.GetObjectList()
	if err != nil {
		return nil, err
	}

	//result := make([]*DbObjectVersion, 0, len(*entries))

	for _, dbObject := range entries {
		fmt.Print(dbObject.Name, ":")
		correctVersion := 0
		for _, version := range dbObject.Versions {
			if version > correctVersion && ((upOrDown == DirectionUp && step >= version) || (upOrDown == DirectionDown && step > version)) {
				correctVersion = version
			} else {
				break
			}
		}
		if correctVersion > 0 {
			result = append(result, &DbObjectVersion{
				DbObject: dbObject,
				Version:  correctVersion,
			})

		}
		fmt.Println(correctVersion)
	}

	return
}

func (m *Migrator) GetDbObjectVersionPath(dov *DbObjectVersion) string {
	folder := m.config.DbObjectPath
	objectFolder := LPAD(dov.DbObject.Priority, m.config.PriorityLpad) + "_" + dov.DbObject.Name
	filename := LPAD(dov.Version, m.config.VersionLpad) + "_" + dov.DbObject.Name + ".sql"
	return folder + "/" + objectFolder + "/" + filename
}

func (m *Migrator) RunFile(path string) error {
	sqlFileContents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = m.config.DB.Exec(string(sqlFileContents))
	if err != nil {
		return err
	}

	fmt.Println("Success!")

	return nil
}

func parseDir(dir os.DirEntry, config *Config) (*DbObject, error) {
	splitEntry := strings.Split(dir.Name(), "_")
	if config.PriorityLpad == 0 {
		config.PriorityLpad = len(splitEntry[0])
	} else if config.PriorityLpad != len(splitEntry[0]) {
		return nil, ErrDifferentPriorityLength
	}

	priority, err := strconv.Atoi(splitEntry[0])
	if err != nil {
		return nil, ErrInvalidPriority
	}

	versions, err := parseVersionFiles(dir.Name(), config)

	return &DbObject{
		Name:     strings.Join(splitEntry[1:], "_"),
		Priority: priority,
		Versions: versions,
	}, nil
}

func parseVersionFiles(dirName string, config *Config) ([]int, error) {
	versionFilePath := filepath.Join(config.DbObjectPath, dirName)
	versionFiles, err := os.ReadDir(versionFilePath)
	if err != nil {
		return nil, err
	}

	orderDirEntries(versionFiles)

	versions := make([]int, 0, len(versionFiles))
	for _, file := range versionFiles {
		versionString := strings.Split(file.Name(), "_")[0]
		if config.VersionLpad == 0 {
			config.VersionLpad = len(versionString)
		} else if config.VersionLpad != len(versionString) {
			return nil, ErrDifferentVersionLength
		}

		version, err := strconv.Atoi(versionString)
		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	return versions, nil
}

func orderDirEntries(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
}

func New(migratorConfig *Config, migrateInstance *migrate.Migrate) (*Migrator, error) {
	return &Migrator{
		migrate: migrateInstance,
		config:  migratorConfig,
	}, nil
}

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
