import {
  RoleServiceApi,
  RoleServiceApiRoleServiceCreateRoleRequest,
  RoleServiceApiRoleServiceDeleteRoleRequest,
  RoleServiceApiRoleServiceUpdateRoleRequest,
} from '/@/api-gen';
import { V1Role } from '/@/api-gen';
import { BasicFetchResult, BasicPageParams } from '/@/api/model/baseModel';
import { FormSchema } from '/@/components/Form/src/types/form';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';
const { t } = useI18n();

export function getRoleColumns(): BasicColumn[] {
  return [
    {
      title: t('role.role.name'),
      dataIndex: 'name',
      slots: { customRender: 'name' },
    },
  ];
}

export const formSchema: FormSchema[] = [
  {
    field: 'name',
    label: t('role.role.name'),
    component: 'Input',
    required: true,
  },
];
// 编辑
export async function putEditRoleData(param: RoleServiceApiRoleServiceUpdateRoleRequest) {
  return await new RoleServiceApi().roleServiceUpdateRole({
    roleId: param.roleId,
    body: param.body,
  });
}

// 增加
export async function postcreateRoleData(param: RoleServiceApiRoleServiceCreateRoleRequest) {
  return await new RoleServiceApi().roleServiceCreateRole({
    body: param.body,
  });
}
// // 删除
export async function postDeleteRoleData(param: RoleServiceApiRoleServiceDeleteRoleRequest) {
  return await new RoleServiceApi().roleServiceDeleteRole({
    id: param.id,
  });
}
// 分页
export async function getRoleData(param: BasicPageParams): Promise<BasicFetchResult<V1Role>> {
  return await new RoleServiceApi()
    .roleServiceListRoles2({
      body: {
        pageSize: param.pageSize,
        pageOffset: (param.page - 1) * param.pageSize,
      },
    })
    .then((response) => {
      const data: BasicFetchResult<V1Role> = {
        total: response.data.totalSize!,
        items: response.data.items!,
      };
      return data;
    });
}
export function getRoleName(record: V1Role): string {
  if (record.isPreserved) {
    return t(`role.role.${record.name?.toLowerCase()}`);
  } else {
    return record.name ?? '';
  }
}
