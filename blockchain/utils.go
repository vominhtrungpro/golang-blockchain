package blockchain

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
)

// return slice byte from int 64
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, num)
	Handle(err)

	return buff.Bytes()
}

// serialize block to byte
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

// deserialize byte to block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

// handle all error
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
