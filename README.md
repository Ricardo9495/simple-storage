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


## Requirements

**Note**: As there are no explicit requirements, all of the following requirements are based on my assumptions.


Features:
- Upload:
  + Upload file with unique name.
  + Size limit 50MB

- Delete:
  + Delete file by name.

- List flie:
  + List all files.

## Architecture
  ![image](https://github.com/user-attachments/assets/f10daf0f-89df-4e0c-9d87-4b0b3d668fe7)


### Component
**1. Cli**
 - a cli developed by cobra cli framework to interact with storage server.

**2. storage server**
 - a go-based application handle business logic, recive request and

**3. storage**
 - For simplicity, I will choose local storage.
 - Architecture design and code design are easy to extend/adapt new kind of storage, for e.g aws S3, azure storage,etc..

**4. metadata DB**
 - Store metadata of file, such as name, filepath, etc..
 - Relational DB woul be a good choice in this case for f

This is a simple hello, world demonstration web server.

It serves version information on /version and answers any other request like /name by saying "Hello, name!".


migrate -path migrations -database '$(PG_URL)?sslmode=disable' up


migrate -source file://./internal/migrations -database "postgres://user:pass@localhost:5432/postgres?sslmode=disable" up


curl -X POST -F "file=@./internal/interfaces/sample.txt" http://localhost:8080/files

