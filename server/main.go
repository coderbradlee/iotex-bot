// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

// Usage:
//   make build
//   ./bin/server -config-file=./config.yaml
//

package main

import (
	"flag"
	"fmt"
	glog "log"
	"os"

	"github.com/lzxm160/iotex-bot/config"
	"github.com/lzxm160/iotex-bot/pkg/log"

	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr,
			"usage: server -config-path=[string]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
}

func main() {
	cfg, err := config.New()
	if err != nil {
		glog.Fatalln("Failed to new config.", zap.Error(err))
	}
	initLogger(cfg)

	// liveness start
	//probeSvr := bot.New(cfg.System.HTTPStatsPort)
	//if err := probeSvr.Start(ctx); err != nil {
	//	log.L().Fatal("Failed to start probe server.", zap.Error(err))
	//}
	log.L().Info("okkkkk")
}

func initLogger(cfg config.Config) {
	if err := log.InitLoggers(cfg.Log); err != nil {
		glog.Println("Cannot config global logger, use default one: ", err)
	}
}
