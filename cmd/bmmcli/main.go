// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// bmmcli is the Carbide Bare Metal Manager CLI.
//
// Commands are generated at runtime from the embedded OpenAPI spec, so the CLI
// stays in sync with API changes without any code generation.  An interactive
// REPL mode (bmmcli interactive) provides autocomplete, scope filtering, and
// guided wizards on top of the same API.
package main

import (
	"context"
	"fmt"
	"os"

	bmmcli "github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/pkg"
	"github.com/nvidia/bare-metal-manager-rest/openapi"
)

func main() {
	ctx := context.Background()

	app, err := bmmcli.NewApp(openapi.Spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
