package fabricapi

import (
	"fmt"
)

func GetItem[TItemResponseModel any](itemId string, apiObjectName string, workspaceId string, c FabricClient) (*TItemResponseModel, error) {

	var err error
	item := new(TItemResponseModel)
	var url string

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for GetGroups: %v", err)
	}
	if workspaceId == "" {
		url = fmt.Sprintf("/v1/%s/%s", apiObjectName, itemId)
	} else if workspaceId != "" {
		url = "/v1/workspaces/" + workspaceId + "/" + apiObjectName + "/" + itemId
	}
	resp, err := client.SetResult(item).Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %v", apiObjectName, err)
	}

	if resp.IsError() {
		message := url + " " + workspaceId
		return nil, fmt.Errorf("failed to get %s: %v", apiObjectName, message)
	}

	return item, nil
}
