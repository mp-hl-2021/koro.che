package usecases

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
}

func (AccountUseCases) CreateAccount(login string, password string) (Account, error) {
	panic("not implemented method")
}

func (AccountUseCases) GetAccountById(id string) (Account, error) {
	panic("not implemented method")
}

func (AccountUseCases) LoginToAccount(login string, password string) (string, error) {
	panic("not implemented method")
}

func (AccountUseCases) Authenticate(token string) (string, error) {
	panic("not implemented method")
}

func (AccountUseCases) Logout() () {
	panic("not implemented method")
}