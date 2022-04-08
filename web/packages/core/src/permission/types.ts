export interface PermissionAcl {
  namespace: string;
  resource: string;
  action: string;
  effect: PermissionEffect;
}

export interface PermissionRequirement {
  namespace: string;
  resource: string;
  action: string;
}

export enum PermissionEffect {
  UNKNOWN = 'UNKNOWN',
  GRANT = 'GRANT',
  FORBIDDEN = 'FORBIDDEN',
}
