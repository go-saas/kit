package common

import "fmt"

type Requirement struct {
	Resource Resource
	Action   Action
	Subject  Subject
	Info     string
}

func NewRequirement(resource Resource, action Action, subject Subject, info string) Requirement {
	return Requirement{
		Resource: resource,
		Action:   action,
		Subject:  subject,
		Info:     info,
	}
}

func (r *Requirement) GetFriendlyString() string {
	if r.Info != "" {
		return r.Info
	}
	return fmt.Sprintf("%s do not have permission to %s %s %s", r.Subject.GetName(), r.Action.GetIdentity(), r.Resource.GetNamespace(), r.Resource.GetIdentity())
}

const (
	AuthenticationRequirement string = "Authentication required"
)
