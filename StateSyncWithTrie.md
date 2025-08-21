# Verifying State Synchronization and Merkle Proofs in a Merkle Patricia Trie Implementation



## Overview of the Implementation



The program implements a Merkle Patricia Trie (MPT) with the following components:

- **Trie Package**: Defines the MPT structure with node types (Empty, Leaf, Extension, Branch) and operations (update, proof generation, root hash computation).

- **State Package**: Manages state synchronization, proof verification, snapshots (full trie state), milestones (metadata with root hash and timestamp), and checkpoints (rollback-capable states).

- **Utils Package**: Provides Keccak-256 hashing and file I/O for persistence.

- **Main Program**: Demonstrates state synchronization, proof verification, snapshot/checkpoint creation, and rollback using example data.



The MPT uses Keccak-256 for hashing (Ethereum-compatible) and converts keys to nibbles (4-bit units) for efficient prefix-based storage. Snapshots and checkpoints are stored as JSON files, and milestones track significant state updates.



## Example Data



The example data consists of two state synchronization operations:

1. **Initial State Sync**:

   - Key: `0x1a`, Value: `1000` (account balance)

   - Key: `0x1b`, Value: `2000`

   - Key: `0x2a`, Value: `3000`

2. **Example data - State Sync**:

   - Key: `0x1a`, Value: `1500` (update existing key)

   - Key: `0x3c`, Value: `4000` (new key)



This data simulates blockchain account states (e.g., addresses and balances). The keys are short (2 bytes in nibble form) for simplicity, and values are ASCII-encoded strings representing balances.





### Step 1: Run the Program

1. Execute the program: `go run main.go`.

2. The program will:

   - Synchronize the initial states.

   - Create a milestone (`initial-sync`), snapshot (`snapshot_initial.json`), and checkpoint (`checkpoint_initial.json`).

   - Synchronize additional states.

   - Generate and verify Merkle proofs for keys `0x1a`, `0x3c`, and `0x2a`.

   - Roll back to the checkpoint.

   - Retrieve the milestone.



### Step 2: Verify State Synchronization

- **Expected Output**:

  ```

  Synchronizing initial states...

  Root Hash after sync: bea0b55d3c043c99a834fb08677065215e20b66922c998ffc80719ccfbbbc758

  Created milestone: initial-sync

  Created snapshot: snapshot_initial.json

  Created checkpoint: checkpoint_initial.json

  ```

- **Verification**:

  - Check that the root hash matches `bea0b55d...` (exact hash may vary slightly due to serialization or system differences).

  - Verify that files `snapshot_initial.json`, `checkpoint_initial.json`, and `milestone_initial-sync.json` are created in the `state-sync/` directory.

  - Open `snapshot_initial.json` to confirm it contains a serialized trie structure with nodes corresponding to the keys `0x1a`, `0x1b`, and `0x2a`. Example (simplified):



    ```json

    {

      "Root": {

        "Type": 2,

        "Key": [3],

        "Child": {

          "Type": 3,

          "Children": [...]

        }

      }

    }



  ```

  Check milestone_initial-sync.json for the root hash and timestamp:



    {

      "Name": "initial-sync",

      "RootHash": "bea0b55d3c043c99a834fb08677065215e20b66922c998ffc80719ccfbbbc758",

      "Timestamp": "2025-08-21T17:28:00+05:30"

    }

  ```



### Step 3: Verify Additional State Synchronization

- **Expected Output**:

  ```

  Synchronizing additional states...

  Root Hash after update: 3570e15f75b4c8b31599a7461b21cde0b5bd292f4467d7493051f9e41b069358

  ```

- **Verification**:

  - Confirm the new root hash (`3570e15f...`) differs from the initial hash, indicating the trie updated correctly for `0x1a:1500` and `0x3c:4000`.

  - The trie now contains keys `0x1a:1500`, `0x1b:2000`, `0x2a:3000`, and `0x3c:4000`.



### Step 4: Verify Merkle Proofs

- **Expected Output**:

  ```

  Generating proof for key 0x1a...

  Prove: Added node type 2, serialized: <hex>

  Prove: Added node type 3, serialized: <hex>

  Prove: Added node type 2, serialized: <hex>

  Prove: Added node type 3, serialized: <hex>

  Prove: Added node type 1, serialized: <hex>

  VerifyProof: Processing node 0, type 2

  VerifyProof: Processing node 1, type 3

  VerifyProof: Processing node 2, type 2

  VerifyProof: Processing node 3, type 3

  VerifyProof: Found leaf with value 1500

  Verified value for key 0x1a: 1500



  Generating proof for key 0x3c...

  ...

  VerifyProof: Found leaf with value 4000

  Verified value for key 0x3c: 4000



  Generating proof for key 0x2a...

  ...

  VerifyProof: Found leaf with value 3000

  Verified value for key 0x2a: 3000

  ```

- **Verification**:

  - **Proof Generation**: The `Prove` logs show a sequence of nodes (extension, branch, leaf) forming the path to each key. For `0x1a`, the path includes multiple nodes due to shared prefixes (e.g., `0x1a` and `0x1b` share a nibble prefix `1`).

  - **Proof Verification**: The `VerifyProof` logs confirm that each node’s hash matches the expected child hash, and the final leaf node contains the correct value (`1500` for `0x1a`, `4000` for `0x3c`, `3000` for `0x2a`).

  - **Manual Check**:

    - Compute the Keccak-256 hash of the first proof node’s serialized data (from `Prove` logs) using a tool like `github.com/ethereum/go-ethereum/crypto` or an online Keccak-256 calculator.

    - Verify it matches the root hash (`3570e15f...`).

    - For each subsequent node in the proof, compute its Keccak-256 hash and confirm it matches the child hash in the parent node’s serialized data.

    - Check that the final leaf node’s key (in nibbles) matches the input key (e.g., `0x1a` → `[0, 1, 4, 10]`) and its value is correct.



### Step 5: Verify Checkpoint Rollback

- **Expected Output**:

  ```

  Rolling back to checkpoint...

  Root Hash after rollback: bea0b55d3c043c99a834fb08677065215e20b66922c998ffc80719ccfbbbc758

  Rollback successful: root hash matches initial state

  ```

- **Verification**:

  - Confirm the root hash after rollback matches the initial hash (`bea0b55d...`).

  - Re-run proof verification for `0x1a` after rollback to ensure it returns `1000` (initial value) instead of `1500`:

    ```go

    // Add to main.go after rollback

    proof, err := sm.GenerateProof([]byte("0x1a"))

    if err != nil {

        fmt.Printf("Failed to generate proof: %v\n", err)

    } else {

        value, valid := sm.VerifyProof(restoredRoot, []byte("0x1a"), proof)

        if valid {

            fmt.Printf("Post-rollback verified value for key 0x1a: %s\n", value)

        } else {

            fmt.Println("Post-rollback proof verification failed")

        }

    }

    ```

  - Expected output: `Post-rollback verified value for key 0x1a: 1000`.



### Step 6: Verify Milestone

- **Expected Output**:

  ```

  Milestone initial-sync: Root Hash bea0b55d3c043c99a834fb08677065215e20b66922c998ffc80719ccfbbbc758, Timestamp <timestamp>

  ```

- **Verification**:

  - Check `milestone_initial-sync.json` to confirm the root hash and timestamp match the output.

  - Verify the milestone’s root hash matches the initial state’s root hash.





