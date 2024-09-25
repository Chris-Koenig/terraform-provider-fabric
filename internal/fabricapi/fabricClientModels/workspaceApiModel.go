package fabricClientModels

// This is a special object need, because the Fabric Capacity is not part
// of the API Request Model and will be handled in a sepperate call
// type WorkspaceCreateRequestModelWithCapacity struct {
// 	DisplayName      string
// 	Description      string
// 	FabricCapacityId string
// }

type WorkspaceCreateRequestModel struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type WorkspaceReadModel struct {
	Id               string `json:"id"`
	DisplayName      string `json:"displayName"`
	Description      string `json:"description"`
	FabricCapacityId string `json:"capacityId"`
}

type WorkspaceUpdateRequestModel struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}
