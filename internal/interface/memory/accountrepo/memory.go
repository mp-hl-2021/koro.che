package accountrepo

import (
	"koro.che/internal/domain/account"
	"strconv"
	"sync"
)

type Memory struct {
	accountsById    map[string]account.Account
	accountsByLogin map[string]account.Account
	nextId          uint64
	mu              *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		accountsById:    make(map[string]account.Account),
		accountsByLogin: make(map[string]account.Account),
		mu:              &sync.Mutex{},
	}
}

func (m *Memory) CreateAccount(cred account.Credentials) (account.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.accountsByLogin[cred.Login]; ok {
		return account.Account{}, account.ErrAlreadyExist
	}
	a := account.Account{
		Id: strconv.FormatUint(m.nextId, 16),
		Credentials: cred,
	}
	m.accountsById[a.Id] = a
	m.accountsByLogin[a.Login] = a
	m.nextId++
	return a, nil
}

func (m *Memory) GetAccountById(id string) (account.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsById[id]
	if !ok {
		return a, account.ErrNotFound
	}
	return a, nil
}

func (m *Memory) GetAccountByLogin(login string) (account.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsByLogin[login]
	if !ok {
		return a, account.ErrNotFound
	}
	return a, nil
}