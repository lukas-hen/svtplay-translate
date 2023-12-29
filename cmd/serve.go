package cmd

import (
	"log"

	"github.com/lukas-hen/svtplay-translate/internal/utils"
	"github.com/lukas-hen/svtplay-translate/pkg/fileserver"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a file over http.",
	Long:  `Serve a file over http.`,
	Run: func(cmd *cobra.Command, args []string) {

		f, _ := cmd.Flags().GetString("file")
		i, _ := cmd.Flags().GetString("interface")
		ipaddr, err := utils.GetInterfaceIpv4Addr(i)
		if err != nil {
			log.Fatalf("Could not get ip from interface \"%s\"\n", i)
		}

		fileserver.Run(ipaddr, f)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringP("file", "f", "", "Path to the file you want to serve.")
	serveCmd.PersistentFlags().StringP("interface", "i", "", "The network interface you want to serve from.")
	serveCmd.MarkPersistentFlagRequired("file")
	serveCmd.MarkPersistentFlagRequired("interface")
}
