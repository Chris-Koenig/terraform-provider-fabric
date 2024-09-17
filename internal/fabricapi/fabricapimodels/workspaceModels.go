package fabricapimodels

// Group represents a response detailing a Power BI group (workspace).
type WorkspaceReadModel struct {
	// CapacityId            string `json:"capacityId"`            // The capacity ID
	// DataflowStorageId     string `json:"dataflowStorageId"`     // The Power BI dataflow storage account ID
	Id string `json:"id"` // The workspace ID
	// IsOnDedicatedCapacity bool   `json:"isOnDedicatedCapacity"` // Whether the group (workspace) is assigned to a dedicated capacity
	IsReadOnly bool `json:"isReadOnly"` // Whether the group (workspace) is read-only
	// LogAnalyticsWorkspace string `json:"logAnalyticsWorkspace"` // The Log Analytics workspace assigned to the group. This is returned only when retrieving a single group.
	DisplayName string `json:"displayName"` // The group name
	// DefaultDatasetStorageFormat DefaultDatasetStorageFormat `json:"defaultDatasetStorageFormat"` // The default dataset storage format in the workspace. Returned only when isOnDedicatedCapacity is true
}

type WorkspaceCreateModel struct {
	DisplayName string `json:"displayname"` // The name
	// 	CapacityId  string `json:"capacityId"`  // The name
	// 	Description string `json:"description"` // The name
}

type WorkspaceDeleteModel struct {
	Id string `json:"id"` // The group (workspace) id
}

type WorkspaceUpdateModel struct {
	Name string `json:"name"` // The group (workspace) name
}

// Groups represents a response detailing a list of Power BI groups (workspaces).
type WorkspacesReadModel struct {
	ODataContext string               `json:"@odata.context"` // The OData context
	ODataCount   int                  `json:"@odata.count"`   // The OData count
	Value        []WorkspaceReadModel `json:"value"`          // The list of groups
}
