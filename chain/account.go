package chain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/dustinxie/ecc"
)

type Account struct {
	prv  *ecdsa.PrivateKey
	addr Address
}

// gerei uma key usando as funcoes criadas, e o ecdsa base, é pego com generateKey
func NewAccount() (Account, error) {
	key, err := ecdsa.GenerateKey(ecc.P256k1(), rand.Reader)
	if err != nil {
		return Account{}, err
	}
	newAddr := New_address(&key.PublicKey)
	return Account{
		prv:  key,
		addr: newAddr,
	}, nil
}

func (a *Account) GetAddr() []byte {
	return []byte(a.addr)
}

func (a *Account) encodeKeyPair() ([]byte, error) {
	fmt.Println("a.prv.D.Bytes():", a.prv.D.Bytes())
	return a.prv.D.Bytes(), nil
}
