package workspace

// Groups represents a response detailing a list of Power BI groups (workspaces).
type WorkspacesModel struct {
	ODataContext string           `json:"@odata.context"` // The OData context
	ODataCount   int              `json:"@odata.count"`   // The OData count
	Value        []WorkspaceModel `json:"value"`          // The list of groups
}
