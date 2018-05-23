package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/leosunmo/consularis/pkg/config"
	c "github.com/leosunmo/consularis/pkg/controller"
)

var namespace string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "consularis",
	Short: "Manage Consul kv in Kubernetes",
	Long:  `Manage Consul Key/Value pairs in Kubernetes using CRDs`,

	Run: func(cmd *cobra.Command, args []string) {
		config := config.New()
		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		if err == nil {
			config.Kubeconfig = kubeconfig
		} else {
			log.WithField("err", err.Error()).Fatal("Error building config")
		}
		masterurl, err := cmd.Flags().GetString("master")
		if err == nil {
			config.MasterURL = masterurl
		} else {
			log.WithField("err", err.Error()).Fatal("Error building config")
		}
		namespace, err := cmd.Flags().GetString("namespace")
		if err == nil {
			if len(namespace) > 0 {
				config.Namespace = namespace
			}
		} else {
			log.WithField("err", err.Error()).Fatal("Error building config")
		}
		consul, err := cmd.Flags().GetString("consul")
		if err == nil {
			if len(consul) > 0 {
				config.Consul = consul
			} else {
				config.Consul = "consul"
			}
		} else {
			log.WithField("err", err.Error()).Fatal("Error building config")
		}
		port, err := cmd.Flags().GetString("port")
		if err == nil {
			if len(port) > 0 {
				config.Port = port
			} else {
				config.Port = "8500"
			}
		} else {
			log.WithField("err", err.Error()).Fatal("Error building config")
		}

		c.Run(config)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().StringP("kubeconfig", "k", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	RootCmd.Flags().StringP("master", "m", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	RootCmd.Flags().StringP("namespace", "n", "", "Namespace to watch (default is all)")
	RootCmd.Flags().StringP("consul", "c", "", "Consul to manage (default \"consul\")")
	RootCmd.Flags().StringP("port", "p", "", "Consul port (default \"8500\")")
}
