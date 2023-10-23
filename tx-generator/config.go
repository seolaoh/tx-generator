package tx_generator

import (
	"crypto/ecdsa"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

type Config struct {
	privKey    *ecdsa.PrivateKey
	from       common.Address
	txType     uint64
	txInterval time.Duration
	chainID    *big.Int
	ethclient  *ethclient.Client
	signer     types.Signer
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	dummyTxType := ctx.Uint64(DummyTransactionTypeFlag.Name)
	dummyTxAccPrivKey := ctx.String(DummyTransactionAccPrivateKeyFlag.Name)
	dummyTxSendInterval := ctx.Duration(DummyTransactionSendIntervalFlag.Name)
	chainID := ctx.Uint64(ChainIDFlag.Name)
	l2EthRpc := ctx.String(L2EthRpcFlag.Name)

	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(dummyTxAccPrivKey, "0x"))
	if err != nil {
		return nil, err
	}
	from := crypto.PubkeyToAddress(privKey.PublicKey)

	intChainID := new(big.Int).SetUint64(chainID)

	client, err := ethclient.Dial(l2EthRpc)
	if err != nil {
		return nil, err
	}

	signer := types.NewCancunSigner(big.NewInt(int64(chainID)))

	return &Config{
		privKey:    privKey,
		from:       from,
		txType:     dummyTxType,
		txInterval: dummyTxSendInterval,
		chainID:    intChainID,
		ethclient:  client,
		signer:     signer,
	}, nil
}
