package fabricapi

// Group represents a response detailing a Power BI group (workspace).
type WorkspaceReadModel struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type WorkspaceCreateModel struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type WorkspaceDeleteModel struct {
	Id string `json:"id"`
}

type WorkspaceUpdateModel struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

// Groups represents a response detailing a list of Power BI groups (workspaces).
type WorkspacesReadModel struct {
	ODataContext string               `json:"@odata.context"` // The OData context
	ODataCount   int                  `json:"@odata.count"`   // The OData count
	Value        []WorkspaceReadModel `json:"value"`          // The list of groups
}
