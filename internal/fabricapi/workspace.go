package fabricapi

import (
	"fmt"
	"strconv"
)

func (c FabricClient) CreateWorkspace(workspaceToCreate WorkspaceCreateModel) (*WorkspaceReadModel, error) {

	var err error
	ws := &WorkspaceReadModel{}

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for CreateGroup: %v", err)
	}
	baseURL := "https://api.fabric.microsoft.com"
	url := fmt.Sprintf("%s/v1/workspaces", baseURL)
	resp, err := client.SetResult(ws).
		// SetQueryParam("workspaceV2", "True").
		SetBody(workspaceToCreate).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to create group: %v", resp.Error())
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

	baseURL := "https://api.fabric.microsoft.com"
	url := fmt.Sprintf("%s/v1/workspaces/%s", baseURL, workspaceId)
	resp, err := client.SetResult(workspace).Get(url)

	//resp, err := client.SetResult(workspace).Get(fmt.Sprintf("/v1/workspaces/%s", workspaceId))
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get workspace: %v", resp.Error())
	}

	return workspace, nil
}

func (c *FabricClient) GetWorkspaces(filter string, top int, skip int) (*WorkspacesReadModel, error) {

	var err error
	groups := &WorkspacesReadModel{}

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for GetGroups: %v", err)
	}

	if filter != "" {
		client.SetQueryParam("$filter", filter)
	}
	if top > 0 {
		client.SetQueryParam("$top", strconv.Itoa(top))
	}
	if skip > 0 {
		client.SetQueryParam("$skip", strconv.Itoa(skip))
	}

	resp, err := client.SetResult(&groups).Get("/v1/workspaces")
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get groups: %v", resp.Error())
	}

	return groups, nil
}

func (c *FabricClient) UpdateWorkspace(workspaceIdToUpdate string, workspaceToUpdate WorkspaceUpdateModel) error {

	var err error

	client, err := c.prepRequest()
	if err != nil {
		return fmt.Errorf("failed to prepare the request for GetGroups: %v", err)
	}

	body := workspaceToUpdate

	resp, err := client.
		SetBody(body).
		Patch(fmt.Sprintf("/v1/workspaces/%s", workspaceIdToUpdate))

	if err != nil {
		return fmt.Errorf("failed to update group: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to get group: [%v] %s", resp.StatusCode(), resp.String())
	}

	return nil
}
