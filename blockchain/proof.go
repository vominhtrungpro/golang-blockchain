package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The first few bytes must contain 0s

// diff high make mining take longer
const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// taking a block and create a new hash to proof of work
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

// init data for calculating hash in proof of work
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0

	// infinite loop
	for nonce < math.MaxInt64 {
		// prepare data
		data := pow.InitData(nonce)
		// hash into sha256 format
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		// convert hash into big int
		intHash.SetBytes(hash[:])

		// compare Target with new big int
		if intHash.Cmp(pow.Target) == -1 {
			// mean hash is last than the target
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// check if the nonce on proof of work is valid
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// return slice byte from int 64
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
