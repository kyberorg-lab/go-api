package app

//Envs
const (
	EnvJwtSecret                = "JWT_SECRET"
	EnvEncryptSalt              = "ENCRYPT_SALT"
	EnvEncryptSecretKeyPassword = "ENCRYPT_SECRET_KEY_PASSWORD"
	EnvDatabaseFile             = "DATABASE_FILE"
	EnvFirstUserName            = "FIRST_USER_NAME"
	EnvFirstUserPassword        = "FIRST_USER_PASSWORD"
)

//Defaults
const (
	DefaultDBFile            = "go-rest.db"
	DefaultFirstUserScope    = ScopeSuperAdmin
	DefaultFirstUserName     = "toor"
	DefaultFirstUserPassword = "OpenItIsToor"
	DefaultSalt              = "Please_Use_ENCRYPT_SALT_Env_Instead"
	DefaultSecretKeyPassword = "Please_Use_ENCRYPT_SECRET_KEY_PASSWORD_Env_Instead"
)

//Database-related
const (
	DBDialect = "sqlite3"
)

//scopes
const (
	ScopeSuperAdmin = "SUPER_ADMIN"
	ScopeUser       = "USER"
)
