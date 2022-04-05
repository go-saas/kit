<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"
          >{{ t('permission.permission.create') }}
        </a-button>
      </template>
      <template #action="{ record }">
        <TableAction
          :actions="[
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
    <PermissionDrawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { BasicTable, useTable, TableAction } from '/@/components/Table';
  import { useDrawer } from '/@/components/Drawer';
  import { getPermissiontColumns, getPermissiontData, postDeletePermissiontData } from './data';
  import PermissionDrawer from './PermissionDrawer.vue';
  import { V1RemoveSubjectPermissionRequest } from '/@/api-gen';
  import { useI18n } from '/@/hooks/web/useI18n';
  export default defineComponent({
    components: { BasicTable, TableAction, PermissionDrawer },
    setup() {
      const { t } = useI18n();
      const [registerDrawer, { openDrawer }] = useDrawer();
      const [registerTable, { reload }] = useTable({
        title: t('permission.permission.list'),
        api: getPermissiontData,
        columns: getPermissiontColumns(),
        isTreeTable: true,
        striped: false,
        bordered: true,
        canResize: false,
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
        const deleteParameter: V1RemoveSubjectPermissionRequest = {
          namespace: record.namespace,
          resource: record.resource,
          action: record.action,
          subject: record.subject,
          effects: [record.effect],
          tenantId: record.tenantId,
        };
        // // 删除
        postDeletePermissiontData(deleteParameter as V1RemoveSubjectPermissionRequest);
        reload();
      }
      function handleSuccess() {
        reload();
      }
      // onMounted(() => {
      // });
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
