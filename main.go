package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/kshyr/pact/config"
	"github.com/kshyr/pact/demon"
	"github.com/kshyr/pact/invoker"
	"github.com/kshyr/pact/logger"
	"github.com/kshyr/pact/script"
	"github.com/kshyr/pact/timekeeper"
	"github.com/spf13/cobra"
)

func demonLoop(
	d *demon.Demon,
	tk *timekeeper.Timekeeper,
	cfg *config.Config,
	reg *script.ScriptRegistry,
) {
}

func main() {
	if demon.IsDemonExecutable() {
		d, err := demon.New()
		if err != nil {
			fmt.Printf("Error initializing demon: %v\n", err)
			os.Exit(1)
		}

		cfg, err := config.NewConfig()
		if err != nil {
			d.Log.Errorf("Failed to load config: %v", err)
			return
		}

		registryPath := filepath.Join(cfg.ScriptsDir, "scripts.toml")
		reg, err := script.NewRegistry(registryPath)
		if err != nil {
			d.Log.Errorf("Failed to load script registry: %v", err)
			return
		}

		tk := timekeeper.New(reg.Scripts, *cfg)
		tk.Start()
		defer tk.Stop()

		err = d.Run(func() {
			demonLoop(d, tk, cfg, reg)
		})

		if err := demon.RemovePIDFile(); err != nil {
			d.Log.Error(err)
		}

		if err != nil {
			d.Log.Errorf("demon couldn't start: %v", err)
			os.Exit(1)
		}

		for {
		}

	}

	var rootCmd = &cobra.Command{
		Use:   "pact",
		Short: "Pact CLI Application",
	}

	var demonCmd = &cobra.Command{
		Use:   "demon",
		Short: "Manage the demon process",
	}

	var demonStartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the demon",
		Run: func(_ *cobra.Command, args []string) {
			if demon.IsRunning() {
				fmt.Println("demon is already running")
				os.Exit(1)
			}

			exePath, err := os.Executable()
			if err != nil {
				fmt.Printf("failed to get executable path: %w\n", err)
				os.Exit(1)
			}

			cmd := exec.Command(exePath)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setsid: true,
			}

			logFilePath, err := logger.CreateLogFile("demon")
			if err != nil {
				fmt.Printf("failed to create log file: %w\n", err)
				os.Exit(1)
			}
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Printf("failed to open log file: %w\n", err)
			}

			cmd.Env = append(os.Environ(), "DEMON=1")
			cmd.Stdout = file
			cmd.Stderr = file

			if err := cmd.Start(); err != nil {
				fmt.Printf("failed to start demon process: %w\n", err)
				os.Exit(1)
			}

			if err = demon.WritePIDFile(cmd.Process.Pid); err != nil {
				fmt.Printf("failed to write pid file", err)
			}

			os.Exit(0)
		},
	}

	var demonStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check demon status",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := demon.New()
			if err != nil {
				fmt.Printf("Error initializing demon: %v\n", err)
				os.Exit(1)
			}

			running := demon.IsRunning()
			if err != nil {
				d.Log.Errorf("Error checking status: %v\n", err)
				os.Exit(1)
			}

			if running {
				d.Log.Println("Demon is running")
			} else {
				d.Log.Println("Demon is not running")
			}
		},
	}

	var demonStopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop the demon",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := demon.New()
			if err != nil {
				fmt.Printf("Error initializing demon: %v\n", err)
				os.Exit(1)
			}

			err = d.Stop()
			if err != nil {
				d.Log.Errorf("Error stopping demon: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Demon stopped")
		},
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage Pact config",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.NewConfig()
			if err != nil {
				fmt.Printf("config error: %v\n", err)
				os.Exit(1)
			}

			tomlContents, err := toml.Marshal(cfg)
			if err != nil {
				fmt.Printf("marshaling to toml error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("config reads: \n%s", string(tomlContents))
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all the items",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.NewConfig()
			if err != nil {
				fmt.Printf("config error: %v\n", err)
				os.Exit(1)
			}

			dirEntries, err := os.ReadDir(cfg.ScriptsDir)
			if err != nil {
				fmt.Printf("failed to read scripts directory: %v\n", err)
				os.Exit(1)
			}

			registryPath := filepath.Join(cfg.ScriptsDir, "scripts.toml")
			if _, err := os.Stat(registryPath); os.IsNotExist(err) {
				fmt.Println("scripts.toml file does not exist")
				os.Exit(1)
			}

			registry, err := script.NewRegistry(registryPath)
			if err != nil {
				fmt.Printf("failed to load metadata: %v\n", err)
				os.Exit(1)
			}

			for _, dirEntry := range dirEntries {
				script.GetByFile(registry.Scripts, dirEntry.Name())
				fmt.Println(dirEntry.Name())
			}
		},
	}

	var invokeCmd = &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a script",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("script name is required")
				os.Exit(1)
			}

			cfg, err := config.NewConfig()
			if err != nil {
				fmt.Printf("config error: %v\n", err)
				os.Exit(1)
			}

			registryPath := filepath.Join(cfg.ScriptsDir, "scripts.toml")
			if _, err := os.Stat(registryPath); os.IsNotExist(err) {
				fmt.Println("scripts.toml file does not exist")
				os.Exit(1)
			}

			registry, err := script.NewRegistry(registryPath)
			if err != nil {
				fmt.Printf("failed to load metadata: %v\n", err)
				os.Exit(1)
			}

			scriptName := args[0]
			for _, s := range registry.Scripts {
				if s.Name == scriptName {
					err = invoker.InvokeScript(s, *cfg)
					if err != nil {
						fmt.Printf("failed to invoke script: %v\n", err)
						os.Exit(1)
					}
					return
				}
			}
		},
	}

	demonCmd.AddCommand(demonStartCmd, demonStatusCmd, demonStopCmd)
	rootCmd.AddCommand(demonCmd, configCmd, listCmd, invokeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
