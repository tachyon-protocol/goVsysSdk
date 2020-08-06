package goVsysSdk

import (
	"crypto/rand"
	"encoding/binary"
	"strings"
)

// IsValidatePhrase checks if phrase valid
func IsValidatePhrase(phrase string) bool {
	wordSet := map[string]bool{}
	for i := range wordList {
		wordSet[wordList[i]] = true
	}
	words := strings.Split(phrase, " ")
	for i := range words {
		_, ok := wordSet[words[i]]
		if !ok {
			return false
		}
	}
	return true
}

// IsValidateAddress checks if address valid
func IsValidateAddress(address string, network NetType) bool {
	data := MustBase58Decode(address)
	if len(data) != 26 {
		return false
	}
	if data[0] != addrVersion || data[1] != byte(network) {
		return false
	}
	key := data[0:22]
	check := data[22:26]
	keyHash := HashChain(key)[0:4]
	for i := 0; i < 4; i++ {
		if check[i] != keyHash[i] {
			return false
		}
	}
	return true
}

// PublicKeyToAddress return address with base58 encoded
func PublicKeyToAddress(publicKey string, network NetType) string {
	return publicKeyToAddress(MustBase58Decode(publicKey), network)
}

// MustGenerateSeed generates seed string
func MustGenerateSeed() string {
	var (
		words     string
		wordCount int64 = 2048
		w1        int64
		w2        int64
		w3        int64
		x         int64
	)
	r := make([]byte, 4)
	for i := 1; i <= 5; i++ {
		_, err := rand.Read(r)
		if err != nil {
			panic(err)
		}
		x = (int64(r[3]) & 0xff) + (int64(r[2])&0xff)<<8 + (int64(r[1])&0xff)<<16 + (int64(r[0])&0xff)<<24
		w1 = x % wordCount
		w2 = (((x / wordCount) >> 0) + w1) % wordCount
		w3 = ((((((x / wordCount) >> 0) + w1) % wordCount) >> 0) + w2) % wordCount
		words += wordList[w1] + " "
		words += wordList[w2] + " "
		words += wordList[w3] + " "
	}
	words = words[:len(words)-1]
	return words
}

func genRandomBytes(n int) []byte {
	retBytes:=make([]byte,n)
	nr,err:=rand.Read(retBytes)
	if err!=nil || nr!=n{
		panic(err)
	}
	return retBytes
}

func bytesToByteArrayWithSize(bytes []byte) (result []byte) {
	result = append(result, uint16ToByte(uint16(len(bytes)))...)
	result = append(result, bytes...)
	return
}

func uint64ToByte(data uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(data))
	return b
}

func uint32ToByte(data int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(data))
	return b
}

func uint16ToByte(data uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(data))
	return b
}
