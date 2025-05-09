param (
    [Parameter(Position=0)]
    [string]$Command
)

switch ($Command) {
    "up" {
        docker-compose up -d
    }
    "down" {
        docker-compose down
    }
    "createdb" {
        docker-compose exec postgres createdb --username=root --owner=root go-chat
    }
    "dropdb" {
        docker-compose exec postgres dropdb go-chat
    }
    "postgres" {
        docker-compose exec postgres psql -U root
    }
    "migrateup" {
        migrate -path server/db/migrations -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose up
    }
    "migratedown" {
        migrate -path server/db/migrations -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose down
    }
    default {
        Write-Host "Usage: .\db.ps1 [up|down|createdb|dropdb|postgres|migrateup|migratedown]"
    }
}
