package fabricapi

import (
	"fmt"
)

func (c FabricClient) CreateWorkspace(workspaceToCreate WorkspaceCreateModel) (*WorkspaceReadModel, error) {

	var err error
	ws := &WorkspaceReadModel{}

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for CreateWorkspace: %v", err)
	}

	url := "/v1/workspaces"
	resp, err := client.SetResult(ws).
		SetBody(workspaceToCreate).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %v %s", err, url)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to create workspace: %v %s", resp.Error(), url)
	}

	return ws, nil
}

func (c FabricClient) DeleteWorkspace(workspaceTodelete WorkspaceDeleteModel) error {

	var err error

	client, err := c.prepRequest()
	if err != nil {
		return fmt.Errorf("failed to prepare the request for DeleteGroup: %v", err)
	}

	resp, err := client.Delete(fmt.Sprintf("/v1/workspaces/%s", workspaceTodelete.Id))
	if err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to delete group: %v", resp.Error())
	}

	return nil
}

func (c *FabricClient) GetWorkspace(workspaceId string) (*WorkspaceReadModel, error) {

	var err error
	workspace := &WorkspaceReadModel{}

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for GetGroups: %v", err)
	}

	url := fmt.Sprintf("/v1/workspaces/%s", workspaceId)
	resp, err := client.SetResult(workspace).Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get workspace: %v", resp.Error())
	}

	return workspace, nil
}

func (c *FabricClient) UpdateWorkspace(workspaceIdToUpdate string, workspaceToUpdate WorkspaceUpdateModel) error {

	var err error

	client, err := c.prepRequest()

	if err != nil {
		return fmt.Errorf("failed to prepare the request update workspace: %v", err)
	}

	body := workspaceToUpdate

	// baseURL := "https://api.fabric.microsoft.com"
	url := fmt.Sprintf("/v1/workspaces/%s", workspaceIdToUpdate)
	resp, err := client.
		SetBody(body).
		Patch(url)

	if err != nil {
		return fmt.Errorf("failed to update workspace: %v %s", err, url)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to update workspace: [%v] %s", resp.StatusCode(), resp.String()+" URL: "+url)
	}

	return nil
}
