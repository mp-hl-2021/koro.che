package linkstorage

import (
	"math/rand"
	"sync"
)

type Memory struct {
	linkByKey  map[string]string
	StatsByKey map[string]uint64
	userToLinksKeys map[string]map[string]bool
	mu         *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		linkByKey:  make(map[string]string),
		StatsByKey: make(map[string]uint64),
		userToLinksKeys: make(map[string]map[string]bool),
		mu:         &sync.Mutex{},
	}
}

const allowedLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const linkSize = 6

func RandString() string {
	b := make([]byte, linkSize)
	for i := range b {
		b[i] = allowedLetters[rand.Intn(len(allowedLetters))]
	}
	return string(b)
}

func (m *Memory) CreateShortLink(link string, userId string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var shortLink string
	for {
		shortLink = RandString()
		if m.linkByKey[shortLink] == "" {
			m.linkByKey[shortLink] = link
			m.StatsByKey[shortLink] = 0
			break
		}
	}
	if userId != "" {
		m.userToLinksKeys[userId][shortLink] = true
	}
	return shortLink, nil
}

func (m *Memory) GetLinkByKey(key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var link = m.linkByKey[key]
	var err error = nil
	if link == "" {
		err = ErrNotExist
	}
	return link, err
}

func (m *Memory) MakeRedirect(key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var link = m.linkByKey[key]
	var err error = nil
	//todo refactor mb use GetLinkByKey
	if link == "" {
		err = ErrNotExist
	}
	if err!= nil {
		return "", err
	}
	m.StatsByKey[key] += 1
	return link, err
}

func (m *Memory) DeleteLink(key string, userId string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var link = m.linkByKey[key]
	var err error = nil
	if link == "" {
		err = ErrNotExist
	}
	if err != nil {
		return "", err
	}
	delete(m.linkByKey, key)
	delete(m.StatsByKey, key)
	if userId != "" {
		delete(m.userToLinksKeys[userId], link)
	}
	return link, err
}

func (m *Memory) GetUserLinks(userId string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var links, ok = m.userToLinksKeys[userId]
	if !ok {
		return []string{}, ErrNotExist
	}
	keys := make([]string, 0, len(links))
	for k, _ := range links {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m *Memory) GetLinkStat(link string) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var useCounter, ok = m.StatsByKey[link]
	if !ok {
		return 0, ErrNotExist
	}
	return useCounter, nil
}

func (m *Memory) CreateUserLinksStorage(userId string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userToLinksKeys[userId] = map[string]bool{}
	return "", nil
}