package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type BlockChain struct {
	Blocks []*Block
}

// block class
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// create hash based on previous hash and data
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// func create actual block from data and previous hash
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	//block.DeriveHash()
	pow := NewProof(block)

	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce

	return block
}

// func to add a block to the chain
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

// generate a genesis block with first block data
func Genesis(fbData string) *Block {
	return CreateBlock(fbData, []byte{})
}

// initialize block chain with first block data
func InitBlockChain(fbData string) *BlockChain {
	return &BlockChain{[]*Block{Genesis(fbData)}}
}
