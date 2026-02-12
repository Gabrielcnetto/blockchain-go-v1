package chain

import (
	"fmt"
	"testing"
)

func Test_CreateAccount(t *testing.T) {
	password := "Senha"
	account, err := NewAccount()
	if err != nil {
		t.Error("error to create blockchainAccount:", err.Error())
	}
	//persistance on db
	err = account.Persistence_account([]byte(password))
	if err != nil {
		t.Error("error to persist account:", err.Error())
	}
	t.Log("Account created and persisted with sucessefully!")
	t.Log("Your address:", account.addr)
}

/*
endereço para testes:
9795d272d608c55aa4195e61c7cfc31197d85879676588bf12bf024e98ea6645
*/
func Test_ReadAccount(t *testing.T) {
	pass := "Senha"
	userAddr := "9795d272d608c55aa4195e61c7cfc31197d85879676588bf12bf024e98ea6645"
	key := []byte(userAddr)

	Account, err := ReadAccount([]byte(pass), key)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(Account.prv)
	fmt.Println("Got account with sucessfully")

}
