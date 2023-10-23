package tx_generator

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	DummyTransactionTypeFlag = &cli.Uint64Flag{
		Name:     "dummy-tx-type",
		Usage:    "generate dummy transaction of given type repeatedly",
		EnvVars:  []string{"DUMMY_TX_TYPE"},
		Required: true,
	}
	DummyTransactionAccPrivateKeyFlag = &cli.StringFlag{
		Name:     "dummy-tx-account",
		Usage:    "EOA to generate dummy transaction",
		EnvVars:  []string{"DUMMY_TX_ACC_PRIV_KEY"},
		Required: true,
	}
	ChainIDFlag = &cli.Uint64Flag{
		Name:     "chain-id",
		Usage:    "chain id to send transaction",
		EnvVars:  []string{"CHAIN_ID"},
		Required: true,
	}
	L2EthRpcFlag = &cli.StringFlag{
		Name:     "l2-eth-rpc",
		Usage:    "L2 ETH RPC",
		EnvVars:  []string{"L2_ETH_ROC"},
		Required: true,
	}
	DummyTransactionSendIntervalFlag = &cli.DurationFlag{
		Name:     "dummy-tx-send-interval",
		Usage:    "Interval second to send dummy transaction",
		EnvVars:  []string{"DUMMY_TX_SEND_INTERVAL"},
		Required: false,
		Value:    1 * time.Second,
	}
)

var requiredFlags = []cli.Flag{
	DummyTransactionTypeFlag,
	DummyTransactionAccPrivateKeyFlag,
	ChainIDFlag,
	L2EthRpcFlag,
}

var optionalFlags = []cli.Flag{
	DummyTransactionSendIntervalFlag,
}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag
