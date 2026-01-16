package cmd

import (
	"context"
	"fmt"
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
	Short: "command line tool to receive JSON of SBOM Report from stdin",
	Long: `Send Aqua Security Trivy Operator's SBOM Report to OWASP Dependency-Track.

two ways to use:

1. command line tool to receive JSON of SBOM Report from stdin

  $ kubectl get sbom hoge -o json | sbomreport-to-dependencytrac

2. http server that receives JSON of SBOM Report from Trivy Operator webhook

  $ sbomreport-to-dependencytrack server --port 8080

Templates with the SBOM Report as a variable can be used for the following items to be registered in the Dependency-Track.

  * Project Name
  * Project Version
  * Project Tags
  
  $ kubectl get sbom hoge -o json | sbomreport-to-dependencytrack \
      --base-url http://127.0.0.1:8081/ \
      --api-key 1234567890 \
      --project-name "[[.sbomReport.report.artifact.repository]]" \                  # e.g. "library/alpine"
      --project-version "[[.sbomReport.report.artifact.tag]]" \                      # e.g. "3.13.5"
      --project-tags "tag1,kube_namespace:[[.sbomReport.metadata.namespace]]" # e.g. ["tag1", "kube_namespace:default"]

  For template, go-template and sprig functions can be used.
  The delimiter of template is "[[" "]]". This is to avoid conflicts with other tools such as Helm.

Environment variables can be used instead of command line arguments, which may be useful when running on Kubernetes.

  $ kubectl get sbom hoge -o json | \
      DT_BASE_URL=http://127.0.0.1:8081 \
      DT_API_KEY=1234567890 \
      DT_PROJECT_NAME="[[.sbomReport.report.artifact.repository]]" \
      DT_PROJECT_VERSION="[[.sbomReport.report.artifact.tag]]" \
      DT_PROJECT_TAGS="tag1,kube_namespace:[[.sbomReport.metadata.namespace]]" \
      sbomreport-to-dependencytrack

Dependency-Track APK key permissions required:

  * BOM_UPLOAD 
  * PORTFOLIO_MANAGEMENT 
  * PROJECT_CREATION_UPLOAD 
  * VIEW_PORTFOLIO 
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		c := config.New(
			viper.GetString("base-url"),
			viper.GetString("api-key"),
			viper.GetString("project-name"),
			viper.GetString("project-version"),
			viper.GetStringSlice("project-tags"),
			viper.GetString("parent-name"),
			viper.GetString("parent-version"),
			viper.GetFloat64("dtrack-client-timeout"),
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

	flags.StringP(
		"base-url",
		"u",
		"http://127.0.0.1:8081/",
		"Dependency Track base URL (env: DT_BASE_URL)",
	)
	flags.StringP("api-key", "k", "", "Dependency Track API key (env: DT_API_KEY)")
	flags.StringP(
		"project-name",
		"",
		"[[.sbomReport.report.artifact.repository]]",
		"Project name template (env: DT_PROJECT_NAME)",
	)
	flags.StringP(
		"project-version",
		"",
		"[[.sbomReport.report.artifact.tag]]",
		"Project version template (env: DT_PROJECT_VERSION)",
	)
	flags.StringSliceP(
		"project-tags",
		"t",
		[]string{},
		"Project tags template (env: DT_PROJECT_TAGS (comma separated))",
	)
	flags.StringP("parent-name", "", "", "Parent project name template (env: DT_PARENT_NAME)")
	flags.StringP(
		"parent-version",
		"",
		"",
		"Parent project version template (env: DT_PARENT_VERSION)",
	)
	flags.Float64P("dtrack-client-timeout", "", 10, "Dependency Track client timeout seconds")
	flags.Float64P(
		"sbom-upload-timeout-sec",
		"",
		30,
		"Seconds to timeout waiting for completion of SBOM upload of Dependency Track",
	)
	flags.Float64P(
		"sbom-upload-check-interval-sec",
		"",
		1,
		"Interval seconds to check for completion of SBOM upload of Dependency Track",
	)
	flags.StringP(
		"sbom-delete-action",
		"",
		"ignore",
		"Action to perform when a SBOMReport deletion is reported via webhook. (ignore|deactivate|delete) (env: DT_SBOM_DELETE_ACTION)",
	)

	if err := viper.BindPFlags(flags); err != nil {
		panic(fmt.Errorf("bind flags to viper: %w", err))
	}
}

func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	return rootCmd.Execute()
}
