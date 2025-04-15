# Simple Storage

A simple application with cli as client and Go-based server.
Allows file uploads, storing file with a unique name to ensure no conflicts.
### features
- upload file
- delete file
- list file

## Prerequisites
docker
 
## Usage
Start
```
Docker compose up
```
Use
```
// upload file
docker compose run upload-file --name filename --file filepath

// delete file
docker compose run delete-file --name

// list file
docker compose run list-file
```

# Document

## Architecture

### Requirements

**Note**: As there are no explicit requirements, all of the following requirements are based on my assumptions.


Features:
- Upload: 
upload file with unique name, size limit 50MB

- Delete:
  delete file by name.

- List flie:
  List all files.

  

This is a simple hello, world demonstration web server.

It serves version information on /version and answers any other request like /name by saying "Hello, name!".


migrate -path migrations -database '$(PG_URL)?sslmode=disable' up


migrate -source file://./internal/migrations -database "postgres://user:pass@localhost:5432/postgres?sslmode=disable" up


curl -X POST -F "file=@./internal/interfaces/sample.txt" http://localhost:8080/files

