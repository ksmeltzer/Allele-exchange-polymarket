//go:build wasm
package main

import "fmt"

func main() {
	// Adapter listens on STDIN/STDOUT via WASI for engine instructions,
	// and dials the WS network natively.
	fmt.Println("allele-exchange-polymarket loaded")
}
