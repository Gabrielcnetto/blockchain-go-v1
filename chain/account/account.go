package account

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/dustinxie/ecc"
)

type account struct {
	prv  *ecdsa.PrivateKey
	addr Address
}

// gerei uma key usando as funcoes criadas, e o ecdsa base, é pego com generateKey
func NewAccount() (account, error) {
	key, err := ecdsa.GenerateKey(ecc.P256k1(), rand.Reader)
	if err != nil {
		return account{}, err
	}
	newAddr := New_address(&key.PublicKey)
	return account{
		prv:  key,
		addr: newAddr,
	}, nil
}

func (a *account) GetAddr() []byte {
	return []byte(a.addr)
}

func (a *account) encodeKeyPair() ([]byte, error) {
	fmt.Println("a.prv.D.Bytes():", a.prv.D.Bytes())
	return a.prv.D.Bytes(), nil
}
