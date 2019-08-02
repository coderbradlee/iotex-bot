module github.com/lzxm160/iotex-bot

go 1.12

require (
	github.com/ethereum/go-ethereum v1.8.27
	github.com/iotexproject/iotex-address v0.2.0
	github.com/iotexproject/iotex-core v0.8.0
	github.com/pkg/errors v0.8.1
	go.uber.org/automaxprocs v1.2.0
	go.uber.org/config v1.3.1
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.21.0
)

replace github.com/ethereum/go-ethereum => github.com/iotexproject/go-ethereum v0.2.1-0.20190723221211-a9fbf57d1cb7
