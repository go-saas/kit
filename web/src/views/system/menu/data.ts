import { MenuServiceApi } from '/@/api-gen/api/menu-service-api';
import { V1Menu } from '/@/api-gen/models/v1-menu';
import { BasicFetchResult, BasicPageParams } from '/@/api/model/baseModel';
import { FormSchema } from '/@/components/Form/src/types/form';
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

export const formSchema: FormSchema[] = [];

export async function getMenuData(param: BasicPageParams): Promise<BasicFetchResult<V1Menu>> {
  return await new MenuServiceApi()
    .menuServiceListMenu2({
      body: {
        pageSize: param.pageSize,
        pageOffset: (param.page - 1) * param.pageSize,
      },
    })
    .then((response) => {
      const menuData: BasicFetchResult<V1Menu> = {
        total: response.data.items!.length,
        items: response.data.items!,
      };
      return menuData;
    });
}
