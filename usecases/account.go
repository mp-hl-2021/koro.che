package usecases

import (
	"golang.org/x/crypto/bcrypt"
	"koro.che/auth"

	"errors"
	"koro.che/accountstorage"
	"unicode"
)

var (
	ErrInvalidLoginString    = errors.New("login string contains invalid character")
	ErrInvalidPasswordString = errors.New("password string contains invalid character")
	ErrTooShortLogin       = errors.New("too short login")
	ErrTooLongLogin        = errors.New("too long login")
	ErrTooShortPassword       = errors.New("too short password")
	ErrTooLongPassword         = errors.New("too long password")
)

const (
	minLoginLength    = 6
	maxLoginLength    = 20
	minPasswordLength = 8
	maxPasswordLength = 48
)

type Account struct {
	Id string
}

type AccountUseCasesInterface interface {
	CreateAccount(login, password string) (Account, error)
	GetAccountById(id string) (Account, error)
	LoginToAccount(login, password string) (string, error)
	Authenticate(token string) (string, error)
	Logout() ()
}

type AccountUseCases struct {
	AccountStorage accountstorage.Interface
	Auth           auth.Interface
}

func (a* AccountUseCases) CreateAccount(login string, password string) (Account, error) {
	if err := validateLogin(login); err != nil {
		return Account{}, err
	}
	if err := validatePassword(password); err != nil {
		return Account{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}
	acc, err := a.AccountStorage.CreateAccount(accountstorage.Credentials{
		Login:    login,
		Password: string(hashedPassword),
	})
	if err != nil {
		return Account{}, err
	}
	return Account{Id: acc.Id}, nil
}

func (a* AccountUseCases) GetAccountById(id string) (Account, error) {
	acc, err := a.AccountStorage.GetAccountById(id)
	if err != nil {
		return Account{}, err
	}
	return Account{Id: acc.Id}, err
}

func (a* AccountUseCases) LoginToAccount(login string, password string) (string, error) {
	if err := validateLogin(login); err != nil {
		return "", err
	}
	if err := validatePassword(password); err != nil {
		return "", err
	}
	acc, err := a.AccountStorage.GetAccountByLogin(login)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password)); err != nil {
		return "", err
	}
	token, err := a.Auth.IssueToken(acc.Id)
	if err != nil {
		return "", err
	}
	return token, err
}

func (a* AccountUseCases) Authenticate(token string) (string, error) {
	return a.Auth.UserIdByToken(token)
}

func (a* AccountUseCases) Logout() () {
	panic("not implemented method")
}

func validateLogin(login string) error {
	loginLength := 0
	for _, r := range login {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ErrInvalidLoginString
		}
		loginLength++
	}
	if loginLength < minLoginLength {
		return ErrTooShortLogin
	}
	if loginLength > maxLoginLength {
		return ErrTooLongLogin
	}
	return nil
}

func validatePassword(password string) error {
	passwordLength := 0
	for _, r := range password {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return ErrInvalidPasswordString
		}
		passwordLength++
	}
	if passwordLength < minPasswordLength {
		return ErrTooShortPassword
	}
	if passwordLength > maxPasswordLength {
		return ErrTooLongPassword
	}
	return nil
}