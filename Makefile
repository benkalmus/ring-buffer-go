VERSION=0.1.0
PACKAGE=trie

.PHONY: test
test: 
	go test -timeout 10s -failfast -v -race ./...
