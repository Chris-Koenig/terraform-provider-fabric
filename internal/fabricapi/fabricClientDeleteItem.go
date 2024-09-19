package fabricapi

import (
	"fmt"
)

func DeleteItem(idToDelete string, itemName string, c FabricClient) error {

	var err error

	client, err := c.prepRequest()
	if err != nil {
		return fmt.Errorf("failed to prepare the request for DeleteGroup: %v", err)
	}

	resp, err := client.Delete(fmt.Sprintf("/v1/%s/%s", itemName, idToDelete))
	if err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	if resp.IsError() {
		return fmt.Errorf("failed to delete group: %v", resp.Error())
	}

	return nil
}
