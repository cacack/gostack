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
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/jwisard/goos"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var getCmd = &cobra.Command{
	Use:   "get (flavor|image|instance)",
	Short: "Get an object of type flavor, image, or instance",
}

var flavorCmd = &cobra.Command{
	Use:   "flavor (--id <id>|--name <name>)",
	Short: "gets a flavor by id or name",
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetString("id")
		name, _ := cmd.Flags().GetString("name")

		flavor, err := getFlavor(client, id, name)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("%+v\n", flavor)
		}
	},
}

var imageCmd = &cobra.Command{
	Use:   "image (--id <id>|--name <name>)",
	Short: "gets an image by id or name",
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetString("id")
		name, _ := cmd.Flags().GetString("name")

		image, err := getImage(client, id, name)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("%+v\n", image)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(flavorCmd)
	getCmd.AddCommand(imageCmd)

	flavorCmd.Flags().String("id", "", "the ID of the desired flavor")
	flavorCmd.Flags().String("name", "", "the name of the desired flavor")

	imageCmd.Flags().String("id", "", "the ID of the desired image")
	imageCmd.Flags().String("name", "", "the name of the desired image")
}

func getFlavor(osClient goos.OSClient, id string, name string) (*flavors.Flavor, error) {

	if id != "" {

		flavor, err := osClient.RetrieveFlavorByID(id)

		if err != nil {
			return nil, err
		}

		return flavor, nil
	}

	if name != "" {

		flavor, err := osClient.RetrieveFlavorByName(name)

		if err != nil {
			return nil, err
		}

		return flavor, nil
	}

	return nil, errors.New("one of flavor id or name is required")
}

func getImage(osClient goos.OSClient, id string, name string) (*images.Image, error) {

	if id != "" {

		image, err := osClient.RetrieveImageByID(id)

		if err != nil {
			return nil, err
		}

		return image, nil
	}

	if name != "" {

		image, err := osClient.RetrieveImageByName(name)

		if err != nil {
			return nil, err
		}

		return image, nil
	}

	return nil, errors.New("one of flavor id or name is required")
}
