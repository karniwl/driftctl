// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package repository

import (
	armnetwork "github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	mock "github.com/stretchr/testify/mock"
)

// MockNetworkRepository is an autogenerated mock type for the NetworkRepository type
type MockNetworkRepository struct {
	mock.Mock
}

// ListAllFirewalls provides a mock function with given fields:
func (_m *MockNetworkRepository) ListAllFirewalls() ([]*armnetwork.AzureFirewall, error) {
	ret := _m.Called()

	var r0 []*armnetwork.AzureFirewall
	if rf, ok := ret.Get(0).(func() []*armnetwork.AzureFirewall); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*armnetwork.AzureFirewall)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAllPublicIPAddresses provides a mock function with given fields:
func (_m *MockNetworkRepository) ListAllPublicIPAddresses() ([]*armnetwork.PublicIPAddress, error) {
	ret := _m.Called()

	var r0 []*armnetwork.PublicIPAddress
	if rf, ok := ret.Get(0).(func() []*armnetwork.PublicIPAddress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*armnetwork.PublicIPAddress)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAllRouteTables provides a mock function with given fields:
func (_m *MockNetworkRepository) ListAllRouteTables() ([]*armnetwork.RouteTable, error) {
	ret := _m.Called()

	var r0 []*armnetwork.RouteTable
	if rf, ok := ret.Get(0).(func() []*armnetwork.RouteTable); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*armnetwork.RouteTable)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAllSecurityGroups provides a mock function with given fields:
func (_m *MockNetworkRepository) ListAllSecurityGroups() ([]*armnetwork.NetworkSecurityGroup, error) {
	ret := _m.Called()

	var r0 []*armnetwork.NetworkSecurityGroup
	if rf, ok := ret.Get(0).(func() []*armnetwork.NetworkSecurityGroup); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*armnetwork.NetworkSecurityGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAllSubnets provides a mock function with given fields: virtualNetwork
func (_m *MockNetworkRepository) ListAllSubnets(virtualNetwork *armnetwork.VirtualNetwork) ([]*armnetwork.Subnet, error) {
	ret := _m.Called(virtualNetwork)

	var r0 []*armnetwork.Subnet
	if rf, ok := ret.Get(0).(func(*armnetwork.VirtualNetwork) []*armnetwork.Subnet); ok {
		r0 = rf(virtualNetwork)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*armnetwork.Subnet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*armnetwork.VirtualNetwork) error); ok {
		r1 = rf(virtualNetwork)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAllVirtualNetworks provides a mock function with given fields:
func (_m *MockNetworkRepository) ListAllVirtualNetworks() ([]*armnetwork.VirtualNetwork, error) {
	ret := _m.Called()

	var r0 []*armnetwork.VirtualNetwork
	if rf, ok := ret.Get(0).(func() []*armnetwork.VirtualNetwork); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*armnetwork.VirtualNetwork)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
