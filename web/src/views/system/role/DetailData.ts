import {
  RoleServiceApi,
  RoleServiceApiRoleServiceGetRolePermissionRequest,
  RoleServiceApiRoleServiceGetRoleRequest,
  RoleServiceApiRoleServiceUpdateRolePermissionRequest,
  RoleServiceApiRoleServiceUpdateRoleRequest,
} from '/@/api-gen';
import { FormSchema } from '/@/components/Form/src/types/form';
import { useI18n } from '/@/hooks/web/useI18n';
const { t } = useI18n();
export const formSchema: FormSchema[] = [
  {
    field: 'action',
    label: t('routes.system.permission.action'),
    component: 'Input',
  },
  {
    field: 'effect',
    label: t('routes.system.permission.effect'),
    component: 'Select',
    componentProps: {
      options: [
        {
          label: 'GRANT',
          value: 'GRANT',
        },
        {
          label: 'FORBIDDEN',
          value: 'FORBIDDEN',
        },
      ],
    },
  },
  {
    field: 'namespace',
    label: t('routes.system.permission.namespace'),
    component: 'Input',
  },
  {
    field: 'resource',
    label: t('routes.system.permission.resource'),
    component: 'Input',
  },
];

export async function getRoleDetail(params: RoleServiceApiRoleServiceGetRoleRequest) {
  return await new RoleServiceApi().roleServiceGetRole(params);
}
export async function putRole(params: RoleServiceApiRoleServiceUpdateRoleRequest) {
  return await new RoleServiceApi().roleServiceUpdateRole(params);
}

export async function getRolePermission(param: RoleServiceApiRoleServiceGetRolePermissionRequest) {
  return await new RoleServiceApi().roleServiceGetRolePermission(param);
}

export async function putRolePermission(
  param: RoleServiceApiRoleServiceUpdateRolePermissionRequest,
) {
  return await new RoleServiceApi().roleServiceUpdateRolePermission({
    id: param.id,
    body: param.body,
  });
}
