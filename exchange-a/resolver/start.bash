cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver
go run resolver.go --listen localhost:7001 --peer localhost:7002 --lnd-rpc-ltc localhost:10001 --lnd-rpc-btc localhost:10002
