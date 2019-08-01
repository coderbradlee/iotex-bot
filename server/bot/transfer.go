// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package bot

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/lzxm160/iotex-bot/config"
	"github.com/lzxm160/iotex-bot/pkg/util/grpcutil"
)

// Transfer
type Transfer struct {
	cfg    config.Config
	ctx    context.Context
	cancel context.CancelFunc
	name   string
}

// NewTransfer
func NewTransfer(cfg config.Config, name string) (Service, error) {
	return newTransfer(cfg, name)
}

func newTransfer(cfg config.Config, name string) (Service, error) {
	svr := Transfer{
		cfg:  cfg,
		name: name,
	}
	return &svr, nil
}

// Start starts the server
func (s *Transfer) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	return s.startTransfer()
}

// Stop stops the server
func (s *Transfer) Stop() error {
	s.cancel()
	return nil
}

// Name
func (s *Transfer) Name() string {
	return s.name
}

func (s *Transfer) startTransfer() error {
	// load keystore
	pri, err := s.getPrivateKey()
	if err != nil {
		return err
	}
	err = s.transfer(pri)
	if err != nil {
		return err
	}
	// check if timeout
	return nil
}
func (s *Transfer) transfer(pri crypto.PrivateKey) error {
	conn, err := grpcutil.ConnectToEndpoint(s.cfg.API.Url, false)
	if err != nil {
		return err
	}
	defer conn.Close()
	cli := iotexapi.NewAPIServiceClient(conn)

	from, err := address.FromBytes(pri.PublicKey().Hash())
	if err != nil {
		return err
	}
	request := iotexapi.GetAccountRequest{Address: from.String()}
	response, err := cli.GetAccount(context.Background(), &request)
	if err != nil {
		return err
	}
	nonce := response.AccountMeta.PendingNonce

	gasprice := big.NewInt(0).SetUint64(s.cfg.Transfer.GasPrice)
	fmt.Println(gasprice.String())
	tx, err := action.NewTransfer(nonce, big.NewInt(1000000000000000000),
		s.cfg.Transfer.To[0], nil, s.cfg.Transfer.GasLimit, gasprice)
	if err != nil {
		return err
	}
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetGasLimit(s.cfg.Transfer.GasLimit).
		SetGasPrice(gasprice).
		SetAction(tx).Build()
	selp, err := action.Sign(elp, pri)
	if err != nil {
		return err
	}
	req := &iotexapi.SendActionRequest{Action: selp.Proto()}
	if _, err = cli.SendAction(context.Background(), req); err != nil {
		return err
	}
	shash := hash.Hash256b(byteutil.Must(proto.Marshal(selp.Proto())))
	txhash := hex.EncodeToString(shash[:])
	log.L().Info("transfer:", zap.String("transfer hash", txhash), zap.Uint64("nonce", nonce), zap.String("from", s.cfg.Transfer.From[0]), zap.String("to", s.cfg.Transfer.To[0]))
	return nil
}
func (s *Transfer) getPrivateKey() (crypto.PrivateKey, error) {
	ks := keystore.NewKeyStore(s.cfg.Wallet,
		keystore.StandardScryptN, keystore.StandardScryptP)

	from, err := address.FromString(s.cfg.Transfer.From[0])
	if err != nil {
		return nil, err
	}
	for _, account := range ks.Accounts() {
		if bytes.Equal(from.Bytes(), account.Address.Bytes()) {
			return crypto.KeystoreToPrivateKey(account, s.cfg.Transfer.From[1])
		}
	}
	return nil, errors.New("src address not found")
}