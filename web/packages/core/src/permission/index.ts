import { PermissionAcl, PermissionRequirement, PermissionEffect } from './types';
import { isArray } from '../utils';
export function isGrant(
  requirement: PermissionRequirement | PermissionRequirement[],
  aclList: PermissionAcl[],
): boolean {
  if (!isArray(requirement)) {
    requirement = [requirement];
  }

  for (const r of requirement) {
    let effect = PermissionEffect.UNKNOWN;
    for (const acl of aclList) {
      if (
        keyMatch(r.action, acl.action) &&
        keyMatch(r.namespace, acl.namespace) &&
        keyMatch(r.resource, acl.resource)
      ) {
        if (acl.effect == PermissionEffect.FORBIDDEN) {
          return false;
        }
        if (acl.effect == PermissionEffect.GRANT) {
          effect = PermissionEffect.GRANT;
        }
      }
    }
    if (effect != PermissionEffect.GRANT) {
      return false;
    }
  }
  return true;
}

function keyMatch(source: string, target: string): boolean {
  const i = target.indexOf('*');
  if (i == -1) {
    return source == target;
  }

  if (source.length > i) {
    return source.substring(0, i) == target.substring(0, i);
  }
  return source == target.substring(0, i);
}

export * from './types';
