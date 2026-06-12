.PHONY: game test

game:
	go run ./cmd/game $(filter-out $@,$(MAKECMDGOALS))

test:
	$${GOCACHE:=.gocache} go test ./...

%:
	@:
