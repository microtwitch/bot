set dotenv-load

run: 
    go run cmd/bot/main.go

compose:
    docker compose pull
    docker compose up -d
