package main

import (
	"context"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/simar7/trivy-wasm-bundle/pkg/wasm/sdk/opa"
)

func prepareModules(paths []string, data string) map[string]*ast.Module {
	modules := make(map[string]*ast.Module)
	for _, path := range paths {
		module, err := ast.ParseModuleWithOpts(path, data, ast.ParserOptions{
			ProcessAnnotation: true,
		})
		if err != nil {
			panic(err)
		}
		modules[path] = module
	}
	return modules
}

func LoadRego(inputRegoFile string, inputVal bool) string {
	return LoadRegoWithPrecompile(inputRegoFile, inputVal, false)
}

func LoadRegoWithPrecompile(inputRegoFile string, inputVal bool, precompile bool) string {
	policy, err := os.ReadFile(inputRegoFile)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	var regoOpts []func(*rego.Rego)
	if precompile {
		compiler := ast.NewCompiler().
			WithUseTypeCheckAnnotations(true).
			WithCapabilities(ast.CapabilitiesForThisVersion())

		compiler.Compile(prepareModules([]string{inputRegoFile}, string(policy)))
		regoOpts = append(regoOpts, rego.Compiler(compiler))
	}

	r := rego.New(
		rego.Query("data.example.allow"),
		rego.Module("policy.rego", string(policy)),
		rego.Input(map[string]interface{}{"foo": inputVal}))
	for _, opts := range regoOpts {
		opts(r)
	}

	resultSet, err := r.Eval(ctx)
	if err != nil {
		panic(err)
	}

	if len(resultSet) > 0 {
		return fmt.Sprintf("Policy result from Rego direct: %v\n", resultSet.Allowed())
	} else {
		panic("No results found")
	}

}

func LoadWASM(inputWASMFile string, inputVal bool) string {
	policy, err := os.ReadFile(inputWASMFile)
	if err != nil {
		panic(err)
	}

	rwasm, err := opa.New().WithPolicyBytes(policy).Init()
	if err != nil {
		panic(err)
	}

	defer rwasm.Close()

	var input interface{} = map[string]interface{}{
		"foo": inputVal,
	}

	ctx := context.Background()
	eps, err := rwasm.Entrypoints(ctx)
	if err != nil {
		panic(err)
	}

	entrypointID, ok := eps["example/allow"]
	if !ok {
		panic(err)
	}

	result, err := rwasm.Eval(ctx, opa.EvalOpts{Entrypoint: entrypointID, Input: &input})
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Policy result from WASM: %v\n", string(result.Result))
}

func main() {
	fmt.Println(LoadWASM("example-check-rego.wasm", true))
	//LoadWASM("example-check-go.wasm")
	fmt.Println(LoadRego("example-check.rego", true))
	fmt.Println(LoadRegoWithPrecompile("example-check.rego", true, true))
}
