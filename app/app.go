package app

import "time"

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

//user-agent related
const (
	UserAgentUnknown = "Unknown"
	IPUnknown        = "0.0.0.0"
	IPUADelimiter    = "---"
)

//token timeouts
const (
	TimeoutAccessToken  = 15 * time.Minute
	TimeoutRefreshToken = 24 * time.Hour
)

//errors
const (
	GeneralError = "Something went wrong at out side"
	AccessDenied = "AccessDenied"
)
