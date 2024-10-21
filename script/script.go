package script

import (
	"fmt"
	"time"
)

const (
	iso8601 = "2006-01-02T15:04:05Z"
)

type Script struct {
	Name        string
	Description string
	File        string
	Args        []string
	Schedule    string
	At          time.Time
	Active      bool
	Tags        []string
}

type scriptTOML struct {
	Name        string   `toml:"name"`
	File        string   `toml:"file"`
	Args        []string `toml:"args"`
	Schedule    string   `toml:"schedule"`
	At          string   `toml:"at"`
	Description string   `toml:"description"`
	Active      bool     `toml:"active"`
	Tags        []string `toml:"tags,omitempty"`
}

func (s *scriptTOML) fromToml() Script {
	at, err := time.Parse(iso8601, s.At)
	if err != nil {
		at = time.Time{}
	}

	return Script{
		Name:        s.Name,
		File:        s.File,
		Args:        s.Args,
		Schedule:    s.Schedule,
		At:          at,
		Description: s.Description,
		Active:      s.Active,
		Tags:        s.Tags,
	}
}

func (s *Script) toToml() scriptTOML {
	at := s.At.Format(iso8601)
	return scriptTOML{
		Name:        s.Name,
		File:        s.File,
		Args:        s.Args,
		Schedule:    s.Schedule,
		At:          at,
		Description: s.Description,
		Active:      s.Active,
		Tags:        s.Tags,
	}
}

func (s *Script) String() string {
	return fmt.Sprintf("%s (%s). Active: %t. Schedule: %s", s.Name, s.Description, s.Active, s.Schedule)
}

func GetByFile(scripts []Script, file string) Script {
	for _, s := range scripts {
		if s.File == file {
			return s
		}
	}
	return Script{}
}
