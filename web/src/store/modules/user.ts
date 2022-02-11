import type { UserInfo } from '/#/store';
import type { ErrorMessageMode } from '/#/axios';
import { defineStore } from 'pinia';
import { store } from '/@/store';
import { PageEnum } from '/@/enums/pageEnum';
import { TOKEN_KEY, USER_INFO_KEY } from '/@/enums/cacheEnum';
import { getAuthCache, setAuthCache } from '/@/utils/auth';
import { AuthWebApi, AccountApi } from '/@/api-gen/api';
import { V1LoginAuthRequest } from '/@/api-gen/models';
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
  }),
  getters: {
    getUserInfo(): UserInfo {
      return this.userInfo || getAuthCache<UserInfo>(USER_INFO_KEY) || {};
    },
    getIsLogin(): boolean {
      return (this.userInfo || getAuthCache<UserInfo>(USER_INFO_KEY) || {})?.id != null;
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
      params: V1LoginAuthRequest & {
        redirect?: string;
        mode?: ErrorMessageMode;
      },
    ): Promise<Nullable<UserInfo>> {
      try {
        const { redirect = '/', mode, ...loginParams } = params;
        const data = await new AuthWebApi().authWebWebLogin(
          { body: loginParams },
          { data: { errorMessageMode: mode } },
        );

        const { accessToken } = data.data;
        if (accessToken) {
          // save token
          this.setToken(accessToken);
        }
        return this.afterLoginAction(redirect);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    async afterLoginAction(redirect = '/'): Promise<Nullable<UserInfo>> {
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
        await router.replace(redirect);
      }
      return userInfo;
    },
    async getUserInfoAction(): Promise<UserInfo | null> {
      const userInfo = (await new AccountApi().accountGetProfile()).data;
      const { roles = [] } = userInfo;

      const converted: UserInfo = {
        id: userInfo.id!,
        username: userInfo.username ?? '',
        name: userInfo.name ?? '',
        //TODO avatar
        avatar: '',
        roles: roles.map((item) => {
          return { name: item.name ?? '', isPreserved: item.isPreserved ?? false, id: item.id! };
        }),
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
