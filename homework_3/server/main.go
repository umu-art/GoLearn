package main

import (
	"GoLearn/homework_3/proto"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"sync"
)

type server struct {
	proto.UnimplementedAccountServiceServer
}

var secretKey = uuid.New().String()
var accounts = make(map[string]*proto.Account)
var guard = &sync.RWMutex{}

func (server) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.Account, error) {
	guard.RLock()

	account, ok := accounts[req.Name]

	guard.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "account "+req.Name+" not found")
	}

	return account, nil
}

func (server) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.Account, error) {
	if len(req.Name) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}

	guard.Lock()

	if _, ok := accounts[req.Name]; ok {
		guard.Unlock()

		return nil, status.Errorf(codes.PermissionDenied, "account already exists")
	}

	newAccount := &proto.Account{
		Name:   req.Name,
		Amount: req.Amount,
	}

	accounts[req.Name] = newAccount

	guard.Unlock()

	return newAccount, nil
}

func (server) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*emptypb.Empty, error) {
	guard.Lock()

	_, ok := accounts[req.Name]

	if !ok {
		guard.Unlock()

		return nil, status.Errorf(codes.NotFound, "account "+req.Name+" not found")
	}

	delete(accounts, req.Name)

	guard.Unlock()

	return nil, nil
}

func (server) PatchAccount(ctx context.Context, req *proto.PatchAccountRequest) (*proto.Account, error) {
	guard.Lock()

	account, ok := accounts[req.Name]

	if !ok {
		guard.Unlock()

		return nil, status.Errorf(codes.NotFound, "account "+req.Name+" not found")
	}

	account.Amount = req.Amount

	guard.Unlock()

	return account, nil
}

func (server) RenameAccount(ctx context.Context, req *proto.RenameAccountRequest) (*proto.Account, error) {
	if len(req.NewName) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}

	guard.Lock()

	account, ok := accounts[req.OldName]

	if !ok {
		guard.Unlock()

		return nil, status.Errorf(codes.NotFound, "account "+req.OldName+" not found")
	}

	_, ok = accounts[req.NewName]

	if ok {
		guard.Unlock()

		return nil, status.Errorf(codes.PermissionDenied, "account "+req.NewName+" already exists")
	}

	account.Name = req.NewName
	accounts[req.NewName] = account
	delete(accounts, req.OldName)

	guard.Unlock()

	return account, nil
}

func (server) GetAllAccounts(ctx context.Context, req *proto.GetAllAccountsRequest) (*proto.GetAllAccountsResponse, error) {
	if req.SecretKey != secretKey {
		return nil, status.Errorf(codes.PermissionDenied, "мимо")
	}

	accountsArray := make([]*proto.Account, 0, len(accounts))

	for _, account := range accounts {
		accountsArray = append(accountsArray, &proto.Account{
			Name:   account.Name,
			Amount: account.Amount,
		})
	}

	return &proto.GetAllAccountsResponse{
		Accounts: accountsArray,
	}, nil
}

func main() {
	println("Секретный ключ:", secretKey)

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
