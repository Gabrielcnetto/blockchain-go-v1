package chain

import "testing"

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
