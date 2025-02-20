package main

import "testing"

func BenchmarkLoadWASM(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		LoadWASM("example-check-rego.wasm", true)
	}
}
