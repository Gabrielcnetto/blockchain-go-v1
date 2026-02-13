package chain

import "time"

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
