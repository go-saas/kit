package util

import (
	"github.com/go-saas/kit/pkg/authz/authz"
	pb "github.com/go-saas/kit/user/api/permission/v1"
)

func MapPermissionBeanToPb(bean authz.PermissionBean, t *pb.Permission) {
	t.Subject = bean.Subject
	t.Namespace = bean.Namespace
	t.Resource = bean.Resource
	t.Action = bean.Action
	t.TenantId = bean.TenantID
	t.Effect = MapAuthEffect2PbEffect(bean.Effect)
}
func MapPbPermissionToBean(t *pb.Permission, bean *authz.PermissionBean) {
	bean.Subject = t.Subject
	bean.Namespace = t.Namespace
	bean.Resource = t.Resource
	bean.Action = t.Action
	bean.TenantID = t.TenantId
	bean.Effect = MapPbEffect2AuthEffect(t.Effect)

}

func MapPbEffect2AuthEffect(eff pb.Effect) authz.Effect {
	effect := authz.EffectUnknown
	switch eff {
	case pb.Effect_GRANT:
		effect = authz.EffectGrant
		break
	case pb.Effect_FORBIDDEN:
		effect = authz.EffectForbidden
		break
	}
	return effect
}
func MapAuthEffect2PbEffect(eff authz.Effect) pb.Effect {
	switch eff {
	case authz.EffectUnknown:
		return pb.Effect_UNKNOWN
	case authz.EffectGrant:
		return pb.Effect_GRANT
	case authz.EffectForbidden:
		return pb.Effect_FORBIDDEN

	}
	return pb.Effect_UNKNOWN
}
