new:
	migrate create -dir ./migrations -ext .sql -seq -digits 5 init
clean_mg:
	rm -rf ./migrations
clean_db:
	rm -rf ./db
