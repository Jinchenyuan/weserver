package servicehandler

import (
	"context"
	"fmt"
	"server/core"
	"server/model"
	pb "server/protobuf/gen"
	protocol "server/submodule/protocol/gen/golang"
)

type S3 struct{}

func (a *S3) PutKey(ctx context.Context, req *pb.PutKeyReq, rsp *pb.PutKeyResp) error {
	fmt.Printf("PutKey request received: data=%s, key=%s\n", req.GetData(), req.GetKey())
	m := core.GetGlobalMesa()
	if m == nil {
		rsp.Code = int32(protocol.ErrorCode_INTERNAL)
		rsp.Message = "failed to get global mesa"
		return nil
	}

	s3kv := &model.S3KV{
		Key:   req.GetKey(),
		Value: req.GetData(),
	}
	s3kv.SetDB(m.DB)

	err := s3kv.Create(ctx)
	if err != nil {
		rsp.Code = int32(protocol.ErrorCode_INTERNAL)
		rsp.Message = fmt.Sprintf("failed to put key: %v", err)
		return nil
	}

	rsp.Code = int32(protocol.ErrorCode_OK)
	rsp.Message = "PutKey successful"
	return nil
}

func (a *S3) GetKey(ctx context.Context, req *pb.GetKeyReq, rsp *pb.GetKeyResp) error {
	fmt.Printf("GetKey request received: key=%s\n", req.GetKey())
	m := core.GetGlobalMesa()
	if m == nil {
		rsp.Code = int32(protocol.ErrorCode_INTERNAL)
		rsp.Message = "failed to get global mesa"
		return nil
	}

	s3kv := &model.S3KV{}
	s3kv.SetDB(m.DB)

	s3kv, err := s3kv.FindByKey(ctx, req.GetKey())
	if err != nil {
		rsp.Code = int32(protocol.ErrorCode_INTERNAL)
		rsp.Message = fmt.Sprintf("failed to get key: %v", err)
		return nil
	}

	rsp.Code = int32(protocol.ErrorCode_OK)
	rsp.Data = s3kv.Value
	rsp.Message = "GetKey successful"
	return nil
}

func (a *S3) DeleteKey(ctx context.Context, req *pb.DeleteKeyReq, rsp *pb.DeleteKeyResp) error {
	fmt.Printf("DeleteKey request received: key=%s\n", req.GetKey())
	m := core.GetGlobalMesa()
	if m == nil {
		rsp.Code = int32(protocol.ErrorCode_INTERNAL)
		rsp.Message = "failed to get global mesa"
		return nil
	}

	s3kv := &model.S3KV{}
	s3kv.SetDB(m.DB)
	err := s3kv.Delete(ctx, req.GetKey())
	if err != nil {
		rsp.Code = int32(protocol.ErrorCode_INTERNAL)
		rsp.Message = fmt.Sprintf("failed to delete key: %v", err)
		return nil
	}

	rsp.Code = int32(protocol.ErrorCode_OK)
	rsp.Message = "DeleteKey successful"
	return nil
}
