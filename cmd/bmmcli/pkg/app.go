// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package bmmcli

import (
	"context"
	"fmt"
	"os"

	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/pkg/interactive"
	"github.com/nvidia/bare-metal-manager-rest/sdk/standard"
	cli "github.com/urfave/cli/v2"
)

// NewApp builds a cli.App from the embedded OpenAPI spec data.
func NewApp(specData []byte) (*cli.App, error) {
	spec, err := ParseSpec(specData)
	if err != nil {
		return nil, fmt.Errorf("parsing embedded spec: %w", err)
	}

	defaultBaseURL := ""
	if len(spec.Servers) > 0 {
		defaultBaseURL = spec.Servers[0].URL
	}

	cfg, _ := LoadConfig()

	commands := BuildCommands(spec)
	commands = append(commands, LoginCommand())
	commands = append(commands, interactiveCommand())
	commands = append(commands, initCommand())
	commands = append(commands, completionCommand())

	app := &cli.App{
		Name:                 "bmmcli",
		Usage:                spec.Info.Title,
		Version:              spec.Info.Version,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Path to config file",
				EnvVars: []string{"BMM_CONFIG"},
				Value:   ConfigPath(),
			},
			&cli.StringFlag{
				Name:    "base-url",
				Usage:   "API base URL",
				EnvVars: []string{"BMM_BASE_URL"},
				Value:   configDefault(cfg.API.Base, defaultBaseURL),
			},
			&cli.StringFlag{
				Name:    "org",
				Usage:   "Organization name",
				EnvVars: []string{"BMM_ORG"},
				Value:   cfg.API.Org,
			},
			&cli.StringFlag{
				Name:    "token",
				Usage:   "API bearer token",
				EnvVars: []string{"BMM_TOKEN"},
			},
			&cli.StringFlag{
				Name:  "token-command",
				Usage: "Shell command that prints a bearer token",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output format: json, yaml, table",
				Value: "json",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
			&cli.StringFlag{
				Name:    "token-url",
				Usage:   "OIDC token endpoint URL",
				EnvVars: []string{"BMM_TOKEN_URL"},
			},
			&cli.StringFlag{
				Name:    "keycloak-url",
				Usage:   "Keycloak base URL (constructs token-url if --token-url is not set)",
				EnvVars: []string{"BMM_KEYCLOAK_URL"},
			},
			&cli.StringFlag{
				Name:    "realm",
				Usage:   "Keycloak realm (used with --keycloak-url)",
				EnvVars: []string{"BMM_REALM"},
				Value:   "carbide-dev",
			},
			&cli.StringFlag{
				Name:    "client-id",
				Usage:   "OAuth client ID",
				EnvVars: []string{"BMM_CLIENT_ID"},
				Value:   "carbide-api",
			},
		},
		Commands: commands,
	}

	return app, nil
}

// interactiveCommand launches the TUI REPL.
// It first shows a config file selector (when multiple configs exist in ~/.bmm/),
// then starts the session. A token is not required — the user can type 'login'.
func interactiveCommand() *cli.Command {
	return &cli.Command{
		Name:    "interactive",
		Aliases: []string{"i"},
		Usage:   "Start interactive REPL mode with autocomplete and scope filtering",
		Action: func(c *cli.Context) error {
			// If --config was explicitly set, use it directly.
			// Otherwise show the picker for multiple configs in ~/.bmm/.
			var cfgPath string
			var cfg *ConfigFile
			var err error
			if c.IsSet("config") {
				cfgPath = c.String("config")
				cfg, err = LoadConfigFromPath(cfgPath)
				if err != nil {
					return fmt.Errorf("loading config %s: %w", cfgPath, err)
				}
				fmt.Printf("Using: %s\n\n", cfgPath)
			} else {
				cfgPath, cfg, err = ChooseConfigFile()
				if err != nil {
					return fmt.Errorf("loading config: %w", err)
				}
			}

			// Resolve base URL: explicit flag > selected config > flag default (spec URL).
			baseURL := cfg.API.Base
			if c.IsSet("base-url") {
				baseURL = c.String("base-url")
			}
			if baseURL == "" {
				baseURL = c.String("base-url") // flag default (spec server URL)
			}

			// Resolve org: explicit flag > selected config.
			org := cfg.API.Org
			if c.IsSet("org") {
				org = c.String("org")
			}
			if org == "" {
				return fmt.Errorf("org is required: set api.org in %s or pass --org", cfgPath)
			}

			// Resolve token: flag → token-command → stored config token → silent API-key exchange.
			token := c.String("token")
			if token == "" {
				token, _ = AutoRefreshToken(cfg)
			}
			if t, err := ResolveToken(token, c.String("token-command")); err == nil && t != "" {
				token = t
			}
			// For API-key configs, silently exchange on startup so the user
			// doesn't have to type 'login' for the first command of the session.
			if token == "" && HasAPIKeyConfig(cfg) {
				if t, err := doInteractiveAPIKeyLogin(cfg, cfgPath); err == nil {
					token = t
				}
			}

			configuration := standard.NewConfiguration()
			configuration.Servers = standard.ServerConfigurations{{URL: baseURL}}
			configuration.HTTPClient = interactive.NewHTTPClient(cfg.API.Name, cfg.Debug.IsVerbose())
			apiClient := standard.NewAPIClient(configuration)

			ctx := context.Background()
			if token != "" {
				ctx = context.WithValue(ctx, standard.ContextAccessToken, token)
			}

			session := interactive.NewSession(apiClient, ctx, org)
			session.Token = token
			session.ConfigPath = cfgPath
			session.LoginFn = InteractiveLoginFn(cfg, cfgPath)
			session.Verbose = cfg.Debug.IsVerbose()

			if token == "" {
				fmt.Fprintf(os.Stderr, "\033[33mWarning:\033[0m No auth token found. Type \033[1mlogin\033[0m to authenticate.\n\n")
			}

			return interactive.RunREPL(session)
		},
	}
}

