package chain

import (
	"fmt"
	"sort"
)

// 添加链类型枚举
type ChainType string

const (
	ChainTypeEVM    ChainType = "EVM"
	ChainTypeTron   ChainType = "TRON"
	ChainTypeSolana ChainType = "SOLANA"
)

type ChainInfo struct {
	ID            string
	Name          string
	Type          ChainType
	IsMainnet     bool
	Confirmations int64
}

var SupportedChains = map[string]ChainInfo{
	// EVM主网
	"1":     {"1", "Ethereum", ChainTypeEVM, true, 12},
	"56":    {"56", "BNB Chain", ChainTypeEVM, true, 15},
	"137":   {"137", "Polygon", ChainTypeEVM, true, 6},
	"42161": {"42161", "Arbitrum", ChainTypeEVM, true, 6},
	"10":    {"10", "Optimism", ChainTypeEVM, true, 6},
	"43114": {"43114", "Avalanche", ChainTypeEVM, true, 6},

	// EVM测试网
	"5":        {"5", "Goerli", ChainTypeEVM, false, 3},
	"11155111": {"11155111", "Sepolia", ChainTypeEVM, false, 3},
	"97":       {"97", "BSC Testnet", ChainTypeEVM, false, 3},
	"80001":    {"80001", "Mumbai", ChainTypeEVM, false, 3},
	"421613":   {"421613", "Arbitrum Goerli", ChainTypeEVM, false, 3},

	// Tron
	"TRX":  {"TRX", "Tron Mainnet", ChainTypeTron, true, 1},
	"nile": {"nile", "Tron Nile", ChainTypeTron, false, 1},

	// Solana
	"SOL":  {"SOL", "Solana Mainnet", ChainTypeSolana, true, 32},
	"SOLT": {"SOLT", "Solana Testnet", ChainTypeSolana, false, 16},
}

func IsSupportedChain(chainID string) bool {
	_, exists := SupportedChains[chainID]
	return exists
}
func GetChainInfo(chainID string) (ChainInfo, bool) {
	info, exists := SupportedChains[chainID]
	return info, exists
}

func GetChainList(mainnetOnly bool) []string {
	var chains []string
	for _, info := range SupportedChains {
		if mainnetOnly && !info.IsMainnet {
			continue
		}
		chains = append(chains, fmt.Sprintf("%s (%s)", info.ID, info.Name))
	}
	sort.Strings(chains)
	return chains
}

type VerifierFactory struct {
	evmVerifier    *EVMVerifier
	tronVerifier   *TronVerifier
	solanaVerifier *SolanaVerifier
}

func NewVerifierFactory() *VerifierFactory {
	return &VerifierFactory{
		evmVerifier:    &EVMVerifier{},
		tronVerifier:   &TronVerifier{},
		solanaVerifier: &SolanaVerifier{},
	}
}

func (f *VerifierFactory) GetVerifier(chainType ChainType) (ChainVerifier, error) {
	switch chainType {
	case ChainTypeEVM:
		return f.evmVerifier, nil
	case ChainTypeTron:
		return f.tronVerifier, nil
	case ChainTypeSolana:
		return f.solanaVerifier, nil
	default:
		return nil, fmt.Errorf("unsupported chain type: %s", chainType)
	}
}
