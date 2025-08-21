package state

import (
	"encoding/json"
	"state-sync/trie"
	"state-sync/utils"
)

type Snapshot struct {
	Root *trie.Node
}

func (sm *StateManager) CreateSnapshot(path string) error {
	snapshot := Snapshot{Root: sm.Trie.Root()}
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	return utils.WriteFile(path, data)
}

func (sm *StateManager) RestoreSnapshot(path string) error {
	data, err := utils.ReadFile(path)
	if err != nil {
		return err
	}
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}
	sm.Trie = trie.NewTrie()
	sm.Trie.SetRoot(snapshot.Root)
	return nil
}