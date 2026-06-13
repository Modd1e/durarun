run:
	go run cmd/durarun/main.go

local:
	docker compose -f docker-compose.locl.yaml up --build -d

atlas-apply:
	set -a; . ./.env; set +a; atlas migrate apply --env local
