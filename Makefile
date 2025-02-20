bundle:
	opa build -t wasm -e example/allow ./example-check.rego

clean:
	rm -rf bundle.tar.gz
	rm -rf *.wasm

extract:
	tar -xzvf bundle.tar.gz /policy.wasm
	mv policy.wasm example-check-rego.wasm

run:
	go run main.go .

bench:
	go test -v -bench=. -benchtime=100x -count=10 .

report:
	go test -v -bench=. -benchtime=100x -count=10 . > benchmark.txt
	go run pkg/reporter/benchmark-output-generator.go benchmark.txt 
