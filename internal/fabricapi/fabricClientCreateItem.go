package fabricapi

import (
	"fmt"
)

func CreateItem[TCreate any, TRead any](itemToCreate TCreate, apiObjectName string, workspaceId string, c FabricClient) (*TRead, error) {

	var err error
	var url string
	result := new(TRead) // Create a new instance of the generic TRead type

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for CreateItem: %v", err)
	}

	if workspaceId == "" {
		url = "/v1/" + apiObjectName
	} else if workspaceId != "" {
		url = "/v1/workspaces/" + workspaceId + "/" + apiObjectName
	}

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
