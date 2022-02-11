import type { RouteRecordRaw } from 'vue-router';

import { usePermissionStore } from '/@/store/modules/permission';

import { useTabs } from './useTabs';

import { router, resetRouter } from '/@/router';
// import { RootRoute } from '/@/router/routes';

import { useMultipleTabStore } from '/@/store/modules/multipleTab';

import { PermissionRequirement } from '/#/store';
import { isGrant } from '/@/utils/permission';
// User permissions related operations
export function usePermission() {
  const permissionStore = usePermissionStore();
  const { closeAll } = useTabs(router);

  /**
   * Reset and regain authority resource information
   * @param id
   */
  async function resume() {
    const tabStore = useMultipleTabStore();
    tabStore.clearCacheTabs();
    resetRouter();
    const routes = await permissionStore.buildRoutesAction();
    routes.forEach((route) => {
      router.addRoute(route as unknown as RouteRecordRaw);
    });
    permissionStore.setLastBuildMenuTime();
    closeAll();
  }

  /**
   * Determine whether there is permission
   */
  function hasPermission(
    requirement: Nullable<PermissionRequirement | PermissionRequirement[]>,
    def = true,
  ): boolean {
    // Visible by default
    if (!requirement) {
      return def;
    }
    return isGrant(requirement, permissionStore.getPermissionList);
  }

  /**
   * refresh menu data
   */
  async function refreshMenu() {
    resume();
  }

  return { hasPermission, refreshMenu };
}
