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
