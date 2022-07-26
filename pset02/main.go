package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// A hash is a sha256 hash, as in pset01
type Hash [32]byte

// ToString gives you a hex string of the hash
func (self Hash) ToString() string {
	return fmt.Sprintf("%x", self)
}

// Blocks are what make the chain in this pset; different than just a 32 byte array
// from last time.  Has a previous block hash, a name and a nonce.
type Block struct {
	PrevHash Hash
	Name     string
	Nonce    string
}

// ToString turns a block into an ascii string which can be sent over the
// network or printed to the screen.
func (self Block) ToString() string {
	return fmt.Sprintf("%x %s %s", self.PrevHash, self.Name, self.Nonce)
}

// Hash returns the sha256 hash of the block.  Hopefully starts with zeros!
func (self Block) Hash() Hash {
	return sha256.Sum256([]byte(self.ToString()))
}

// BlockFromString takes in a string and converts it to a block, if possible
func BlockFromString(s string) (Block, error) {
	var bl Block

	// check string length
	if len(s) < 66 || len(s) > 100 {
		return bl, fmt.Errorf("Invalid string length %d, expect 66 to 100", len(s))
	}
	// split into 3 substrings via spaces
	subStrings := strings.Split(s, " ")

	if len(subStrings) != 3 {
		return bl, fmt.Errorf("got %d elements, expect 3", len(subStrings))
	}

	hashbytes, err := hex.DecodeString(subStrings[0])
	if err != nil {
		return bl, err
	}
	if len(hashbytes) != 32 {
		return bl, fmt.Errorf("got %d byte hash, expect 32", len(hashbytes))
	}

	copy(bl.PrevHash[:], hashbytes)

	bl.Name = subStrings[1]

	// remove trailing newline if there; the blocks don't include newlines, but
	// when transmitted over TCP there's a newline to signal end of block
	bl.Nonce = strings.TrimSpace(subStrings[2])

	// TODO add more checks on name/nonce ...?

	return bl, nil
}

func main() {

	fmt.Printf("NameChain Miner v0.1\n")
	targetBits := uint8(20)
	THREAD_COUNT := 8
	// Your code here!

	// Basic idea:
	// Get tip from server, mine a block pointing to that tip,
	// then submit to server.
	// To reduce stales, poll the server every so often and update the
	// tip you're mining off of if it has changed.
	kill := make(chan bool)
	out := make(chan uint64)
	tip := Block{sha256.Sum256([]byte("")), "", ""}
	block := Block{sha256.Sum256([]byte("")), "", ""}
	first := true

	for {
		select {
		case <-time.After(5 * time.Second):
			// request tip from server
			new_tip, err := GetTipFromServer()
			if err != nil {
				fmt.Println(err)
			}
			if new_tip.Nonce != tip.Nonce {
				tip = new_tip
				// new block has been mined
				if !first {
					for i := 0; i < THREAD_COUNT; i++ {
						kill <- true
					}
				} else {
					first = false
				}
				block = Block{
					new_tip.Hash(),
					"Ben135892",
					"",
				}
				go block.Mine(targetBits, kill, out)
			}
		case nonce := <-out:
			// successfully mined a block, send to server
			for i := 0; i < THREAD_COUNT-1; i++ {
				kill <- true
			}
			first = true
			block.Nonce = fmt.Sprintf("%d", nonce)
			SendBlockToServer(block)
			fmt.Printf("Mined a new block!")
		}
	}

}
