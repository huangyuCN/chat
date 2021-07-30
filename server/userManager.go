package main

import (
	"errors"
	"net"
	"sync"
)

// userManager 管理所有在线玩家
type userManager struct {
	users map[string]*user //所有玩家
	lock  *sync.Mutex
}

// NewUserManager 初始化玩家管理对象
func NewUserManager() *userManager {
	return &userManager{
		users: make(map[string]*user),
		lock:  new(sync.Mutex),
	}
}

// Register 注册一个新玩家
func (m *userManager) Register(name string, conn *net.TCPConn) (error, *user) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.users[name]; ok {
		return errors.New(UserNameRepeated), nil
	}
	user := NewUser(name, conn, m)
	m.users[name] = user
	return nil, user
}

// Unregister 注销
func (m *userManager) Unregister(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.users, name)
}
