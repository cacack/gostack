package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulule/deepcopier"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
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
			id:       "d8e4ff70-f3d6-49ef-872c-65455735005a",
			name:     "",
			output:   &output1,
			expected: expected1,
			err:      nil,
		},
		{
			id:       "",
			name:     "m1.xlarge",
			output:   &output1,
			expected: expected1,
			err:      nil,
		},
		{
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

		if c.err != nil {
			assert.NotNil(err, "getFlavor should have returned an error but didn't")
		} else {
			assert.Nil(err, "getFlavor should not have returned an error but did")
		}

		assert.Equal(c.expected, flavor, "getFlavor should have returned the same Flavor but didn't")

		mockClient.AssertExpectations(t)
	}
}
