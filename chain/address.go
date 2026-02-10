package chain

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"

	"golang.org/x/crypto/sha3"
)

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
