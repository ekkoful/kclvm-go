// Copyright 2021 The KCL Authors. All rights reserved.

package service

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"kusionstack.io/kclvm-go/pkg/spec/gpyrpc"
)

var _ = fmt.Sprint

func RunGrpcServer(address string) error {
	grpcServer := grpc.NewServer()
	gpyrpc.RegisterKclvmServiceServer(grpcServer, newKclvmServiceImpl())

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer.Serve(lis)
	return nil
}

type _KclvmServiceImpl struct {
	c *KclvmServiceClient
}

func newKclvmServiceImpl() *_KclvmServiceImpl {
	return &_KclvmServiceImpl{
		c: NewKclvmServiceClient(),
	}
}

func (p *_KclvmServiceImpl) Ping(ctx context.Context, args *gpyrpc.Ping_Args) (*gpyrpc.Ping_Result, error) {
	return p.c.Ping(args)
}
func (p *_KclvmServiceImpl) ExecProgram(ctx context.Context, args *gpyrpc.ExecProgram_Args) (*gpyrpc.ExecProgram_Result, error) {
	return p.c.ExecProgram(args)
}
func (p *_KclvmServiceImpl) FormatCode(ctx context.Context, args *gpyrpc.FormatCode_Args) (*gpyrpc.FormatCode_Result, error) {
	return p.c.FormatCode(args)
}
func (p *_KclvmServiceImpl) FormatPath(ctx context.Context, args *gpyrpc.FormatPath_Args) (*gpyrpc.FormatPath_Result, error) {
	return p.c.FormatPath(args)
}
func (p *_KclvmServiceImpl) LintPath(ctx context.Context, args *gpyrpc.LintPath_Args) (*gpyrpc.LintPath_Result, error) {
	return p.c.LintPath(args)
}
func (p *_KclvmServiceImpl) OverrideFile(ctx context.Context, args *gpyrpc.OverrideFile_Args) (*gpyrpc.OverrideFile_Result, error) {
	return p.c.OverrideFile(args)
}
func (p *_KclvmServiceImpl) GetSchemaType(ctx context.Context, args *gpyrpc.GetSchemaType_Args) (*gpyrpc.GetSchemaType_Result, error) {
	return p.c.GetSchemaType(args)
}
func (p *_KclvmServiceImpl) GetSchemaTypeMapping(ctx context.Context, args *gpyrpc.GetSchemaTypeMapping_Args) (*gpyrpc.GetSchemaTypeMapping_Result, error) {
	return p.c.GetSchemaTypeMapping(args)
}
func (p *_KclvmServiceImpl) ValidateCode(ctx context.Context, args *gpyrpc.ValidateCode_Args) (*gpyrpc.ValidateCode_Result, error) {
	return p.c.ValidateCode(args)
}
func (p *_KclvmServiceImpl) ListDepFiles(ctx context.Context, args *gpyrpc.ListDepFiles_Args) (*gpyrpc.ListDepFiles_Result, error) {
	return p.c.ListDepFiles(args)
}
func (p *_KclvmServiceImpl) LoadSettingsFiles(ctx context.Context, args *gpyrpc.LoadSettingsFiles_Args) (*gpyrpc.LoadSettingsFiles_Result, error) {
	return p.c.LoadSettingsFiles(args)
}
