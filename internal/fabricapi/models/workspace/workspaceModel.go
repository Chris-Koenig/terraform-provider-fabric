package workspace

// Group represents a response detailing a Power BI group (workspace).
type WorkspaceModel struct {
	CapacityId            string `json:"capacityId"`            // The capacity ID
	DataflowStorageId     string `json:"dataflowStorageId"`     // The Power BI dataflow storage account ID
	Id                    string `json:"id"`                    // The workspace ID
	IsOnDedicatedCapacity bool   `json:"isOnDedicatedCapacity"` // Whether the group (workspace) is assigned to a dedicated capacity
	IsReadOnly            bool   `json:"isReadOnly"`            // Whether the group (workspace) is read-only
	LogAnalyticsWorkspace string `json:"logAnalyticsWorkspace"` // The Log Analytics workspace assigned to the group. This is returned only when retrieving a single group.
	Name                  string `json:"name"`                  // The group name
	// DefaultDatasetStorageFormat DefaultDatasetStorageFormat `json:"defaultDatasetStorageFormat"` // The default dataset storage format in the workspace. Returned only when isOnDedicatedCapacity is true
}
