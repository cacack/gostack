// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list (flavors|images|instances)",
	Short: "List objects of type flavor, image, or instance",
}

var flavorsCmd = &cobra.Command{
	Use:     "flavors",
	Short:   "lists the flavors that are available within a tenant",
	Aliases: []string{"flavor"},
	Run: func(cmd *cobra.Command, args []string) {
		listFlavors(cmd, args)
	},
}

var imagesCmd = &cobra.Command{
	Use:     "images",
	Short:   "lists the images that are available within a tenant",
	Aliases: []string{"image"},
	Run: func(cmd *cobra.Command, args []string) {
		listImages(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(flavorsCmd)
	listCmd.AddCommand(imagesCmd)

	imagesCmd.Flags().BoolP("all", "a", false, "if set, all fields of the image will be printed")
}

func listFlavors(cmd *cobra.Command, args []string) {
	allFlavors, err := retrieveFlavors(provider)

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, flavor := range allFlavors {
		fmt.Printf("%+v\n", flavor)
	}
}

func listImages(cmd *cobra.Command, args []string) {
	allImages, err := retrieveImages(provider)

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, image := range allImages {
		fmt.Printf("%+v\n", image)
	}
}

func retrieveFlavors(provider *gophercloud.ProviderClient) ([]flavors.Flavor, error) {

	compute, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return nil, err
	}

	listOpts := flavors.ListOpts{
		Limit:      20,
		AccessType: flavors.PublicAccess,
	}

	allPages, err := flavors.ListDetail(compute, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		return nil, err
	}

	return allFlavors, nil
}

func retrieveImages(provider *gophercloud.ProviderClient) ([]images.Image, error) {

	compute, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return nil, err
	}

	listOpts := images.ListOpts{
		Limit: 20,
	}

	allPages, err := images.ListDetail(compute, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return nil, err
	}

	return allImages, nil
}
