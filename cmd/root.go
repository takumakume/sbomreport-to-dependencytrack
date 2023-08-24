package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/takumakume/sbomreport-to-dependencytrack/dependencytrack"
)

var rootCmd = &cobra.Command{
	Use:   "sbomreport-to-dependencytrack",
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

		sbomReportJSON, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		if err := upload(config, sbomReportJSON); err != nil {
			return err
		}

		return nil
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
}

func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	return rootCmd.Execute()
}

func upload(c *config, sbomReportJSON []byte) error {
	sbomReport := map[string]interface{}{}
	err := json.Unmarshal(sbomReportJSON, &sbomReport)
	if err != nil {
		return err
	}
	t := newTemplateEngine(sbomReport)
	projectName, err := t.render(c.projectName)
	if err != nil {
		return err
	}
	projectVersion, err := t.render(c.projectVersion)
	if err != nil {
		return err
	}
	projectTags := []string{}
	for _, tag := range c.projectTags {
		tag, err := t.render(tag)
		if err != nil {
			return err
		}
		projectTags = append(projectTags, tag)
	}

	dtrackClient, err := dependencytrack.New(c.baseURL, c.apiKey, c.timeout)
	if err != nil {
		return err
	}

	bom, err := getBOM(sbomReportJSON)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := dtrackClient.UploadBOM(ctx, projectName, projectVersion, bom); err != nil {
		return err
	}

	if len(c.projectTags) > 0 {
		if err := dtrackClient.AddTagsToProject(ctx, projectName, projectVersion, projectTags); err != nil {
			return err
		}
	}

	return nil
}

func getBOM(sbomReportJSON []byte) ([]byte, error) {
	var data map[string]interface{}

	if err := json.Unmarshal(sbomReportJSON, &data); err != nil {
		return nil, err
	}

	kind, ok := data["kind"].(string)
	if !ok || kind != "SbomReport" {
		return nil, errors.New("kind is not SbomReport")
	}

	apiVersion, ok := data["apiVersion"].(string)
	if !ok {
		return nil, fmt.Errorf("apiVersion %q is not found", apiVersion)
	}

	report, ok := data["report"].(map[string]interface{})
	if !ok {
		return nil, errors.New("report is not found")
	}

	bom, ok := report["components"].(map[string]interface{})
	if !ok {
		return nil, errors.New("bom is not found")
	}

	jsonBytes, err := json.Marshal(bom)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
