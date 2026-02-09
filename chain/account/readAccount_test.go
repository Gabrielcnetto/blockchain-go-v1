package account

import (
	"fmt"
	"testing"
)

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
