package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"github.com/vominhtrungpro/golang-blockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	clearDatabaseCmd := flag.NewFlagSet("clear", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "clear":
		err := clearDB(cli.blockchain.Database)
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if clearDatabaseCmd.Parsed() {
		cli.printChain()
	}
}

func clearDB(db *badger.DB) error {
	err := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			err := txn.Delete(it.Item().Key())
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()
}
