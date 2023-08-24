package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := newConfig(
			viper.GetString("base-url"),
			viper.GetString("api-key"),
			viper.GetString("project-name"),
			viper.GetString("project-version"),
			viper.GetStringSlice("project-tags"),
			viper.GetInt("timeout"),
		)

		return server(config, viper.GetInt("port"))
	},
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	flags := serverCmd.Flags()
	flags.IntP("port", "p", 80, "Port number to listen http to receive webhook from trivy-operator (env: DT_PORT)")
	viper.BindPFlag("port", flags.Lookup("port"))

	rootCmd.AddCommand(serverCmd)
}

func server(c *config, port int) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ResponseWriter.Write(w, []byte(`healthy`))
			return
		case "POST":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := upload(c, body); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			http.ResponseWriter.Write(w, []byte(`ok`))
			return
		default:
			log.Println("Invalid request method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}

	return nil
}
