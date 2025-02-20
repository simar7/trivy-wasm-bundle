package main

//
//import (
//	"encoding/json"
//	"syscall/js"
//)
//
//type Input struct {
//	Foo bool `json:"foo"`
//}
//
//// allow is the WASM entry point for "example/allow"
//func allow(this js.Value, args []js.Value) interface{} {
//	if len(args) == 0 {
//		return js.ValueOf("error: missing input")
//	}
//
//	var input Input
//	err := json.Unmarshal([]byte(args[0].String()), &input)
//	if err != nil {
//		return js.ValueOf("error: invalid JSON input")
//	}
//
//	if input.Foo {
//		return js.ValueOf("allowed")
//	}
//	return js.ValueOf("denied")
//}
//
//func main() {
//	js.Global().Set("example/allow", js.FuncOf(allow))
//	select {} // Keeps the WASM module running
//}
//
