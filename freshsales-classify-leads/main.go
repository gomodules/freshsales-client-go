package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	freshsalesclient "gomodules.xyz/freshsales-client-go"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"kmodules.xyz/client-go/logs"
)

func main() {
	var (
		filename string
	)
	var rootCmd = &cobra.Command{
		Use:   "freshsales-classify-leads",
		Short: "Classify Freshsales leads",
		RunE: func(cmd *cobra.Command, args []string) error {
			return classify(filename)
		},
	}
	flags := rootCmd.Flags()

	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&filename, "in", filename, "Path to input leads.json file")

	logs.ParseFlags()

	utilruntime.Must(rootCmd.Execute())
}

type Product string

const (
	ProductKubeDB  Product = "kubedb"
	ProductStash   Product = "stash"
	ProductKubeVault   Product = "kubevault"
	ProductKubeform   Product = "kubeform"
	ProductVoyager Product = "voyager"
	ProductUnknown Product = "unknown"
)

type Criterion struct {
	Product  Product
	Keywords []string
}

var criteria = []Criterion{
	{
		Product:  ProductStash,
		Keywords: []string{"backup", "stash"},
	},
	{
		Product:  ProductKubeDB,
		Keywords: []string{"database", "db", "kubedb", "rds", "postgres", "postgresql", "mysql", "elasticsearch", "mariadb", "redis", "memcached", "mongodb"},
	},
	{
		Product:  ProductVoyager,
		Keywords: []string{"voyager", "ingress", "haproxy", "nginx"},
	},
	{
		Product:  ProductKubeVault,
		Keywords: []string{"vault", "kubevault"},
	},
	{
		Product:  ProductKubeform,
		Keywords: []string{"kubeform", "terraform"},
	},
}

func classify(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var leads []freshsalesclient.Lead
	err = json.Unmarshal(data, &leads)
	if err != nil {
		return err
	}

	groups := map[Product][]freshsalesclient.Lead{}
	for _, l := range leads {
		q := strings.ToLower(l.CustomField.KubernetesSetup) + " " + strings.ToLower(l.CustomField.CalendlyMeetingAgenda)

		prod := ProductUnknown
		for _, c := range criteria {
			if strings.EqualFold(l.CustomField.Interest, string(c.Product)) {
				prod = c.Product
			} else {
				for _, w := range c.Keywords {
					if strings.Contains(q, w) {
						prod = c.Product
						break
					}
				}
			}

			if prod != ProductUnknown {
				break
			}
		}
		groups[prod] = append(groups[prod], l)
	}

	dir := filepath.Dir(filename)
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	for prod, ls := range groups {
		sort.Slice(ls, func(i, j int) bool {
			return ls[i].Email < ls[j].Email
		})

		data, err := json.MarshalIndent(ls, "", "  ")
		if err != nil {
			return err
		}
		f := filepath.Join(dir, fmt.Sprintf("%s_%s.json", prod, base))
		err = ioutil.WriteFile(f, data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
