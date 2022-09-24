package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	freshsalesclient "gomodules.xyz/freshsales-client-go"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"kmodules.xyz/client-go/logs"
)

func main() {
	var (
		freshsalesHost     = "https://appscode.freshsales.io"
		freshsalesAPIToken = os.Getenv("CRM_API_TOKEN")
		filename           string
	)
	var rootCmd = &cobra.Command{
		Use:   "freshsales-export-leads",
		Short: "Export Freshsales leads",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := freshsalesclient.New(freshsalesHost, freshsalesAPIToken)
			return export(client, filename)
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&freshsalesHost, "freshsales.host", freshsalesHost, "Freshsales host url")
	flags.StringVar(&freshsalesAPIToken, "freshsales.token", freshsalesAPIToken, "Freshsales api token")
	flags.StringVar(&filename, "out", filename, "Path to output json file")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
}

func export(c *freshsalesclient.Client, filename string) error {
	leads, err := c.ListAllLeads()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(leads, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
