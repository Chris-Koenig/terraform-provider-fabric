package fabricapi

import (
	"fmt"
)

func GetItem[TItemResponseModel any](itemId string, itemName string, c FabricClient) (*TItemResponseModel, error) {

	var err error
	item := new(TItemResponseModel)

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for GetGroups: %v", err)
	}

	url := fmt.Sprintf("/v1/%s/%s", itemName, itemId)
	resp, err := client.SetResult(item).Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get workspace: %v", resp.Error())
	}

	return item, nil
}
