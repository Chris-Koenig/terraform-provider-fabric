package fabricClientModels

type Principal struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type RoleAssignmentCreateRequestModel struct {
	Principal Principal `json:"principal"`
	Role      string    `json:"role"`
}

type RoleAssignmentReadModel struct {
	Id        string    `json:"id"`
	Principal Principal `json:"principal"`
	Role      string    `json:"role"`
}

type RoleAssignmentUpdateRequestModel struct {
	Principal Principal `json:"principal"`
	Role      string    `json:"role"`
}
