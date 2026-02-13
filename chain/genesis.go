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
			genesis, _ := json.Marshal(g)
			return txn.Set([]byte("chain"), genesis)
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

func ReadGenesis() (SigGenesis, error) {
	db, err := clients.StartBadger()
	if err != nil {
		return SigGenesis{}, err
	}
	defer db.Close()
	var gen SigGenesis
	err = db.View(func(txn *badger.Txn) error {
		stateChain := VerifyChainExist(txn)
		if !stateChain {
			return fmt.Errorf("chain not created")
		}
		item, err := txn.Get([]byte("chain"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			fmt.Println("Item pego:", val)
			data := append([]byte{}, val...)
			if err := json.Unmarshal(data, &gen); err != nil {
				return nil
			}
			return nil
		})
		return err
	})
	return gen, err
}
