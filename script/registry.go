package script

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type ScriptRegistry struct {
	Scripts []Script
}

type scriptRegistryTOML struct {
	Scripts []scriptTOML `toml:"scripts"`
}

func NewRegistry(registryPath string) (*ScriptRegistry, error) {
	var registryTOML scriptRegistryTOML
	if _, err := toml.DecodeFile(registryPath, &registryTOML); err != nil {
		return nil, fmt.Errorf("error loading script metadata: %v", err)
	}

	var registry ScriptRegistry
	for _, s := range registryTOML.Scripts {
		registry.Scripts = append(registry.Scripts, s.fromToml())
	}

	return &registry, nil
}
