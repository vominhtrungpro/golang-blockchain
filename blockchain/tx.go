package blockchain

type TxOutput struct {
	Value  int
	Pubkey string
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlock(data string) bool {
	return out.Pubkey == data
}
