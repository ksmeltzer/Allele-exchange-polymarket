//go:build wasm

package main

import (
	"fmt"
	"unsafe"
)

//go:wasmexport Manifest
func Manifest() uint64 {
	manifestJSON := `{"name": "allele-exchange-polymarket", "version": "v1.0.0", "description": "Polymarket Exchange Adapter", "author": "Allele Org", "dependencies": [], "config": [{"key": "WALLET_ADDRESS", "type": "string", "description": "Public Polygon Wallet Address (0x...)", "required": true}, {"key": "WALLET_PRIVATE_KEY", "type": "secret", "description": "Polygon Wallet Private Key (0x...)", "required": true}]}`
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
