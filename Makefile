run:
	(cd cmd/genuuidv7; go run genuuidv7.go)
build:
	(cd cmd/genuuidv7; go build)
bench:
	(cd pkg/uuidv7; go test -bench=.)

