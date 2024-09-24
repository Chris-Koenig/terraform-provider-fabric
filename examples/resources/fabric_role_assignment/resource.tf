resource "fabric_role_assignment" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000" # Required: Replace with the actual workspace ID
  role         = "Member"                               # Required: Role assigned to the principal, e.g., "Member" or "Admin"

  principal {
    id   = "00000000-0000-0000-0000-000000000000" # Required: Replace with the actual principal ID (user or group)
    type = "User"                                 # Required: Type of principal (e.g., "User" or "Group")
  }
}
