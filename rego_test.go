package main

import "testing"

func BenchmarkLoadRego(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		LoadRego("example-check.rego", true)
	}
}
func BenchmarkLoadRegoWithPrecompile(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		LoadRegoWithPrecompile("example-check.rego", true, true)
	}
}
