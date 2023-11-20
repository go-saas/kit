package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	kapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/authz/authz"
	"github.com/go-saas/kit/pkg/localize"
	pb "github.com/go-saas/kit/user/api/permission/v1"
	"github.com/go-saas/kit/user/util"
	"github.com/go-saas/saas"
	"github.com/samber/lo"
)

// CheckForSubjects internal api for remote permission checker
func (s *PermissionService) CheckForSubjects(ctx context.Context, req *pb.CheckSubjectsPermissionRequest) (*pb.CheckSubjectsPermissionReply, error) {
	if err := kapi.ErrIfUntrusted(ctx, s.trust); err != nil {
		return nil, err
	}

	subjects := make([]authz.Subject, len(req.Subjects))
	for i, subject := range req.Subjects {
		subjects[i] = authz.SubjectStr(subject)
	}
	grantList, err := s.auth.BatchCheckForSubjects(ctx, lo.Map(req.Requirements, func(t *pb.PermissionRequirement, _ int) *authz.Requirement {
		return authz.NewRequirement(authz.NewEntityResource(t.Namespace, t.Resource), authz.ActionStr(t.Action))
	}), subjects...)
	if err != nil {
		//other error
		return nil, err
	}
	effList := lo.Map(grantList, func(t *authz.Result, _ int) pb.Effect {

		effect := pb.Effect_FORBIDDEN
		if t.Allowed {
			effect = pb.Effect_GRANT
		}
		return effect
	})

	return &pb.CheckSubjectsPermissionReply{EffectList: effList}, err
}

func (s *PermissionService) AddSubjectPermission(ctx context.Context, req *pb.AddSubjectPermissionRequest) (*pb.AddSubjectPermissionResponse, error) {
	if err := kapi.ErrIfUntrusted(ctx, s.trust); err != nil {
		return nil, err
	}

	if err := s.findAnyValidateModifyPermissionDef(ctx, req.Namespace, req.Action); err != nil {
		return nil, err
	}
	if err := s.permissionMgr.AddGrant(ctx, authz.NewEntityResource(req.Namespace, req.Resource),
		authz.ActionStr(req.Action), authz.SubjectStr(req.Subject), req.TenantId, util.MapPbEffect2AuthEffect(req.Effect)); err != nil {
		return nil, err
	}
	return &pb.AddSubjectPermissionResponse{}, nil
}

func (s *PermissionService) ListSubjectPermission(ctx context.Context, req *pb.ListSubjectPermissionRequest) (*pb.ListSubjectPermissionResponse, error) {
	if err := kapi.ErrIfUntrusted(ctx, s.trust); err != nil {
		return nil, err
	}
	subs := make([]authz.Subject, len(req.Subjects))
	for i, subject := range req.Subjects {
		subs[i] = authz.SubjectStr(subject)
	}
	acl, err := s.permissionMgr.ListAcl(ctx, subs...)
	if err != nil {
		return nil, err
	}
	resItems := make([]*pb.Permission, len(acl))
	for i, bean := range acl {
		r := &pb.Permission{}
		util.MapPermissionBeanToPb(bean, r)
		resItems[i] = r
	}
	res := &pb.ListSubjectPermissionResponse{
		Acl: resItems,
	}
	ti, _ := saas.FromCurrentTenant(ctx)
	var groups []*pb.PermissionDefGroup
	authz.WalkGroups(len(ti.GetId()) == 0, true, func(group *authz.PermissionDefGroup) {
		g := &pb.PermissionDefGroup{}
		mapGroupDef2Pb(ctx, group, g)
		groups = append(groups, g)
		var defs []*pb.PermissionDef
		group.Walk(len(ti.GetId()) == 0, true, func(def *authz.PermissionDef) {
			d := &pb.PermissionDef{}
			mapDef2Pb(ctx, def, d)
			defs = append(defs, d)
		})
		g.Def = defs
	})
	res.DefGroups = groups
	return res, nil
}

