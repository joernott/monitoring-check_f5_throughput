package cmd

import (
	"github.com/joernott/check_f5_throughput/checker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile           string
	statsFile         string
	warningThreshold  string
	criticalThreshold string
	host              string
	port              int
	community         string

	rootCmd = &cobra.Command{
		Use:   "check_f5_throughput",
		Short: "An Icinga2/Nagios plugin to monitor F5 throughput",
		Long: `This Icinga 2 /Nagios plugin monitors the throughput of
F5 BIGIP loadbalancers via SNMP. The procedure is outlined in
https://support.f5.com/csp/article/K50309321.`,
		Run: func(cmd *cobra.Command, args []string) {
			checker.Check(
				viper.GetString("host"),
				uint16(viper.GetInt("port")),
				viper.GetString("community"),
				viper.GetString("warningThreshold"),
				viper.GetString("criticalThreshold"),
				viper.GetString("file"))
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "/etc/icinga2/check_f5_throughput.yaml", "config file (default is /etc/icinga2/check_f5_throughput.yaml)")
	rootCmd.PersistentFlags().StringVarP(&warningThreshold, "warning", "w", "", "warning range")
	rootCmd.PersistentFlags().StringVarP(&criticalThreshold, "critical", "c", "", "critical range")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "host/ip address of the laod balancer (default 127.0.0.1)")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "P", 161, "SNMP port (default 161)")
	rootCmd.PersistentFlags().StringVarP(&community, "community", "C", "public", "SNMP community (default public")
	rootCmd.PersistentFlags().StringVarP(&statsFile, "file", "F", " /var/lib/icinga2/check_f5_throughput_stats.json", "statistics file (default /var/lib/icinga2/check_f5_throughput_stats.json)")

	viper.BindPFlag("warning", rootCmd.PersistentFlags().Lookup("warning"))
	viper.BindPFlag("critical", rootCmd.PersistentFlags().Lookup("critical"))
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("community", rootCmd.PersistentFlags().Lookup("community"))
	viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", 161)
	viper.SetDefault("community", "public")
	viper.SetDefault("warning", "")
	viper.SetDefault("critical", "")
	viper.SetDefault("file", "/var/lib/icinga2/check_f5_throughput_stats.json")
}

func initConfig() {
	if cfgFile == "" {
		cfgFile = "/etc/icinga2/check_f5_throughput.yaml"
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
