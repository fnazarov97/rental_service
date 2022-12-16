
run: 
	go run main.go
migrate-up:
	migrate -path ./storage/migrations -database 'postgres://farhod:f@rhod666997@localhost:5432/rentals?sslmode=disable' up
migrate-down:
	migrate -path ./storage/migrations -database 'postgres://farhod:f@rhod666997@localhost:5432/rentals?sslmode=disable' down
pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule foreach --recursive git pull
	# git submodule update --remote --merge
pull-proto:
	git submodule update --init --recursive