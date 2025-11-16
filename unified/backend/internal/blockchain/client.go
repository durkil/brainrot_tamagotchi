package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client represents a blockchain client
type Client struct {
	Eth        *ethclient.Client
	ChainID    *big.Int
	PrivateKey *ecdsa.PrivateKey

	// Contract addresses
	NFTAddress        common.Address
	CaseAddress       common.Address
	MarketplaceAddress common.Address
	BurnAddress       common.Address
}

// NewClient creates a new blockchain client
func NewClient() (*Client, error) {
	rpcURL := os.Getenv("BASE_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://sepolia.base.org" // Default to testnet
	}

	// Connect to blockchain
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Load private key (optional, for transactions)
	var privateKey *ecdsa.PrivateKey
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex != "" {
		privateKey, err = crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key: %w", err)
		}
	}

	// Load contract addresses
	nftAddr := common.HexToAddress(os.Getenv("CONTRACT_NFT_ADDRESS"))
	caseAddr := common.HexToAddress(os.Getenv("CONTRACT_CASE_ADDRESS"))
	marketAddr := common.HexToAddress(os.Getenv("CONTRACT_MARKETPLACE_ADDRESS"))
	burnAddr := common.HexToAddress(os.Getenv("CONTRACT_BURN_ADDRESS"))

	return &Client{
		Eth:                client,
		ChainID:            chainID,
		PrivateKey:         privateKey,
		NFTAddress:         nftAddr,
		CaseAddress:        caseAddr,
		MarketplaceAddress: marketAddr,
		BurnAddress:        burnAddr,
	}, nil
}

// GetTransactor creates a transaction auth object
func (c *Client) GetTransactor(ctx context.Context) (*bind.TransactOpts, error) {
	if c.PrivateKey == nil {
		return nil, fmt.Errorf("private key not set")
	}

	nonce, err := c.Eth.PendingNonceAt(ctx, crypto.PubkeyToAddress(c.PrivateKey.PublicKey))
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.Eth.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.PrivateKey, c.ChainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	return auth, nil
}

// WaitForTransaction waits for a transaction to be mined
func (c *Client) WaitForTransaction(ctx context.Context, txHash common.Hash) error {
	receipt, err := bind.WaitMined(ctx, c.Eth, &types.Transaction{})
	if err != nil {
		return err
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

