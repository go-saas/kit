import { h } from 'vue';
import {
  MenuServiceApi,
  MenuServiceApiMenuServiceCreateMenuRequest,
  MenuServiceApiMenuServiceDeleteMenuRequest,
  MenuServiceApiMenuServiceUpdateMenuRequest,
} from '/@/api-gen';
import { V1UpdateMenu } from '/@/api-gen';
import { BasicFetchResult, BasicPageParams } from '/@/api/model/baseModel';
import { FormSchema } from '/@/components/Form/src/types/form';
import Icon from '/@/components/Icon';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';
import { transformObjToAppRouteRecordRaw } from '/@/router/helper/routeHelper';
import { AppRouteRecordRaw } from '/@/router/types';
const { t } = useI18n();

export function getMenuColumns(): BasicColumn[] {
  return [
    {
      title: t('routes.system.menu.name'),
      dataIndex: 'name',
      align: 'left',
    },
    {
      title: t('routes.system.menu.icon'),
      dataIndex: 'icon',
      width: 50,
      customRender: ({ record }) => {
        return h(Icon, { icon: record.meta.icon });
      },
    },
    {
      title: t('routes.system.menu.component'),
      dataIndex: 'component',
    },
    {
      title: t('routes.system.menu.path'),
      dataIndex: 'path',
    },
    {
      title: t('routes.system.menu.redirect'),
      dataIndex: 'redirect',
    },
  ];
}

export const formSchema: FormSchema[] = [
  {
    field: 'name',
    label: t('routes.system.menu.name'),
    component: 'Input',
    required: true,
  },
  {
    field: 'component',
    label: t('routes.system.menu.component'),
    required: true,
    component: 'AutoComplete',
    componentProps: {
      options: [
        {
          label: 'LAYOUT',
          value: 'LAYOUT',
          key: 'LAYOUT',
        },
        {
          label: 'IFRAME',
          value: 'IFRAME',
          key: 'IFRAME',
        },
        {
          label: 'MICROAPP',
          value: 'MICROAPP',
          key: 'MICROAPP',
        },
      ],
    },
    defaultValue: 'LAYOUT',
  },
  {
    field: 'path',
    label: t('routes.system.menu.path'),
    component: 'Input',
    required: true,
  },
  {
    field: 'iframe',
    label: t('routes.system.menu.iframe'),
    component: 'Input',
    required: (callback) => {
      return callback.values['component'] == 'IFRAME';
    },
    ifShow: (callback) => {
      return callback.values['component'] == 'IFRAME';
    },
  },
  {
    field: 'microApp',
    label: 'MicroApp',
    component: 'Input',
    required: (callback) => {
      return callback.values['component'] == 'MICROAPP';
    },
    ifShow: (callback) => {
      return callback.values['component'] == 'MICROAPP';
    },
  },
  {
    field: 'redirect',
    label: t('routes.system.menu.redirect'),
    component: 'Input',
  },
  {
    field: 'icon',
    label: t('routes.system.menu.icon'),
    component: 'IconPicker',
  },
  {
    field: 'title',
    label: t('routes.system.menu.title'),
    component: 'Input',
    required: true,
  },

  {
    field: 'priority',
    label: t('routes.system.menu.priority'),
    component: 'InputNumber',
  },
  {
    field: 'parent',
    label: t('routes.system.menu.parent'),
    component: 'TreeSelect',
    componentProps: {
      fieldNames: {
        label: 'name',
        key: 'id',
        value: 'id',
      },
    },
    colProps: { span: 8 },
  },
];

// 编辑
export async function putEditMenuData(param: V1UpdateMenu) {
  const request: MenuServiceApiMenuServiceUpdateMenuRequest = {
    menuId: param.id,
    body: {
      menu: {
        ...param,
        meta: {
          icon: param.icon,
          microApp: param.microApp,
          title: param.title,
        },
      },
    },
  };
  return await new MenuServiceApi().menuServiceUpdateMenu(request);
}

// 增加
export async function postcreateMenuData(param: MenuServiceApiMenuServiceCreateMenuRequest) {
  return await new MenuServiceApi().menuServiceCreateMenu({
    body: param.body,
  });
}
// 删除
export async function postDeleteMenuData(param: MenuServiceApiMenuServiceDeleteMenuRequest) {
  return await new MenuServiceApi().menuServiceDeleteMenu({
    id: param.id,
  });
}
// 分页
export async function getMenuData(
  param?: BasicPageParams,
): Promise<BasicFetchResult<AppRouteRecordRaw>> {
  return await new MenuServiceApi()
    .menuServiceListMenu2({
      body: {
        pageSize: param ? param.pageSize : -1,
        pageOffset: param ? (param.page - 1) * param.pageSize : -1,
      },
    })
    .then((response) => {
      const menuData: BasicFetchResult<AppRouteRecordRaw> = {
        total: response.data.items!.length,
        items: transformObjToAppRouteRecordRaw(response.data.items!),
      };
      return menuData;
    });
}
