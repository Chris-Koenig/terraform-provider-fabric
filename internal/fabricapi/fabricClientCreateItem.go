package fabricapi

import (
	"fmt"
)

// Define a generic CreateItem method
func CreateItem[TCreate any, TRead any](itemToCreate TCreate, apiObjectName string, c FabricClient) (*TRead, error) {

	var err error
	result := new(TRead) // Create a new instance of the generic TRead type

	client, err := c.prepRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request for CreateItem: %v", err)
	}

	url := "/v1/" + apiObjectName
	resp, err := client.SetResult(result).
		SetBody(itemToCreate).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to create item: %v %s", err, url)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to create item: %v %s", resp.Error(), url)
	}

	return result, nil
}
