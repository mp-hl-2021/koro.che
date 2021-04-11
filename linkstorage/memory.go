package linkstorage

import (
	"math/rand"
	"sync"
)

type Memory struct {
	linkByKey  map[string]string
	StatsNyKey map[string]uint64
	mu         *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		linkByKey:  make(map[string]string),
		StatsNyKey: make(map[string]uint64),
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

func (m *Memory) CreateShortLink(link string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var shortLink string
	for {
		shortLink = RandString()
		if m.linkByKey[shortLink] == "" {
			m.linkByKey[shortLink] = link
			m.StatsNyKey[shortLink] = 0
			break
		}
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
	m.StatsNyKey[key] += 1
	return link, err
}