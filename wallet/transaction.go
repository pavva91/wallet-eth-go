package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ctx         = context.Background()
	url         = "https://sepolia.infura.io/v3/e22263da932241b4a652df190751bfc0"
	client, err = ethclient.DialContext(ctx, url)
)

func CurrentBlock() {
	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(block.Number())
}

func CheckBalance(addressString string) {
	address := common.HexToAddress(addressString)
	balanceWei, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		log.Print("There was an error", err)
	}
	fmt.Println("The balance in WEI at the current block number is", balanceWei)

	var convertionRate = big.NewInt(1000000000)
	var bigFloat = new(big.Float).SetInt(convertionRate)
	var balanceGwei = new(big.Float)
	balanceGwei.Quo(new(big.Float).SetInt(balanceWei), bigFloat)
	fmt.Println("The balance in GWEI at the current block number is", balanceGwei)
	fmt.Println("The balance in GWEI at the current block number is", balanceGwei.String())
	var balanceEth = new(big.Float)
	balanceEth.Quo(balanceGwei, bigFloat)
	fmt.Println("The balance in ETH at the current block number is", balanceEth)
	fmt.Println("The balance in ETH at the current block number is", balanceEth.String())
}

func QueryTransactions() {
	block, _ := client.BlockByNumber(ctx, nil)
	for _, transaction := range block.Transactions() {
		fmt.Println("-----------------------------------------------")
		fmt.Println("Transaction Value: ", transaction.Value().String())
		fmt.Println("Gas: ", transaction.Gas())
		fmt.Println("Gas Price: ", transaction.GasPrice().Uint64())
		fmt.Println("Nonce: ", transaction.Nonce())
		recipietAddressTx := transaction.To()
		if recipietAddressTx != nil {
			fmt.Println("To: ", transaction.To().Hex())
		}
	}
}

func GetTransactionsPerBlock() {
	block, err := client.BlockByNumber(ctx, nil)
	figure, err := client.TransactionCount(ctx, block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transactions per block: ", figure)
}

func CreateTransactionWrap() {
	// Recipient Info
	recipientAddress := common.HexToAddress("0x1Dabfed3934E7ab52ebeAbdd9153fe0003ee648f")

	// Sender Info
	senderPrivateKeyBytes, _ := hexutil.Decode("0x911510851e683067cbd3a0ed25d484dddfb829c537f272df53342098c3a04763")
	totalBudget := big.NewInt(592245011523813826)
	fraction := big.NewInt(100)
	fmt.Println("TotalBudget: ", totalBudget)

	amount := big.NewInt(0)
	amount.Div(totalBudget, fraction)

	CreateTransaction(senderPrivateKeyBytes, recipientAddress, *amount)

	// Retrieve Sender Address from string public key
	senderPublicKeyString := "0x04b86b83ab380657d59c999ff8d5d0647338f1971f6a31d43cc15445c7a6c5e584776206188793f102459efc722efb42db27d0abf0e37e30ff44f9aa9fb9e813eb"
	fmt.Println("Public Key direct String: ", senderPublicKeyString)
	senderPublicKeyBytes, _ := hexutil.Decode(senderPublicKeyString)
	senderPublicKeyECDSA, _ := crypto.UnmarshalPubkey(senderPublicKeyBytes)
	SenderAddress := crypto.PubkeyToAddress(*senderPublicKeyECDSA)
	fmt.Println(SenderAddress)

}

// recipient is wallet2.json
// sender is wallet1.json
func CreateTransaction(senderPrivateKeyBytes []byte, recipientAddressString common.Address, amount big.Int) {
	// Recipient Info
	RecipientAddress := common.HexToAddress("0x1Dabfed3934E7ab52ebeAbdd9153fe0003ee648f")

	// Sender Info
	// private key sender
	senderPrivateKeyBytes, err := hexutil.Decode("0x911510851e683067cbd3a0ed25d484dddfb829c537f272df53342098c3a04763")
	if err != nil {
		log.Println(err)
	}
	senderPrivateKey, err := crypto.ToECDSA(senderPrivateKeyBytes)
	if err != nil {
		log.Println(err)
	}

	// public key sender
	senderPublicKey := senderPrivateKey.Public()
	senderPublicKeyECDSA, ok := senderPublicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	senderPublicKeyBytes := crypto.FromECDSAPub(senderPublicKeyECDSA)
	senderPublicKeyString := hexutil.Encode(senderPublicKeyBytes)
	fmt.Println("Public Key Sender From Private Key:", senderPublicKeyString)

	// Sender Address
	SenderAddress := crypto.PubkeyToAddress(*senderPublicKeyECDSA)
	fmt.Println(SenderAddress)

	// You’ll need to declare variables for the amount of Ether you’re sending (in Gwei), the nonce (the number of transactions from the address), the gas price, the gas limit, and the ChainID.
	nonce, err := client.PendingNonceAt(ctx, SenderAddress)
	if err != nil {
		log.Println(err)
	}

	gasLimit := 3600
	gasLimit = 999999999
	gasLimit = 21000
	gas, err := client.SuggestGasPrice(ctx)

	if err != nil {
		log.Println(err)
	}

	ChainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Nonce : ", nonce)
	fmt.Println("RecipientAddress: ", RecipientAddress)
	fmt.Println("Amount Tx: ", amount)
	fmt.Println("Gas Limit: ", gasLimit)
	fmt.Println("Suggested gas: ", gas)
	fmt.Println("Chain ID: ", ChainID)

	// Create Tx
	transaction := types.NewTransaction(nonce, RecipientAddress, &amount, uint64(gasLimit), gas, nil)

	// Sign Tx
	signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(ChainID), senderPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Send Tx
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("transaction sent: %s\n", signedTx.Hash().Hex())
}
