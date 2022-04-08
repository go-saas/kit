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
  import { formSchema, postcreateMenuData, putEditMenuData } from './data';
  import { BasicDrawer, useDrawerInner } from '/@/components/Drawer';
  import { V1UpdateMenu, MenuServiceApiMenuServiceCreateMenuRequest } from '/@/api-gen';

  export default defineComponent({
    name: 'MenuDrawer',
    components: { BasicDrawer, BasicForm },
    emits: ['success', 'register'],
    setup(_, { emit }) {
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
            iframe: data.record?.meta?.frameSrc,
          });
        }
        const treeData = data.rawData.items;
        updateSchema({
          field: 'parent',
          componentProps: { treeData },
        });
      });
      const getTitle = computed(() => (!unref(isUpdate) ? '新增菜单' : '编辑菜单'));
      async function handleSubmit() {
        try {
          let values = await validate();
          setDrawerProps({ confirmLoading: true });

          values = { ...values, id: rowId.value };
          if (unref(isUpdate)) {
            const params = values as V1UpdateMenu;
            await putEditMenuData(params);
          } else if (!unref(isUpdate)) {
            // 增加
            const param: MenuServiceApiMenuServiceCreateMenuRequest = {
              body: values,
            };
            await postcreateMenuData(param as MenuServiceApiMenuServiceCreateMenuRequest);
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
