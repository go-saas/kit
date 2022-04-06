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
  import { formSchema, postcreateUser } from './data';
  import { BasicDrawer, useDrawerInner } from '/@/components/Drawer';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { UserServiceApiUserServiceCreateUserRequest } from '/@/api-gen';
  export default defineComponent({
    name: 'MenuDrawer',
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
        !unref(isUpdate) ? t('user.user.create') : t('user.user.edit'),
      );
      async function handleSubmit() {
        try {
          let values = await validate();
          setDrawerProps({ confirmLoading: true }); // TODO custom api

          values = { ...values, id: rowId.value };
          const createUserParamas: UserServiceApiUserServiceCreateUserRequest = {
            body: {
              avatar: values.avatar,
              password: values.password,
              confirmPassword: values.confirmPassword,
              name: values.name,
              username: values.username,
              phone: values.phone,
              email: values.email,
              gender: values.gender,
              birthday: values.birthday,
              rolesId: [...(values.roles ?? [])],
            },
          };
          console.log(createUserParamas);
          await postcreateUser(createUserParamas);
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
