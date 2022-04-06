import { ErrorTypeEnum } from '/@/enums/exceptionEnum';
import { MenuModeEnum, MenuTypeEnum } from '/@/enums/menuEnum';
import { PermissionEffect } from '/@/enums/permissionEnum';
// Lock screen information
export interface LockInfo {
  // Password required
  pwd?: string | undefined;
  // Is it locked?
  isLock?: boolean;
}

// Error-log information
export interface ErrorLogInfo {
  // Type of error
  type: ErrorTypeEnum;
  // Error file
  file: string;
  // Error name
  name?: string;
  // Error message
  message: string;
  // Error stack
  stack?: string;
  // Error detail
  detail: string;
  // Error url
  url: string;
  // Error time
  time?: string;
}

export interface RoleInfo {
  id: string;
  name: string;
  isPreserved: boolean;
}

export interface TenantInfo {
  id: string;
  name: string;
  displayName: string;
  region: string;
  logo?: BlobFile;
}

export interface BlobFile {
  url?: string;
}
export interface UserTenantInfo {
  isHost: bool;
  tenant?: TenantInfo;
}

export interface UserInfo {
  id: string | number;
  username: string;
  name: string;
  avatar: string;
  roles: RoleInfo[];
  tenants: UserTenantInfo[] = [];
  currentTenant: UserTenantInfo;
}

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

export interface BeforeMiniState {
  menuCollapsed?: boolean;
  menuSplit?: boolean;
  menuMode?: MenuModeEnum;
  menuType?: MenuTypeEnum;
}
