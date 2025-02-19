package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

var key = "29b75c2edac447aea792135fdf065d86"
var endpoint = "https://mainnet.infura.io/v3/29b75c2edac447aea792135fdf065d86"

const keystoreDir = "./wallets" // Keystore 文件存放目录

func main() {
	fmt.Println("1️⃣ 生成新的 Keystore 并存储")
	mnemonic, privateKey, address, err := createKeystore("your-secure-password")
	if err != nil {
		log.Fatal("Keystore 生成失败:", err)
	}
	fmt.Println("Keystore 生成成功，钱包地址:", address)
	fmt.Println("助记词:", mnemonic)
	fmt.Println("私钥:", hex.EncodeToString(privateKey.D.Bytes()))

	fmt.Println("\n2️⃣ 生成子钱包")
	subMnemonic, subPrivateKey, subAddress, err := createSubWallet(mnemonic)
	if err != nil {
		log.Fatal("子钱包生成失败:", err)
	}
	fmt.Println("子钱包生成成功，子钱包地址:", subAddress)
	fmt.Println("子钱包助记词:", subMnemonic)
	fmt.Println("子钱包私钥:", hex.EncodeToString(subPrivateKey.D.Bytes()))
}

// 创建 Keystore 并存储
func createKeystore(password string) (string, *ecdsa.PrivateKey, string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", nil, "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", nil, "", err
	}
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", nil, "", err
	}
	derivedKey, err := deriveETHKey(masterKey)
	if err != nil {
		return "", nil, "", err
	}
	privateKey, err := crypto.ToECDSA(derivedKey.Key)
	if err != nil {
		return "", nil, "", err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return mnemonic, privateKey, address, nil
}

// 生成子钱包（BIP44规范：m/44'/60'/0'/0/1）
func createSubWallet(parentMnemonic string) (string, *ecdsa.PrivateKey, string, error) {
	seed := bip39.NewSeed(parentMnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", nil, "", err
	}
	// 使用 m/44'/60'/0'/0/1 作为子钱包派生路径
	subKey, err := deriveETHKeyWithIndex(masterKey, 1)
	if err != nil {
		return "", nil, "", err
	}
	privateKey, err := crypto.ToECDSA(subKey.Key)
	if err != nil {
		return "", nil, "", err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return parentMnemonic, privateKey, address, nil
}

// 按照 BIP44 规范派生 ETH 钱包（路径：m/44'/60'/0'/0/0）
func deriveETHKey(masterKey *bip32.Key) (*bip32.Key, error) {
	purpose, _ := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	coinType, _ := purpose.NewChildKey(60 + bip32.FirstHardenedChild)
	account, _ := coinType.NewChildKey(0 + bip32.FirstHardenedChild)
	change, _ := account.NewChildKey(0)
	return change.NewChildKey(0)
}

// 按索引派生 ETH 子钱包（路径：m/44'/60'/0'/0/index）
func deriveETHKeyWithIndex(masterKey *bip32.Key, index uint32) (*bip32.Key, error) {
	purpose, _ := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	coinType, _ := purpose.NewChildKey(60 + bip32.FirstHardenedChild)
	account, _ := coinType.NewChildKey(0 + bip32.FirstHardenedChild)
	change, _ := account.NewChildKey(0)
	return change.NewChildKey(index)
}