// initCommand writes a sample config file to ~/.bmm/config.yaml.
func initCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Create a sample config file at ~/.bmm/config.yaml",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "Write to this path instead of the default",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Overwrite existing config file",
			},
		},
		Action: func(c *cli.Context) error {
			cfgPath := c.String("config")
			if cfgPath == "" {
				cfgPath = ConfigPath()
			}
			force := c.Bool("force")

			if !force {
				if _, err := os.Stat(cfgPath); err == nil {
					return fmt.Errorf("config already exists at %s (use --force to overwrite)", cfgPath)
				}
			}

			if err := os.MkdirAll(ConfigDir(), 0700); err != nil {
				return fmt.Errorf("creating config directory: %w", err)
			}

			if err := os.WriteFile(cfgPath, []byte(sampleConfigContent), 0600); err != nil {
				return fmt.Errorf("writing config: %w", err)
			}

			fmt.Printf("Config written to %s\n", cfgPath)
			return nil
		},
	}
}

func completionCommand() *cli.Command {
	return &cli.Command{
		Name:  "completion",
		Usage: "Output shell completion script",
		Subcommands: []*cli.Command{
			{
				Name:  "bash",
				Usage: "Output bash completion script",
				Action: func(c *cli.Context) error {
					fmt.Print(bashCompletion)
					return nil
				},
			},
			{
				Name:  "zsh",
				Usage: "Output zsh completion script",
				Action: func(c *cli.Context) error {
					fmt.Print(zshCompletion)
					return nil
				},
			},
			{
				Name:  "fish",
				Usage: "Output fish completion script",
				Action: func(c *cli.Context) error {
					fmt.Print(fishCompletion)
					return nil
				},
			},
		},
	}
}

const bashCompletion = `# bash completion for bmmcli
# Add to ~/.bashrc:  eval "$(bmmcli completion bash)"
_bmmcli_complete() {
    local cur opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$(${COMP_WORDS[0]} --generate-bash-completion "${COMP_WORDS[@]:1:$COMP_CWORD}")
    COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
    return 0
}
complete -o default -F _bmmcli_complete bmmcli
`

const zshCompletion = `# zsh completion for bmmcli
# Add to ~/.zshrc:  eval "$(bmmcli completion zsh)"
_bmmcli_complete() {
    local -a opts
    opts=(${(f)"$(${words[1]} --generate-bash-completion ${words:1:$CURRENT-1})"})
    _describe 'bmmcli' opts
}
compdef _bmmcli_complete bmmcli
`

const fishCompletion = `# fish completion for bmmcli
# Add to ~/.config/fish/completions/bmmcli.fish or run:
#   bmmcli completion fish > ~/.config/fish/completions/bmmcli.fish
complete -c bmmcli -f -a '(bmmcli --generate-bash-completion (commandline -cop))'
`

func configDefault(cfgValue, fallback string) string {
	if cfgValue != "" {
		return cfgValue
	}
	return fallback
}
