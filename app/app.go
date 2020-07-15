package app

//Envs
const (
	EnvJwtSecret    = "JWT_SECRET"
	EnvEncryptSalt  = "ENCRYPT_SALT"
	EnvDatabaseFile = "DATABASE_FILE"
)

//Defaults
const (
	DefaultDBFile         = "go-rest.db"
	DefaultSuperUserScope = "SUPER_ADMIN"
	DefaultSalt           = "Please_Use_ENCRYPT_SALT_Env_Instead"
)

//Database-related
const (
	DBDialect = "sqlite3"
)
