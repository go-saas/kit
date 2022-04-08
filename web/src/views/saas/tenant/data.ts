import {
  TenantServiceApi,
  TenantServiceApiTenantServiceCreateTenantRequest,
  TenantServiceApiTenantServiceDeleteTenantRequest,
  TenantServiceApiTenantServiceUpdateTenantRequest,
} from '/@/api-gen';
import { V1Tenant } from '/@/api-gen';
import { BasicFetchResult, BasicPageParams } from '/@/api/model/baseModel';
import { FormSchema } from '/@/components/Form/src/types/form';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';
import { uploadApi } from '/@/api/sys/upload';
import { UploadFileParams } from '/#/axios';
import { FileItem } from '/@/components/Upload/src/typing';
const { t } = useI18n();

export function getTenantColumns(): BasicColumn[] {
  return [
    // {
    //   title: 'Logo',
    //   dataIndex: 'logo',
    // },
    {
      title: t('saas.tenant.logo'),
      dataIndex: 'logo',
      width: 100,
      slots: { customRender: 'logo' },
    },
    {
      title: t('saas.tenant.name'),
      dataIndex: 'name',
    },
    {
      title: t('saas.tenant.region'),
      dataIndex: 'region',
    },
    {
      title: t('saas.tenant.displayName'),
      dataIndex: 'displayName',
    },
    {
      title: t('saas.tenant.createdAt'),
      dataIndex: 'createdAt',
      format: 'date|YYYY-MM-DD HH:mm:ss',
    },
    {
      title: t('saas.tenant.updatedAt'),
      dataIndex: 'updatedAt',
      format: 'date|YYYY-MM-DD HH:mm:ss',
    },
  ];
}

export const formSchema: FormSchema[] = [
  {
    field: 'name',
    label: t('saas.tenant.name'),
    component: 'Input',
    required: true,
    rules: [{ pattern: new RegExp('^[A-Za-z0-9](?:[A-Za-z0-9\\-]{1,61}[A-Za-z0-9])?$') }],
  },
  {
    field: 'region',
    label: t('saas.tenant.region'),
    component: 'Input',
    required: true,
  },
  {
    field: 'displayName',
    label: t('saas.tenant.displayName'),
    component: 'Input',
    required: true,
  },
  {
    field: 'logo',
    label: t('saas.tenant.logo'),
    component: 'Upload',
    required: true,
    componentProps: {
      api: uploadTenantLogo,
      maxNumber: 1,
    },
    valueTransformer: (files: FileItem[]) => {
      return files.find((_) => true)?.id;
    },
    displayTransformer: (blob) => [blob],
  },
  {
    field: 'separateDb',
    label: t('saas.tenant.separateDb'),
    defaultValue: false,
    component: 'Switch',
    required: true,
  },
  {
    field: 'adminUsername',
    label: t('saas.tenant.adminUsername'),
    component: 'Input',
    required: (callback) => {
      return callback.values['separateDb'];
    },
  },
  {
    field: 'adminPassword',
    label: t('saas.tenant.adminPassword'),
    component: 'StrengthMeter',
    required: (callback) => {
      return callback.values['separateDb'];
    },
  },
];

export const updateFormSchema: FormSchema[] = [
  {
    field: 'name',
    label: t('saas.tenant.name'),
    component: 'Input',
    required: true,
    rules: [{ pattern: new RegExp('^[A-Za-z0-9](?:[A-Za-z0-9\\-]{1,61}[A-Za-z0-9])?$') }],
  },
  {
    field: 'displayName',
    label: t('saas.tenant.displayName'),
    component: 'Input',
    required: true,
  },
  {
    field: 'logo',
    label: t('saas.tenant.logo'),
    component: 'Upload',
    required: true,
    componentProps: {
      api: uploadTenantLogo,
      maxNumber: 1,
    },
    valueTransformer: (files: FileItem[]) => {
      return files.find((_) => true)?.id;
    },
    displayTransformer: (blob) => [blob],
  },
];

// 编辑
export async function putEditTenantData(param: TenantServiceApiTenantServiceUpdateTenantRequest) {
  return await new TenantServiceApi().tenantServiceUpdateTenant({
    tenantId: param.tenantId,
    body: {
      tenant: {
        id: param.tenantId,
        name: param.body.tenant?.name == undefined ? '' : param.body.tenant?.name,
        displayName: param.body.tenant?.displayName,
        logo: param.body.tenant?.logo,
      },
    },
  });
}

export async function getTenantDataDetail(id: string) {
  return await new TenantServiceApi().tenantServiceGetTenant({ idOrName: id });
}

// 增加
export async function postcreateTenantData(
  param: TenantServiceApiTenantServiceCreateTenantRequest,
) {
  return await new TenantServiceApi().tenantServiceCreateTenant({
    body: param.body,
  });
}
// 删除
export async function postDeleteTenantData(
  param: TenantServiceApiTenantServiceDeleteTenantRequest,
) {
  return await new TenantServiceApi().tenantServiceDeleteTenant({
    id: param.id,
  });
}

export async function uploadTenantLogo(
  params: UploadFileParams,
  onUploadProgress: (progressEvent: ProgressEvent) => void,
) {
  return uploadApi('/v1/saas/tenant/logo', params, onUploadProgress);
}

export async function getTenantData(param: BasicPageParams): Promise<BasicFetchResult<V1Tenant>> {
  return await new TenantServiceApi()
    .tenantServiceListTenant2({
      body: {
        pageSize: param.pageSize,
        pageOffset: (param.page - 1) * param.pageSize,
      },
    })
    .then((response) => {
      const tenantData: BasicFetchResult<V1Tenant> = {
        total: response.data.items!.length,
        items: response.data.items!,
      };
      return tenantData;
    });
}
