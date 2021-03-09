package drivers

// Driver are the names of supported driver names
type Driver string

const (
	// MySQL database driver
	MySQL = "mysql"
	// Postgres database driver
	Postgres = "postgres"
	// SQLServer database driver
	SQLServer = "sqlserver"
	// SQLite database driver
	SQLite = "sqlite"
	// SQLite3 database driver
	SQLite3 = "sqlite3"
)
