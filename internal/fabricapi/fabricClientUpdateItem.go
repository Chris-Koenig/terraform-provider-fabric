package fabricapi

import (
	"fmt"
)

func UpdateItem[TUpdate any](itemIdToUpdate string, apiObjectName string, itemUpdateRequestModel TUpdate, workspaceId string, c FabricClient) error {

	var err error
	var url string
	client, err := c.prepRequest()

	if err != nil {
		return fmt.Errorf("failed to prepare the request update workspace: %v", err)
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
		return fmt.Errorf("failed to update workspace: %v %s", err, url)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to update workspace: [%v] %s", resp.StatusCode(), resp.String())
	}

	return nil
}
