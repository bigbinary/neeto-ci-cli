package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"
	"github.com/renderedtext/sem/config"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		RunInit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

const semaphore_yaml_template = `
version: "v1.0"
name: My first pipeline
semaphore_image: standard
blocks:
  - name: "Stage 1"
    build:
      prologue:
        commands:
          - checkout
      epilogue:
        commands:
          - echo "Yay job finished"
          - echo $SEMAPHORE_JOB_RESULT
      env_vars:
        - name: VAR1
          value: Environment Variable 1
        - name: PI
          value: "3.14159"
      jobs:
      - name: Just ls
        commands:
          - pwd
          - echo "test"
          - ls /etc

      - name: List files
        commands:
          - echo "First env var -> $VAR1"
          - echo "My files:"
          - ls -lah

  - name: "Stage 2"
    build:
      jobs:
      - name: Echo job
        commands:
          - checkout
          - pwd
          - echo $SEMAPHORE_PIPELINE_ID
          - echo "Hello from $SEMAPHORE_JOB_ID"
`

func RunInit(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		utils.Fail("not a git repository")
	}

	if _, err := os.Stat(".semaphore/semaphore.yml"); err == nil {
		utils.Fail(".semaphore/semaphore.yml is already present in the repository")
	}

	createInitialSemaphoreYaml()

	repo_url, err := gitconfig.OriginURL()

	utils.CheckWithMessage(err, "failed to extract remote from git configuration")

	name := constructProjectName(repo_url)
	host := config.GetHost()

	project_url := registerProjectOnSemaphore(name, host, repo_url)

	fmt.Printf("Project is created. You can find it at %s.\n", project_url)
	fmt.Println("")
	fmt.Printf("To run your first pipeline execute:\n")
	fmt.Println("")
	fmt.Printf("  git add .semaphore/semaphore.yml && git commit -m \"First pipeline\" && git push\n")
	fmt.Println("")
}

func constructProjectName(repo_url string) string {
	re := regexp.MustCompile(`git\@github\.com:.*\/(.*).git`)
	match := re.FindStringSubmatch(repo_url)

	if len(match) < 2 {
		utils.Fail("unrecognized git remote format")
	}

	return match[1]
}

func createInitialSemaphoreYaml() {
	if _, err := os.Stat(".semaphore"); err != nil {
		err := os.Mkdir(".semaphore", 0755)

		utils.Check(err)
	}

	err := ioutil.WriteFile(".semaphore/semaphore.yml", []byte(semaphore_yaml_template), 0755)

	utils.CheckWithMessage(err, "failed to create .semaphore/semaphore.yml")
}

const project_template = `
apiVersion: v1alpha
kind: Project
metadata:
  name: %s
spec:
  repository:
    url: "%s"
`

func registerProjectOnSemaphore(name string, host string, repo_url string) string {
	c := client.FromConfig()

	project, err := yaml.YAMLToJSON([]byte(fmt.Sprintf(project_template, name, repo_url)))

	utils.CheckWithMessage(err, "connecting project to Semaphore failed")

	body, status, err := c.Post("projects", project)

	utils.CheckWithMessage(err, "connecting project to Semaphore failed")

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return fmt.Sprintf("https://%s/projects/%s", host, name)
}
