package cmd

import (
	"os"

	"github.com/hacker301et/wifipen/cmd/logic"
	"github.com/spf13/cobra"
)

var iface string

var rootCmd = &cobra.Command{
	Use:   "wifipen",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		w := logic.Init(iface)
		w.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&iface, "iface", "i", "", "wifi inteface used to creack wifi")
	if err := rootCmd.MarkFlagRequired("iface"); err != nil {
		return
	}
}
