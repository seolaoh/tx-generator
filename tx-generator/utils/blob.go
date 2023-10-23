package utils

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
)

func BlobToVHash(commit kzg4844.Commitment) common.Hash {
	hasher := sha256.New()
	hasher.Write(commit[:])
	hash := hasher.Sum(nil)

	var vhash common.Hash
	vhash[0] = params.BlobTxHashVersion
	copy(vhash[1:], hash[1:])

	return vhash
}
