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
	UAUserAgentUnknown = "Unknown"
	UAIPUnknown        = "0.0.0.0"
	UAIPDelimiter      = "---"
)

//token lifetime
const (
	LifetimeAccessToken  = 15 * time.Minute
	LifetimeRefreshToken = 24 * time.Hour
)

//errors
const (
	ErrInvalidJson    = "Invalid json provided"
	ErrGeneralError   = "Something went wrong at out side"
	ErrAccessDenied   = "AccessDenied"
	ErrNotImplemented = "Not implemented yet"
)

type MessageJson struct {
	Message string `json:"message"`
}

type ErrJson struct {
	Err string `json:"err"`
}
