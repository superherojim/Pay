package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

var (
	ErrInvalidMnemonic  = errors.New("invalid mnemonic")
	ErrWalletGeneration = errors.New("failed to generate wallet")
	ErrDerivationPath   = errors.New("invalid derivation path")
)

// HDWallet 分层确定性钱包
type HDWallet struct {
	mu            sync.Mutex
	masterWallets []*MasterWalletInfo // 新增主钱包列表
	mnemonic      string
	passphrase    string
	basePath      accounts.DerivationPath
	index         uint32
	wallet        *hdwallet.Wallet
}

type Wallet struct {
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Path       string
}

// MasterWalletInfo 主钱包信息
type MasterWalletInfo struct {
	Mnemonic     string   // 助记词
	Address      string   // 主钱包地址
	PrivateKey   string   // 主钱包私钥(Hex)
	PublicKey    string   // 主钱包公钥(Hex)
	DerivePath   string   // 派生路径
	ChildWallets []Wallet // 已生成的子钱包
}

// NewMasterWallet 创建主钱包
func NewMasterWallet(mnemonic, passphrase string) (*HDWallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, ErrInvalidMnemonic
	}

	// 以太坊标准路径：m/44'/60'/0'/0
	basePath, err := accounts.ParseDerivationPath("m/44'/60'/0'/0")
	if err != nil {
		return nil, ErrDerivationPath
	}

	// 从助记词创建HD钱包
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to create HD wallet: %v", err)
	}

	return &HDWallet{
		mnemonic:   mnemonic,
		passphrase: passphrase,
		basePath:   basePath,
		wallet:     wallet,
	}, nil
}

// GenerateChildWallet 生成子钱包
func (w *HDWallet) GenerateChildWallet() (*Wallet, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 生成派生路径：m/44'/60'/0'/0/index
	path := make(accounts.DerivationPath, len(w.basePath))
	copy(path, w.basePath)
	path = append(path, w.index)

	// 派生账户
	account, err := w.wallet.Derive(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %v", err)
	}

	// 获取私钥
	privateKey, err := w.wallet.PrivateKey(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %v", err)
	}

	w.index++

	return &Wallet{
		Address:    account.Address,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Path:       path.String(),
	}, nil
}

// GenerateMnemonic 生成新助记词（256位熵，24个单词）
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// GetMasterWalletInfo 获取主钱包完整信息
func (w *HDWallet) GetMasterWalletInfo() (*MasterWalletInfo, error) {
	// 生成主钱包路径 m/44'/60'/0'/0/0
	basePath := w.basePath.String()
	firstChildPath := basePath + "/0"

	// 生成第一个子钱包作为主钱包账户
	child, err := w.GenerateChildWallet()
	if err != nil {
		return nil, err
	}

	return &MasterWalletInfo{
		Mnemonic:     w.mnemonic,
		Address:      child.Address.Hex(),
		PrivateKey:   hex.EncodeToString(crypto.FromECDSA(child.PrivateKey)),
		PublicKey:    hex.EncodeToString(crypto.FromECDSAPub(child.PublicKey)),
		DerivePath:   firstChildPath,
		ChildWallets: []Wallet{*child},
	}, nil
}

// AddMainWallet 新增AddMainWallet方法
func (w *HDWallet) AddMainWallet(info *MasterWalletInfo) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.masterWallets = append(w.masterWallets, info)
}

// CreateSystemMasterWallet 修改NewMasterWallet为创建系统主钱包
func CreateSystemMasterWallet() (*MasterWalletInfo, error) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		return nil, err
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("failed to create HD wallet: %v", err)
	}

	// 生成主账户路径 m/44'/60'/0'/0/0
	path, _ := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	masterInfo := &MasterWalletInfo{
		Mnemonic:   mnemonic,
		Address:    account.Address.Hex(),
		PrivateKey: hex.EncodeToString(crypto.FromECDSA(privateKey)),
		PublicKey:  hex.EncodeToString(crypto.FromECDSAPub(&privateKey.PublicKey)),
		DerivePath: path.String(),
	}

	return masterInfo, nil
}

// DeriveChildWallet 新增从主钱包派生指定路径子钱包的方法
func DeriveChildWallet(masterInfo *MasterWalletInfo, index uint32) (*Wallet, error) {
	wallet, err := hdwallet.NewFromMnemonic(masterInfo.Mnemonic)
	if err != nil {
		return nil, err
	}

	basePath, err := accounts.ParseDerivationPath(masterInfo.DerivePath)
	if err != nil {
		return nil, err
	}

	// 派生路径格式: 主路径/index
	childPath := make(accounts.DerivationPath, len(basePath))
	copy(childPath, basePath)
	childPath = append(childPath, index)

	account, err := wallet.Derive(childPath, false)
	if err != nil {
		return nil, err
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Address:    account.Address,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		Path:       childPath.String(),
	}, nil
}

// GetMasterWallets 新增获取所有主钱包的方法
func (w *HDWallet) GetMasterWallets() []*MasterWalletInfo {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.masterWallets
}
