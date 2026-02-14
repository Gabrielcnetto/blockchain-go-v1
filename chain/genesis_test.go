package chain

import (
	"fmt"
	"testing"
)

func Test_create_genesis(t *testing.T) {
	//acount pq precisaremos assinar
	pass := "Senha"
	key := []byte(from)
	account, err := ReadAccount([]byte(pass), key)
	if err != nil {
		t.Error("error to get chain account:", err.Error())
		return
	}
	not_sig_gen := NewGenesis(Address(from), 100, "chain")
	sigGen, err := account.SigGenesis(not_sig_gen)
	if err != nil {
		t.Error("Error to sign new genesis:", err.Error())
		return
	}
	state, err := VerifyGenesis(sigGen)
	if err != nil {
		t.Error("error to verify sign genesis:", err)
		return
	}
	if !state {
		t.Error("state false")
		return
	}
	//salvar chain
	stateSave := sigGen.SaveGenesis([]byte(pass))
	if stateSave != nil {
		t.Error("error to verify state save from blockchain:", stateSave.Error())
		return
	}
	fmt.Println("Genesis created with sucessfully!")

}

func Test_readGenesis(t *testing.T) {
	chain, err := ReadGenesis()
	if err != nil {
		t.Error("Error to read genesis:", err.Error())
	}
	fmt.Println("authority:", chain.Authority)
	fmt.Println("balances:", chain.Balances)
	fmt.Println("Sucess to get blockchain")
	fmt.Println("___________________-")
}
