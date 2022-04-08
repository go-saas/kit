import type { AppRouteRecordRaw, Menu } from '/@/router/types';

import { defineStore } from 'pinia';
import { store } from '/@/store';
import { useI18n } from '/@/hooks/web/useI18n';
import {
  flatMultiLevelRoutes,
  transformObjToAppRouteRecordRaw,
  transformObjToRoute,
} from '/@/router/helper/routeHelper';
import { transformRouteToMenu } from '/@/router/helper/menuHelper';
import { PermissionAcl } from '@kit/core';
import { PermissionEffect } from '@kit/core';

import { getAuthCache, setAuthCache } from '/@/utils/auth';
import { ERROR_LOG_ROUTE, PAGE_NOT_FOUND_ROUTE } from '/@/router/routes/basic';

import { filter } from '/@/utils/helper/treeHelper';

import { PermissionServiceApi, MenuServiceApi } from '/@/api-gen';
import { PERMISSION_KEY } from '/@/enums/cacheEnum';
import { useMessage } from '/@/hooks/web/useMessage';
import { PageEnum } from '/@/enums/pageEnum';
import { isGrant } from '@kit/core';

interface PermissionState {
  permissionList: PermissionAcl[];
  // Whether the route has been dynamically added
  isDynamicAddedRoute: boolean;
  // To trigger a menu update
  lastBuildMenuTime: number;
  // Backstage menu list
  backMenuList: Menu[];
  frontMenuList: Menu[];
}
export const usePermissionStore = defineStore({
  id: 'app-permission',
  state: (): PermissionState => ({
    permissionList: [],
    // Whether the route has been dynamically added
    isDynamicAddedRoute: false,
    // To trigger a menu update
    lastBuildMenuTime: 0,
    // Backstage menu list
    backMenuList: [],
    // menu List
    frontMenuList: [],
  }),
  getters: {
    getPermissionList(): PermissionAcl[] {
      return this.permissionList.length > 0
        ? this.permissionList
        : getAuthCache<PermissionAcl[]>(PERMISSION_KEY);
    },
    getBackMenuList(): Menu[] {
      return this.backMenuList;
    },
    getLastBuildMenuTime(): number {
      return this.lastBuildMenuTime;
    },
    getIsDynamicAddedRoute(): boolean {
      return this.isDynamicAddedRoute;
    },
  },
  actions: {
    setPermissionList(permissionList: PermissionAcl[]) {
      this.permissionList = permissionList;
      setAuthCache(PERMISSION_KEY, permissionList);
    },
    setBackMenuList(list: Menu[]) {
      this.backMenuList = list;
      list?.length > 0 && this.setLastBuildMenuTime();
    },

    setLastBuildMenuTime() {
      this.lastBuildMenuTime = new Date().getTime();
    },

    setDynamicAddedRoute(added: boolean) {
      this.isDynamicAddedRoute = added;
    },
    resetState(): void {
      this.isDynamicAddedRoute = false;
      this.permissionList = [];
      this.backMenuList = [];
      this.lastBuildMenuTime = 0;
    },
    async changePermissionCode() {
      const codeList = await new PermissionServiceApi().permissionServiceGetCurrent();
      this.setPermissionList(
        (codeList.data?.acl ?? []).map((p) => {
          return {
            namespace: p.namespace ?? '',
            resource: p.resource ?? '',
            action: p.action ?? '',
            effect: p.effect as PermissionEffect,
          };
        }),
      );
    },
    async buildRoutesAction(): Promise<AppRouteRecordRaw[]> {
      const { t } = useI18n();

      let routes: AppRouteRecordRaw[] = [];

      const routeFilter = (route: AppRouteRecordRaw) => {
        const { meta } = route;
        const { requirement } = meta || {};
        if (!requirement) return true;
        return isGrant(requirement, this.permissionList);
      };

      const routeRemoveIgnoreFilter = (route: AppRouteRecordRaw) => {
        const { meta } = route;
        const { ignoreRoute } = meta || {};
        return !ignoreRoute;
      };

      /**
       * @description 根据设置的首页path，修正routes中的affix标记（固定首页）
       * */
      const patchHomeAffix = (routes: AppRouteRecordRaw[]) => {
        if (!routes || routes.length === 0) return;
        let homePath: string = PageEnum.BASE_HOME;
        function patcher(routes: AppRouteRecordRaw[], parentPath = '') {
          if (parentPath) parentPath = parentPath + '/';
          routes.forEach((route: AppRouteRecordRaw) => {
            const { path, children, redirect } = route;
            const currentPath = path.startsWith('/') ? path : parentPath + path;
            if (currentPath === homePath) {
              if (redirect) {
                homePath = route.redirect! as string;
              } else {
                route.meta = Object.assign({}, route.meta, { affix: true });
                throw new Error('end');
              }
            }
            children && children.length > 0 && patcher(children, currentPath);
          });
        }
        try {
          patcher(routes);
        } catch (e) {
          // 已处理完毕跳出循环
        }
        return;
      };

      const { createMessage } = useMessage();

      createMessage.loading({
        content: t('sys.app.menuLoading'),
        duration: 1,
      });

      // !Simulate to obtain permission codes from the background,
      // this function may only need to be executed once, and the actual project can be put at the right time by itself
      let routeList: AppRouteRecordRaw[] = [];
      try {
        await this.changePermissionCode();
        routeList = transformObjToAppRouteRecordRaw(
          (await new MenuServiceApi().menuServiceGetAvailableMenus()).data?.items ?? [],
        );
      } catch (error) {
        console.error(error);
      }

      routeList = transformObjToRoute(routeList);
      //  Background routing to menu structure
      const backMenuList = transformRouteToMenu(routeList);

      this.setBackMenuList(backMenuList);

      // remove meta.ignoreRoute item
      routeList = filter(routeList, routeRemoveIgnoreFilter);
      routeList = routeList.filter(routeRemoveIgnoreFilter);

      routeList = filter(routeList, routeFilter);
      routeList = routeList.filter(routeFilter);

      routeList = flatMultiLevelRoutes(routeList);
      routes = [PAGE_NOT_FOUND_ROUTE, ...routeList];

      routes.push(ERROR_LOG_ROUTE);
      patchHomeAffix(routes);
      return routes;
    },
  },
});

// Need to be used outside the setup
export function usePermissionStoreWithOut() {
  return usePermissionStore(store);
}
