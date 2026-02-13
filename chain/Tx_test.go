package chain

import (
	"fmt"
	"testing"
)

var from = "3591f0a62b229788301c315a8daefbac8cd69fd478e46f02b5703532a489ecb9"
var to = "9426f0773c4169f388e1680cd2d6af6ddd800df3bdae1e8c44118338301ae948"

func Test_newTx(t *testing.T) {
	pass := "Senha"
	newTx := NewTx(Address(from), Address(to), 100, 1)
	newTx.Hash()
	key := []byte(from)
	Account, err := ReadAccount([]byte(pass), key)
	if err != nil {
		t.Error(err.Error())
	}
	sigTx, err := Account.SignTx(newTx)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println("Sucesso ao criar sigTx:", sigTx)

}
