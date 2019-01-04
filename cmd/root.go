package cmd

import (
	"errors"
	"fmt"

	"os"
	"strings"

	"github.com/jwisard/goos"
	homedir "github.com/mitchellh/go-homedir"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile = ""

var viperCfg *viper.Viper

var client goos.OSClient

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
	rootCmd.PersistentFlags().StringP("domain", "d", "", "OpenStack authentication domain. If not provided, will be pulled from the OS_USER_DOMAIN_NAME env variable. Default is 'Default'")
	rootCmd.PersistentFlags().StringP("cloud", "c", "", "OpenStack cloud name (e.g. bob, alice, watchtower). If not provided, will be pulled from the OS_AUTH_URL env variable")
	rootCmd.PersistentFlags().StringP("tenant", "t", "", "OpenStack tenant name (project) for which the cmd will be executed. If not provided, will be pulled from the OS_PROJECT_NAME env variable")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display verbose logging")

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

	v, err := initializeViper(rootCmd, cfgFile)

	if err != nil {
		logger.Fatal("Failed to initialize configuration (Viper): " + err.Error())
	}

	viperCfg = v
}

func initializeViper(cmd *cobra.Command, file string) (*viper.Viper, error) {

	// Note that Viper maintains its data globally, but that makes testing
	// difficult.  Instead, we'll create a new Viper instance and maintain that.
	v := viper.New()

	if file != "" {
		// Use config file from the flag.
		v.SetConfigFile(file)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		v.SetConfigName(".cloudconfig")
		// Search in the current working directory for the ".cloudconfig" file
		v.AddConfigPath(".")
		// ...or search in home directory
		v.AddConfigPath(home)
	}

	v.BindPFlag("user", cmd.PersistentFlags().Lookup("user"))
	v.BindPFlag("password", cmd.PersistentFlags().Lookup("password"))
	v.BindPFlag("domain", cmd.PersistentFlags().Lookup("domain"))
	v.BindPFlag("cloud", cmd.PersistentFlags().Lookup("cloud"))
	v.BindPFlag("tenant", cmd.PersistentFlags().Lookup("tenant"))
	v.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))
	v.BindEnv("user", "OS_USERNAME")
	v.BindEnv("password", "OS_PASSWORD")
	v.BindEnv("domain", "OS_USER_DOMAIN_NAME")
	v.BindEnv("tenant", "OS_PROJECT_NAME")
	v.BindEnv("OS_AUTH_URL")

	err := v.ReadInConfig()

	// If a config file is found, read it in.
	if err == nil {
		logger.Info("Using cloud config file:", viper.ConfigFileUsed())
	}

	return v, err
}

func initializeProvider(cmd *cobra.Command, args []string) {

	cfg, err := validatedAuthConfig(viperCfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Collected AuthConfig: " + fmt.Sprintf("%+v", cfg))

	osClient, err := goos.CreateOSClient(cfg)

	if err != nil {
		logger.Fatal("Failed to establish an authenticated OpenStack client " + err.Error())
	}

	client = osClient
}

func validatedAuthConfig(v *viper.Viper) (*goos.AuthConfig, error) {

	user := v.GetString("user")

	if user == "" {
		return &goos.AuthConfig{}, errors.New("No userId was provided")
	}

	pass := v.GetString("password")

	if pass == "" {
		return &goos.AuthConfig{}, errors.New("No password was provided")
	}

	authDomain, err := validAuthDomain(v)
	if err != nil {
		return &goos.AuthConfig{}, err
	}

	authURL, err := validAuthURL(v)
	if err != nil {
		return &goos.AuthConfig{}, err
	}

	ten := v.GetString("tenant")

	if ten == "" {
		return &goos.AuthConfig{}, errors.New("No tenant name was provided")
	}

	return &goos.AuthConfig{
		User:       user,
		Password:   pass,
		AuthURL:    authURL,
		AuthDomain: authDomain,
		TenantName: ten,
	}, nil
}

func validAuthURL(v *viper.Viper) (string, error) {

	// cloud contains the value provided by the user via a flag
	if cloud := v.GetString("cloud"); cloud != "" {
		all := v.AllSettings()

		_, ok := all[cloud]

		if !ok {
			return "", errors.New("unknown cloud name: " + cloud + ". valid values are " + strings.Join(viper.GetStringSlice("clouds"), ","))
		}

		cloudConfig := v.GetStringMapString(cloud)

		return cloudConfig["authurl"], nil
	}

	// otherwise, extract the authURL from defaults and environment
	if authurl := v.GetString("OS_AUTH_URL"); authurl != "" {
		return authurl, nil
	}

	return "", errors.New("unable to determine cloud configuration from environment (no OS_AUTH_URL set). set it, or use the --cloud setting")
}

func validAuthDomain(v *viper.Viper) (string, error) {

	// authDomain contains the value provided by the user, either by flag or env variable
	if authDomain := v.GetString("domain"); authDomain != "" {

		// now validate the input
		validDomains := v.GetStringSlice("authdomains")

		for _, dom := range validDomains {
			if authDomain == dom {
				return authDomain, nil
			}
		}

		return "", errors.New(authDomain + " is not a valid domain. valid values are " + strings.Join(validDomains, ","))
	}

	// no domain specified, return the default
	return "Default", nil
}
