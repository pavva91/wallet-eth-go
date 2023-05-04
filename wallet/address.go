package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Wallets are composed of three main components; the public key, the private key, and the public address.
func CreateWallet() (thePublicAddress string, thePublicKey string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("SAVE BUT DO NOT SHARE THIS (Private Key):", hexutil.Encode(privateKeyBytes))

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyString := hexutil.Encode(publicKeyBytes)
	fmt.Println("Public Key:", publicKeyString)
	backToBytes, err := hexutil.Decode(publicKeyString)
	fmt.Println("Public Key:", hexutil.Encode(backToBytes))
	backToECDSA, _ := crypto.UnmarshalPubkey(backToBytes)
	fmt.Println("Public Key:", hexutil.Encode(crypto.FromECDSAPub(backToECDSA)))

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address 1:", address)

	return address, publicKeyString
}
