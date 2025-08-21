package state

import (
	"encoding/json"
	"time"
	"state-sync/utils"
)

type Milestone struct {
	Name      string
	RootHash  []byte
	Timestamp time.Time
}

func (sm *StateManager) CreateMilestone(name string, rootHash []byte) error {
	milestone := Milestone{
		Name:      name,
		RootHash:  rootHash,
		Timestamp: time.Now(),
	}
	data, err := json.MarshalIndent(milestone, "", "  ")
	if err != nil {
		return err
	}
	return utils.WriteFile("milestone_"+name+".json", data)
}

func (sm *StateManager) GetMilestone(name string) (*Milestone, error) {
	data, err := utils.ReadFile("milestone_" + name + ".json")
	if err != nil {
		return nil, err
	}
	var milestone Milestone
	if err := json.Unmarshal(data, &milestone); err != nil {
		return nil, err
	}
	return &milestone, nil
}