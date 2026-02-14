package chain

import (
	"fmt"
	"testing"
)

func Test_NewBlock(t *testing.T) {
	Senha := "Senha"
	SigBlock, err := Coinbase(Address(from), []byte(Senha))
	if err != nil {
		t.Error("Error to create coinbase:", err.Error())
	}
	err = SaveBlock(*SigBlock)
	if err != nil {
		t.Error("Error to save block:", err.Error())
	}
	fmt.Println("Sucess to create and save block")
}

func Test_ReadBSIglock(t *testing.T) {
	block, err := ReadBlock()
	if err != nil {
		t.Error("error to get block:", err.Error())
	}
	fmt.Println("Block get with sucessfully")
	fmt.Println("Tree:", block.String())
	for _, item := range block.Txs {
		fmt.Println("Item tx:", item)
	}
}
