package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"orion-cli/internal/ai"
	"orion-cli/internal/config"
	"orion-cli/internal/core"
	githubadapter "orion-cli/internal/github"
	"orion-cli/internal/gitrepo"
)

func init() {
	rootCmd.AddCommand(taskCmd)
}

var taskCmd = &cobra.Command{
	Use:   "task [gorev_metni]",
	Short: "Doğal dil görevini analiz edip GitHub issue'ları açar",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := config.LoadEnv()
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("çalışma dizini alınamadı: %w", err)
		}

		analyzer := ai.NewGeminiAdapter(env.GeminiAPIKey, env.GeminiModel, nil)
		creator := githubadapter.NewIssuesAdapter(env.GitHubToken)
		locator := gitrepo.NewDetector(wd)

		facade := core.NewTaskFacade(analyzer, creator, locator)
		createdIssues, err := facade.Execute(args[0])
		if err != nil {
			return err
		}

		fmt.Println("✅ Issue'lar başarıyla oluşturuldu:")
		for _, issue := range createdIssues {
			fmt.Printf("- #%d %s -> %s\n", issue.Number, issue.Title, issue.URL)
		}
		return nil
	},
}
