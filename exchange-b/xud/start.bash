cd $GOPATH/src/github.com/offerm/swap-resolver
go run xud.go --listen localhost:8896 --peer localhost:8886 --lnd-rpc-ltc localhost:20001 --lnd-rpc-btc localhost:20002
