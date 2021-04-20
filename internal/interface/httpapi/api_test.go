package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"koro.che/internal/usecases/account"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AccountUseCasesFake struct{}

func (AccountUseCasesFake) CreateAccount(login, password string) (account.Account, error) {
	switch login {
	case "test":
		return account.Account{
			Id: "test_id",
		}, nil
	default:
		return account.Account{}, errors.New("failed to create an account")
	}
}

func (AccountUseCasesFake) GetAccountById(id string) (account.Account, error) {
	panic("implement me")
}

func (AccountUseCasesFake) LoginToAccount(login, password string) (string, error) {
	if login == "test" && password == "test" {
		return "token", nil
	}
	return "", errors.New("invalid login or password")
}

func (a *AccountUseCasesFake) Authenticate(token string) (string, error) {
	panic("implement me")
}

func (a AccountUseCasesFake) Logout() {
	panic("implement me")
}

func Test_postSignup(t *testing.T) {
	service := NewApi(&AccountUseCasesFake{}, nil)
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/api/register")
		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("failed to create account", func(t *testing.T) {
		m := registrationModel{
			Login:    "bob",
			Password: "test",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal("failed to marshal struct")
		}
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(b))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assertStatusCode(t, resp.Code, http.StatusInternalServerError)
	})

	t.Run("successful account creation", func(t *testing.T) {
		m := registrationModel{
			Login:    "test",
			Password: "test",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal("failed to marshal struct")
		}
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(b))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assertStatusCode(t, resp.Code, http.StatusCreated)

		location := resp.Header().Get("Location")
		if location != "/accounts/test_id" {
			t.Errorf("Server MUST return %s Location header, but %s given", "/accounts/test_id", location)
		}
	})
}

func Test_postSignin(t *testing.T) {
	service := NewApi(&AccountUseCasesFake{}, nil)
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/api/register")
		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("failed login with incorrect login or password", func(t *testing.T) {
		m := registrationModel{
			Login:    "bob",
			Password: "test",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal("failed to marshal struct")
		}
		req := httptest.NewRequest(http.MethodPut, "/api/login", bytes.NewReader(b))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("successful login with correct password", func(t *testing.T) {
		m := registrationModel{
			Login:    "test",
			Password: "test",
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal("failed to marshal struct")
		}
		req := httptest.NewRequest(http.MethodPut, "/api/login", bytes.NewReader(b))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assertStatusCode(t, resp.Code, http.StatusOK)
	})
}

func assertStatusCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Server MUST return %d (%s) status code, but %d (%s) given",
			expectedCode, http.StatusText(expectedCode), actualCode, http.StatusText(actualCode))
	}
}

func invalidJsonTest(router http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{a:")))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}
