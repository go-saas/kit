<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"> {{ t('saas.tenant.create') }} </a-button>
      </template>
      <template #logo="{ record }">
        <img :src="record.logo?.url" alt="" />
      </template>
      <template #action="{ record }">
        <TableAction
          :actions="[
            { icon: 'clarity:note-edit-line', onClick: handleEdit.bind(null, record) },
            {
              icon: 'clarity:switch-line',
              onClick: handleChangeTenant.bind(null, record),
              tooltip: t('saas.tenant.change'),
              auth: [{ namespace: '*', resource: '*', action: '*', hostOnly: true }],
            },
            {
              icon: 'ant-design:delete-outlined',
              color: 'error',
              popConfirm: {
                title: t('common.confirmDelete'),
                confirm: handleDelete.bind(null, record),
              },
            },
          ]"
        />
      </template>
    </BasicTable>
    <TenantDrawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { BasicTable, useTable, TableAction } from '/@/components/Table';
  import { useDrawer } from '/@/components/Drawer';
  import { getTenantColumns, getTenantData, postDeleteTenantData } from './data';
  import { TenantServiceApiTenantServiceDeleteTenantRequest } from '/@/api-gen';
  import TenantDrawer from './TenantDrawer.vue';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useUserStore } from '/@/store/modules/user';
  export default defineComponent({
    components: { BasicTable, TableAction, TenantDrawer },
    setup() {
      const { t } = useI18n();
      const userStore = useUserStore();
      const defaultRequirement = {
        namespace: '*',
        resource: '*',
        action: '*',
      };
      const [registerDrawer, { openDrawer }] = useDrawer();
      const [registerTable, { reload, updateTableDataRecord }] = useTable({
        title: t('saas.tenant.list'),
        api: getTenantData,
        columns: getTenantColumns(),
        isTreeTable: true,
        striped: false,
        bordered: true,
        canResize: false,
        rowKey: 'id',
        actionColumn: {
          width: 120,
          title: t('common.operating'),
          requirement: defaultRequirement,
          slots: { customRender: 'action' },
          fixed: undefined,
        },
        handleSearchInfoFn(info) {
          console.log(info);
        },
      });
      function handleCreate(record: Recordable) {
        openDrawer(true, {
          record,
          isUpdate: false,
        });
      }
      function handleEdit(record: Recordable) {
        openDrawer(true, {
          record,
          isUpdate: true,
        });
      }
      function handleDelete(record: Recordable) {
        const deleteId: TenantServiceApiTenantServiceDeleteTenantRequest = {
          id: record.id,
        };
        // 删除
        postDeleteTenantData(deleteId as TenantServiceApiTenantServiceDeleteTenantRequest);
        reload();
      }
      function handleSuccess({ isUpdate, values }) {
        if (isUpdate) {
          updateTableDataRecord(values.id, values);
        }
        reload();
      }
      function handleChangeTenant(record: Recordable) {
        userStore.changeTenant(record.id!);
      }
      return {
        handleCreate,
        handleEdit,
        handleDelete,
        handleSuccess,
        registerDrawer,
        registerTable,
        t,
        handleChangeTenant,
      };
    },
  });
</script>
