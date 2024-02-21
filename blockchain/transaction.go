package blockchain

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/vominhtrungpro/golang-blockchain/wallet"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	Handle(err)

	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txIn := TxInput{[]byte{}, -1, nil, []byte(data)}
	txOut := NewTXOutput(100, to)

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{*txOut}}
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallet.CreateWallets()
	Handle(err)
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validOutputs := chain.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, to))

	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	chain.SignTransaction(&tx, w.PrivateKey)

	return &tx
}

func (tx *Transaction) Sign(privKey rsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		prevTXs := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].Pubkey = prevTXs.Outputs[in.Out].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].Pubkey = nil

		hashed := sha256.Sum256(txCopy.ID)
		signature, err := rsa.SignPKCS1v15(rand.Reader, &privKey, crypto.SHA256, hashed[:])
		Handle(err)

		tx.Inputs[inId].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var input []TxInput
	var output []TxOutput

	for _, in := range tx.Inputs {
		input = append(input, TxInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		output = append(output, TxOutput{out.Value, out.PubkeyHash})
	}

	txCopy := Transaction{tx.ID, input, output}

	return txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Previous transaction does not exist")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].Pubkey = prevTx.Outputs[in.Out].PubkeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].Pubkey = nil

		hashed := sha256.Sum256(txCopy.ID)

		rawPubKey := rsa.PublicKey{
			N: new(big.Int).SetBytes(in.Pubkey),
			E: 65537,
		}

		err := rsa.VerifyPKCS1v15(&rawPubKey, crypto.SHA256, hashed[:], in.Signature)
		Handle(err)
	}

	return true
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Out))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.Pubkey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubkeyHash))
	}

	return strings.Join(lines, "\n")
}
