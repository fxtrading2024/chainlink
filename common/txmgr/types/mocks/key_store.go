// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/smartcontractkit/chainlink/v2/common/types"
)

// KeyStore is an autogenerated mock type for the KeyStore type
type KeyStore[ADDR types.Hashable, CHAIN_ID types.ID, SEQ types.Sequence] struct {
	mock.Mock
}

// CheckEnabled provides a mock function with given fields: address, chainID
func (_m *KeyStore[ADDR, CHAIN_ID, SEQ]) CheckEnabled(address ADDR, chainID CHAIN_ID) error {
	ret := _m.Called(address, chainID)

	if len(ret) == 0 {
		panic("no return value specified for CheckEnabled")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(ADDR, CHAIN_ID) error); ok {
		r0 = rf(address, chainID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EnabledAddressesForChain provides a mock function with given fields: chainId
func (_m *KeyStore[ADDR, CHAIN_ID, SEQ]) EnabledAddressesForChain(chainId CHAIN_ID) ([]ADDR, error) {
	ret := _m.Called(chainId)

	if len(ret) == 0 {
		panic("no return value specified for EnabledAddressesForChain")
	}

	var r0 []ADDR
	var r1 error
	if rf, ok := ret.Get(0).(func(CHAIN_ID) ([]ADDR, error)); ok {
		return rf(chainId)
	}
	if rf, ok := ret.Get(0).(func(CHAIN_ID) []ADDR); ok {
		r0 = rf(chainId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ADDR)
		}
	}

	if rf, ok := ret.Get(1).(func(CHAIN_ID) error); ok {
		r1 = rf(chainId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeToKeyChanges provides a mock function with given fields:
func (_m *KeyStore[ADDR, CHAIN_ID, SEQ]) SubscribeToKeyChanges() (chan struct{}, func()) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SubscribeToKeyChanges")
	}

	var r0 chan struct{}
	var r1 func()
	if rf, ok := ret.Get(0).(func() (chan struct{}, func())); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan struct{})
		}
	}

	if rf, ok := ret.Get(1).(func() func()); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	return r0, r1
}

// NewKeyStore creates a new instance of KeyStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKeyStore[ADDR types.Hashable, CHAIN_ID types.ID, SEQ types.Sequence](t interface {
	mock.TestingT
	Cleanup(func())
}) *KeyStore[ADDR, CHAIN_ID, SEQ] {
	mock := &KeyStore[ADDR, CHAIN_ID, SEQ]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
