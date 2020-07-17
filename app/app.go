package app

//Envs
const (
	EnvJwtSecret         = "JWT_SECRET"
	EnvDatabaseFile      = "DATABASE_FILE"
	EnvFirstUserName     = "FIRST_USER_NAME"
	EnvFirstUserPassword = "FIRST_USER_PASSWORD"
)

//Defaults
const (
	DefaultDBFile            = "go-rest.db"
	DefaultFirstUserScope    = ScopeSuperAdmin
	DefaultFirstUserName     = "toor"
	DefaultFirstUserPassword = "OpenItIsToor"
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

type Token struct {
	Token string `json:"token"  binding:"required"`
}
