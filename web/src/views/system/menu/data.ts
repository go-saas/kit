import { MenuServiceApi } from '/@/api-gen/api/menu-service-api';
import { V1Menu } from '/@/api-gen/models/v1-menu';
import { BasicFetchResult } from '/@/api/model/baseModel';
import { BasicColumn } from '/@/components/Table/src/types/table';

const defaultRequirement = {
  namespace: '*',
  resource: '*',
  action: '*',
};

export function getMenuColumns(): BasicColumn[] {
  return [
    {
      title: 'Name',
      dataIndex: 'name',
      requirement: defaultRequirement,
    },
    {
      title: 'Component',
      dataIndex: 'component',
      requirement: defaultRequirement,
    },
    {
      title: 'Path',
      dataIndex: 'path',
      requirement: defaultRequirement,
    },
    {
      title: 'Redirect',
      dataIndex: 'redirect',
      requirement: defaultRequirement,
    },
    {
      title: 'Parent',
      dataIndex: 'parent',
      requirement: defaultRequirement,
    },
  ];
}

export async function getMenuData(): Promise<BasicFetchResult<V1Menu>> {
  return await new MenuServiceApi().menuServiceGetAvailableMenus().then((response) => {
    const menuData: BasicFetchResult<V1Menu> = {
      total: response.data.items!.length,
      items: response.data.items!,
    };
    return menuData;
  });
}
