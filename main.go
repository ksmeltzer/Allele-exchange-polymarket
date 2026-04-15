//go:build wasm
package main

import (
	"fmt"
	"unsafe"
)

//go:wasmexport Manifest
func Manifest() uint64 {
	manifestJSON := `{"name": "allele-exchange-polymarket", "version": "v1.0.0", "description": "Polymarket Exchange Adapter", "author": "Allele Org", "dependencies": [], "config": [{"key": "POLY_API_KEY", "type": "secret", "description": "Polymarket API Key", "required": true}, {"key": "POLY_API_SECRET", "type": "secret", "description": "Polymarket API Secret", "required": true}, {"key": "POLY_API_PASSPHRASE", "type": "secret", "description": "Polymarket API Passphrase", "required": true}]}`
	outBytes := []byte(manifestJSON)
	outPtr := uint32(uintptr(unsafe.Pointer(&outBytes[0])))
	outLen := uint32(len(outBytes))
	return (uint64(outPtr) << 32) | uint64(outLen)
}

func main() {
	// Adapter listens on STDIN/STDOUT via WASI for engine instructions,
	// and dials the WS network natively.
	fmt.Println("allele-exchange-polymarket loaded")
	
}
