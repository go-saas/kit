import type { UserInfo, UserTenantInfo } from '/#/store';
import type { ErrorMessageMode } from '/#/axios';
import { defineStore } from 'pinia';
import { store } from '/@/store';
import { PageEnum } from '/@/enums/pageEnum';
import { TOKEN_KEY, USER_INFO_KEY, TENANT_INFO_KEY } from '/@/enums/cacheEnum';
import { getAuthCache, setAuthCache } from '/@/utils/auth';
import { AuthWebApi, AccountApi } from '/@/api-gen/api';
import { V1WebLoginAuthRequest } from '/@/api-gen/models';
import { useI18n } from '/@/hooks/web/useI18n';
import { useMessage } from '/@/hooks/web/useMessage';
import { router } from '/@/router';
import { usePermissionStore } from '/@/store/modules/permission';
import { RouteRecordRaw } from 'vue-router';
import { PAGE_NOT_FOUND_ROUTE } from '/@/router/routes/basic';
import { h } from 'vue';

interface UserState {
  userInfo: Nullable<UserInfo>;
  sessionTimeout?: boolean;
  lastUpdateTime: number;
  selectionTenantId?: Nullable<string>;
}

export const useUserStore = defineStore({
  id: 'app-user',
  state: (): UserState => ({
    // user info
    userInfo: null,
    // Whether the login expired
    sessionTimeout: false,
    // Last fetch time
    lastUpdateTime: 0,
    selectionTenantId: null,
  }),
  getters: {
    getUserInfo(): UserInfo {
      return this.userInfo || getAuthCache<UserInfo>(USER_INFO_KEY) || {};
    },
    getIsLogin(): boolean {
      return (this.userInfo || getAuthCache<UserInfo>(USER_INFO_KEY) || {})?.id != null;
    },
    getCurrentTenant(): UserTenantInfo {
      return (this.userInfo || getAuthCache<UserInfo>(USER_INFO_KEY) || {})?.currentTenant || {};
    },
    getCurrentIsHost(): boolean {
      return (
        (this.userInfo || getAuthCache<UserInfo>(USER_INFO_KEY) || {})?.currentTenant?.isHost ??
        false
      );
    },
    getCurrentSettingTenant(): Nullable<string> {
      return this.selectionTenantId || getAuthCache<string>(TENANT_INFO_KEY);
    },
    getSessionTimeout(): boolean {
      return !!this.sessionTimeout;
    },
    getLastUpdateTime(): number {
      return this.lastUpdateTime;
    },
  },
  actions: {
    setToken(info: string | undefined) {
      setAuthCache(TOKEN_KEY, info);
    },

    setUserInfo(info: UserInfo | null) {
      this.userInfo = info;
      this.lastUpdateTime = new Date().getTime();
      setAuthCache(USER_INFO_KEY, info);
    },
    setSessionTimeout(flag: boolean) {
      this.sessionTimeout = flag;
    },
    resetState() {
      this.userInfo = null;

      this.sessionTimeout = false;
    },
    /**
     * @description: login
     */
    async login(
      params: V1WebLoginAuthRequest & {
        redirect?: string;
        mode?: ErrorMessageMode;
      },
    ): Promise<Nullable<UserInfo>> {
      try {
        const { mode, ...loginParams } = params;
        // const data = await new AuthWebApi().authWebWebLogin(
        //   { body: loginParams },
        //   { data: { errorMessageMode: mode } },
        // );
        const data = await new AuthWebApi().authWebWebLogin(
          { body: loginParams },
          { data: { errorMessageMode: mode } },
        );

        const { accessToken, redirect } = data.data;
        if (accessToken) {
          // save token
          this.setToken(accessToken);
        }
        return this.afterLoginAction(redirect);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    async afterLoginAction(redirect): Promise<Nullable<UserInfo>> {
      if (redirect == null || redirect == '') {
        redirect = '/';
      }
      // get user info
      const userInfo = await this.getUserInfoAction();

      const sessionTimeout = this.sessionTimeout;
      if (sessionTimeout) {
        this.setSessionTimeout(false);
      } else {
        const permissionStore = usePermissionStore();
        if (!permissionStore.isDynamicAddedRoute) {
          const routes = await permissionStore.buildRoutesAction();
          routes.forEach((route) => {
            router.addRoute(route as unknown as RouteRecordRaw);
          });
          router.addRoute(PAGE_NOT_FOUND_ROUTE as unknown as RouteRecordRaw);
          permissionStore.setDynamicAddedRoute(true);
        }
        if (redirect.indexOf('http://') === 0 || redirect.indexOf('https://') === 0) {
          location.replace(redirect);
        } else {
          await router.replace(redirect);
        }
      }
      return userInfo;
    },
    async getUserInfoAction(): Promise<UserInfo | null> {
      const userInfo = (await new AccountApi().accountGetProfile()).data;
      const { roles = [], tenants = [] } = userInfo;

      const converted: UserInfo = {
        id: userInfo.id!,
        username: userInfo.username ?? '',
        name: userInfo.name ?? '',
        avatar: userInfo.avatar?.url ?? '',
        roles: roles.map((item) => {
          return { name: item.name ?? '', isPreserved: item.isPreserved ?? false, id: item.id! };
        }),
        tenants: tenants.map((item) => item as UserTenantInfo),
        currentTenant: userInfo.currentTenant as UserTenantInfo,
      };
      this.setUserInfo(converted);
      return converted;
    },
    /**
     * @description: logout
     */
    async logout(goLogin = false) {
      if (this.getIsLogin) {
        await new AuthWebApi().authWebWebLogout({ body: {} });
      }
      this.setToken(undefined);
      this.setSessionTimeout(false);
      this.setUserInfo(null);
      goLogin && router.push(PageEnum.BASE_LOGIN);
    },
    async changeTenant(userTenant?: string | null) {
      this.selectionTenantId = userTenant;
      setAuthCache(TENANT_INFO_KEY, userTenant);
      window.location.reload();
    },
    /**
     * @description: Confirm before logging out
     */
    confirmLoginOut() {
      const { createConfirm } = useMessage();
      const { t } = useI18n();
      createConfirm({
        iconType: 'warning',
        title: () => h('span', t('sys.app.logoutTip')),
        content: () => h('span', t('sys.app.logoutMessage')),
        onOk: async () => {
          await this.logout(true);
        },
      });
    },
  },
});

// Need to be used outside the setup
export function useUserStoreWithOut() {
  return useUserStore(store);
}
