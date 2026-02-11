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
5ba289b463a8638e1c5ae72f8d7ff309ea5bf8c0dcb4b6f93041026a46ed94f3
*/
func Test_ReadAccount(t *testing.T) {
	pass := "Senha"
	userAddr := "5ba289b463a8638e1c5ae72f8d7ff309ea5bf8c0dcb4b6f93041026a46ed94f3"
	key := []byte(userAddr)

	Account, err := ReadAccount([]byte(pass), key)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(Account.prv)
	fmt.Println("Got account with sucessfully")

}
