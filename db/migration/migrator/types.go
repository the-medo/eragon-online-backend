package migrator

type Config struct {
	DbObjectPath string
	PriorityLpad int
	VersionLpad  int
}

type DbObject struct {
	Name     string
	Priority int
	Versions []int
}
