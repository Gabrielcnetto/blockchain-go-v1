package chain

import (
	"fmt"
	"testing"
)

var from = "9bdc139130ccf97ac1d19098c39fafae62c9f0de821ed31688933fef28a9d952"
var to = "881b4a56cb1932d520ef2ef88111d603471f9a30c61f1ae1e808e1756b27bdf9"

func Test_newTx(t *testing.T) {
	pass := "Senha"
	newTx := NewTx(Address(from), Address(to), 300, 1)

	newTx.Hash()
	key := []byte(from)
	Account, err := ReadAccount([]byte(pass), key)
	if err != nil {
		t.Error(err.Error())
		return
	}
	sigTx, err := Account.SignTx(newTx)
	if err != nil {
		t.Error(err.Error())
		return
	}
	actualBlock, err := ReadBlock()
	if err != nil {
		t.Error("error to get block:", err.Error())
		return
	}
	fmt.Println("Block:", actualBlock)
	newBlock, err := NewBlock([]SigTx{sigTx}, actualBlock.Number+1, actualBlock.BlockToHash())
	if err != nil {
		t.Error("error to create new block:", err.Error())
		return
	}
	sigBlock, err := Account.SignBlock(newBlock)
	if err != nil {
		t.Error("error to sign block:", err.Error())
		return
	}
	verifyStatus, err := VerifyBlock(sigBlock, Account.addr)
	if err != nil || !verifyStatus {
		t.Error("error to verify status with authority:", err.Error())
		return
	}
	err = SaveBlock(sigBlock)
	if err != nil {
		t.Error("Error to save new block:", err.Error())
	}
	fmt.Println("Sucesso ao criar sigTx:", sigTx)

}
