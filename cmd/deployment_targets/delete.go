package deployment_targets

import (
	"fmt"

	client "github.com/bigbinary/neeto-ci-cli/api/client"
	"github.com/bigbinary/neeto-ci-cli/cmd/utils"
)

func Delete(targetId string) {
	c := client.NewDeploymentTargetsV1AlphaApi()

	err := c.Delete(targetId)
	utils.Check(err)

	fmt.Printf("Deployment target '%s' deleted.\n", targetId)
}
