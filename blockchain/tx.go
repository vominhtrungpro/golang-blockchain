package blockchain

import (
	"bytes"

	"github.com/vominhtrungpro/golang-blockchain/wallet"
)

type TxOutput struct {
	Value      int
	PubkeyHash []byte
}

type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	Pubkey    []byte
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (in *TxInput) UsesKey(pubkeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.Pubkey)

	return bytes.Equal(lockingHash, pubkeyHash)
}

func (out *TxOutput) Lock(address []byte) {
	pubkeyHash := wallet.Base58Decode(address)
	pubkeyHash = pubkeyHash[1 : len(pubkeyHash)-4]
	out.PubkeyHash = pubkeyHash
}

func (out *TxOutput) IsLockedWithKey(pubkeyHash []byte) bool {
	return bytes.Equal(out.PubkeyHash, pubkeyHash)
}
