package cmd

import (
	"github.com/stretchr/testify/mock"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
)

// MockOSClient is a mock object for OSClient invocations and is used for testing.
type MockOSClient struct {
	mock.Mock
}

func (m *MockOSClient) RetrieveFlavors() ([]flavors.Flavor, error) {
	args := m.Called()
	return args.Get(0).([]flavors.Flavor), args.Error(1)
}

func (m *MockOSClient) RetrieveFlavorByID(flavorID string) (*flavors.Flavor, error) {
	args := m.Called(flavorID)
	return args.Get(0).(*flavors.Flavor), args.Error(1)
}

func (m *MockOSClient) RetrieveFlavorByName(flavorName string) (*flavors.Flavor, error) {
	args := m.Called(flavorName)
	return args.Get(0).(*flavors.Flavor), args.Error(1)
}

func (m *MockOSClient) RetrieveImages() ([]images.Image, error) {
	args := m.Called()
	return args.Get(0).([]images.Image), args.Error(1)
}

func (m *MockOSClient) RetrieveImageByID(imageID string) (*images.Image, error) {
	args := m.Called(imageID)
	return args.Get(0).(*images.Image), args.Error(1)
}

func (m *MockOSClient) RetrieveImageByName(imageName string) (*images.Image, error) {
	args := m.Called(imageName)
	return args.Get(0).(*images.Image), args.Error(1)
}
