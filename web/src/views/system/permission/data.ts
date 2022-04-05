import {
  PermissionServiceApi,
  V1AddSubjectPermissionRequest,
  V1ListSubjectPermissionRequest,
  V1Permission,
  V1RemoveSubjectPermissionRequest,
} from '/@/api-gen';
import { PermissionResult } from './PermissionType.data';
import { FormSchema } from '/@/components/Form/src/types/form';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';
const { t } = useI18n();
const defaultRequirement = {
  namespace: '*',
  resource: '*',
  action: '*',
};
export function getPermissiontColumns(): BasicColumn[] {
  return [
    {
      title: t('routes.system.permission.namespace'),
      dataIndex: 'namespace',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.subject'),
      dataIndex: 'subject',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.resource'),
      dataIndex: 'resource',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.tenantId'),
      dataIndex: 'tenantId',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.effect'),
      dataIndex: 'effect',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.permission.createdAt'),
      dataIndex: 'createdAt',
      requirement: defaultRequirement,
      format: 'date|YYYY-MM-DD HH:mm:ss',
    },
    {
      title: t('routes.system.permission.updatedAt'),
      dataIndex: 'updatedAt',
      requirement: defaultRequirement,
      format: 'date|YYYY-MM-DD HH:mm:ss',
    },
    {
      width: 60,
      title: t('routes.system.permission.action'),
      requirement: defaultRequirement,
      dataIndex: 'action',
      slots: { customRender: 'action' },
      fixed: undefined,
    },
  ];
}

export const formSchema: FormSchema[] = [
  {
    field: 'namespace',
    label: t('routes.system.permission.namespace'),
    component: 'Input',
    required: true,
  },
  {
    field: 'resource',
    label: t('routes.system.permission.resource'),
    component: 'Input',
    required: true,
  },
  {
    field: 'subject',
    label: t('routes.system.permission.subject'),
    component: 'Input',
    required: true,
  },
  {
    field: 'action',
    label: t('routes.system.permission.action'),
    component: 'Input',
    required: true,
  },
  {
    field: 'effect',
    label: t('routes.system.permission.effect'),
    component: 'Select',
    componentProps: {
      options: [
        {
          label: t('permission.grant'),
          value: 'GRANT',
        },
        {
          label: t('permission.fobidden'),
          value: 'FORBIDDEN',
        },
      ],
    },
    colProps: { span: 8 },
  },
];

// 增加
export async function postcreatePermissiontData(param: V1AddSubjectPermissionRequest) {
  return await new PermissionServiceApi().permissionServiceAddSubjectPermission({
    body: param,
  });
}
// 删除
export async function postDeletePermissiontData(param: V1RemoveSubjectPermissionRequest) {
  return await new PermissionServiceApi().permissionServiceRemoveSubjectPermission({
    body: param,
  });
}

// 权限列表
export async function getPermissiontData(
  param: V1ListSubjectPermissionRequest,
): Promise<PermissionResult<V1Permission>> {
  return await new PermissionServiceApi()
    .permissionServiceListSubjectPermission2({
      body: { subjects: param.subjects },
    })
    .then((response) => {
      const menuData: PermissionResult<V1Permission> = {
        total: response.data.acl!.length,
        items: response.data.acl!,
      };
      return menuData;
    });
}
