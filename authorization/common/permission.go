package common

type PermissionManagementService interface {
	EnsureGrant(resource Resource, action Action, subject Subject) error
	EnsureDisallow(resource Resource, action Action, subject Subject) error
	IsGrant(resource Resource, action Action, subject Subject) (bool, error)
}
