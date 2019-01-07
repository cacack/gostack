package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulule/deepcopier"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
)

func TestGetFlavor(t *testing.T) {

	assert := assert.New(t)

	output1 := flavors.Flavor{
		ID:         "d8e4ff70-f3d6-49ef-872c-65455735005a",
		Disk:       50,
		RAM:        37880,
		Name:       "m1.xlarge",
		RxTxFactor: 1,
		Swap:       0,
		VCPUs:      24,
		IsPublic:   true,
		Ephemeral:  0,
	}

	expected1 := &flavors.Flavor{}
	deepcopier.Copy(output1).To(expected1)

	cases := []struct {
		id       string
		name     string
		output   *flavors.Flavor
		expected *flavors.Flavor
		err      error
	}{
		{
			// Matching ID
			id:       "d8e4ff70-f3d6-49ef-872c-65455735005a",
			name:     "",
			output:   &output1,
			expected: expected1,
			err:      nil,
		},
		{
			// not a matching ID
			id:       "not-a-matching-id",
			name:     "",
			output:   &output1,
			expected: expected1,
			err:      nil,
		},
		{
			// matching name
			id:       "",
			name:     "m1.xlarge",
			output:   &output1,
			expected: expected1,
			err:      nil,
		},
		{
			// not a matching name (case sensitive)
			id:       "",
			name:     "M1.XLarge",
			output:   &output1,
			expected: expected1,
			err:      nil,
		},
		{
			// both inputs are empty
			id:       "",
			name:     "",
			output:   &flavors.Flavor{},
			expected: nil,
			err:      errors.New("not nil"),
		},
	}

	for _, c := range cases {
		mockClient := new(MockOSClient)
		if c.id != "" {
			mockClient.On("RetrieveFlavorByID", c.id).Return(c.output, c.err)
		} else if c.name != "" {
			mockClient.On("RetrieveFlavorByName", c.name).Return(c.output, c.err)
		}

		// actual invocation of get command method
		flavor, err := getFlavor(mockClient, c.id, c.name)

		assert.Equal(c.err != nil, err != nil, "getFlavor's returned error didn't match the expected bahavior")

		assert.Equal(c.expected, flavor, "getFlavor should have returned the same Flavor but didn't")

		mockClient.AssertExpectations(t)
	}
}

func TestGetImage(t *testing.T) {

	assert := assert.New(t)

	output1 := images.Image{
		ID:       "5a562b14-49a9-425d-af57-089b466ed226",
		Created:  "2016-10-21T13:11:06Z",
		MinDisk:  0,
		MinRAM:   0,
		Name:     "Ubuntu 16.04.1 LTS xenial (cloudimg)",
		Progress: 100,
		Status:   "ACTIVE",
		Updated:  "2016-10-21T13:14:02Z",
		Metadata: make(map[string]interface{}),
	}

	expected1 := &images.Image{}
	deepcopier.Copy(output1).To(expected1)

	cases := []struct {
		id       string
		name     string
		output   *images.Image
		expected *images.Image
		err      error
	}{
		{
			// Matching ID
			id:       "5a562b14-49a9-425d-af57-089b466ed226",
			name:     "",
			output:   &output1,
			expected: expected1,
			err:      nil,
		},
		{
			// not a matching ID
			id:       "not-a-matching-id",
			name:     "",
			output:   &images.Image{},
			expected: nil,
			err:      errors.New("not nil"),
		},
		{
			// matching name
			id:       "",
			name:     "Ubuntu 16.04.1 LTS xenial (cloudimg)",
			output:   &output1,
			expected: expected1,
			err:      nil,
		},
		{
			// not a matching name (case sensitive)
			id:       "",
			name:     "ubuntu 16.04.1 lts xenial (cloudimg)",
			output:   &images.Image{},
			expected: nil,
			err:      errors.New("not nil"),
		},
		{
			// both inputs are empty
			id:       "",
			name:     "",
			output:   &images.Image{},
			expected: nil,
			err:      errors.New("not nil"),
		},
	}

	for _, c := range cases {
		mockClient := new(MockOSClient)
		if c.id != "" {
			mockClient.On("RetrieveImageByID", c.id).Return(c.output, c.err)
		} else if c.name != "" {
			mockClient.On("RetrieveImageByName", c.name).Return(c.output, c.err)
		}

		// actual invocation of get command method
		flavor, err := getImage(mockClient, c.id, c.name)

		assert.Equal(c.err != nil, err != nil, "getImage's returned error didn't match the expected bahavior")

		assert.Equal(c.expected, flavor, "getImage should have returned the same Image but didn't")

		mockClient.AssertExpectations(t)
	}
}
