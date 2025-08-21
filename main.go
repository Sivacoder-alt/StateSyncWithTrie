package main

import (
	"fmt"
	"state-sync/state"
	"state-sync/trie"
	"time"
	"bytes"
)

func main() {
	t := trie.NewTrie()
	sm := state.NewStateManager(t)

	states := map[string]string{
		"0x1a": "1000",
		"0x1b": "2000",
		"0x2a": "3000",
	}
	fmt.Println("synchronizing initial states...")
	sm.SyncStates(states)
	rootHash := sm.Trie.RootHash()
	fmt.Printf("root Hash after sync: %x\n", rootHash)

	milestoneName := "initial-sync"
	err := sm.CreateMilestone(milestoneName, rootHash)
	if err != nil {
		fmt.Printf("failed to create milestone: %v\n", err)
	} else {
		fmt.Printf("created milestone: %s\n", milestoneName)
	}

	snapshotPath := "snapshot_initial.json"
	err = sm.CreateSnapshot(snapshotPath)
	if err != nil {
		fmt.Printf("failed to create snapshot: %v\n", err)
	} else {
		fmt.Printf("created snapshot: %s\n", snapshotPath)
	}

	checkpointPath := "checkpoint_initial.json"
	err = sm.CreateCheckpoint(checkpointPath)
	if err != nil {
		fmt.Printf("failed to create checkpoint: %v\n", err)
	} else {
		fmt.Printf("created checkpoint: %s\n", checkpointPath)
	}

	moreStates := map[string]string{
		"0x1a": "1500", 
		"0x3c": "4000", 
	}
	fmt.Println("\nsynchronizing additional states...")
	sm.SyncStates(moreStates)
	newRootHash := sm.Trie.RootHash()
	fmt.Printf("root Hash after update: %x\n", newRootHash)

	keys := [][]byte{[]byte("0x1a"), []byte("0x3c"), []byte("0x2a")}
	for _, key := range keys {
		fmt.Printf("\ngenerating proof for key %s...\n", key)
		proof, err := sm.GenerateProof(key)
		if err != nil {
			fmt.Printf("failed to generate proof for %s: %v\n", key, err)
		} else {
			value, valid := sm.VerifyProof(newRootHash, key, proof)
			if valid {
				fmt.Printf("verified value for key %s: %s\n", key, value)
			} else {
				fmt.Printf("p  roof verification failed for key %s\n", key)
			}
		}
	}

	fmt.Println("\nRolling back to checkpoint...")
	err = sm.RestoreCheckpoint(checkpointPath)
	if err != nil {
		fmt.Printf("Failed to restore checkpoint: %v\n", err)
	} else {
		restoredRoot := sm.Trie.RootHash()
		fmt.Printf("Root Hash after rollback: %x\n", restoredRoot)
		if bytes.Equal(restoredRoot, rootHash) {
			fmt.Println("Rollback successful: root hash matches initial state")
		}
	}

	milestone, err := sm.GetMilestone(milestoneName)
	if err != nil {
		fmt.Printf("Failed to get milestone: %v\n", err)
	} else {
		fmt.Printf("Milestone %s: Root Hash %x, Timestamp %s\n", milestoneName, milestone.RootHash, milestone.Timestamp.Format(time.RFC3339))
	}
}