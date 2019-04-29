generate: bindec_bin generate-bench generate-test

bindec_bin:
	go build -o bindec_bin ./cmd/bindec/main.go

generate-bench: bindec_bin
	rm -f bench/foo_bindec.go
	./bindec_bin -type=Foo bench

generate-test: bindec_bin
	rm -f bindec_test.go
	go generate .

test: generate-test
	go test -cover -coverprofile=coverage.txt -covermode="atomic" . -v

bench: generate-bench
	go test -bench=./bench

clean:
	rm -f bindec_bin

.PHONY: test
