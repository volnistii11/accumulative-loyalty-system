new:
	migrate create -dir ./internal/gophermart/storage/migrations -ext .sql -seq -digits 5 init
clean_mg:
	rm -rf ./internal/gophermart/storage/migrations
clean_db:
	rm -rf ./db
