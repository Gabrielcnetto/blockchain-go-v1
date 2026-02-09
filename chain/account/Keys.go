package account

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/dustinxie/ecc"
)

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
