dependencies:
	docker compose up

run-e2e:
	cd ./script-e2e/ && go run ./main.go

run-load:
	cd ./script-load-test && go run ./main.go
