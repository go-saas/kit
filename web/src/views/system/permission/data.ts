import { PermissionServiceApi } from '/@/api-gen/api/permission-service-api';
import { V1Permission } from '/@/api-gen/models/v1-permission';
import { BasicFetchResult } from '/@/api/model/baseModel';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';

const defaultRequirement = {
  namespace: '*',
  resource: '*',
  action: '*',
};

const { t } = useI18n();

export function getPremissionColumns(): BasicColumn[] {
  return [
    {
      title: t('routes.system.permission.action'),
      dataIndex: 'action',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.effect'),
      dataIndex: 'effect',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.namespace'),
      dataIndex: 'namespace',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.resource'),
      dataIndex: 'resource',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.subject'),
      dataIndex: 'subject',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.tenantId'),
      dataIndex: 'tenantId',
      requirement: defaultRequirement,
    },
  ];
}

export async function getPremissionData(): Promise<BasicFetchResult<V1Permission>> {
  return await new PermissionServiceApi()
    .permissionServiceListSubjectPermission()
    .then((response) => {
      const premissionData: BasicFetchResult<V1Permission> = {
        total: response.data.acl!.length,
        items: response.data.acl!,
      };
      return premissionData;
    });
}
