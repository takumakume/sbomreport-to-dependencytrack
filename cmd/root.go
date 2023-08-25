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
	Short: "",
	Long:  ``,
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
	flags.StringP("project-name", "", "[[.input.report.artifact.repository]]", "Project name template (env: DT_PROJECT_NAME)")
	flags.StringP("project-version", "", "[[.input.report.artifact.tag]]", "Dependency Track API key (env: DT_PROJECT_VERSION)")
	flags.StringSliceP("project-tags", "t", []string{}, "Dependency Track project tags template (env: DT_PROJECT_TAGS comma separated)")
	flags.IntP("timeout", "", 10, "Dependency Track timeout seconds (env: DT_TIMEOUT)")

	viper.BindPFlag("base-url", flags.Lookup("base-url"))
	viper.BindPFlag("api-key", flags.Lookup("api-key"))
	viper.BindPFlag("project-name", flags.Lookup("project-name"))
	viper.BindPFlag("project-version", flags.Lookup("project-version"))
	viper.BindPFlag("project-tags", flags.Lookup("project-tags"))
	viper.BindPFlag("timeout", flags.Lookup("timeout"))
}

func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	return rootCmd.Execute()
}
