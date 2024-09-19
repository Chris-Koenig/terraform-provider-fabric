package fabricClientModels

type Principal struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type WorkspaceRoleAssignmentCreateRequestModel struct {
	Principal Principal `json:"principal"`
	Role      string    `json:"role"`
}

type WorkspaceRoleAssignmentReadModel struct {
	Id        string    `json:"id"`
	Principal Principal `json:"principal"`
	Role      string    `json:"role"`
}

type WorkspaceRoleAssignmentUpdateRequestModel struct {
	Principal Principal `json:"principal"`
	Role      string    `json:"role"`
}
