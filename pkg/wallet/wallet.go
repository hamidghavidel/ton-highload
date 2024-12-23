package wallet

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"golang.org/x/crypto/ed25519"
	"time"
)

type Info struct {
	W           *wallet.Wallet
	Address     string
	Testnet     bool
	Balance     tlb.Coins
	SubWalletID uint32
	PrivateKey  ed25519.PrivateKey
}

func NewClient(ctx context.Context) (*ton.APIClient, error) {
	connection := liteclient.NewConnectionPool()
	configUrl := "https://ton-blockchain.github.io/testnet-global.config.json"
	err := connection.AddConnectionsFromConfigUrl(ctx, configUrl)
	if err != nil {
		return nil, err
	}
	client := ton.NewAPIClient(connection)
	return client, nil
}

func HighLoadV3(ctx context.Context, seed []string, testnet bool) (Info, error) {
	var i Info
	client, err := NewClient(ctx)
	if err != nil {
		return i, err
	}
	w, err := wallet.FromSeed(client, seed, wallet.ConfigHighloadV3{
		MessageTTL: 3000,
		MessageBuilder: func(ctx context.Context, subWalletId uint32) (id uint32, createdAt int64, err error) {
			createdAt = time.Now().Unix() - 30
			return uint32(createdAt % (1 << 23)), createdAt, nil
		},
	})
	fmt.Println("hi", w.WalletAddress().String())

	block, err := client.CurrentMasterchainInfo(ctx)
	if err != nil {
		return i, err
	}

	i.Address = w.WalletAddress().Testnet(testnet).String()
	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		return i, err
	}
	i.Balance = balance
	i.Testnet = testnet
	i.SubWalletID = w.GetSubwalletID()
	i.PrivateKey = w.PrivateKey()
	i.W = w
	return i, nil
}
