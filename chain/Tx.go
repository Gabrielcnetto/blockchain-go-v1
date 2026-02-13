package chain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dustinxie/ecc"
	"golang.org/x/crypto/sha3"
)

type Hash [32]byte

func NewHash(val any) Hash {
	byte_val, _ := json.Marshal(val)
	//instancia a ferramenta que cria hash
	state := sha3.NewLegacyKeccak256()
	//adciona o hash no byte gerado com marshal
	_, _ = state.Write(byte_val)
	//pra salvar o write faz o sum
	state.Sum(nil)
	return Hash(byte_val)
}

// retornar o hexadecimal do hash gerado
func (h *Hash) Hash_to_hex_string() string {
	return hex.EncodeToString(h[:])
}

// retornar o hash em byte puro
// *Lembrando que quando cria [32] ele fica com array, ai o : converte pra slice
func (h *Hash) Hash_to_bytes() []byte {
	return h[:]
}

//funcao antes de ser assinada

//aqui vamos criar esse item, antes de fazer qualquer transacao, criamos um tx, assinamos > enviamos

type Tx struct {
	From  Address
	To    Address
	Value int64
	Nonce int64
	Time  time.Time
}

func NewTx(from, to Address, value, nonce int64) Tx {
	return Tx{
		From:  from,
		To:    to,
		Value: value,
		Nonce: nonce,
		Time:  time.Now(),
	}
}

// cada transação gera um hash da transacao criada
func (t *Tx) Hash() Hash {
	return NewHash(t)
}

type SigTx struct {
	Tx
	Sig []byte `json:"sig"`
}

func NewSig(tx Tx, sig []byte) SigTx {
	return SigTx{
		Tx:  tx,
		Sig: sig,
	}
}

func SigTxToHash(tx SigTx) Hash {
	return NewHash(tx)
}

// a assinatura assinada tmb gera um hash dela
func (s *SigTx) Hash() Hash {
	return NewHash(s)
}

func (t SigTx) String() string {
	return fmt.Sprintf(
		"tx %.7s: %.7s -> %.7s %8d %8d", t.Hash(), t.From, t.To, t.Value, t.Nonce,
	)
}

//vamos gerar uma funcao, que recebe uma sigTx e gera um hash dela(precisaremos no futuro)
//objetivo: Receber via construtor e nao via dependecia, e retornar a sigTX em Hash

func NewSigHash(sigTx SigTx) Hash {
	return NewHash(sigTx)
}

// junta em nós, para a arvore de markle facilitar a busca nos blocos
/*
        Root
       /    \
     AB     CD    ← Você recebe CD (pai do galho)
    /  \   /  \
   A    B C    D  ← Você tem C, recebe D (folha irmã)
         ↑
       Sua TX
*/
func TxPairHash(left, right Hash) Hash {
	var nilHash Hash
	if right == nilHash {
		//retorna apenas 1 para a arvore
		return left
	}
	//se nao junta os 2
	return NewHash(left.Hash_to_hex_string() + right.Hash_to_hex_string())
}

// verificar a assinatura, engenharia reversa do processo feito para assinar
// verificar se assinatura é referente a quem ta tentando enviar
// isso comprova que o "dono" mesmo que ta tentando enviar
// *O Ecc.RecoveryPubKey retorna a chave publica > que gerará o endereço do dono
func VerifyTx(sigTx SigTx) (bool, error) {
	hash := sigTx.Hash()
	data, err := ecc.RecoverPubkey("P-256k1", hash[:], sigTx.Sig)
	if err != nil {
		return false, err
	}
	accountAddr := New_address(data) //1 pub sempre gera a mesma saída
	return accountAddr == sigTx.From, nil
}
