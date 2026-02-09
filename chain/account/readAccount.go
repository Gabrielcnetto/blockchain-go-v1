package account

import (
	"app/clients"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"math/big"

	"github.com/dgraph-io/badger/v4"
	"github.com/dustinxie/ecc"
	"golang.org/x/crypto/argon2"
)

func ReadAccount(password []byte, address []byte) (account, error) {
	//precisamos receber o endereço, pra conseguir gerar
	AccountFromDb, err := GetAcountFromDb(address)
	if err != nil {
		return account{}, err
	}
	bitPrivateKey, err := decryptWithPassword(AccountFromDb, password)
	if err != nil {
		return account{}, err
	}
	privKey := reacreatePrivate(bitPrivateKey)
	return account{
		prv:  privKey,
		addr: Address(address),
	}, nil
}

func GetAcountFromDb(addr []byte) ([]byte, error) {
	db, err := clients.StartBadger()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	AccountEncriptedData := []byte{}
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(addr))
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			AccountEncriptedData = append([]byte{}, val...)
			return nil
		})
		return nil
	})
	return AccountEncriptedData, err

}

func decryptWithPassword(data, password []byte) ([]byte, error) {
	//como o salt é de 32 bits, eu puxo os primerios 32 bits que foram colocados no valor total
	salt := data[:32]

	encritpPass := data[32:]
	pass := argon2.IDKey(password, salt, 1, 256, 1, 32)
	aesInstance, err := aes.NewCipher(pass)
	if err != nil {
		return nil, err
	}
	gcmInstance, err := cipher.NewGCM(aesInstance)
	if err != nil {
		return nil, err
	}
	nonceSize := gcmInstance.NonceSize()
	nounce := encritpPass[:nonceSize]
	keypairEncoded := encritpPass[nonceSize:]
	byteKeyPair, err := gcmInstance.Open(
		nil, nounce, keypairEncoded, nil,
	)
	if err != nil {
		return nil, err
	}
	return byteKeyPair, nil
}

func reacreatePrivate(data []byte) *ecdsa.PrivateKey {
	d := new(big.Int).SetBytes(data)
	privKey := new(ecdsa.PrivateKey)
	privKey.D = d
	privKey.PublicKey.Curve = ecc.P256k1()
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(d.Bytes())
	return privKey
}
