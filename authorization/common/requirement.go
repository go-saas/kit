package common


type Requirement interface {
	GetRequireName() string
}

type RequirementStr string

func (r RequirementStr) GetRequireName() string {
	return string(r)
}

const(
	AuthenticationRequirement RequirementStr = "Authentication"
)