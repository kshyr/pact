package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/kshyr/pact/demon"
	"github.com/kshyr/pact/logger"
	"github.com/spf13/cobra"
)

func demonLoop(d *demon.Demon) {
	d.Log.Infof("Demon is alive: It's %d.", d.Pid)
	time.Sleep(time.Second)
}

func main() {
	if demon.IsDemonExecutable() {
		d, err := demon.New()
		if err != nil {
			fmt.Printf("Error initializing demon: %v\n", err)
			os.Exit(1)
		}

		err = d.Run(func() {
			demonLoop(d)
		})

		if err := demon.RemovePIDFile(); err != nil {
			d.Log.Error(err)
		}

		if err != nil {
			d.Log.Errorf("demon couldn't start: %v", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	var rootCmd = &cobra.Command{
		Use:   "pact",
		Short: "Pact CLI Application",
	}

	var demonCmd = &cobra.Command{
		Use:   "demon",
		Short: "Manage the demon process",
	}

	var startCmd = &cobra.Command{
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

	var statusCmd = &cobra.Command{
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

	var stopCmd = &cobra.Command{
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

	demonCmd.AddCommand(startCmd, statusCmd, stopCmd)
	rootCmd.AddCommand(demonCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
