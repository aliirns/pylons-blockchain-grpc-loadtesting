package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"google.golang.org/grpc"

	pylonstypes "github.com/Pylons-tech/pylons/x/pylons/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/goombaio/namegenerator"
)

const (
	recipeProductPrefix = "recipe"
	tradeProductPrefix  = "trade"
	productSeparator    = "/"
)

func getTrade(grpcConn *grpc.ClientConn, tradeID uint64) (pylonstypes.Trade, error) {
	pylonsClient := pylonstypes.NewQueryClient(grpcConn)
	tradeRes, err := pylonsClient.Trade(context.Background(), &pylonstypes.QueryGetTradeRequest{Id: tradeID})
	if err != nil {
		return pylonstypes.Trade{}, err
	}
	return tradeRes.GetTrade(), nil
}

func dialGrpc() (*grpc.ClientConn, error) {
	grpcConn, err := grpc.Dial(
		"192.168.88.86:9090",
		grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
	)

	if err != nil {
		log.Printf("grpc.Dial: %v", err)
		return &grpc.ClientConn{}, err
	}

	return grpcConn, nil
}

func getBalance(address string, denom string) (sdk.Int, error) {
	grpcConn, err := dialGrpc()
	if err != nil {
		return sdk.Int{}, err
	}

	defer grpcConn.Close()

	// This creates a gRPC client to query the x/bank service.
	bankClient := banktypes.NewQueryClient(grpcConn)
	bankRes, err := bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: address, Denom: denom})
	if err != nil {
		return sdk.Int{}, err
	}
	return bankRes.Balance.Amount, nil
}

func createAccount() (bool, error) {
	grpcConn, err := dialGrpc()
	if err != nil {
		return false, err
	}

	defer grpcConn.Close()

	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	name := nameGenerator.Generate()

	cmd := exec.Command("pylonsd", "keys", "add", name)
	_, err = cmd.Output()
	if err != nil {
		return false, err
	}

	// pylonsClient := pylonstypes.NewMsgClient(grpcConn)
	referral := ""

	// creatorAddress := name
	//username := name
	token := "eyJraWQiOiJsWUJXVmciLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxOjYyODM2NTMzODM4MzphbmRyb2lkOmNiYTIxZjViY2I3ZTA5MzRkMDQwN2EiLCJhdWQiOlsicHJvamVjdHNcLzYyODM2NTMzODM4MyIsInByb2plY3RzXC9weWxvbnMtYjkyM2MiXSwiaXNzIjoiaHR0cHM6XC9cL2ZpcmViYXNlYXBwY2hlY2suZ29vZ2xlYXBpcy5jb21cLzYyODM2NTMzODM4MyIsImV4cCI6MTY1ODkyOTg5MiwiaWF0IjoxNjU4OTI2MjkyfQ.loB6bqJh-8H4FFuFyfWmAWmbNkzcS3Ce4goCdKIhuDzNf_mUnqrJizTTSyvGphG3GpbOPeqPUCxycJ2Xrg3t72f71fAUsQYeukpBFU2B38V1rUngHgIBnz9QTn26xl0-yGGqf7rjFZUbNjMdIXJ7zuGV4ZW0-z1jVFmaFLi7TH_XTxZ4JwmKkBSaDIM5P43e8Wll7iyckWNfPuyxLm81lQEsZGmyIJ-HzRNmle5ejqzvis_7rRyR97DcIxQoEfjPvkAXLOv48VhrU_U2_yTzK3vXJ_ZOGWVM0aUuPalquRnHiSKmgceQ2tdSQKAFlXnpAuEIlFxQr3x6mKVXJrZri6ycpZIYovs3A8ZC182VRWlWP3N-8wp4A8PY8iBbSCs-93Oi1WMDOhbEDz0SCbKfFfyLetNnIJFqlgZiWJr80cwGgDrLauLy-qeOH0tIBXVvuX59Aq4Rd_4SmGKwF9nmW7-a8Bg7qfGyv8v7q34g_Sd3yJUWctveIR5t27TSEacx"
	// createAccountMsg := pylonstypes.NewMsgCreateAccount(creatorAddress, username, token, referral)
	// res, err := pylonsClient.CreateAccount(context.Background(), createAccountMsg)
	// if err != nil {
	// 	return false, err
	// }
	cmd = exec.Command("pylonsd", "tx", "pylons", "create-account", name, token, referral, "--from", name, "-y")
	_, err = cmd.Output()

	if err != nil {
		fmt.Println(name)
		fmt.Println("Error creating account")
		return false, err
	}
	return true, nil

}

func GenerateAddressesInKeyring(ring keyring.Keyring, n int) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		var err error
		info, _, _ := ring.NewMnemonic("NewUser"+strconv.Itoa(i), keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
		addrs[i], err = info.GetAddress()
		if err != nil {
			panic(err)
		}
		_, err = info.GetPubKey()
		if err != nil {
			panic(err)
		}
	}
	return addrs
}
