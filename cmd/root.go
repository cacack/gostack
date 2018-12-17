package cmd

import (
	"errors"
	"fmt"

	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"

	homedir "github.com/mitchellh/go-homedir"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile = ""
var verbose bool

type Cloud struct {
	name    string
	authUrl string
	domain  string
	region  string
}

type ValidatedConfig struct {
	cloud    Cloud
	user     string
	password string
	tenant   string
}

var provider *gophercloud.ProviderClient

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "gostack",
	Short:            "Demonstrating OpenStack API useing Go",
	PersistentPreRun: initializeProvider,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "cloudconfig", "", "OpenStack cloud config file (default is $HOME/.cloudconfig.yaml, then current dir)")
	rootCmd.PersistentFlags().StringP("user", "u", "", "OpenStack userId. If not provided, will be pulled from the OS_USERNAME env variable")
	rootCmd.PersistentFlags().StringP("password", "p", "", "OpenStack password. If not provided, will be pulled from the OS_PASSWORD env variable")
	rootCmd.PersistentFlags().StringP("cloud", "c", "", "OpenStack cloud name (e.g. bob, alice, watchtower). If not provided, will be pulled from the OS_AUTH_URL env variable")
	rootCmd.PersistentFlags().StringP("tenant", "t", "", "OpenStack tenant name (project) for which the cmd will be executed. If not provided, will be pulled from the OS_PROJECT_NAME env variable")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display verbose logging")
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("cloud", rootCmd.PersistentFlags().Lookup("cloud"))
	viper.BindPFlag("tenant", rootCmd.PersistentFlags().Lookup("tenant"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	logger.SetFormatter(&logger.TextFormatter{})

	if viper.GetBool("verbose") {
		logger.SetLevel(logger.DebugLevel)
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetLevel(logger.WarnLevel)
		logger.SetOutput(os.Stderr)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigName(".cloudconfig")
		// Search in the current working directory for the ".cloudconfig" file
		viper.AddConfigPath(".")
		// ...or search in home directory
		viper.AddConfigPath(home)
	}

	viper.BindEnv("user", "OS_USERNAME")
	viper.BindEnv("password", "OS_PASSWORD")
	viper.BindEnv("tenant", "OS_PROJECT_NAME")
	viper.BindEnv("OS_AUTH_URL")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Info("Using cloud config file:", viper.ConfigFileUsed())
	}
}

func initializeProvider(cmd *cobra.Command, args []string) {

	// This grabs a reference to the global Viper object.
	// This allows us to write more unit-testable code, rather than
	// only referring to the global Viper via the viper package.
	v := viper.GetViper()

	cfg, err := validatedConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Collected config: " + fmt.Sprintf("%+v", cfg))

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: cfg.cloud.authUrl,
		Username:         cfg.user,
		Password:         cfg.password,
		TenantName:       cfg.tenant,
		DomainName:       cfg.cloud.domain,
	}

	authClient, err := openstack.AuthenticatedClient(authOpts)

	if err != nil {
		logger.Fatal("Failed to establish an authenticated OpenStack client " + err.Error())
	}

	provider = authClient
}

func validatedConfig(v *viper.Viper) (ValidatedConfig, error) {

	cloud, err := cloudConfig(v)
	if err != nil {
		return ValidatedConfig{}, err
	}

	user := v.GetString("user")

	if user == "" {
		return ValidatedConfig{}, errors.New("No userId was provided")
	}

	pass := v.GetString("password")

	if pass == "" {
		return ValidatedConfig{}, errors.New("No password was provided")
	}

	ten := v.GetString("tenant")

	if ten == "" {
		return ValidatedConfig{}, errors.New("No tenant name was provided")
	}

	return ValidatedConfig{cloud, user, pass, ten}, nil
}

func cloudConfig(v *viper.Viper) (Cloud, error) {

	if cloud := viper.GetString("cloud"); cloud != "" {
		all := viper.AllSettings()

		_, ok := all[cloud]

		if !ok {
			return Cloud{}, errors.New("unknown cloud name: " + cloud + ". valid values are " + strings.Join(viper.GetStringSlice("clouds"), ","))
		}

		cloudConfig := viper.GetStringMapString(cloud)

		c := Cloud{}
		c.name = cloud
		c.domain = cloudConfig["domain"]
		c.authUrl = cloudConfig["authurl"]
		c.region = cloudConfig["region"]

		return c, nil
	}

	// otherwise, build the cloudConfig from defaults and environment
	if authurl := viper.GetString("OS_AUTH_URL"); authurl != "" {
		return Cloud{"", authurl, "Default", "RegionOne"}, nil
	}

	return Cloud{}, errors.New("unable to determine cloud configuration from environment (no OS_AUTH_URL set). set it, or use the --cloud setting")
}
