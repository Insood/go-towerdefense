.PHONY: harness game test

harness:
	go run ./cmd/harness $(filter-out $@,$(MAKECMDGOALS))

game:
	go run ./cmd/game $(filter-out $@,$(MAKECMDGOALS))

test:
	$${GOCACHE:=.gocache} go test ./...

%:
	@:
