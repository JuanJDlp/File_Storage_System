## Endpoints

### File Related

**POST** "/api/v1/files" -> Create a new file

**DELETE** "/api/v1/files" -> Deletes all files

**DELETE** "/api/v1/files/:name" -> Gets a file by name
**GET** "/api/v1/files/:name" -> Start the dowload of an specific file

**GET** "/api/v1/file/all" -> Will list all files owned by a user but not start the dowload.

### User Related

**POST** "/api/v1/users/login" -> Logs in and return the jwt token

**POST** "/api/v1/users/register" -> Creates a new User

**POST** "/api/v1/users/update" -> Updates the username or password of an alredy existing user
