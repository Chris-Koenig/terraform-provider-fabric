package fabricClientModels

type WorkspaceCreateRequestModel struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type WorkspaceReadModel struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type WorkspaceDeleteRequestModel struct {
	Id string `json:"id"`
}

type WorkspaceUpdateRequestModel struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}
