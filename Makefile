.PHONY: db
db:
	docker run -v bench:/var/lib/postgresql/data -p 54321:5432 --rm --name bench-db -e POSTGRES_USER=bench -e POSTGRES_PASSWORD=bench -e POSTGRES_DB=bench postgres:15
