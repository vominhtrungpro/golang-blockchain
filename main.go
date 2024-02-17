package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type BlockChain struct {
	blocks []*Block
}

// block class
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

// create hash based on previous hash and data
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// func create actual block from data and previous hash
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

// func to add a block to the chain
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

// generate a genesis block with first block data
func Genesis(fbData string) *Block {
	return CreateBlock(fbData, []byte{})
}

// initialize block chain with first block data
func InitBlockChain(fbData string) *BlockChain {
	return &BlockChain{[]*Block{Genesis(fbData)}}
}

func main() {
	fbData := "Genesis"
	chain := InitBlockChain(fbData)
	chain.AddBlock("First block after genesis")
	chain.AddBlock("Second block after genesis")
	chain.AddBlock("Third block after genesis")

	for _, block := range chain.blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %x\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
