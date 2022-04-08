import {
  UserServiceApi,
  UserServiceApiUserServiceCreateUserRequest,
  V1User,
  RoleServiceApi,
} from '/@/api-gen';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';
import { BasicFetchResult, BasicPageParams } from '/@/api/model/baseModel';
import { FormSchema } from '/@/components/Table';
import { FileItem } from '/@/components/Upload/src/typing';
import { uploadApi } from '/@/api/sys/upload';
import { UploadFileParams } from '/#/axios';

const { t } = useI18n();

export function getUsertColumns(): BasicColumn[] {
  return [
    {
      title: t('routes.system.user.avatar'),
      dataIndex: 'avatar',

      slots: { customRender: 'avatar' },
    },

    {
      title: t('routes.system.user.username'),
      dataIndex: 'username',
    },
    {
      title: t('routes.system.user.phone'),
      dataIndex: 'phone',
    },
    {
      title: t('routes.system.user.email'),
      dataIndex: 'email',
    },
    {
      title: t('user.user.role'),
      dataIndex: 'roles',

      slots: { customRender: 'roles' },
    },
    {
      title: t('routes.system.user.name'),
      dataIndex: 'name',
    },
    {
      title: t('user.gender.gender'),
      dataIndex: 'gender',

      format: (text: string, _record: Recordable, _index: number) => {
        return t(`user.gender.${text.toLowerCase()}`);
      },
    },
    {
      title: t('routes.system.user.birthday'),
      dataIndex: 'birthday',

      format: 'date|YYYY-MM-DD HH:mm:ss',
    },
    {
      title: t('common.createdAt'),
      dataIndex: 'createdAt',

      format: 'date|YYYY-MM-DD HH:mm:ss',
    },
    {
      title: t('common.updatedAt'),
      dataIndex: 'updatedAt',

      format: 'date|YYYY-MM-DD HH:mm:ss',
    },
  ];
}

export const formSchema: FormSchema[] = [
  {
    field: 'avatar',
    label: t('routes.system.user.avatar'),
    component: 'Upload',
    componentProps: {
      api: uploadUserAvatar,
      maxNumber: 1,
    },
    valueTransformer: (files: FileItem[]) => {
      return files ? files.find((_) => true)?.id : undefined;
    },
    displayTransformer: (blob) => [blob],
  },
  {
    field: 'name',
    label: t('routes.system.user.name'),
    component: 'Input',
  },
  {
    field: 'username',
    label: t('routes.system.user.username'),
    component: 'Input',
    required: true,
  },
  {
    field: 'phone',
    label: t('routes.system.user.phone'),
    component: 'Input',
  },
  {
    field: 'email',
    label: t('routes.system.user.email'),
    component: 'Input',
    rules: [{ type: 'email' }],
  },
  {
    field: 'birthday',
    label: t('routes.system.user.birthday'),
    component: 'DatePicker',
  },
  {
    field: 'gender',
    label: t('user.gender.gender'),
    component: 'Select',
    componentProps: {
      options: [
        {
          label: t('user.gender.male'),
          value: 'MALE',
          key: 'MALE',
        },
        {
          label: t('user.gender.female'),
          value: 'FEMALE',
          key: 'FEMALE',
        },
        {
          label: t('user.gender.other'),
          value: 'OTHER',
          key: 'OTHER',
        },
      ],
    },
    defaultValue: 'OTHER',
  },
  {
    field: 'roles',
    label: t('routes.system.user.roles'),
    component: 'ApiSelect',
    componentProps: {
      mode: 'multiple',
      api: getRoles,
      labelField: 'name',
      valueField: 'id',
    },
  },
  {
    field: 'divider',
    label: '',
    component: 'Divider',
  },
  {
    field: 'password',
    label: t('routes.system.user.password'),
    component: 'StrengthMeter',
    required: true,
  },
  {
    field: 'confirmPassword',
    label: t('routes.system.user.confirmPassword'),
    component: 'Input',
    required: true,
    dynamicRules: ({ values }) => {
      return [
        {
          required: true,
          validator: (_, value) => {
            if (!value) {
              return Promise.reject('not null');
            }
            if (value !== values.password) {
              return Promise.reject('not equal');
            }
            return Promise.resolve();
          },
        },
      ];
    },
  },
];

export async function getRoles() {
  return await new RoleServiceApi().roleServiceListRoles().then((response) => {
    return response.data.items;
  });
}

export async function uploadUserAvatar(
  params: UploadFileParams,
  onUploadProgress: (progressEvent: ProgressEvent) => void,
) {
  return uploadApi('/v1/user/avatar', params, onUploadProgress);
}

export async function postcreateUser(param: UserServiceApiUserServiceCreateUserRequest) {
  return await new UserServiceApi().userServiceCreateUser(param).then((response) => {
    return response.data;
  });
}

// user列表
export async function getUsertData(param: BasicPageParams): Promise<BasicFetchResult<V1User>> {
  return await new UserServiceApi()
    .userServiceListUsers2({
      body: {
        pageSize: param.pageSize,
        pageOffset: (param.page - 1) * param.pageSize,
      },
    })
    .then((response) => {
      const userData: BasicFetchResult<V1User> = {
        total: response.data.totalSize!,
        items: response.data.items!,
      };
      console.log(userData);
      return userData;
    });
}
