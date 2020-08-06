package goVsysSdk

import (
	"crypto/sha256"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/sha3"
	"strconv"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func HashChain(nonceSecret []byte) []byte {
	blake2bHash, err := blake2b.New256(nil)
	if err != nil {
		panic(err.Error())
	}
	blake2bHash.Write(nonceSecret)
	return Keccak256(blake2bHash.Sum(nil))
}

func BuildSeedHash(seed string, nonce int) []byte {
	nonceSeed := strconv.Itoa(nonce) + seed
	return HashChain([]byte(nonceSeed))
}

// GenerateKeyPair generate Account using seed byte array
func GenerateKeyPair(seed []byte) *Account {
	var originPublicKey = new([32]byte)
	originPrivateKey := sha256.Sum256([]byte(seed))
	curve25519.ScalarBaseMult(originPublicKey, &originPrivateKey)
	originPrivateKey[0] &= 248
	originPrivateKey[31] &= 127
	originPrivateKey[31] |= 64
	return &Account{publicKey: originPublicKey[:], privateKey: originPrivateKey[:]}
}
