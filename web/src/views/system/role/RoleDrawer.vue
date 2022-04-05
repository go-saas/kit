<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="registerDrawer"
    showFooter
    :title="getTitle"
    width="50%"
    @ok="handleSubmit"
  >
    <BasicForm @register="registerForm" />
  </BasicDrawer>
</template>
<script lang="ts">
  import { useI18n } from '/@/hooks/web/useI18n';
  import { defineComponent, ref, computed, unref } from 'vue';
  import { BasicForm, useForm } from '/@/components/Form/index';
  import { formSchema, postcreateRoleData, putEditRoleData } from './data';
  import { BasicDrawer, useDrawerInner } from '/@/components/Drawer';
  import {
    RoleServiceApiRoleServiceCreateRoleRequest,
    RoleServiceApiRoleServiceUpdateRoleRequest,
  } from '/@/api-gen';
  export default defineComponent({
    name: 'RoleDrawer',
    components: { BasicDrawer, BasicForm },
    emits: ['success', 'register'],
    setup(_, { emit }) {
      const { t } = useI18n();
      const isUpdate = ref(true);
      const rowId = ref('');
      const [registerForm, { resetFields, setFieldsValue, validate }] = useForm({
        labelWidth: 100,
        schemas: formSchema,
        showActionButtonGroup: false,
        baseColProps: { lg: 12, md: 24 },
      });
      const [registerDrawer, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
        resetFields();
        setDrawerProps({ confirmLoading: false });
        isUpdate.value = !!data?.isUpdate;
        rowId.value = data.record.id;
        if (unref(isUpdate)) {
          setFieldsValue({
            ...data.record,
          });
        }
      });
      const getTitle = computed(() =>
        !unref(isUpdate) ? t('role.role.create') : t('role.role.edit'),
      );
      async function handleSubmit() {
        try {
          let values = await validate();
          setDrawerProps({ confirmLoading: true });
          // TODO custom api
          values = { ...values, id: rowId.value };
          if (unref(isUpdate)) {
            const paramEdit: RoleServiceApiRoleServiceUpdateRoleRequest = {
              roleId: rowId.value,
              body: {
                role: values,
              },
            };
            // 编辑
            await putEditRoleData(paramEdit as RoleServiceApiRoleServiceUpdateRoleRequest);
          } else if (!unref(isUpdate)) {
            const paramCreate: RoleServiceApiRoleServiceCreateRoleRequest = {
              body: values,
            };
            // 增加
            await postcreateRoleData(paramCreate as RoleServiceApiRoleServiceCreateRoleRequest);
          }
          closeDrawer();
          emit('success', { isUpdate: unref(isUpdate), values: values });
        } finally {
          setDrawerProps({ confirmLoading: false });
        }
      }
      return { registerDrawer, registerForm, getTitle, handleSubmit };
    },
  });
</script>
