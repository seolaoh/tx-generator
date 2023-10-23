package tx_generator

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
	"github.com/urfave/cli/v2"

	"github.com/seolaoh/tx-generator/tx-generator/utils"
)

var zeroAddr = common.Address{0}

func Main(cliCtx *cli.Context) error {
	fmt.Println("initializing tx generator")

	generator, err := NewTxGenerator(cliCtx)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to create tx generator"))
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generator.Start(ctx)
	<-utils.WaitInterrupt()
	generator.Stop()

	return nil
}

type TxGenerator struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *Config
}

func NewTxGenerator(ctx *cli.Context) (*TxGenerator, error) {
	config, err := NewConfig(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to create config"))
		return nil, err
	}

	return &TxGenerator{
		cfg: config,
	}, nil
}

func (g TxGenerator) Start(ctx context.Context) {
	g.ctx, g.cancel = context.WithCancel(ctx)

	ticker := time.NewTicker(g.cfg.txInterval)
	defer ticker.Stop()

	fmt.Printf("tx generator started: tx type %d\n", g.cfg.txType)

	for ; ; <-ticker.C {
		select {
		case <-g.ctx.Done():
			fmt.Println("stopping tx generator")
			return
		default:
			if err := g.generateTx(g.cfg.txType); err != nil {
				fmt.Println(fmt.Errorf("failed to generate dummy tx: %w", err))
			}
		}
	}
}

func (g TxGenerator) Stop() {
	g.cancel()
}

func (g TxGenerator) generateTx(txType uint64) error {
	fmt.Println("generating dummy tx...")

	tx, err := g.getTx(txType)
	if err != nil {
		return fmt.Errorf("failed to get dummy tx: %w", err)
	}

	signedTx, err := types.SignTx(tx, g.cfg.signer, g.cfg.privKey)
	if err != nil {
		return fmt.Errorf("failed to sign tx: %w", err)
	}
	fmt.Printf("signed transaction hash: %#x\n", signedTx.Hash())

	if err = g.cfg.ethclient.SendTransaction(g.ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	return nil
}

func (g TxGenerator) getTx(txType uint64) (*types.Transaction, error) {
	str := "hello kroma"
	blob := kzg4844.Blob{}
	copy(blob[:], str)
	blobCommit, _ := kzg4844.BlobToCommitment(blob)
	blobProof, _ := kzg4844.ComputeBlobProof(blob, blobCommit)
	blobVHash := utils.BlobToVHash(blobCommit)

	nonce, err := g.cfg.ethclient.NonceAt(g.ctx, g.cfg.from, nil)
	//nonce, err := g.cfg.ethclient.PendingNonceAt(g.ctx, g.cfg.from)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	accessList := types.AccessList{
		types.AccessTuple{
			Address:     zeroAddr,
			StorageKeys: []common.Hash{common.BigToHash(common.Big0)},
		},
	}

	var txData types.TxData
	switch txType {
	case 0:
		txData = &types.LegacyTx{
			Nonce:    nonce,
			GasPrice: hexutil.MustDecodeBig("0x600f1f"),
			Gas:      21000,
			To:       &zeroAddr,
			Value:    common.Big1,
			Data:     hexutil.MustDecode("0x"),
		}
	case 1:
		txData = &types.AccessListTx{
			ChainID:    g.cfg.chainID,
			Nonce:      nonce,
			GasPrice:   hexutil.MustDecodeBig("0x600f1f"),
			Gas:        30000,
			To:         &zeroAddr,
			Value:      common.Big1,
			Data:       hexutil.MustDecode("0x"),
			AccessList: accessList,
		}
	case 2:
		txData = &types.DynamicFeeTx{
			ChainID:    g.cfg.chainID,
			Nonce:      nonce,
			Gas:        30000,
			GasFeeCap:  new(big.Int).SetInt64(1000000005),
			GasTipCap:  new(big.Int).SetInt64(1000000005),
			To:         &zeroAddr,
			Value:      common.Big1,
			Data:       hexutil.MustDecode("0x"),
			AccessList: accessList,
		}
	case 3:
		txData = &types.BlobTx{
			ChainID:    uint256.MustFromBig(g.cfg.chainID),
			Nonce:      nonce,
			GasTipCap:  uint256.NewInt(1000000005),
			GasFeeCap:  uint256.NewInt(1000000005),
			Gas:        21000,
			BlobFeeCap: uint256.NewInt(20000000000),
			BlobHashes: []common.Hash{blobVHash},
			Value:      uint256.MustFromBig(common.Big1),
			To:         zeroAddr,
			Sidecar: &types.BlobTxSidecar{
				Blobs:       []kzg4844.Blob{blob},
				Commitments: []kzg4844.Commitment{blobCommit},
				Proofs:      []kzg4844.Proof{blobProof},
			},
		}
	}

	return types.NewTx(txData), nil
}
