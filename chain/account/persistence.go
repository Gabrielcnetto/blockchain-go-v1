package account

import (
	"app/clients"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/dgraph-io/badger/v4"
	"golang.org/x/crypto/argon2"
)

/*
O processo de persistência da conta
	1)Codifique o par de chaves da conta.
	2)Criptografe o par de chaves codificadas com a senha fornecida pelo proprietário.
	3)Escreva o par de chaves criptografadas em um arquivo com acesso restrito.
*/
//função gerente da persistencia
func (a *account) Persistence_account(password []byte) error {
	// inicia o DB
	DB, err := clients.StartBadger()
	if err != nil {
		return err
	}
	defer DB.Close()
	// pega o keypair em bytes
	keyPairEncoded, err := a.encodeKeyPair()
	if err != nil {
		return err
	}

	// criptografa com a senha (keypairBytes + Senha recebida do usuario)
	password_encript, err := CriptografarKeyPair_mais_senha(password, keyPairEncoded)
	if err != nil {
		return err
	}

	// salva no Badger: key = address, value = keypair criptografado
	err = DB.Update(func(txn *badger.Txn) error {
		return txn.Set(a.GetAddr(), password_encript)
	})

	return err
}

func CriptografarKeyPair_mais_senha(pass []byte, keypayrEncoded []byte) ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	//usando a senha vamos criar uma senha forte com base no salt
	key := argon2.IDKey(
		pass, salt, 1, 256, 1, 32,
	)
	//instancie o AES(padrao usado em criptografias) com a chave forte criada em argon
	aesInstance, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//agora com o aes, instancie o GCM
	gcmInstance, err := cipher.NewGCM(aesInstance)
	if err != nil {
		return nil, err
	}
	//gere um nounce com base no gcm nounce
	Nounce := make([]byte, gcmInstance.NonceSize())
	if _, err := rand.Read(Nounce); err != nil {
		return nil, err
	}
	//agora junte o nounce, a keypair e o salt em um pacotão unico
	data := gcmInstance.Seal(Nounce, Nounce, keypayrEncoded, nil)
	//junta com o salt iniciado la em cima
	data = append(salt, data...)
	return data, nil
}
