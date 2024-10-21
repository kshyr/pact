package invoker

import (
	"os"
	"os/exec"

	"github.com/kshyr/pact/config"
	"github.com/kshyr/pact/script"
)

func InvokeScript(s script.Script, cfg config.Config) error {
	scriptPath := cfg.ScriptsDir + "/" + s.File
	args := append([]string{scriptPath}, s.Args...)
	cmd := exec.Command("/bin/sh", args...)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
