package chain

import (
	"app/clients"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/dustinxie/ecc"
)

type Genesis struct {
	Chain     string
	Authority Address
	Balances  map[Address]uint64
	Time      time.Time
}

func NewGenesis(addr Address, balance uint64, name string) Genesis {
	balances := make(map[Address]uint64, 1)
	balances[addr] = balance
	return Genesis{
		Chain:     name,
		Authority: addr,
		Balances:  balances,
		Time:      time.Now(),
	}
}

//gerando o hash do genesis

func (g *Genesis) Hash() Hash {
	return NewHash(g)
}

// para o genesis valer, precisa assinar ele, seguindo a ideia de TX criamos um sigGenesis
type SigGenesis struct {
	Genesis
	Sig []byte `json:"sig"`
}

func NewSigGenesis(gen Genesis, sig []byte) SigGenesis {
	return SigGenesis{
		Genesis: gen,
		Sig:     sig,
	}
}

// hash da sig genesis
func (g *SigGenesis) Hash() Hash {
	return NewHash(g)
}

// verificar a assinatura da genesis
func VerifyGenesis(gen SigGenesis) (bool, error) {
	hash := gen.Hash()
	pubKey, err := ecc.RecoverPubkey("P-256k1", hash[:], gen.Sig)
	if err != nil {
		return false, err
	}
	addr := New_address(pubKey)
	return addr == Address(gen.Authority), nil
}

// persistencia da blockchain com badgers
// verificando se a blockchain ja foi criada, antes de setar coisa nova

func VerifyChainExist(txn *badger.Txn) bool {
	_, err := txn.Get([]byte("chain"))
	if err == badger.ErrKeyNotFound {
		fmt.Println("Chain não encontrada:", err.Error())
		return false
	} else if err != nil {
		fmt.Println("erro ao ler a chain, retorno false")
		return false
	}
	return true

}

func (g *SigGenesis) SaveGenesis() error {
	db, err := clients.StartBadger()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		stateChain := VerifyChainExist(txn)
		if !stateChain {
			//gerar com a genesis
			hashed_chain := g.Hash()
			return txn.Set([]byte("chain"), hashed_chain[:])
		} else {
			return fmt.Errorf("Genesis ja definido")
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func ChainByte(sig *SigGenesis) []byte {
	sigByte, _ := json.Marshal(sig)
	return sigByte
}
