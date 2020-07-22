package global

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
	ErrInvalidJson      = "invalid json provided"
	ErrGeneralError     = "something went wrong at out side"
	ErrAccessDenied     = "accessDenied"
	ErrNotImplemented   = "not implemented yet"
	ErrEmptyTokenString = "got empty string instead of token"
	ErrEmptyToken       = "got empty token"
	ErrMalformedToken   = "malformed token"
	ErrTokenExpired     = "token is either expired or not active yet"
	ErrMalformedClaims  = "got malformed claims"
)
