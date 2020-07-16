# GO API
REST with JWT written in go

# Endpoints

### POST /auth 
Scope: no auth needed

* 422 - when send malformed JSON + ErrorJson
* 401 - login details wrong + ErrorJson
* 500 - server-side error + ErrorJson
* 200 - Token issued + Token

