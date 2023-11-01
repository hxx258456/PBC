package biz

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	pb "pyramid/pyramid-manage/backend/api/backend/v1"
	sdk "pyramid/pyramid-manage/backend/app/backend/internal/biz/chainmaker-sdk"
)

type SysChainUseCase struct {
}

type SysChainRepo interface {
}

func NewSysChainUseCase() *SysChainUseCase {
	return &SysChainUseCase{}
}

func (uc *SysChainUseCase) AddOrganization(ctx context.Context, org *pb.Organization) (*pb.BaseReply, error) {
	client, _ := sdk.NewChainMakerClient(ctx, "../../configs/chain2_config.yml")
	o, _ := json.Marshal(org)

	rParams := []*common.KeyValuePair{
		{Key: "type", Value: []byte("0")},
		{Key: "organization", Value: o},
	}

	resp, err := client.Invoke("pyramid2", "save", "", rParams, true)
	if err != nil {
	}
	if resp.Code != common.TxStatusCode_SUCCESS {

	}
	return &pb.BaseReply{}, nil
}

func (uc *SysChainUseCase) QueryOrganization(ctx context.Context, org *pb.Organization) (*pb.Organization, error) {
	client, err := sdk.NewChainMakerClient(ctx, "../../configs/chain2_config.yml")
	if err != nil {
		fmt.Println(err)
	}
	var o = pb.Organization{}
	rParams := []*common.KeyValuePair{
		{Key: "type", Value: []byte("0")},
		{Key: "org_id", Value: []byte(org.OrgId)},
	}
	resp, err := client.Invoke("pyramid2", "query", "", rParams, true)
	if err != nil {
	}
	if resp.Code != common.TxStatusCode_SUCCESS {
	}
	json.Unmarshal(resp.ContractResult.Result, &o)
	return &o, nil
}
