This is a simple hello, world demonstration web server.

It serves version information on /version and answers any other request like /name by saying "Hello, name!".


migrate -path migrations -database '$(PG_URL)?sslmode=disable' up


migrate -source file://./internal/migrations -database "postgres://user:pass@localhost:5432/postgres?sslmode=disable" up


curl -X POST -F "file=@./internal/interfaces/sample.txt" http://localhost:8080/files

