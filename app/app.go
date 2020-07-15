package app

//Envs
const (
	EnvJwtSecret    = "JWT_SECRET"
	EnvDatabaseFile = "DATABASE_FILE"
)

//Defaults
const (
	DefaultDBFile         = "go-rest.db"
	DefaultSuperUserScope = "SUPER_ADMIN"
)

//Database-related
const (
	DBDialect = "sqlite3"
)
