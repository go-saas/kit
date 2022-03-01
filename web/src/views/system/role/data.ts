import { RoleServiceApi } from '/@/api-gen/api/role-service-api';
import { V1Role } from '/@/api-gen/models/v1-role';
import { BasicFetchResult } from '/@/api/model/baseModel';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';

const defaultRequirement = {
  namespace: '*',
  resource: '*',
  action: '*',
};

const { t } = useI18n();

export function getRoleColumns(): BasicColumn[] {
  return [
    {
      title: t('routes.system.role.name'),
      dataIndex: 'name',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.role.isPreserved'),
      dataIndex: 'isPreserved',
      requirement: defaultRequirement,
    },
  ];
}

export async function getRoleData(): Promise<BasicFetchResult<V1Role>> {
  return await new RoleServiceApi().roleServiceListRoles().then((response) => {
    const roleData: BasicFetchResult<V1Role> = {
      total: response.data.totalSize!,
      items: response.data.items!,
    };
    console.log(response);
    return roleData;
  });
}
