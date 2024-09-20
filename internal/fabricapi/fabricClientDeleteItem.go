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
	finalUrl := fmt.Sprintf("%s/%s", url, idToDelete)
	resp, err := client.Delete(finalUrl)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %v", apiObjectName, err)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to delete %s: %v", apiObjectName, finalUrl)
	}

	return nil
}
