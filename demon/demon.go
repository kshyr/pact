package demon

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"os/signal"

	"github.com/adrg/xdg"
	"github.com/kshyr/pact/logger"
)

const (
	logFileRelPath = "pact/logs/demon.log"
	pidFileRelPath = "pact/demon.pid"
)

type Demon struct {
	Log     *logger.Logger
	pidFile string
	Pid     int
}

func IsDemonExecutable() bool {
	return os.Getenv("DEMON") == "1"
}

func New() (*Demon, error) {
	logFile, err := logger.CreateLogFile("demon")
	if err != nil {
		return nil, err
	}

	log, err := logger.New(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	pidFile, err := PIDFile()

	d := &Demon{
		Log:     log,
		pidFile: pidFile,
	}

	pid, err := ReadPID()
	if err != nil {
		return nil, err
	}
	d.Pid = pid

	return d, nil
}

func (d Demon) Run(callback func()) error {
	if err := setProcessName(); err != nil {
		return d.Log.Errorf("failed to set process name: %v", err)
	}

	pid := os.Getpid()
	if err := os.WriteFile(d.pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return d.Log.Errorf("failed to write pid file: %v", err)
	}

	d.Log.Info("Demon started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case sig := <-sigs:
			d.Log.Infof("Received signal: %v. Shutting down.", sig)
			if err := RemovePIDFile(); err != nil {
				return err
			}
			os.Exit(0)
		default:
			callback()
		}
	}
}

func IsRunning() bool {
	pid, err := ReadPID()
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		RemovePIDFile()
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func ReadPID() (int, error) {
	pidFile, err := PIDFile()
	if err != nil {
		return 0, err
	}

	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, fmt.Errorf("failed to read pid file: %w", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, fmt.Errorf("invalid pid in pid file: %w", err)
	}

	return pid, nil
}

func (d Demon) Stop() error {
	process, err := os.FindProcess(d.Pid)
	if err != nil {
		return d.Log.Errorf("failed to find demon process: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return d.Log.Errorf("failed to terminate demon process: %w", err)
	}

	if err := RemovePIDFile(); err != nil {
		return err
	}

	return nil
}

func PIDFile() (string, error) {
	pidFilePath, err := xdg.CacheFile(pidFileRelPath)
	if err != nil {
		return "", fmt.Errorf("failed to create pid file: %w", err)
	}
	return pidFilePath, nil
}

func WritePIDFile(pid int) error {
	pidFile, err := PIDFile()
	if err != nil {
		return err
	}

	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("failed to write pid file: %w", err)
	}

	return nil
}

func RemovePIDFile() error {
	pidFile, err := PIDFile()
	if err != nil {
		return err
	}

	if err := os.Remove(pidFile); err != nil {
		return fmt.Errorf("failed to remove pid file: %w", err)
	}

	return nil
}
