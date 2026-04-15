package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run keygen.go <PUBLIC_ADDRESS> <POLYGON_PRIVATE_KEY>")
		fmt.Println("Example: go run keygen.go 0xYourAddress 0xYourPrivateKey")
		os.Exit(1)
	}

	address := os.Args[1]
	privKeyHex := os.Args[2]

	privKeyHex = strings.TrimPrefix(privKeyHex, "0x")
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		fmt.Printf("Invalid private key: %v\n", err)
		os.Exit(1)
	}

	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := "0"

	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
			},
			"ClobAuth": []apitypes.Type{
				{Name: "address", Type: "address"},
				{Name: "timestamp", Type: "string"},
				{Name: "nonce", Type: "uint256"},
				{Name: "message", Type: "string"},
			},
		},
		PrimaryType: "ClobAuth",
		Domain: apitypes.TypedDataDomain{
			Name:    "ClobAuthDomain",
			Version: "1",
			ChainId: math.NewHexOrDecimal256(137), // Polygon Mainnet
		},
		Message: apitypes.TypedDataMessage{
			"address":   address,
			"timestamp": timestampStr,
			"nonce":     nonce,
			"message":   "This message attests that I control the given wallet",
		},
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		fmt.Printf("Failed to hash domain: %v\n", err)
		os.Exit(1)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		fmt.Printf("Failed to hash message: %v\n", err)
		os.Exit(1)
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := crypto.Keccak256Hash(rawData)

	sig, err := crypto.Sign(hash.Bytes(), privKey)
	if err != nil {
		fmt.Printf("Failed to sign: %v\n", err)
		os.Exit(1)
	}

	sig[64] += 27 // V adjust
	signature := fmt.Sprintf("0x%x", sig)

	req, err := http.NewRequest("POST", "https://clob.polymarket.com/auth/api-key", nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("POLY_ADDRESS", address)
	req.Header.Set("POLY_SIGNATURE", signature)
	req.Header.Set("POLY_TIMESTAMP", timestampStr)
	req.Header.Set("POLY_NONCE", nonce)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Polymarket API rejected signature. Status: %d\nBody: %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	var result struct {
		ApiKey     string `json:"apiKey"`
		Secret     string `json:"secret"`
		Passphrase string `json:"passphrase"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=======================================================")
	fmt.Println("🎉 POLYMARKET CLOB API KEYS SUCCESSFULLY GENERATED 🎉")
	fmt.Println("=======================================================")
	fmt.Printf("POLY_API_KEY        : %s\n", result.ApiKey)
	fmt.Printf("POLY_API_SECRET     : %s\n", result.Secret)
	fmt.Printf("POLY_API_PASSPHRASE : %s\n", result.Passphrase)
	fmt.Println("=======================================================")
	fmt.Println("You can now safely paste these into the Polymarket plugin config UI.")
}
