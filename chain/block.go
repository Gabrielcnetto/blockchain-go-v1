package chain

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustinxie/ecc"
)

// tem o bloco, e tem o block Assinado(que contem o hash do anterior)
// tem a lista de Transactions validadas e assinadas
// arvore de markle
// time da criação
type Block struct {
	Number     uint64    `json:"number"`
	Parent     Hash      `json:"parent"` //hash do ultimo block
	Txs        []SigTx   `json:"txs"`
	merkleTree []Hash    //arvore inteira (em prod geralmente nao vai ter)
	MerkleRoot Hash      `json:"merkleRoot"`
	Time       time.Time `json:"time"`
}

func NewBlock(transactions []SigTx, numberBLock uint64, parent Hash) (Block, error) {
	merkleTree, err := MerkleHash(transactions, SigTxToHash, TxPairHash)
	if err != nil {
		return Block{}, err
	}
	return Block{
		merkleTree: merkleTree,
		Number:     numberBLock,
		Parent:     parent,
		Time:       time.Now(),
		Txs:        transactions,
		MerkleRoot: merkleTree[0],
	}, nil
}
func (b *Block) BlockToHash() Hash {
	return NewHash(b)
}

//agora, vamos ter o block assinado que é o bloco realmente atribuido a rede

type SigBlock struct {
	Block
	Sig []byte `json:"sig"`
}

func NewSigBlock(b Block, sig []byte) SigBlock {
	return SigBlock{
		Block: b,
		Sig:   sig,
	}
}

// ele tmb gera um hash
func (b *SigBlock) SigBlockToHash() Hash {
	return NewHash(b)
}

// funcao pega da internet pra printar o bloco
func (b SigBlock) String() string {
	var bld strings.Builder
	bld.WriteString(
		fmt.Sprintf(
			"blk %7d: %.7s -> %.7s   mrk %.7s\n",
			b.Number, b.SigBlockToHash(), b.Parent, b.MerkleRoot,
		),
	)
	for _, tx := range b.Txs {
		bld.WriteString(fmt.Sprintf("%v\n", tx))
	}
	return bld.String()
}

// assinando com ECDSA, padrão usado em todo o código (Criamos a private key com ela), quem assina é a conta
func (a *Account) SignBlock(blk Block) (SigBlock, error) {
	hash := blk.BlockToHash()
	hashBytes := hash[:] //deixando de array > slice Bytes
	signature, err := ecc.SignBytes(a.prv, hashBytes, ecc.LowerS|ecc.RecID)
	if err != nil {
		return SigBlock{}, err
	}
	return NewSigBlock(
		blk, signature,
	), nil
}

func VerifyBlock(blk SigBlock, authority Address) (bool, error) {
	HashBlock := blk.BlockToHash()
	hashBytes := HashBlock[:]
	publicKey, err := ecc.RecoverPubkey("P-256k1", hashBytes, blk.Sig)
	if err != nil {
		return false, err
	}
	acc := New_address(publicKey)
	return acc == authority, nil
}
