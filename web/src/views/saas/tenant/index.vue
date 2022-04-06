<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <Authority :value="[{ namespace: 'saas.tenant', resource: '*', action: 'create' }]">
          <a-button type="primary" @click="handleCreate">
            {{ t('saas.tenant.create') }}
          </a-button></Authority
        >
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
              disabled: record.separateDb,
              auth: [{ namespace: '*', resource: '*', action: '*' }],
            },
            {
              icon: 'ant-design:delete-outlined',
              color: 'error',
              popConfirm: {
                title: t('common.confirmDelete'),
                confirm: handleDelete.bind(null, record),
              },
              auth: [{ namespace: 'saas.tenant', resource: '*', action: 'delete' }],
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
  import { Authority } from '/@/components/Authority';
  export default defineComponent({
    components: { BasicTable, TableAction, TenantDrawer, Authority },
    setup() {
      const { t } = useI18n();
      const userStore = useUserStore();
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
      async function handleSuccess({ isUpdate, values }) {
        if (isUpdate) {
          updateTableDataRecord(values.id, values);
        } else {
          await reload();
        }
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