func (s *PermissionService) UpdateSubjectPermission(ctx context.Context, req *pb.UpdateSubjectPermissionRequest) (*pb.UpdateSubjectPermissionResponse, error) {
	if err := kapi.ErrIfUntrusted(ctx, s.trust); err != nil {
		return nil, err
	}

	var acl []authz.UpdateSubjectPermission
	for _, a := range req.Acl {
		if err := s.findAnyValidateModifyPermissionDef(ctx, a.Namespace, a.Action); err != nil {
			return nil, err
		}
		effect := util.MapPbEffect2AuthEffect(a.Effect)
		acl = append(acl, authz.UpdateSubjectPermission{
			Resource: authz.NewEntityResource(a.Namespace, a.Resource),
			Action:   authz.ActionStr(a.Action),
			TenantID: a.TenantId,
			Effect:   effect,
		})
	}

	if err := s.permissionMgr.UpdateGrant(ctx, authz.SubjectStr(req.Subject), acl); err != nil {
		return nil, err
	}
	return &pb.UpdateSubjectPermissionResponse{}, nil
}

func (s *PermissionService) RemoveSubjectPermission(ctx context.Context, req *pb.RemoveSubjectPermissionRequest) (*pb.RemoveSubjectPermissionReply, error) {
	if err := kapi.ErrIfUntrusted(ctx, s.trust); err != nil {
		return nil, err
	}
	effList := make([]authz.Effect, len(req.Effects))
	for i, effect := range req.Effects {
		effList[i] = util.MapPbEffect2AuthEffect(effect)
	}
	var opts []authz.FilterFunc
	opts = append(opts, authz.WithEffectsFilter(effList...))
	if req.Action != nil {
		opts = append(opts, authz.WithActionFilter(authz.ActionStr(*req.Action)))
	}
	if req.Resource != nil && req.Namespace != nil {
		opts = append(opts, authz.WithResourceFilter(authz.NewEntityResource(*req.Namespace, *req.Resource)))
	}
	if req.TenantId != nil {
		opts = append(opts, authz.WithTenantFilter(*req.TenantId))
	}
	if err := s.permissionMgr.RemoveGrant(ctx, authz.SubjectStr(req.Subject), opts...); err != nil {
		return nil, err
	}
	return &pb.RemoveSubjectPermissionReply{}, nil
}

func (s *PermissionService) findAnyValidateModifyPermissionDef(ctx context.Context, namespace string, action string) error {
	//find def
	def, err := authz.FindDef(namespace, authz.ActionStr(action), true)
	if err != nil {
		return err
	}
	ti, _ := saas.FromCurrentTenant(ctx)
	if (def.Side == authz.PermissionAllowSide_HOST_ONLY && len(ti.GetId()) != 0) || (def.Side == authz.PermissionAllowSide_TENANT_ONLY && len(ti.GetId()) == 0) {
		return errors.New(400, authz.DefNotFoundReason, fmt.Sprintf("action %s in %s side mismatch", action, namespace))
	}
	return nil
}

func mapGroupDef2Pb(ctx context.Context, a *authz.PermissionDefGroup, b *pb.PermissionDefGroup) {
	b.Name = a.Name
	b.DisplayName = localize.GetMsg(ctx, a.Name, a.Name, nil, nil)
	b.Side = mapSide2Pb(a.Side)
	b.Priority = a.Priority
	b.Extra = a.Extra
}

func mapDef2Pb(ctx context.Context, a *authz.PermissionDef, b *pb.PermissionDef) {
	b.Name = a.Name
	b.DisplayName = localize.GetMsg(ctx, a.Name, a.Name, nil, nil)
	b.Side = mapSide2Pb(a.Side)
	b.Extra = a.Extra
	b.Namespace = a.Namespace
	b.Action = a.Action

}

func mapSide2Pb(side authz.PermissionAllowSide) pb.PermissionSide {
	switch side {
	case authz.PermissionAllowSide_HOST_ONLY:
		return pb.PermissionSide_HOST_ONLY
	case authz.PermissionAllowSide_TENANT_ONLY:
		return pb.PermissionSide_TENANT_ONLY
	default:
		return pb.PermissionSide_BOTH
	}
}
