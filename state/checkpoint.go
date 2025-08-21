package state

func (sm *StateManager) CreateCheckpoint(path string) error {
	return sm.CreateSnapshot(path) 
}

func (sm *StateManager) RestoreCheckpoint(path string) error {
	return sm.RestoreSnapshot(path)
}
