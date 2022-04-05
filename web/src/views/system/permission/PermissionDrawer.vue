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
  import { formSchema, getPermissiontData, postcreatePermissiontData } from './data';
  import { BasicDrawer, useDrawerInner } from '/@/components/Drawer';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { V1AddSubjectPermissionRequest } from '/@/api-gen';
  export default defineComponent({
    name: 'MenuDrawer',
    components: { BasicDrawer, BasicForm },
    emits: ['success', 'register'],
    setup(_, { emit }) {
      const { t } = useI18n();
      const isUpdate = ref(true);
      const rowId = ref('');
      const [registerForm, { resetFields, setFieldsValue, updateSchema, validate }] = useForm({
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
        const subjectsdata = {
          subjects: [],
        };
        const treeData = await getPermissiontData(subjectsdata);
        updateSchema({
          field: 'parentMenu',
          componentProps: { treeData },
        });
      });
      const getTitle = computed(() =>
        !unref(isUpdate) ? t('menu.menu.create') : t('menu.menu.edit'),
      );
      async function handleSubmit() {
        try {
          let values = await validate();
          setDrawerProps({ confirmLoading: true });
          // TODO custom api
          values = { ...values, id: rowId.value };
          if (!unref(isUpdate)) {
            const param: V1AddSubjectPermissionRequest = values;
            postcreatePermissiontData(param as V1AddSubjectPermissionRequest);
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
