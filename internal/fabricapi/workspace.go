package fabricapi

import (
	"fmt"
	"strconv"
	"terraform-provider-fabric/internal/fabricapi/models/workspace"
)

// CreateGroup creates a new group.
// https://learn.microsoft.com/en-us/rest/api/power-bi/groups/create-group
func (c *FabricClient) CreateWorkspace(groupName string) (*workspace.WorkspaceModel, error) {
	// POST https://api.fabric.com/v1.0/myorg/groups

	var err error
	group := &workspace.WorkspaceModel{}

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for CreateGroup: %v", err)
	}

	resp, err := client.SetResult(group).
		SetQueryParam("workspaceV2", "True").
		SetBody(&workspace.WorkspaceCreateModel{Name: groupName}).
		Post("/v1.0/myorg/groups")
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to create group: %v", resp.Error())
	}

	return group, nil
}

// DeleteGroup deletes a group by its ID.
// https://learn.microsoft.com/en-us/rest/api/power-bi/groups/delete-group
func (c *FabricClient) DeleteWorkspace(workspaceId string) error {
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
// https://learn.microsoft.com/en-us/rest/api/power-bi/groups/get-group
func (c *FabricClient) GetWorkspace(groupId string) (*workspace.WorkspaceModel, error) {
	// GET https://api.powerbi.com/v1.0/myorg/groups/{groupId}

	var err error
	group := &workspace.WorkspaceModel{}

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for GetGroups: %v", err)
	}

	resp, err := client.SetResult(group).Get(fmt.Sprintf("/v1.0/myorg/groups/%s", groupId))
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get group: %v", resp.Error())
	}

	return group, nil
}

// GetGroups retrieves a list of groups.
// https://learn.microsoft.com/en-us/rest/api/power-bi/groups/get-groups
func (c *FabricClient) GetWorkspaces(filter string, top int, skip int) (*workspace.WorkspaceModel, error) {
	// GET https://api.powerbi.com/v1.0/myorg/groups

	var err error
	groups := &workspace.WorkspaceModel{}

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
