<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <Authority :value="[{ namespace: 'user.role', resource: '*', action: 'create' }]">
          <a-button type="primary" @click="handleCreate">
            {{ t('user.user.create') }}
          </a-button>
        </Authority>
      </template>
      <template #avatar="{ record }">
        <img :src="record.avatar?.url" alt="" />
      </template>
      <template #roles="{ record }">
        <Tag
          v-for="role in record.roles"
          :key="role.id"
          :color="role.isPreserved ? 'green' : 'geekblue'"
        >
          {{ getRoleName(role) }}
        </Tag>
      </template>
      <template #action="{ record }">
        <TableAction
          :actions="[
            {
              icon: 'clarity:note-edit-line',
              onClick: handleEdit.bind(null, record),
              disabled: record.isPreserved ?? false,
            },
            // {
            //   icon: 'ant-design:delete-outlined',
            //   color: 'error',
            //   popConfirm: {
            //     title: t('common.confirmDelete'),
            //     confirm: handleDelete.bind(null, record),
            //   },
            // },
          ]"
        />
      </template>
      <!-- <template #action="{ record }">
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
      </template> -->
    </BasicTable>
    <UserDrawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { BasicTable, useTable, TableAction } from '/@/components/Table';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useDrawer } from '/@/components/Drawer';
  import { getUsertColumns, getUsertData } from './data';
  import { getRoleName } from '../role/data';
  import UserDrawer from './UserDrawer.vue';
  import { Tag } from 'ant-design-vue';
  import { Authority } from '/@/components/Authority';
  export default defineComponent({
    components: { BasicTable, UserDrawer, Tag, TableAction, Authority },
    setup() {
      const { t } = useI18n();
      const [registerDrawer, { openDrawer }] = useDrawer();
      const [registerTable, { reload, updateTableDataRecord }] = useTable({
        title: t('user.user.list'),
        api: getUsertData,
        columns: getUsertColumns(),
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
      function handleEdit(_record: Recordable) {
        // openDrawer(true, {
        //   record,
        //   isUpdate: true,
        // });
      }
      // function handleDelete(record: Recordable) {
      //   // const deleteParameter: V1RemoveSubjectUserRequest = {
      //   //   namespace: record.namespace,
      //   //   resource: record.resource,
      //   //   action: record.action,
      //   //   subject: record.subject,
      //   //   effects: [record.effect],
      //   //   tenantId: record.tenantId,
      //   // };
      //   // // 删除
      // }
      async function handleSuccess({ isUpdate, values }) {
        if (isUpdate) {
          updateTableDataRecord(values.id, values);
        } else {
          await reload();
        }
      }
      // onMounted(() => {
      // });
      return {
        registerDrawer,
        handleCreate,
        handleSuccess,
        registerTable,
        t,
        getRoleName,
        handleEdit,
      };
    },
  });
</script>
