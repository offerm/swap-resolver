cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver
go run resolver.go --listen localhost:7002 --peer localhost:7001 --lnd-rpc-ltc localhost:20001 --lnd-rpc-btc localhost:20002
