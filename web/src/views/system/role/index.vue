<template>
  <div>
    <BasicTable @register="registerTable">
      <template #name="{ record }">
        <Tag :color="record.isPreserved ? 'green' : 'geekblue'">
          {{ getRoleName(record) }}
        </Tag>
      </template>
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"> {{ t('role.role.create') }} </a-button>
      </template>
      <template #action="{ record }">
        <TableAction
          :actions="[
            {
              icon: 'clarity:info-standard-line',
              onClick: handleView.bind(null, record),
            },
            {
              icon: 'clarity:note-edit-line',
              onClick: handleEdit.bind(null, record),
              disabled: record.isPreserved ?? false,
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
    <RoleDrawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { BasicTable, useTable, TableAction } from '/@/components/Table';
  import { useDrawer } from '/@/components/Drawer';
  import { getRoleColumns, getRoleData, postDeleteRoleData, getRoleName } from './data';
  import RoleDrawer from './RoleDrawer.vue';
  import { RoleServiceApiRoleServiceDeleteRoleRequest } from '/@/api-gen';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useGo } from '/@/hooks/web/usePage';
  import { Tag } from 'ant-design-vue';
  export default defineComponent({
    components: { BasicTable, TableAction, RoleDrawer, Tag },
    setup() {
      const go = useGo();
      const { t } = useI18n();
      const defaultRequirement = {
        namespace: '*',
        resource: '*',
        action: '*',
      };
      const [registerDrawer, { openDrawer }] = useDrawer();
      const [registerTable, { reload, updateTableDataRecord }] = useTable({
        title: t('role.role.list'),
        api: getRoleData,
        columns: getRoleColumns(),
        isTreeTable: true,
        striped: false,
        bordered: true,
        canResize: false,
        rowKey: 'id',
        actionColumn: {
          width: 120,
          title: t('common.operating'),
          requirement: defaultRequirement,
          dataIndex: 'action',
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
        const deleteId: RoleServiceApiRoleServiceDeleteRoleRequest = {
          id: record.id,
        };
        // // 删除
        postDeleteRoleData(deleteId as RoleServiceApiRoleServiceDeleteRoleRequest);
        reload();
      }
      async function handleSuccess({ isUpdate, values }) {
        if (isUpdate) {
          updateTableDataRecord(values.id, values);
        }
        reload();
      }
      // /system/role/AccountDetail
      // onMounted(() => {
      // });
      function handleView(record: Recordable) {
        go(`/role/${record.id}/detail`);
      }
      return {
        handleView,
        handleCreate,
        handleEdit,
        handleDelete,
        handleSuccess,
        registerDrawer,
        registerTable,
        t,
        getRoleName,
      };
    },
  });
</script>
