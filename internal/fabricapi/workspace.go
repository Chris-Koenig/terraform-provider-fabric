package fabricapi

import (
	"fmt"
	"strconv"
)

// CreateGroup creates a new group.
// https://learn.microsoft.com/en-us/rest/api/power-bi/groups/create-group
func (c FabricClient) CreateWorkspace(displayName string) (*WorkspaceReadModel, error) {

	var err error
	ws := &WorkspaceReadModel{}
	ws1 := &WorkspaceCreateModel{DisplayName: displayName}

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for CreateGroup: %v", err)
	}
	baseURL := "https://api.fabric.microsoft.com"
	url := fmt.Sprintf("%s/v1/workspaces", baseURL)
	resp, err := client.SetResult(ws).
		// SetQueryParam("workspaceV2", "True").
		SetBody(ws1).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to create group: %v", resp.Error())
	}

	return ws, nil
}

// DeleteGroup deletes a group by its ID.
// https://learn.microsoft.com/en-us/rest/api/power-bi/groups/delete-group
func (c FabricClient) DeleteWorkspace(workspaceId string) error {
	// DELETE https://api.fabric.com/v1.0/myorg/groups/{groupId}
	var err error

	client, err := c.prepRequest()
	if err != nil {
		return fmt.Errorf("failed to prepare the request for DeleteGroup: %v", err)
	}

	resp, err := client.Delete(fmt.Sprintf("/v1.0/myorg/groups/%s", workspaceId))
	if err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to delete group: %v", resp.Error())
	}

	return nil
}

// GetGroup retrieves a group by its ID.
func (c *FabricClient) GetWorkspace(workspaceId string) (*WorkspaceReadModel, error) {
	// GET https://api.fabric.microsoft.com/v1/workspaces/{workspaceid}

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

// GetGroups retrieves a list of groups.
// https://learn.microsoft.com/en-us/rest/api/power-bi/groups/get-groups
func (c *FabricClient) GetWorkspaces(filter string, top int, skip int) (*WorkspacesReadModel, error) {
	// GET https://api.powerbi.com/v1.0/myorg/groups

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

	resp, err := client.SetResult(&groups).Get("/v1.0/myorg/groups")
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get groups: %v", resp.Error())
	}

	return groups, nil
}

// // UpdateGroup updates a specified workspace.
// // https://learn.microsoft.com/en-us/rest/api/power-bi/groups/update-group
// func (c *FabricClient) UpdateWorkspace(groupId string, updateGroupRequest *workspace.WorkspaceModel) error {
// 	// PATCH https://api.powerbi.com/v1.0/myorg/groups/{groupId}

// 	var err error

// 	client, err := c.prepRequest()
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare the request for GetGroups: %v", err)
// 	}

// 	body := updateGroupRequest.Validate()

// 	resp, err := client.
// 		SetBody(body).
// 		Patch(fmt.Sprintf("/v1.0/myorg/groups/%s", groupId))
// 	if err != nil {
// 		return fmt.Errorf("failed to update groups: %v", err)
// 	}

// 	if resp.IsError() {
// 		return fmt.Errorf("failed to get groups: [%v] %s", resp.StatusCode(), resp.String())
// 	}

// 	return nil
// }
