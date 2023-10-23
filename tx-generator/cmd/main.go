package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	txGenerator "github.com/seolaoh/tx-generator/tx-generator"
)

var (
	Version   = ""
	GitCommit = ""
	GitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Flags = txGenerator.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	app.Name = "tx-generator"
	app.Usage = "Dummy tx generator"

	app.Action = curryMain()
	err := app.Run(os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

// curryMain transforms the txGenerator.Main function into an app.Action
// This is done to capture the Version of the tx generator.
func curryMain() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return txGenerator.Main(ctx)
	}
}
