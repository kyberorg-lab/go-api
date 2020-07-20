# GO API
REST with JWT written in go

## Endpoints

### GET /health
* 200 - when app is alive

### POST /auth/login
Scope/Auth: JSON with a username and password

* 200 - Tokens issued + Tokens (access, refresh)
* 422 - when send malformed JSON + ErrorJson
* 401 - login details wrong + ErrorJson
* 500 - server-side error + ErrorJson

### POST /auth/refresh-token
Scope/Auth: Refresh token

* 201 - New Tokens issued + Tokens (access, refresh)
* 401 - when no token or expired token + ErrorJson
* 403 - when refresh token's owner has no permissions to update tokens
* 422 - when send malformed JSON + ErrorJson
* 500 - when server-side issue + ErrorJson

### POST /auth/logout
Scope/Auth: Access token

* 200 - success + MessageJson
* 401 - when no token or expired token + ErrorJson



## Env Vars

Parameter | Description | Value |
--- | --- | --- |
DATABASE_FILE | Location of database file | go-rest.db |
FIRST_USER_NAME | Username for first application user | toor | 
FIRST_USER_PASSWORD | Password for first application user | OpenItIsToor |

## About Users
In order to protect sensible resources, we use Users.

User passwords stored encrypted. For encryption, we use `AES-GCM`

Users can be configured using `/users` API. This API is available only for SUPER_ADMIN users.

When application starts and detects that there are 0 users, it will create First User.

### First User aka root 
First User's username is by default `toor` (reversed root), but it can be customizing by EnvVar `FIRST_USER_USERNAME`.

First User's password can be set by EnvVar `FIRST_USER_PASSWORD`. 
If not set: application uses built-in password `OpenItIsToor`.

## About Encryption
Passwords stored in Database hashed.
For hashing, we use BCrypt with salt.

