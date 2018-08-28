cd $GOPATH/src/github.com/offerm/swap-resolver
go run xud.go --listen localhost:7002 --peer localhost:7001 --lnd-rpc-ltc localhost:20001 --lnd-rpc-btc localhost:20002
