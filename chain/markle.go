package chain

import (
	"fmt"
	"math"
)

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
	//calcular o tamanho total da arvore, somando folha(TX), o nó das folhas(pairs) até o merkle root
	/*
		1. 6 transações
		2. Potência de 2 que cobre 6? 2³ = 8 (precisa de 3 níveis de folhas)
		3. +1 nível para raiz = 4 níveis totais
		4. 2⁴ = 16 posições totais
		5. -1 = 15 nós na árvore
		6. ✅
	*/

	/*
			Nível 3 (raiz):                    [0]                  ← 1 nó (2⁰)
		                              /              \
		Nível 2:                   [1]              [2]             ← 2 nós (2¹)
		                         /     \           /     \
		Nível 1:               [3]     [4]       [5]     [6]        ← 4 nós (2²)
		                      /  \    /  \      /  \    /  \
		Nível 0 (folhas):   [7][8] [9][10] [11][12] [13][14]        ← 8 nós (2³)
	*/
	levelEnd := int(math.Pow(2, math.Ceil(math.Log2(float64(len(hashTransactions))))+1) - 1)
	merkleTree := make([]Hash, levelEnd) // o tamanho nunca muda
	PrimeiraCamada := levelEnd / 2       //começamos pela metade do lengh (isso acessa o nó com as folhas, antes do pairs)
	for i, folhaIndex := 0, PrimeiraCamada; i < len(transactions); i, folhaIndex = i+1, folhaIndex+1 {
		merkleTree[folhaIndex] = hashTransactions[i]
	}

	for index, folhaIndex := 0, PrimeiraCamada; index < len(transactions); index, folhaIndex = index+1, folhaIndex+1 {
		merkleTree[folhaIndex] = hashTransactions[index] //pega o hash, e salva na merkletree
	}
	//✔ Cada nível tem o dobro de nós do nível acima
	levelEnd, par := PrimeiraCamada*2, PrimeiraCamada/2 //"subindo uma camada"
	//o *2 significa que a camada pai que vai vir, tem o dobro de filhos

	for PrimeiraCamada > 0 { //se o end ta maior que 0 quer dizer que ainda tem camada pra subir
		//isso aqui vai formar a camada dos pairs e ir "fechando" até o root
		for index, pairIndex := PrimeiraCamada, par; index < levelEnd; index, pairIndex = index+2, pairIndex+1 {
			merkleTree[pairIndex] = pairHash(merkleTree[index], merkleTree[index+1])
		}
		PrimeiraCamada = PrimeiraCamada / 2 //"pula para o proximo nivel"
		levelEnd, par = PrimeiraCamada*2, PrimeiraCamada
	}
	return merkleTree, nil
}

/*
Você tem um número: 16 (2⁴)

1️⃣ COMEÇA:
   PrimeiraCamada = 16/2 = 8  (metade)
   (aqui você coloca as folhas)

2️⃣ SOBE UM NÍVEL:
   PrimeiraCamada = 8 * 2 = 16  (dobra)
   (aqui você coloca os pais)

3️⃣ SOBE MAIS UM:
   PrimeiraCamada = 16 * 2 = 32
   (aqui os avós)

4️⃣ ATÉ CHEGAR NA RAIZ
*/
