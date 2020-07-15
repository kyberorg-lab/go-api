package app

//Envs
const (
	EnvJwtSecret                = "JWT_SECRET"
	EnvEncryptSalt              = "ENCRYPT_SALT"
	EnvEncryptSecretKeyPassword = "ENCRYPT_SECRET_KEY_PASSWORD"
	EnvDatabaseFile             = "DATABASE_FILE"
	EnvSuperUserName            = "SUPER_USER_NAME"
	EnvSuperUserPassword        = "SUPER_USER_PASSWORD"
)

//Defaults
const (
	DefaultDBFile            = "go-rest.db"
	DefaultSuperUserScope    = "SUPER_ADMIN"
	DefaultSuperUserName     = "toor"
	DefaultSuperUserPassword = "OpenItIsToor"
	DefaultSalt              = "Please_Use_ENCRYPT_SALT_Env_Instead"
	DefaultSecretKeyPassword = "Please_Use_ENCRYPT_SECRET_KEY_PASSWORD_Env_Instead"
)

//Database-related
const (
	DBDialect = "sqlite3"
)
