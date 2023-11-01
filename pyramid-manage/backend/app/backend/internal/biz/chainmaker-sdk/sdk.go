package chainmaker_sdk

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"context"
)

type ChainMakerClient struct {
	client *sdk.ChainClient
}

func NewChainMakerClient(ctx context.Context, path string) (*ChainMakerClient, error) {
	client, err := sdk.NewChainClient(sdk.WithConfPath(path))
	if err != nil {
		return nil, err
	}
	return &ChainMakerClient{client: client}, err
}

func (c *ChainMakerClient) Invoke(contractName, method, txId string, params []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {
	return c.client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
}
