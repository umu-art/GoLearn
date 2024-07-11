package main

import (
	"GoLearn/homework_2/dto"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ConnectionInfo struct {
	Port int
	Host string
}

type CommandInfo struct {
	ConnectionInfo ConnectionInfo
	Command        string
	Execute        func(ConnectionInfo) error
}

var commands = []CommandInfo{
	{
		Command: "Создать новый аккаунт",
		Execute: func(info ConnectionInfo) error {
			var request dto.CreateAccountRequest

			fmt.Print("Введите имя нового аккаунта: ")
			_, _ = fmt.Scan(&request.Name)

			fmt.Print("Введите баланс: ")
			_, _ = fmt.Scan(&request.Amount)

			data, err := json.Marshal(request)
			if err != nil {
				return fmt.Errorf("json marshal failed: %w", err)
			}

			_, err = ExecuteHttp(info, "api/account/create", "POST", data)
			if err != nil {
				return err
			}

			println("Аккаунт создан")
			return nil
		},
	},
	{
		Command: "Получить информацию об аккаунте",
		Execute: func(info ConnectionInfo) error {
			var name string

			fmt.Print("Введите имя аккаунта: ")
			_, _ = fmt.Scan(&name)

			resp, err := ExecuteHttp(info, fmt.Sprintf("api/account?name=%s", name), "GET", nil)
			if err != nil {
				return err
			}

			var response dto.GetAccountResponse
			if err := json.Unmarshal(resp, &response); err != nil {
				return fmt.Errorf("json unmarshal failed: %w", err)
			}

			fmt.Printf("Имя: %s, Баланс: %d\n", response.Name, response.Amount)

			return nil
		},
	},
	{
		Command: "Изменить баланс",
		Execute: func(info ConnectionInfo) error {
			var request dto.PatchAccountRequest

			fmt.Print("Введите имя аккаунта: ")
			_, _ = fmt.Scan(&request.Name)

			fmt.Print("Введите новую сумму: ")
			_, _ = fmt.Scan(&request.Amount)

			data, err := json.Marshal(request)
			if err != nil {
				return fmt.Errorf("json marshal failed: %w", err)
			}

			_, err = ExecuteHttp(info, "api/account", "PATCH", data)
			if err != nil {
				return err
			}

			println("Баланс изменен")
			return nil
		},
	},
	{
		Command: "Удалить аккаунт",
		Execute: func(info ConnectionInfo) error {
			var request dto.DeleteAccountRequest

			fmt.Print("Введите имя аккаунта: ")
			_, _ = fmt.Scan(&request.Name)

			data, err := json.Marshal(request)
			if err != nil {
				return fmt.Errorf("json marshal failed: %w", err)
			}

			_, err = ExecuteHttp(info, "api/account", "DELETE", data)
			if err != nil {
				return err
			}

			println("Аккаунт удален")
			return nil
		},
	},
	{
		Command: "Переименовать аккаунт",
		Execute: func(info ConnectionInfo) error {
			var request dto.ChangeAccountRequest

			fmt.Print("Введите имя аккаунта: ")
			_, _ = fmt.Scan(&request.Name)

			fmt.Print("Введите новое имя: ")
			_, _ = fmt.Scan(&request.NewName)

			data, err := json.Marshal(request)
			if err != nil {
				return fmt.Errorf("json marshal failed: %w", err)
			}

			_, err = ExecuteHttp(info, "api/account/rename", "POST", data)
			if err != nil {
				return err
			}

			println("Аккаунт переименован")
			return nil
		},
	},
}

func ExecuteHttp(info ConnectionInfo, path string, method string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s:%d/%s", info.Host, info.Port, path),
		bytes.NewReader(data))

	if err != nil {
		return nil, fmt.Errorf("http request create failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusCreated {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return body, nil
}

func main() {
	hostVal := flag.String("host", "http://go.http.umu-art.ru", "server api host")
	portVal := flag.Int("port", 80, "server port")
	secretKey := flag.String("secret-key", "", "Ключ админа")

	flag.Parse()

	connectionInfo := ConnectionInfo{
		Host: *hostVal,
		Port: *portVal,
	}

	if len(*secretKey) > 0 {
		println("Попытка получить список всех аккаунтов...")
		resp, err := ExecuteHttp(connectionInfo, fmt.Sprintf("api/accounts?secret-key=%s", *secretKey), "GET", nil)
		if err != nil {
			println(err.Error())
		} else {
			var response []dto.GetAccountResponse
			if err := json.Unmarshal(resp, &response); err != nil {
				println(err.Error())
			} else {
				fmt.Printf("Всего аккаунтов: %d\n", len(response))
				for _, account := range response {
					fmt.Printf("Имя: %s, Баланс: %d\n", account.Name, account.Amount)
				}
			}
		}
		os.Exit(0)
	}

	for {
		play(connectionInfo)
	}
}

func play(connectionInfo ConnectionInfo) {
	println()

	for i, command := range commands {
		println(i+1, ">", command.Command)
	}
	println("0 > Выход")
	print("Выберите команду: ")

	var commandIndex int
	_, _ = fmt.Scan(&commandIndex)
	if commandIndex < 0 || commandIndex > len(commands) {
		println("Неверный индекс")
		return
	}

	if commandIndex == 0 {
		os.Exit(0)
	}

	err := commands[commandIndex-1].Execute(connectionInfo)
	if err != nil {
		println(err.Error())
	}
}
