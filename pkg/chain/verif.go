package chain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"google.golang.org/grpc"
)

type ChainVerifier interface {
	VerifyTransaction(ctx context.Context, rpcURL string, txHash string) (confirmations int64, status int, err error)
}

// 以太坊系验证器
type EVMVerifier struct{}

func (v *EVMVerifier) VerifyTransaction(ctx context.Context, rpcURL string, txHash string) (int64, int, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return 0, 0, fmt.Errorf("连接节点失败: %v", err)
	}
	defer client.Close()

	hash := common.HexToHash(txHash)
	receipt, err := client.TransactionReceipt(ctx, hash)
	if err != nil {
		return 0, 0, fmt.Errorf("获取交易收据失败: %v", err)
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, int(receipt.Status), fmt.Errorf("获取最新区块失败: %v", err)
	}

	return header.Number.Int64() - receipt.BlockNumber.Int64(), int(receipt.Status), nil
}

type TronVerifier struct{}

func (v *TronVerifier) VerifyTransaction(ctx context.Context, rpcURL string, txHash string) (int64, int, error) {
	// 连接到 Tron 节点的 gRPC 接口
	conn, err := grpc.Dial(rpcURL, grpc.WithInsecure())
	if err != nil {
		return 0, 0, fmt.Errorf("连接到 Tron 节点失败: %v", err)
	}
	defer conn.Close()

	// 创建一个新的客户端
	client := api.NewWalletClient(conn)

	// 查询交易信息
	txInfo, err := client.GetTransactionById(ctx, &api.BytesMessage{Value: []byte(txHash)})
	if err != nil {
		return 0, 0, fmt.Errorf("获取交易信息失败: %v", err)
	}

	// 检查交易状态
	if txInfo == nil {
		return 0, 0, fmt.Errorf("交易未找到")
	}

	if txInfo.Ret[0].ContractRet == core.Transaction_Result_SUCCESS {
		return 1, 1, nil // 交易成功
	} else {
		return 0, 0, fmt.Errorf("交易失败，状态码: %s", txInfo.Ret[0].ContractRet)
	}
}

type SolanaVerifier struct{}

func (v *SolanaVerifier) VerifyTransaction(ctx context.Context, rpcURL string, txHash string) (int64, int, error) {
	// 创建 Solana 客户端
	client := rpc.New(rpcURL)

	// 获取交易签名
	sig, err := solana.SignatureFromBase58(txHash)
	if err != nil {
		return 0, 0, fmt.Errorf("无效的交易哈希: %v", err)
	}

	// 获取交易状态
	status, err := client.GetSignatureStatuses(ctx, true, sig)
	if err != nil {
		return 0, 0, fmt.Errorf("获取交易状态失败: %v", err)
	}

	// 检查交易状态
	if status == nil || len(status.Value) == 0 || status.Value[0] == nil {
		return 0, 0, fmt.Errorf("交易不存在或未确认")
	}

	if status.Value[0].ConfirmationStatus != rpc.ConfirmationStatusFinalized {
		return 0, 0, fmt.Errorf("交易未最终确认，状态: %s", status.Value[0].ConfirmationStatus)
	}

	// 获取当前区块高度
	slot, err := client.GetSlot(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return 0, 0, fmt.Errorf("获取当前区块高度失败: %v", err)
	}

	// 计算确认数
	confirmations := slot - status.Value[0].Slot
	return int64(confirmations), 1, nil
}
