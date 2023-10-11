package v1

import "google.golang.org/protobuf/proto"

func (x *LoginAuthRequest) StringWithMask(mask string) string {
	ret := proto.Clone(x).(*LoginAuthRequest)
	ret.Password = mask
	return ret.String()
}

func (x *TokenRequest) StringWithMask(mask string) string {
	ret := proto.Clone(x).(*TokenRequest)
	ret.Password = mask
	return ret.String()
}

func (x *ValidatePasswordRequest) StringWithMask(mask string) string {
	ret := proto.Clone(x).(*ValidatePasswordRequest)
	ret.Password = mask
	return ret.String()
}

func (x *WebLoginAuthRequest) StringWithMask(mask string) string {
	ret := proto.Clone(x).(*WebLoginAuthRequest)
	ret.Password = mask
	return ret.String()
}

func (x *RegisterAuthRequest) StringWithMask(mask string) string {
	ret := proto.Clone(x).(*RegisterAuthRequest)
	ret.Password = mask
	ret.ConfirmPassword = mask
	return ret.String()
}
