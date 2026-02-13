package chain

import "fmt"

//O algoritmo de hash Merkle constrói a árvore Merkle a partir da lista de transações
/*
Geramos um hash de cada transaction
geramos um par com a tx do lado, depois denovo, e denovo até ter um unico hash "pai" que representa todas as TX
tipo, imagna que temos

T1 | T2 | T3 | T4 (Transaction hash)
   \/        \/
   TA 		 TB (Pair hash)
        \/
		TC (main hash) (chamada de merkle root)
*/
//e no bloco, vai apenas o merkle root (referente a arvore)

func MerkleHash[Transaction any, Hash comparable](transactions []Transaction, typehash func(Transaction) Hash, pairHash func(Hash, Hash) Hash) ([]Hash, error) {
	if len(transactions) == 0 {
		return nil, fmt.Errorf("merkle hash: empty transaction list")
	}
	hashTransactions := make([]Hash, len(transactions))
	for index, transacion := range transactions {
		hashTransactions[index] = typehash(transacion)
	}
}
