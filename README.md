# Useful commands

Run application
`go run cmd/app/main.go`

Compiles the Go code in the cmd directory into an executable binary:
`go build ./cmd`

Add Missing Dependencies and Remove Unused Dependencies:
`go mod tidy`

Starts all services defined in the docker-compose.yml file
`docker compose up -d`

Verify Services Are Running
`docker ps`

Install go migrate:
`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

Execute migrations:
`migrate -database "postgresql://postgres:postgres@localhost:5432/gateway?sslmode=disable" -path migrations up`