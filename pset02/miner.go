package main

import "fmt"

// This file is for the mining code.
// Note that "targetBits" for this assignment, at least initially, is 33.
// This could change during the assignment duration!  I will post if it does.

// Mine mines a block by varying the nonce until the hash has targetBits 0s in
// the beginning.  Could take forever if targetBits is too high.
// Modifies a block in place by using a pointer receiver.
func (self *Block) Mine(targetBits uint8, kill <-chan bool, out chan<- uint64) {
	// your mining code here
	// also feel free to get rid of this method entirely if you want to
	// organize things a different way; this is just a suggestion
	nonce := uint64(0)
	THREAD_COUNT := 8
	for i := 0; i < THREAD_COUNT; i++ {
		go func(nonce uint64) {
			for {
				select {

				case <-kill:
					return

				default:
					self.Nonce = fmt.Sprintf("%d", nonce)
					if CheckWork(self, targetBits) {
						// mined a new block
						out <- nonce
						return
					}
					nonce++
				}
			}
		}(nonce)
		nonce += (1 << 30)
	}
}

// CheckWork checks if there's enough work
func CheckWork(bl *Block, targetBits uint8) bool {
	// your checkwork code here
	// feel free to inline this or do something else.  I just did it this way
	// so I'm giving empty functions here.
	hash := bl.Hash()
	for j := 0; uint8(j) < targetBits; j++ {
		bit := (hash[j/8] >> (7 - (j % 8))) & 0x01
		if bit == 1 {
			return false
		}
	}
	return true
}
