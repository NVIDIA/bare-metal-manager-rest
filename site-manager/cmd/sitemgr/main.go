// Package main is the command entry point
package main

import (
	"context"
	"os"

	cli "github.com/urfave/cli/v2"
	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
	"github.com/nvidia/carbide-rest/site-manager/pkg/sitemgr"
)

func main() {
	cmd := sitemgr.NewCommand()
	app := &cli.App{
		Name:    cmd.Name,
		Usage:   cmd.Usage,
		Version: "0.1.0",
		Flags:   cmd.Flags,
		Action:  cmd.Action,
	}

	ctx := core.NewDefaultContext(context.Background())
	log := core.GetLogger(ctx)
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
