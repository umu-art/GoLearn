package main

import (
	"GoLearn/homework_3/proto"
	"GoLearn/homework_4/server/elastic"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type server struct {
	proto.UnimplementedAccountServiceServer
}

var accountStorage elastic.AccountStorage

func (server) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.Account, error) {

	exists, err := accountStorage.IsExistsAccount(req.Name, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования аккаунта, %w", err)
	}

	if !exists {
		return nil, status.Errorf(codes.NotFound, "account "+req.Name+" not found")
	}

	account, err := accountStorage.GetAccount(req.Name, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения аккаунта, %w", err)
	}

	return account, nil
}

func (server) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.Account, error) {

	if len(req.Name) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}

	exists, err := accountStorage.IsExistsAccount(req.Name, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования аккаунта, %w", err)
	}

	if exists {
		return nil, status.Errorf(codes.PermissionDenied, "account already exists")
	}

	newAccount := &proto.Account{
		Name:   req.Name,
		Amount: req.Amount,
	}

	err = accountStorage.PatchAccount(newAccount, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления аккаунта, %w", err)
	}

	return newAccount, nil
}

func (server) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*emptypb.Empty, error) {

	exists, err := accountStorage.IsExistsAccount(req.Name, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования аккаунта, %w", err)
	}

	if !exists {
		return nil, status.Errorf(codes.NotFound, "account "+req.Name+" not found")
	}

	err = accountStorage.DeleteAccount(req.Name, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка удаления аккаунта, %w", err)
	}

	return nil, nil
}

func (server) PatchAccount(ctx context.Context, req *proto.PatchAccountRequest) (*proto.Account, error) {

	exists, err := accountStorage.IsExistsAccount(req.Name, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования аккаунта, %w", err)
	}

	if !exists {
		return nil, status.Errorf(codes.NotFound, "account "+req.Name+" not found")
	}

	patchedAccount := &proto.Account{
		Name:   req.Name,
		Amount: req.Amount,
	}

	err = accountStorage.PatchAccount(patchedAccount, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления аккаунта, %w", err)
	}

	return patchedAccount, nil
}

func (server) RenameAccount(ctx context.Context, req *proto.RenameAccountRequest) (*proto.Account, error) {

	if len(req.NewName) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}

	exists, err := accountStorage.IsExistsAccount(req.OldName, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования аккаунта, %w", err)
	}

	if !exists {
		return nil, status.Errorf(codes.NotFound, "account "+req.OldName+" not found")
	}

	exists, err = accountStorage.IsExistsAccount(req.NewName, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования аккаунта, %w", err)
	}

	if exists {
		return nil, status.Errorf(codes.PermissionDenied, "account "+req.NewName+" already exists")
	}

	oldAccount, err := accountStorage.GetAccount(req.OldName, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения аккаунта, %w", err)
	}

	newAccount := &proto.Account{
		Name:   req.NewName,
		Amount: oldAccount.Amount,
	}

	err = accountStorage.DeleteAccount(req.OldName, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка удаления аккаунта, %w", err)
	}

	err = accountStorage.PatchAccount(newAccount, ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления аккаунта, %w", err)
	}

	return newAccount, nil
}

func main() {
	err := accountStorage.Init()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", "0.0.0.0:5445")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterAccountServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
