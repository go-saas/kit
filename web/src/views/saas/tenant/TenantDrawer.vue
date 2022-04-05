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
  import { defineComponent, ref, computed, unref } from 'vue';
  import { BasicForm, useForm } from '/@/components/Form/index';
  import {
    formSchema,
    updateFormSchema,
    postcreateTenantData,
    putEditTenantData,
    getTenantDataDetail,
  } from './data';
  import { BasicDrawer, useDrawerInner } from '/@/components/Drawer';
  import { useI18n } from '/@/hooks/web/useI18n';
  import {
    TenantServiceApiTenantServiceCreateTenantRequest,
    TenantServiceApiTenantServiceUpdateTenantRequest,
  } from '/@/api-gen';
  export default defineComponent({
    name: 'MenuDrawer',
    components: { BasicDrawer, BasicForm },
    emits: ['success', 'register'],
    setup(_, { emit }) {
      const { t } = useI18n();
      const isUpdate = ref(true);
      const rowId = ref('');
      const [registerForm, { resetFields, setFieldsValue, validate, resetSchema }] = useForm({
        labelWidth: 100,
        schemas: [],
        showActionButtonGroup: false,
        baseColProps: { lg: 12, md: 24 },
      });
      const [registerDrawer, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
        resetFields();

        isUpdate.value = !!data?.isUpdate;
        rowId.value = data.record.id;
        if (unref(isUpdate)) {
          resetSchema(updateFormSchema);
          setFieldsValue({
            ...data.record,
          });
          setDrawerProps({ loading: true });
          try {
            //reload from api
            data = await getTenantDataDetail(data.record.id);
            setFieldsValue({
              ...data.data,
            });
          } finally {
            setDrawerProps({ loading: false });
          }
        } else {
          resetSchema(formSchema);
        }
      });
      const getTitle = computed(() =>
        !unref(isUpdate) ? t('saas.tenant.create') : t('saas.tenant.edit'),
      );
      async function handleSubmit() {
        try {
          let values = await validate();
          setDrawerProps({ confirmLoading: true });
          values = { ...values, id: rowId.value };
          if (unref(isUpdate)) {
            const param: TenantServiceApiTenantServiceUpdateTenantRequest = {
              tenantId: rowId.value,
              body: { tenant: values },
            };
            // 编辑
            await putEditTenantData(param as TenantServiceApiTenantServiceUpdateTenantRequest);
          } else if (!unref(isUpdate)) {
            const param: TenantServiceApiTenantServiceCreateTenantRequest = {
              body: values,
            };
            await postcreateTenantData(param as TenantServiceApiTenantServiceCreateTenantRequest);
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
