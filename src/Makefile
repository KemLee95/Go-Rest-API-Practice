LDFLAGS := -ldflags "-X main.Version=${VERSION}"
PID_FILE := './.pid'

.PHONY: run
run: ## run the API server
	go run ${LDFLAGS} ./main.go

.PHONY: run-dev
run-dev: ## run the API server with live reload support (requires fswatch)
	@go run ${LDFLAGS} ./main.go & echo $$! > $(PID_FILE)
	@fswatch -x -o --event Created --event Updated --event Renamed -r internal pkg cmd config | xargs -n1 -I {} make run-restart