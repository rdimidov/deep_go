package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type callable = func() any
type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	deps map[string]callable
}

func NewContainer() *Container {
	return &Container{deps: make(map[string]callable)}
}

func (c *Container) RegisterType(name string, constructor any) {
	constructorFunc, ok := constructor.(callable)
	if !ok {
		panic(fmt.Sprintf("dependency %s is not valid", name))
	}
	c.deps[name] = constructorFunc
}

func (c *Container) RegisterSingletonType(name string, constructor any) {
	constructorFunc, ok := constructor.(callable)
	if !ok {
		panic(fmt.Sprintf("dependency %s is not valid", name))
	}
	var (
		instance any
		once     sync.Once
	)

	c.deps[name] = func() any {
		once.Do(func() {
			instance = constructorFunc()
		})
		return instance
	}
}

func (c *Container) Resolve(name string) (any, error) {
	constructor, ok := c.deps[name]
	if !ok {
		return nil, fmt.Errorf("dependency %s has not been registred", name)
	}
	return constructor(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() any {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() any {
		return &MessageService{}
	})
	container.RegisterSingletonType("OtherMessageService", func() any {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	OtherMessageService1, err := container.Resolve("OtherMessageService")
	assert.NoError(t, err)
	OtherMessageService2, err := container.Resolve("OtherMessageService")
	assert.NoError(t, err)

	oms1 := OtherMessageService1.(*MessageService)
	oms2 := OtherMessageService2.(*MessageService)
	assert.True(t, oms1 == oms2)
	assert.NotNil(t, oms1)
	assert.NotNil(t, oms2)

}
