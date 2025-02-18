package client

import (
	"errors"
	"fmt"

	models "github.com/bigbinary/neeto-ci-cli/api/models"
)

type ProjectApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewProjectV1AlphaApi() ProjectApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("cli/v1")

	return ProjectApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "projects",
		ResourceNameSingular: "project",
	}
}

func NewProjectV1AlphaApiWithCustomClient(client BaseClient) ProjectApiV1AlphaApi {
	client.SetApiVersion("cli/v1")

	return ProjectApiV1AlphaApi{
		BaseClient:           client,
		ResourceNamePlural:   "projects",
		ResourceNameSingular: "project",
	}
}

func (c *ProjectApiV1AlphaApi) ListProjects() (*models.ProjectListV1Alpha, error) {
	body, status, err := c.BaseClient.List(c.ResourceNamePlural)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to neetoCI failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewProjectListV1AlphaFromJson(body)
}

func (c *ProjectApiV1AlphaApi) GetProject(name string) (*models.ProjectV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to neetoCI failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewProjectV1AlphaFromJson(body)
}

func (c *ProjectApiV1AlphaApi) DeleteProject(name string) error {
	body, status, err := c.BaseClient.Delete(c.ResourceNamePlural, name)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return nil
}

func (c *ProjectApiV1AlphaApi) CreateProject(d *models.ProjectV1Alpha) (*models.ProjectV1Alpha, error) {
	json_body, err := d.ToJson()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to serialize object '%s'", err))
	}

	body, status, err := c.BaseClient.Post(c.ResourceNamePlural, json_body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("creating %s on neetoCI failed '%s'", c.ResourceNameSingular, err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewProjectV1AlphaFromJson(body)
}

func (c *ProjectApiV1AlphaApi) UpdateProject(d *models.ProjectV1Alpha) (*models.ProjectV1Alpha, error) {
	json_body, err := d.ToJson()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to serialize %s object '%s'", c.ResourceNameSingular, err))
	}

	identifier := d.Metadata.Id

	body, status, err := c.BaseClient.Patch(c.ResourceNamePlural, identifier, json_body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("updating %s on neetoCI failed '%s'", c.ResourceNamePlural, err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewProjectV1AlphaFromJson(body)
}
