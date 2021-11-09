# Kucoin API - Golang

Go version **go1.17**
### How to run?
1. `go build`
2. Create new file `config.yaml` with the following info:
```
key:  <YOUR-API-KEY>
secret: <YOUR-API-SECRET>
passphrase: <YOUR-API-PASSPHRASE>
targetTime: <TRAGET-TIME>   ==> target time for first BUY order to be sent to exchange
```
3. `go run .\main.go --config config.yaml`
