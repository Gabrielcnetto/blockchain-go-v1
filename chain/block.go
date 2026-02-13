package chain

import (
	"app/clients"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
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
func Coinbase(addr Address, pass []byte) (*SigBlock, error) {
	account, err := ReadAccount(pass, []byte(addr))
	if err != nil {
		return nil, err
	}
	tx := Tx{
		From:  addr,
		To:    addr,
		Value: 1500,
		Nonce: 1,
		Time:  time.Now(),
	}
	sigTx, err := account.SignTx(tx)
	if err != nil {
		return nil, err
	}
	block, err := NewBlock([]SigTx{sigTx}, 1, NewHash(nil))
	if err != nil {
		return nil, err
	}
	NewSiggBlock, err := account.SignBlock(block)
	if err != nil {
		return nil, err
	}
	return &NewSiggBlock, nil

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

func GetBlock(txn *badger.Txn) (*SigBlock, int, error) {
	fmt.Println("Entrei auqi")
	item, err := txn.Get([]byte("block"))
	if err == badger.ErrKeyNotFound {
		return nil, 0, err
	}
	var SigBlock SigBlock
	err = item.Value(func(val []byte) error {
		savedByte := append([]byte{}, val...)
		if err := json.Unmarshal(savedByte, &SigBlock); err != nil {
			return err
		}
		return nil
	})
	return &SigBlock, 1, err
}

func SaveBlock(sigBlock SigBlock) error {
	db, err := clients.StartBadger()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		parentBlock, state, err := GetBlock(txn)
		switch state {
		case 0:
			byteBlock, _ := json.Marshal(sigBlock)
			return txn.Set([]byte("block"), byteBlock)
		case 1:
			sigBlock.Parent = parentBlock.BlockToHash()
			byteBlock, _ := json.Marshal(sigBlock)
			err := txn.Set([]byte("block"), byteBlock)
			return err
		default:
			return err
		}

	})
	return err
}
