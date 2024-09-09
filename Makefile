include .env
export

migrate-local-up:
	@migrate -path migrations \
	 -database "postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable" \
	 -verbose up

migrate-local-down:
	 @migrate -path migrations \
         -database "postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable" \
         -verbose down

migrate-create:
	migrate create -dir migrations -ext sql $(name)