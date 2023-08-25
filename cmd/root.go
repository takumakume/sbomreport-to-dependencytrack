package cmd

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/takumakume/sbomreport-to-dependencytrack/config"
	"github.com/takumakume/sbomreport-to-dependencytrack/uploader"
)

var rootCmd = &cobra.Command{
	Use:   "sbomreport-to-dependencytrack",
	Short: "Tool to send trivy-operator's SBOM Report to Dependency Track",
	Long: `
		Tool to send trivy-operator's SBOM Report to Dependency Track,
		which can receive webhooks from Stdin and TrivyOperator and send them to Dependency Track.
		Project and tags can be generated with templates using SBOM Report values.
		
		# from stdin
		$ kubectl get sbom hoge -o json | sbomreport-to-dependencytrack

		# from webhook
		$ sbomreport-to-dependencytrack server --port 80

		# set project name, version and tags
		#  - using go template with sprig functions 
		#  - delimiter: "[[" "]]" (no conflict with helm template)
		#  - ".sbomReport" variable: the root of the SBOM Report
		$ kubectl get sbom hoge -o json | sbomreport-to-dependencytrack \
			--base_url http://localhost:8081 \
			--api-key 1234567890 \
			--project-name "[[ .sbomReport.report.artifact.repository ]]"
			--project-version "[[ .sbomReport.report.artifact.tag ]]"
			--project-tags "tag1,kube_cluster_name:production,kube_namespace:[[ .sbomReport.report.metadaga.namespace ]]"

		# set by environment variables
		$ kubectl get sbom hoge -o json | \
			DT_BASE_URL=http://localhost:8081 \
			DT_API_KEY=1234567890 \
			DT_PROJECT_NAME="[[ .sbomReport.report.artifact.repository ]]" \
			DT_PROJECT_VERSION="[[ .sbomReport.report.artifact.tag ]]" \
			DT_PROJECT_TAGS="tag1,kube_cluster_name:production,kube_namespace:[[ .sbomReport.report.metadaga.namespace ]]" \
			sbomreport-to-dependencytrack
		`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		c := config.New(
			viper.GetString("base-url"),
			viper.GetString("api-key"),
			viper.GetString("project-name"),
			viper.GetString("project-version"),
			viper.GetStringSlice("project-tags"),
		)
		if err := c.Validate(); err != nil {
			return err
		}

		u, err := uploader.New(c)
		if err != nil {
			return err
		}

		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		return u.Run(ctx, input)
	},
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("DT")

	flags := rootCmd.PersistentFlags()

	flags.StringP("base-url", "u", "http://lolcalhost:8081/", "Dependency Track base URL (env: DT_BASE_URL)")
	flags.StringP("api-key", "k", "", "Dependency Track API key (env: DT_API_KEY)")
	flags.StringP("project-name", "", "[[.sbomReport.report.artifact.repository]]", "Project name template (env: DT_PROJECT_NAME)")
	flags.StringP("project-version", "", "[[.sbomReport.report.artifact.tag]]", "Project version template (env: DT_PROJECT_VERSION)")
	flags.StringSliceP("project-tags", "t", []string{}, "Project tags template (env: DT_PROJECT_TAGS (comma separated))")

	viper.BindPFlag("base-url", flags.Lookup("base-url"))
	viper.BindPFlag("api-key", flags.Lookup("api-key"))
	viper.BindPFlag("project-name", flags.Lookup("project-name"))
	viper.BindPFlag("project-version", flags.Lookup("project-version"))
	viper.BindPFlag("project-tags", flags.Lookup("project-tags"))
}

func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	return rootCmd.Execute()
}
