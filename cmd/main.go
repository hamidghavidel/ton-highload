package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	w "ton-highload/pkg/wallet"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	// "github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const configPath string = "./testnet-global.config.json"

func main() {
	ctx := context.Background()
	rawMnemonic := "cancel fork visa lend trust skull bread spoon glimpse where pill beach party scene roof coast icon leaf frame knife extra polar twenty edit"
	mnemonic := strings.Split(rawMnemonic, " ")
	//client, err := w.NewClient(ctx)

	//block, err := client.CurrentMasterchainInfo(ctx) // get current block, we will need it in requests to LiteServer
	//if err != nil {
	//	log.Fatalln("CurrentMasterchainInfo err:", err.Error())
	//	return
	//}

	walletInfo, err := w.HighLoadV3(ctx, configPath, mnemonic, true)
	if err != nil {
		log.Fatalln("HighLoadV3 err:", err.Error())
		return
	}
	fmt.Println("walletInfo:", walletInfo)

	destinationAddress := address.MustParseAddr("0QDRGONvd2MFtX3X7qVUdeYtk31KLbuqV6IKuZ67CM5dePnq")
	// sourceAddress := address.MustParseAddr(walletInfo.Address)
	//internalMessageBody := cell.BeginCell().
	//	MustStoreUInt(0, 32).                // write 32 zero bits to indicate that a text comment will follow
	//	MustStoreStringSnake("Hello, TON!"). // write our text comment
	//	EndCell()
	//internalMessage := cell.BeginCell().
	//	MustStoreUInt(0, 1).     // indicate that it is an internal message -> int_msg_info$0
	//	MustStoreBoolBit(true).  // IHR Disabled
	//	MustStoreBoolBit(true).  // bounce
	//	MustStoreBoolBit(false). // bounced
	//	MustStoreUInt(0, 2).     // src -> addr_none
	//	MustStoreAddr(destinationAddress).
	//	MustStoreCoins(tlb.MustFromTON("0.2").Nano().Uint64()). // amount
	//	MustStoreBoolBit(false).                                // Extra currency
	//	MustStoreCoins(0).                                      // IHR Fee
	//	MustStoreCoins(0).                                      // Forwarding Fee
	//	MustStoreUInt(0, 64).                                   // Logical time of creation
	//	MustStoreUInt(0, 32).                                   // UNIX time of creation
	//	MustStoreBoolBit(false).                                // No State Init
	//	MustStoreBoolBit(true).                                 // We store Message Body as a reference
	//	MustStoreRef(internalMessageBody).                      // Store Message Body as a reference
	//	EndCell()

	//getMethodResult, err := client.RunGetMethod(ctx, block, destinationAddress, "seqno") // run "seqno" GET method from your wallet contract
	//if err != nil {
	//	log.Fatalln("RunGetMethod err:", err.Error())
	//	return
	//}
	//seqno := getMethodResult.MustInt(0)
	//mac := hmac.New(sha512.New, []byte(strings.Join(mnemonic, " ")))
	//hash := mac.Sum(nil)
	//
	//k := pbkdf2.Key(hash, []byte("TON default seed"), 100000, 32, sha512.New)
	//
	//privateKey := ed25519.NewKeyFromSeed(k)

	//toSign := cell.BeginCell().
	//	MustStoreUInt(uint64(walletInfo.SubWalletID), 32).     // subwallet_id | We consider this further
	//	MustStoreUInt(uint64(time.Now().UTC().Unix()+60), 32). // Message expiration time, +60 = 1 minute
	//	MustStoreUInt(seqno.Uint64(), 32).                     // store seqno
	//	MustStoreUInt(uint64(3), 8).                           // store mode of our internal message
	//	MustStoreRef(internalMessage)                          // store our internalMessage as a reference
	//
	//signature := ed25519.Sign(walletInfo.PrivateKey, toSign.EndCell().Hash()) // get the hash of our message to wallet smart contract and sign it to get signature
	//log.Println("Message signed successfully.")
	//body := cell.BeginCell().
	//	MustStoreSlice(signature, 512). // store signature
	//	MustStoreBuilder(toSign).       // store our message
	//	EndCell()
	//
	// externalMessage := cell.BeginCell().
	// 	MustStoreUInt(0b10, 2). // 0b10 -> 10 in binary
	// 	//MustStoreUInt(0, 2).               // src -> addr_none
	// 	MustStoreAddr(sourceAddress).      // Source address
	// 	MustStoreAddr(destinationAddress). // Destination address
	// 	MustStoreCoins(0).                 // Import Fee
	// 	MustStoreBoolBit(false).           // No State Init
	// 	MustStoreBoolBit(true).            // We store Message Body as a reference
	// 	MustStoreRef(body).                // Store Message Body as a reference
	// 	EndCell()
	// msg1 := tlb.Message{
	// 	MsgType: "INTERNAL",
	// 	Msg: ,
	// }
	
	payload := cell.BeginCell().MustStoreUInt(0, 32).MustStoreStringSnake("1234567").EndCell()
	msg1 := wallet.Message{
		Mode: wallet.PayGasSeparately + wallet.IgnoreErrors,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      true,
			DstAddr:     destinationAddress,
			Amount:      tlb.MustFromTON("0.01"),
			Body:        payload,
		},
	}

	// tlb.ExternalMessage{}

	// msg := wallet.SimpleMessage(destinationAddress, tlb.MustFromTON("0.01"), payload)

	ext, err := walletInfo.W.PrepareExternalMessageForMany(ctx, false, []*wallet.Message{&msg1})
	if err != nil {
		panic(err)
	}
	boc, err := tlb.ToCell(ext)
	if err != nil {
		panic(err)
	}
	tx := boc.ToBOCWithFlags(false)
	b := base64.StdEncoding.EncodeToString(tx)
	fmt.Println(b)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("sent tx")
	//log.Println("Prepared external message:", base64.StdEncoding.EncodeToString(externalMessage.ToBOCWithFlags(false)))
	//
	//var resp tl.Serializable
	//err = client.Client().QueryLiteserver(ctx, ton.SendMessage{Body: externalMessage.ToBOCWithFlags(false)}, &resp)
	//fmt.Println("response", resp)
	//if err != nil {
	//	log.Fatalln(err.Error())
	//	return
	//}
}
