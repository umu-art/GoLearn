package account

import (
	"GoLearn/homework_2/dto"
	"GoLearn/homework_2/server/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
)

func New(secretKey string) *Handler {
	return &Handler{
		accounts:  make(map[string]*model.Account),
		guard:     &sync.RWMutex{},
		secretKey: secretKey,
	}
}

type Handler struct {
	accounts  map[string]*model.Account
	guard     *sync.RWMutex
	secretKey string
}

// Actuator проверка живо ли оно
func (h *Handler) Actuator(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// GetAll возвращает список всех аккаунтов (для админки)
func (h *Handler) GetAll(c echo.Context) error {
	secretKey := c.QueryParams().Get("secret-key")

	if secretKey != h.secretKey {
		return c.NoContent(http.StatusForbidden)
	}

	h.guard.RLock()

	accounts := make([]dto.GetAccountResponse, 0, len(h.accounts))

	for _, account := range h.accounts {
		accounts = append(accounts, dto.GetAccountResponse{
			Name:   account.Name,
			Amount: account.Amount,
		})
	}

	h.guard.RUnlock()

	return c.JSON(http.StatusOK, accounts)
}

// CreateAccount создает новый аккаунт
// POST /account/create
// {"name": "alice", "amount": 50}
func (h *Handler) CreateAccount(c echo.Context) error {
	var request dto.CreateAccountRequest
	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err)

		return c.String(http.StatusBadRequest, "invalid request")
	}

	if len(request.Name) == 0 {
		return c.String(http.StatusBadRequest, "empty name")
	}

	h.guard.Lock()

	if _, ok := h.accounts[request.Name]; ok {
		h.guard.Unlock()

		return c.String(http.StatusForbidden, "account already exists")
	}

	h.accounts[request.Name] = &model.Account{
		Name:   request.Name,
		Amount: request.Amount,
	}

	h.guard.Unlock()

	return c.NoContent(http.StatusCreated)
}

// GetAccount возвращает информацию об аккаунте
// GET /account?name=alice
func (h *Handler) GetAccount(c echo.Context) error {
	name := c.QueryParams().Get("name")

	h.guard.RLock()

	account, ok := h.accounts[name]

	h.guard.RUnlock()

	if !ok {
		return c.String(http.StatusNotFound, "account "+name+" not found")
	}

	response := dto.GetAccountResponse{
		Name:   account.Name,
		Amount: account.Amount,
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteAccount удаляет аккаунт
// DELETE /account/delete
// {"name": "alice"}
func (h *Handler) DeleteAccount(c echo.Context) error {
	var request dto.DeleteAccountRequest

	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err)

		return c.String(http.StatusBadRequest, "invalid request")
	}

	h.guard.Lock()

	_, ok := h.accounts[request.Name]

	if !ok {
		h.guard.Unlock()

		return c.String(http.StatusNotFound, "account "+request.Name+" not found")
	}

	delete(h.accounts, request.Name)

	h.guard.Unlock()

	return c.NoContent(http.StatusOK)
}

// PatchAccount меняет баланс
// POST /account/patch
// {"name": "alice", "amount": 50}
func (h *Handler) PatchAccount(c echo.Context) error {
	var request dto.PatchAccountRequest

	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err)

		return c.String(http.StatusBadRequest, "invalid request")
	}

	h.guard.Lock()

	account, ok := h.accounts[request.Name]

	if !ok {
		h.guard.Unlock()

		return c.String(http.StatusNotFound, "account "+request.Name+" not found")
	}

	account.Amount = request.Amount

	h.guard.Unlock()

	return c.NoContent(http.StatusOK)
}

// ChangeAccount меняет имя аккаунта
// POST /account/change
// {"name": "alice", "new-name": 50}
func (h *Handler) ChangeAccount(c echo.Context) error {
	var request dto.ChangeAccountRequest

	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err)

		return c.String(http.StatusBadRequest, "invalid request")
	}

	if len(request.NewName) == 0 {
		return c.String(http.StatusBadRequest, "empty new name")
	}

	h.guard.Lock()

	account, ok := h.accounts[request.Name]

	if !ok {
		h.guard.Unlock()

		return c.String(http.StatusNotFound, "account "+request.Name+" not found")
	}

	_, ok = h.accounts[request.NewName]

	if ok {
		h.guard.Unlock()

		return c.String(http.StatusForbidden, "account "+request.NewName+" already exists")
	}

	account.Name = request.NewName
	h.accounts[request.NewName] = account
	delete(h.accounts, request.Name)

	h.guard.Unlock()

	return c.NoContent(http.StatusOK)
}
