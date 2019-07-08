/*
Copyright © 2019 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/portworx/px/pkg/plugin"
	"github.com/portworx/px/pkg/util"
	"github.com/spf13/cobra"
)

const (
	pxDefaultDir        = ".px"
	pxDefaultConfigName = "config.yml"

	Ki = 1024
	Mi = 1024 * Ki
	Gi = 1024 * Mi
	Ti = 1024 * Gi
)

var (
	cfgDir      string
	cfgFile     string
	optEndpoint string
	pm          *plugin.PluginManager

	// The $HOME/.px/plugins dir will be added at runtime
	pxPluginDefaultDirs = []string{
		"/var/lib/px/plugins",
		"/etc/pwx/plugins",
		"/opt/pwx/plugins",
		"/var/lib/porx/plugins",
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "px",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.Eprintf("%v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/"+pxDefaultDir+"/"+pxDefaultConfigName+")")
	rootCmd.PersistentFlags().StringVar(&optEndpoint, "endpoint", "", "Portworx service endpoint")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output in yaml|json|wide")
	rootCmd.PersistentFlags().Bool("show-labels", false, "Show labels in the last column of the output")
	rootCmd.PersistentFlags().StringP("selector", "l", "", "Comma separated label selector of the form 'key=value,key=value'")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Load plugins
	home, _ := homedir.Dir()
	pxPluginDefaultDirs = append(pxPluginDefaultDirs,
		path.Join(home, pxDefaultDir, "plugins"))
	pm = plugin.NewPluginManager(&plugin.PluginManagerConfig{
		PluginDirs: pxPluginDefaultDirs,
		RootCmd:    rootCmd,
	})
	pm.Load()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(cfgFile) == 0 {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			util.Eprintf("Error: %v\n", err)
			os.Exit(1)
		}

		cfgFile = path.Join(home, pxDefaultDir, pxDefaultConfigName)
	}

}

func GetConfigFile() string {
	return cfgFile
}
