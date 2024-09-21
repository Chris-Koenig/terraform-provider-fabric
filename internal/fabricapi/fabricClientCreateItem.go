package fabricapi

import (
	"fmt"
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"
)

func CreateItem[TCreate any, TRead any](itemToCreate TCreate, apiObjectName string, workspaceId string, c FabricClient) (*TRead, error) {

	var err error
	var url string
	result := new(TRead) // Create a new instance of the generic TRead type

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for CreateItem: %v", err)
	}

	url = "/v1/workspaces/" + workspaceId + "/" + apiObjectName

	resp, err := client.SetResult(result).
		SetBody(itemToCreate).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to create item: %v %s", err, url)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to create item: %v %s", resp, url)
	}

	return result, nil
}

// The Creation of the Workspace is special (Capacity handling), for this reason
// we handle the creation in a sepperate call.
func CreateWorkspace(workspaceToCreate fabricClientModels.WorkspaceCreateRequestModel, fabricCapacityId *string, c FabricClient) (*fabricClientModels.WorkspaceReadModel, error) {

	var err error
	var url string
	var workspaceCreated = new(fabricClientModels.WorkspaceReadModel)

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for CreateItem: %v", err)
	}

	// Create the Workspace
	urlWorkspaceCreate := "/v1/workspaces"
	resp, err := client.SetResult(workspaceCreated).
		SetBody(workspaceToCreate).
		Post(urlWorkspaceCreate)

	if err != nil {
		return nil, fmt.Errorf("failed to create item: %v %s", err, url)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to create item: %v %s", resp, url)
	}

	if fabricCapacityId == nil || *fabricCapacityId == "" {
		workspaceCreated.FabricCapacityId = new(string)

	} else if fabricCapacityId != nil && *fabricCapacityId != "" {
		// Assign the Created Workspace to a capacity
		urlAssignCreatedWorkspaceToCapacity := "/v1/workspaces/" + workspaceCreated.Id + "/assignToCapacity"

		capacityAssignModel := new(fabricClientModels.AssignCapacityRequestModel)
		capacityAssignModel.FabricCapacityId = fabricCapacityId

		resp, err := client.
			SetBody(capacityAssignModel).
			Post(urlAssignCreatedWorkspaceToCapacity)

		if err != nil {
			return nil, fmt.Errorf("failed to create item: %v %s", err, url)
		}

		if resp.IsError() {
			return nil, fmt.Errorf("failed to create item: %v %s", resp, url)
		}

		workspaceCreated.FabricCapacityId = fabricCapacityId
	}
	return workspaceCreated, nil
}
