# allele-exchange-polymarket

**Type**: Exchange Adapter Plugin
**Description**: Polymarket CLOB integration (WASI Network Enabled)

## Architecture (DSS ARC Panel Reviewed)
This WASM plugin runs under the engine's `Exchange` execution scope. It is granted network sockets via WASI to:
1. Stream the Polymarket WebSocket (`wss://ws-subscriptions-clob.polymarket.com/ws/market`)
2. Normalize it into the standardized `MarketState` struct.
3. Sign EIP-712 formatted limit orders using the injected `Wallet` private key when `TradeSignals` are dispatched to it.

**Panel Findings Addressed:**
- Enforced checking the `negRisk` boolean flag, as Polygon collateral mechanics completely change when trading mutually exclusive negative-risk markets.
- Removed reliance on `golang/crypto` CGO extensions; using purely math-based EIP-712 libraries that compile safely to `wasip1`.

## Compilation
```bash
GOOS=wasip1 GOARCH=wasm go build -o adapter.wasm main.go
```
