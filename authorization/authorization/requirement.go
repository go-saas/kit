package authorization

type Requirement interface {
	GetRequiredName() string
}

type RequirementStr string

func (r RequirementStr) GetRequiredName() string {
	return string(r)
}

const (
	AuthenticationRequirement RequirementStr = "Authentication"
)
