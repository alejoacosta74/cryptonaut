
# Run unit tests
test:
	go test -v ./...

# Run benchmarks
bench:
	go test -bench=. ./...
