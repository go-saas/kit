import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const system: AppRouteModule = {
  path: '/system',
  name: 'System',
  component: LAYOUT,
  redirect: '/system/users',
  meta: {
    orderNo: 100,
    icon: 'ion:grid-outline',
    title: t('routes.system.title'),
  },
  children: [
    {
      path: 'users',
      name: 'User',
      component: () => import('/@/views/system/user/index.vue'),
      meta: {
        // affix: true,
        title: t('routes.system.user.management'),
      },
    },
    {
      path: 'roles',
      name: 'Role',
      component: () => import('/@/views/system/role/index.vue'),
      meta: {
        // affix: true,
        title: t('routes.system.role.management'),
      },
    },
    {
      path: 'permission',
      name: 'Permission',
      component: () => import('/@/views/system/permission/index.vue'),
      meta: {
        // affix: true,
        title: t('routes.system.permission.management'),
      },
    },
  ],
};

export default system;
