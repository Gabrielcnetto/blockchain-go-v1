package chain

import (
	"app/clients"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/dgraph-io/badger/v4"
	"github.com/dustinxie/ecc"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
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

//O endereço é um hash derivado da chave pública,
// 	que identifica o usuário na blockchain e é usado para receber ou enviar ativos.
/*
Private key (D)  ---> Public key ---> Address ---> Usado em transações / saldo
*/

// O address é derivado da public key, normalmente através de um hash dela.
// precisamos gerar um endereço para as contas poderem se conversarem, nao imporat como sera gerado
type Address string

func New_address(pub *ecdsa.PublicKey) Address {
	pub_bytes, _ := json.Marshal(New_publickKey(pub))
	hash := make([]byte, 32) //um slice de 32 bytes
	sha3.ShakeSum256(hash, pub_bytes)
	return Address(hex.EncodeToString(hash))
}

//a segurança é baseada e gerada as keys em curva eliptica P-256k1

type P256k1PublickKey struct {
	Curve string `json:"curve"`
	Y     *big.Int
	X     *big.Int
}

func New_publickKey(pub *ecdsa.PublicKey) P256k1PublickKey {
	return P256k1PublickKey{
		Curve: "P-256k1",
		Y:     pub.Y,
		X:     pub.X,
	}
}

type P256k1PrivateKey struct {
	P256k1PublickKey
	D *big.Int //numero aleatorio para gerarmos a private e public key
}

func new_P256k1PrivateKey(priv *ecdsa.PrivateKey) P256k1PrivateKey {
	return P256k1PrivateKey{
		P256k1PublickKey: New_publickKey(&priv.PublicKey),
		D:                priv.D,
	}
}

func (p *P256k1PrivateKey) NewPublicKey() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: ecc.P256k1(),
		X:     p.X,
		Y:     p.X,
	}
}

func (p *P256k1PrivateKey) NewPrivateKey() *ecdsa.PrivateKey {
	return &ecdsa.PrivateKey{
		PublicKey: *p.NewPublicKey(),
		D:         p.D,
	}
}

/*
FLUXO
Geração:
1️⃣ Geramos D aleatório (private key) → fica só na wallet do usuário.
2️⃣ Derivamos Public Key de D.
3️⃣ Derivamos Address da Public Key.
4️⃣ Entregamos só o Address pro usuário.

MEU PAPEL:
“Quando a conta é criada, o sistema gera o address e entrega esse address ao usuário.
Esse address é o que o usuário precisa guardar para receber fundos e identificar a conta.”
*/

func ReadAccount(password []byte, address []byte) (Account, error) {
	//precisamos receber o endereço, pra conseguir gerar
	AccountFromDb, err := GetAcountFromDb(address)
	if err != nil {
		return Account{}, err
	}
	bitPrivateKey, err := decryptWithPassword(AccountFromDb, password)
	if err != nil {
		return Account{}, err
	}
	privKey := reacreatePrivate(bitPrivateKey)
	return Account{
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

/*
O processo de persistência da conta
	1)Codifique o par de chaves da conta.
	2)Criptografe o par de chaves codificadas com a senha fornecida pelo proprietário.
	3)Escreva o par de chaves criptografadas em um arquivo com acesso restrito.
*/
//função gerente da persistencia
func (a *Account) Persistence_account(password []byte) error {
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

//assinando o bloco genesis

func (a *Account) SigGenesis(gen Genesis) (SigGenesis, error) {
	hash := gen.Hash()
	sig, err := ecc.SignBytes(a.prv, hash[:], ecc.LowerS|ecc.RecID)
	if err != nil {
		return SigGenesis{}, err
	}
	sigGen := NewSigGenesis(gen, sig)
	return sigGen, nil
}
