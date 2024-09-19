package fabricapi

import (
	"fmt"
)

func DeleteItem(idToDelete string, apiObjectName string, workspaceId string, c FabricClient) error {

	var err error
	var url string

	client, err := c.prepRequest()
	if err != nil {
		return fmt.Errorf("failed to prepare the request for DeleteGroup: %v", err)
	}

	if workspaceId == "" {
		url = "/v1/" + apiObjectName
	} else if workspaceId != "" {
		url = "/v1/workspaces/" + workspaceId + "/" + apiObjectName
	}

	resp, err := client.Delete(fmt.Sprintf("%s/%s", url, idToDelete))
	if err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to delete group: %v", resp.Error())
	}

	return nil
}
