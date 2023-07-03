set dotenv-load

run: 
    go run main.go

compose:
    docker compose pull
    docker compose up -d
