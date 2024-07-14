package elastic

import (
	"GoLearn/homework_4/proto"
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"os"
)

type AccountStorage struct {
	typedClient *elasticsearch.TypedClient
}

var ConnectionError = errors.New("поломалося соединение с elasticsearch")

func (a *AccountStorage) Init() error {
	cert, _ := os.ReadFile("homework_4/server/elastic/http_ca.crt")

	typedClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{
			"https://go-homework-4-elastic.go.svc.cluster.local:9200/",
		},
		APIKey: "U1VMaHNaQUJscWNENUg4RGUwR1Q6NERKOEJUVjVUalNHYTRFRU9Id2hoZw==",
		CACert: cert,
	})
	if err != nil {
		return ConnectionError
	}

	resp, err := typedClient.Indices.Exists("accounts").Do(context.Background())
	if err != nil {
		return ConnectionError
	}

	if !resp {
		_, err = typedClient.Indices.Create("accounts").Do(context.Background())
		if err != nil {
			return ConnectionError
		}
	}

	a.typedClient = typedClient

	return nil
}

func (a *AccountStorage) IsExistsAccount(name string, ctx context.Context) (bool, error) {
	resp, err := a.typedClient.
		Exists("accounts", name).
		Do(ctx)

	if err != nil {
		return false, err
	}

	return resp, nil
}

func (a *AccountStorage) GetAccount(name string, ctx context.Context) (*proto.Account, error) {
	resp, err := a.typedClient.
		Get("accounts", name).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	var account proto.Account

	err = json.Unmarshal(resp.Source_, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountStorage) PatchAccount(account *proto.Account, ctx context.Context) error {
	_, err := a.typedClient.
		Index("accounts").
		Id(account.Name).
		Request(account).
		Do(ctx)

	return err
}

func (a *AccountStorage) DeleteAccount(name string, ctx context.Context) error {
	_, err := a.typedClient.
		Delete("accounts", name).
		Do(ctx)

	return err
}
