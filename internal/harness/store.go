package harness

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Store struct {
	path string
}

func NewStore(path string) Store {
	return Store{path: path}
}

func (s Store) Load() (*State, error) {
	data, err := os.ReadFile(s.path)
	if err == nil {
		var state State
		if err := json.Unmarshal(data, &state); err != nil {
			return nil, err
		}
		return &state, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	return &State{}, nil
}

func (s Store) Save(state *State) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
