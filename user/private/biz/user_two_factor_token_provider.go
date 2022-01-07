package biz

type UserTwoFactorTokenProvider interface {
	// Name of this provider
	Name() string
	Generate(purpose string, user *User) (token string, err error)
	Validate(purpose string, token string, user *User) (bool, error)
	CanGenerate(user *User) (bool, error)
}

type UserTwoFactorTokenProviderFactory interface {
	Resolve(name string) UserTwoFactorTokenProvider
}
