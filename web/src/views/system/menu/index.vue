<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <Authority :value="[{ namespace: 'sys.menu', resource: '*', action: 'create' }]">
          <a-button type="primary" @click="handleCreate">{{ t('menu.menu.create') }} </a-button>
        </Authority>
      </template>
      <template #action="{ record }">
        <TableAction
          :actions="[
            {
              icon: 'clarity:note-edit-line',
              onClick: handleEdit.bind(null, record),
              auth: [{ namespace: 'sys.menu', resource: '*', action: 'update' }],
            },
            {
              icon: 'ant-design:delete-outlined',
              color: 'error',
              popConfirm: {
                title: t('common.confirmDelete'),
                confirm: handleDelete.bind(null, record),
              },
              auth: [{ namespace: 'sys.menu', resource: '*', action: 'delete' }],
            },
          ]"
        />
      </template>
    </BasicTable>
    <MenuDrawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { BasicTable, useTable, TableAction } from '/@/components/Table';
  import { useDrawer } from '/@/components/Drawer';
  import { getMenuColumns, getMenuData, postDeleteMenuData } from './data';
  import MenuDrawer from './MenuDrawer.vue';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { Authority } from '/@/components/Authority';
  export default defineComponent({
    components: { BasicTable, TableAction, MenuDrawer, Authority },
    setup() {
      const { t } = useI18n();
      const [registerDrawer, { openDrawer }] = useDrawer();
      const [registerTable, { reload, updateTableDataRecord, getRawDataSource }] = useTable({
        title: t('menu.menu.list'),
        api: getMenuData,
        columns: getMenuColumns(),
        isTreeTable: true,
        striped: false,
        bordered: true,
        canResize: false,
        rowKey: 'id',
        actionColumn: {
          width: 80,
          title: t('common.operating'),
          dataIndex: 'action',
          slots: { customRender: 'action' },
          fixed: undefined,
        },
        handleSearchInfoFn(info) {
          console.log(info);
        },
      });
      function handleCreate(record: Recordable) {
        const rawData = getRawDataSource();
        openDrawer(true, {
          rawData,
          record,
          isUpdate: false,
        });
      }
      function handleEdit(record: Recordable) {
        const rawData = getRawDataSource();
        openDrawer(true, {
          rawData,
          record,
          isUpdate: true,
        });
      }
      function handleDelete(record: Recordable) {
        // 删除
        postDeleteMenuData({ id: record.id });
        reload();
      }
      async function handleSuccess({ isUpdate, values }) {
        if (isUpdate) {
          updateTableDataRecord(values.id, values);
        } else {
          await reload();
        }
      }
      return {
        handleCreate,
        handleEdit,
        handleDelete,
        handleSuccess,
        registerDrawer,
        registerTable,
        t,
      };
    },
  });
</script>
