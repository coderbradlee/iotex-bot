// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package bot

//
//import (
//	"context"
//	"fmt"
//	"net/http"
//	"net/http/pprof"
//	"runtime"
//	"sync"
//	"time"
//
//	"github.com/pkg/errors"
//	"go.uber.org/zap"
//
//	"github.com/iotexproject/iotex-core/action/protocol"
//	"github.com/iotexproject/iotex-core/action/protocol/account"
//	"github.com/iotexproject/iotex-core/action/protocol/execution"
//	"github.com/iotexproject/iotex-core/action/protocol/multichain/mainchain"
//	"github.com/iotexproject/iotex-core/action/protocol/poll"
//	"github.com/iotexproject/iotex-core/action/protocol/rewarding"
//	"github.com/iotexproject/iotex-core/action/protocol/rolldpos"
//	"github.com/iotexproject/iotex-core/chainservice"
//	"github.com/iotexproject/iotex-core/config"
//	"github.com/iotexproject/iotex-core/dispatcher"
//	"github.com/iotexproject/iotex-core/p2p"
//	"github.com/iotexproject/iotex-core/pkg/ha"
//	"github.com/iotexproject/iotex-core/pkg/log"
//	"github.com/iotexproject/iotex-core/pkg/probe"
//	"github.com/iotexproject/iotex-core/pkg/routine"
//	"github.com/iotexproject/iotex-core/pkg/util/httputil"
//)
//
//// Server is the iotex server instance containing all components.
//type Server struct {
//	cfg                  config.Config
//	rootChainService     *chainservice.ChainService
//	chainservices        map[uint32]*chainservice.ChainService
//	p2pAgent             *p2p.Agent
//	dispatcher           dispatcher.Dispatcher
//	mainChainProtocol    *mainchain.Protocol
//	initializedSubChains map[uint32]bool
//	mutex                sync.RWMutex
//	subModuleCancel      context.CancelFunc
//}
//
//// NewServer creates a new server
//// TODO clean up config, make root config contains network, dispatch and chainservice
//func NewServer(cfg config.Config) (*Server, error) {
//	return newServer(cfg, false)
//}
//
//
//func newServer(cfg config.Config, testing bool) (*Server, error) {
//	// create dispatcher instance
//	dispatcher, err := dispatcher.NewDispatcher(cfg)
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to create dispatcher")
//	}
//	p2pAgent := p2p.NewAgent(cfg, dispatcher.HandleBroadcast, dispatcher.HandleTell)
//	chains := make(map[uint32]*chainservice.ChainService)
//	var cs *chainservice.ChainService
//	var opts []chainservice.Option
//	if testing {
//		opts = []chainservice.Option{
//			chainservice.WithTesting(),
//		}
//	}
//	cs, err = chainservice.New(cfg, p2pAgent, dispatcher, opts...)
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to create chain service")
//	}
//
//	// Add action validators
//	cs.ActionPool().
//		AddActionEnvelopeValidators(
//			protocol.NewGenericValidator(cs.Blockchain()),
//		)
//	cs.Blockchain().Validator().
//		AddActionEnvelopeValidators(
//			protocol.NewGenericValidator(cs.Blockchain()),
//		)
//	// Install protocols
//	if err := registerDefaultProtocols(cs, cfg); err != nil {
//		return nil, err
//	}
//	mainChainProtocol := mainchain.NewProtocol(cs.Blockchain())
//	if err := cs.RegisterProtocol(mainchain.ProtocolID, mainChainProtocol); err != nil {
//		return nil, err
//	}
//	// TODO: explorer dependency deleted here at #1085, need to revive by migrating to api
//	chains[cs.ChainID()] = cs
//	dispatcher.AddSubscriber(cs.ChainID(), cs)
//	svr := Server{
//		cfg:                  cfg,
//		p2pAgent:             p2pAgent,
//		dispatcher:           dispatcher,
//		rootChainService:     cs,
//		chainservices:        chains,
//		mainChainProtocol:    mainChainProtocol,
//		initializedSubChains: map[uint32]bool{},
//	}
//	// Setup sub-chain starter
//	// TODO: sub-chain infra should use main-chain API instead of protocol directly
//	return &svr, nil
//}
//
//// Start starts the server
//func (s *Server) Start(ctx context.Context) error {
//	cctx, cancel := context.WithCancel(context.Background())
//	s.subModuleCancel = cancel
//	if err := s.p2pAgent.Start(cctx); err != nil {
//		return errors.Wrap(err, "error when starting P2P agent")
//	}
//	if err := s.rootChainService.Blockchain().AddSubscriber(s); err != nil {
//		return errors.Wrap(err, "error when starting sub-chain starter")
//	}
//	for _, cs := range s.chainservices {
//		if err := cs.Start(cctx); err != nil {
//			return errors.Wrap(err, "error when starting blockchain")
//		}
//	}
//	if err := s.dispatcher.Start(cctx); err != nil {
//		return errors.Wrap(err, "error when starting dispatcher")
//	}
//
//	return nil
//}
//
//// Stop stops the server
//func (s *Server) Stop(ctx context.Context) error {
//	defer s.subModuleCancel()
//	if err := s.p2pAgent.Stop(ctx); err != nil {
//		return errors.Wrap(err, "error when stopping P2P agent")
//	}
//	if err := s.dispatcher.Stop(ctx); err != nil {
//		return errors.Wrap(err, "error when stopping dispatcher")
//	}
//	if err := s.rootChainService.Blockchain().RemoveSubscriber(s); err != nil {
//		return errors.Wrap(err, "error when unsubscribing root chain block creation")
//	}
//	for _, cs := range s.chainservices {
//		if err := cs.Stop(ctx); err != nil {
//			return errors.Wrap(err, "error when stopping blockchain")
//		}
//	}
//	return nil
//}
//
//
//// StopChainService stops the chain service run in the server.
//func (s *Server) StopChainService(ctx context.Context, id uint32) error {
//	s.mutex.RLock()
//	defer s.mutex.RUnlock()
//	c, ok := s.chainservices[id]
//	if !ok {
//		return errors.New("Chain ID does not match any existing chains")
//	}
//	return c.Stop(ctx)
//}
