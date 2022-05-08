package voo

import "database/sql"

type initPaths struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}

type databaseConfig struct {
	dsn     string
	databse string
}

type Database struct {
	DatabseType string
	Pool        *sql.DB
}
