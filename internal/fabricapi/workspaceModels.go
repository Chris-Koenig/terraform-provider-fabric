package fabricapi

// Group represents a response detailing a Power BI group (workspace).
type WorkspaceReadModel struct {
	Id          string `json:"id"`          // The workspace ID
	DisplayName string `json:"displayName"` // The group name
	CapacityId  string `json:"capacityId"`  // The capacity ID
	Description string `json:"description"` // The description
}

type WorkspaceCreateModel struct {
	DisplayName string `json:"displayname"`
	CapacityId  string `json:"capacityId"`
	Description string `json:"description"`
}

type WorkspaceDeleteModel struct {
	Id string `json:"id"`
}

type WorkspaceUpdateModel struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayname"`
	Description string `json:"description"`
}

// Groups represents a response detailing a list of Power BI groups (workspaces).
type WorkspacesReadModel struct {
	ODataContext string               `json:"@odata.context"` // The OData context
	ODataCount   int                  `json:"@odata.count"`   // The OData count
	Value        []WorkspaceReadModel `json:"value"`          // The list of groups
}
