import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const saas: AppRouteModule = {
  path: '/saas',
  name: 'Saas',
  component: LAYOUT,
  redirect: '/saas/tenants',
  meta: {
    orderNo: 101,
    icon: 'ion:grid-outline',
    title: t('routes.saas.title'),
  },
  children: [
    {
      path: 'tenants',
      name: 'Tenant',
      component: () => import('/@/views/saas/tenant/index.vue'),
      meta: {
        // affix: true,
        title: t('routes.saas.tenant.management'),
      },
    },
  ],
};

export default saas;
