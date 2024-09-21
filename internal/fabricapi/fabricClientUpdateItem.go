package fabricapi

import (
	"fmt"
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"
)

func UpdateItem[TUpdate any](itemIdToUpdate string, apiObjectName string, itemUpdateRequestModel TUpdate, workspaceId string, c FabricClient) error {

	var err error
	var url string
	client, err := c.prepRequest()

	if err != nil {
		return fmt.Errorf("failed to prepare the request update "+apiObjectName+": %v", err)
	}

	if workspaceId == "" {
		url = "/v1/" + apiObjectName + "/" + itemIdToUpdate
	} else if workspaceId != "" {
		url = "/v1/workspaces/" + workspaceId + "/" + apiObjectName + "/" + itemIdToUpdate
	}

	resp, err := client.
		SetBody(itemUpdateRequestModel).
		Patch(url)

	if err != nil {
		return fmt.Errorf("failed to update "+apiObjectName+": %v %s", err, url)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to update "+apiObjectName+": [%v] %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// The update of the Workspace is special (capacity handling), for this reason
// we handle the creation in a sepperate function.
func UpdateWorkspace(workspaceUpdateRequestModel fabricClientModels.WorkspaceUpdateRequestModel,
	workspaceId string,
	capacityPlanValue string,
	capacityStateValue string,
	c FabricClient) error {

	var err error
	var url string
	client, err := c.prepRequest()

	if err != nil {
		return fmt.Errorf("failed to prepare the request update workspace: %v", err)
	}

	url = "/v1/workspaces/" + workspaceId

	// Update the workspace
	resp, err := client.
		SetBody(workspaceUpdateRequestModel).
		Patch(url)

	if err != nil {
		return fmt.Errorf("failed to update workspace: %v %s", err, url)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to update workspace: [%v] %s", resp.StatusCode(), resp.String())
	}

	// update the capacity

	if capacityPlanValue == capacityStateValue {
		// do nothing
	} else if capacityPlanValue == "" {
		// the capacity is removed from the workspace terraform configuration
		urlUnAssignCapacity := "/v1/workspaces/" + workspaceId + "/unassignFromCapacity"

		resp, err := client.Post(urlUnAssignCapacity)

		if err != nil {
			return fmt.Errorf("failed to create item: %v %s", err, urlUnAssignCapacity)
		}

		if resp.IsError() {
			return fmt.Errorf("failed to create item: %v %s", resp, urlUnAssignCapacity)
		}

	} else if capacityPlanValue != "" {
		// capacity is added to the workspace terraform configuration
		urlAssignCapacity := "/v1/workspaces/" + workspaceId + "/assignToCapacity"

		capacityAssignModel := new(fabricClientModels.AssignCapacityRequestModel)
		capacityAssignModel.FabricCapacityId = &capacityPlanValue

		resp, err := client.
			SetBody(capacityAssignModel).
			Post(urlAssignCapacity)

		if err != nil {
			return fmt.Errorf("failed to create item: %v %s", err, urlAssignCapacity)
		}

		if resp.IsError() {
			return fmt.Errorf("failed to create item: %v %s", resp, urlAssignCapacity)
		}

	}

	return nil
}
