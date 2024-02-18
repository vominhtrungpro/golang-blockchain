package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

// handle all error
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)

	return decode
}
