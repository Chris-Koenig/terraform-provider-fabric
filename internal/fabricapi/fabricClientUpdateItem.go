package fabricapi

import (
	"fmt"
)

func UpdateItem[TUpdate any](itemIdToUpdate string, itemNameToUpdate string, itemUpdateRequestModel TUpdate, c FabricClient) error {

	var err error

	client, err := c.prepRequest()

	if err != nil {
		return fmt.Errorf("failed to prepare the request update workspace: %v", err)
	}

	url := fmt.Sprintf("/v1/%s/%s", itemNameToUpdate, itemIdToUpdate)
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
