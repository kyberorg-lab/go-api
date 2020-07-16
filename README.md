# GO API
REST with JWT written in go

## Endpoints

### POST /auth 
Scope: no auth needed

* 422 - when send malformed JSON + ErrorJson
* 401 - login details wrong + ErrorJson
* 500 - server-side error + ErrorJson
* 200 - Token issued + Token

## Env Vars

Parameter | Description | Value |
--- | --- | --- |
DATABASE_FILE | Location of database file | go-rest.db |
FIRST_USER_NAME | Username for first application user | toor | 
FIRST_USER_PASSWORD | Password for first application user | OpenItIsToor |
ENCRYPT_SECRET_KEY_PASSWORD | Shared secret for password encryption | Please_Use_ENCRYPT_SECRET_KEY_PASSWORD_Env_Instead |
ENCRYPT_SALT | Salt for password encryption | Please_Use_ENCRYPT_SALT_Env_Instead

## About Users
In order to protect sensible resources, we use Users.

User passwords are stored encrypted. For encryption we use `AES-GCM`

Users can be configured using `/users` API. This API is available only for SUPER_ADMIN users.

When application starts and detects that there are 0 users, it will create First User.

### First User aka root 
First User's username is by default `toor` (reversed root), but it can be customizing by EnvVar `FIRST_USER_USERNAME`.

First User's password can be set by EnvVar `FIRST_USER_PASSWORD`. 
If not set: application uses built-in password `OpenItIsToor`

## About Encryption
Passwords are stored in Database encrypted. For encryption and decryption we use encryptor.
This encryptor uses password-based bytes encryptor using 256 bits AES encryption with Galois Counter Mode (GCM).

For initialise encryptor we need both:

* Secret key password 
* Salt

Those params can be defined using:
* Environment variable

Environment variable approach is recommended.

If ENV not set, application will use default values.

#### Secret key password 
* EnvVar `ENCRYPT_SECRET_KEY_PASSWORD`
* Default: `Please_Use_ENCRYPT_SECRET_KEY_PASSWORD_Env_Instead`

#### Salt
* EnvVar `ENCRYPT_SALT` 
* Property `Please_Use_ENCRYPT_SALT_Env_Instead`

