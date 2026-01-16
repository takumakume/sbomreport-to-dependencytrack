package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/takumakume/sbomreport-to-dependencytrack/config"
	"github.com/takumakume/sbomreport-to-dependencytrack/server"
	"github.com/takumakume/sbomreport-to-dependencytrack/uploader"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "http server that receives JSON of SBOM Report from Trivy Operator webhook",
	Long: `http server that receives JSON of SBOM Report from Trivy Operator webhook

	$ sbomreport-to-dependencytrack server --port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.New(
			viper.GetString("base-url"),
			viper.GetString("api-key"),
			viper.GetString("project-name"),
			viper.GetString("project-version"),
			viper.GetStringSlice("project-tags"),
			viper.GetString("parent-name"),
			viper.GetString("parent-version"),
			viper.GetFloat64("dtrack-client-timeout-sec"),
			viper.GetFloat64("sbom-upload-timeout-sec"),
			viper.GetFloat64("sbom-upload-check-interval-sec"),
			viper.GetString("sbom-delete-action"),
		)
		if err := c.Validate(); err != nil {
			return err
		}

		u, err := uploader.New(c)
		if err != nil {
			return err
		}

		s := server.New(u, viper.GetInt("port"))

		return s.Run()
	},
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	flags := serverCmd.Flags()
	flags.IntP("port", "p", 8080, "Port number to listen http to receive webhook from trivy-operator (env: DT_PORT)")
	viper.BindPFlag("port", flags.Lookup("port"))

	rootCmd.AddCommand(serverCmd)
}
