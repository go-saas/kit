import { UserServiceApi } from '/@/api-gen/api/user-service-api';
import { V1User } from '/@/api-gen/models/v1-user';
import { BasicFetchResult, BasicPageParams } from '/@/api/model/baseModel';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';

const defaultRequirement = {
  namespace: '*',
  resource: '*',
  action: '*',
};

const { t } = useI18n();

export function getUserColumns(): BasicColumn[] {
  return [
    {
      title: t('routes.system.role.name'),
      dataIndex: 'name',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.role.birthday'),
      dataIndex: 'birthday',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.role.email'),
      dataIndex: 'email',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.role.gender'),
      dataIndex: 'gender',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.role.phone'),
      dataIndex: 'phone',
      requirement: defaultRequirement,
    },
  ];
}

export async function getUserData(param: BasicPageParams): Promise<BasicFetchResult<V1User>> {
  return await new UserServiceApi()
    .userServiceListUsers({
      pageOffset: (param.page - 1) * param.pageSize,
      pageSize: param.pageSize,
    })
    .then((response) => {
      const userData: BasicFetchResult<V1User> = {
        total: response.data.totalSize!,
        items: response.data.items!,
      };
      return userData;
    });
}
