cd $GOPATH/src/github.com/offerm/swap-resolver
go run xud.go --listen localhost:8886 --peer localhost:8896 --lnd-rpc-ltc localhost:10001 --lnd-rpc-btc localhost:10002
