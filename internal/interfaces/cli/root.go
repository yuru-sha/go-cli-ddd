package cli

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/logger"
)

// RootOptions contains root-level CLI flags resolved before command dispatch.
type RootOptions struct {
	ConfigPath string
	Env        string
	Help       bool
}

// NewRootCommand creates the root CLI command dispatcher.
func NewRootCommand() *RootCommand {
	return &RootCommand{
		configPath: "configs/config.yaml",
		modules:    make(map[string]CommandModule),
	}
}

// Register attaches one subcommand to the root command.
func (r *RootCommand) Register(module CommandModule) {
	r.modules[module.Name()] = module
}

// Execute parses global flags, initializes runtime config, and dispatches the subcommand.
func (r *RootCommand) Execute(args []string) error {
	opts, commandName, moduleArgs, err := ParseRootArgs(args, r.configPath)
	if err != nil {
		return err
	}
	r.configPath = opts.ConfigPath
	r.env = opts.Env

	if opts.Help || commandName == "" {
		r.printUsage()
		return nil
	}

	module, ok := r.modules[commandName]
	if !ok {
		return fmt.Errorf("未対応のコマンドです: %s", commandName)
	}

	cfgOpts := config.NewConfigOptions(r.configPath, r.env)
	cfg, err := config.LoadConfig(cfgOpts)
	if err != nil {
		return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	logger.InitLogger(cfg.App.LogLevel, cfg.App.Debug)

	slog.Info(
		"アプリケーションを起動しました",
		"app_name", cfg.App.Name,
		"log_level", cfg.App.LogLevel,
		"debug", cfg.App.Debug,
		"command", module.Name(),
	)

	return module.Execute(context.Background(), moduleArgs)
}

func (r *RootCommand) printUsage() {
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	sort.Strings(names)

	var b strings.Builder
	b.WriteString("Usage: go-cli-ddd [--config path] [--env local|dev|prd] <command> [flags]\n\n")
	b.WriteString("Commands:\n")
	for _, name := range names {
		b.WriteString("  ")
		b.WriteString(name)
		b.WriteString("\n")
	}
	_, _ = io.WriteString(os.Stdout, b.String())
}

// ParseRootArgs extracts root options and the target subcommand from raw CLI args.
func ParseRootArgs(args []string, defaultConfigPath string) (RootOptions, string, []string, error) {
	opts := RootOptions{ConfigPath: defaultConfigPath}
	commandName := ""
	commandArgs := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-h" || arg == "--help":
			opts.Help = true
		case arg == "--config":
			if i+1 >= len(args) {
				return RootOptions{}, "", nil, fmt.Errorf("--config の値が不足しています")
			}
			i++
			opts.ConfigPath = args[i]
		case strings.HasPrefix(arg, "--config="):
			opts.ConfigPath = strings.TrimPrefix(arg, "--config=")
		case arg == "--env":
			if i+1 >= len(args) {
				return RootOptions{}, "", nil, fmt.Errorf("--env の値が不足しています")
			}
			i++
			opts.Env = args[i]
		case strings.HasPrefix(arg, "--env="):
			opts.Env = strings.TrimPrefix(arg, "--env=")
		case commandName == "" && !strings.HasPrefix(arg, "-"):
			commandName = arg
		case commandName == "":
			return RootOptions{}, "", nil, fmt.Errorf("未対応のグローバルフラグです: %s", arg)
		default:
			commandArgs = append(commandArgs, arg)
		}
	}

	return opts, commandName, commandArgs, nil
}
