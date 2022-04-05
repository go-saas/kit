import { h } from 'vue';
import {
  MenuServiceApi,
  MenuServiceApiMenuServiceCreateMenuRequest,
  MenuServiceApiMenuServiceDeleteMenuRequest,
  MenuServiceApiMenuServiceUpdateMenuRequest,
} from '/@/api-gen/api/menu-service-api';
import { V1UpdateMenu } from '/@/api-gen/models/v1-update-menu';
import { BasicFetchResult, BasicPageParams } from '/@/api/model/baseModel';
import { FormSchema } from '/@/components/Form/src/types/form';
import Icon from '/@/components/Icon';
import { BasicColumn } from '/@/components/Table/src/types/table';
import { useI18n } from '/@/hooks/web/useI18n';
import { transformObjToAppRouteRecordRaw } from '/@/router/helper/routeHelper';
import { AppRouteRecordRaw } from '/@/router/types';
const { t } = useI18n();
export const defaultRequirement = {
  namespace: '*',
  resource: '*',
  action: '*',
};

export function getMenuColumns(): BasicColumn[] {
  return [
    {
      title: t('routes.system.menu.name'),
      dataIndex: 'name',
      requirement: defaultRequirement,
      align: 'left',
    },
    {
      title: t('routes.system.menu.icon'),
      dataIndex: 'icon',
      width: 50,
      requirement: defaultRequirement,
      customRender: ({ record }) => {
        return h(Icon, { icon: record.meta.icon });
      },
    },
    {
      title: t('routes.system.menu.component'),
      dataIndex: 'component',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.menu.path'),
      dataIndex: 'path',
      requirement: defaultRequirement,
    },
    {
      title: t('routes.system.menu.redirect'),
      dataIndex: 'redirect',
      requirement: defaultRequirement,
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
    component: 'Input',
    required: true,
  },
  {
    field: 'path',
    label: t('routes.system.menu.path'),
    component: 'Input',
    required: true,
  },
  {
    field: 'redirect',
    label: t('routes.system.menu.redirect'),
    component: 'Input',
  },
  {
    field: 'icon',
    label: '图标',
    component: 'IconPicker',
  },
  {
    field: 'title',
    label: 'title',
    component: 'Input',
  },
  {
    field: 'microApp',
    label: 'microApp',
    component: 'Input',
  },
  {
    field: 'priority',
    label: 'priority',
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
