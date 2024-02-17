package main

import (
	"fmt"
	"strconv"

	"github.com/vominhtrungpro/golang-blockchain/blockchain"
)

func main() {
	fbData := "Genesis"
	chain := blockchain.InitBlockChain(fbData)

	chain.AddBlock("First block after genesis")
	chain.AddBlock("Second block after genesis")
	chain.AddBlock("Third block after genesis")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Data in block: %x\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
