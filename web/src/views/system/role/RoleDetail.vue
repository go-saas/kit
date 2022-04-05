<template>
  <div>
    <PageWrapper
      :content="content"
      :title="t('role.role.permissions')"
      contentBackground
      @back="goBack"
    >
      <template #footer>
        <a-card :bordered="false" class="mt-5">
          <div v-for="(item, index) in permissionList" :key="index">
            <a-card :title="t(item.displayName ?? '')" :bordered="true" class="mt-5">
              <div v-for="(def, defIndex) in item.def" :key="defIndex">
                <a-descriptions :column="3">
                  <a-descriptions-item :label="t('routes.system.permission.namespace')"
                    >{{ def.namespace }}
                  </a-descriptions-item>
                  <a-descriptions-item :label="t('routes.system.permission.action')">
                    {{ def.action }}
                  </a-descriptions-item>
                  <a-descriptions-item :label="t('routes.system.permission.action')">
                    <a-switch
                      :disabled="!editable"
                      :checked="def?.granted ?? false"
                      @change="(checked) => onChange(def, checked)"
                    />
                  </a-descriptions-item>
                </a-descriptions>
              </div>
            </a-card>
          </div>
        </a-card>
        <!-- <a-button style="margin-bottom: 20px" type="primary" @click="handlModify"
          >{{ t('role.role.editPermissions') }}
        </a-button> -->
      </template>
    </PageWrapper>
  </div>
</template>

<script lang="ts">
  import { defineComponent, onMounted, ref } from 'vue';
  import { Card, Descriptions, Switch } from 'ant-design-vue';
  import { PageWrapper } from '/@/components/Page';
  import { useGo } from '/@/hooks/web/usePage';
  import { useTabs } from '/@/hooks/web/useTabs';
  import { getRolePermission, putRolePermission, getRoleDetail } from './DetailData';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useRouter } from 'vue-router';
  import { V1PermissionDefGroup, V1PermissionDef, V1UpdateRolePermissionAcl } from '/@/api-gen';
  export default defineComponent({
    components: {
      PageWrapper,
      [Card.name]: Card,
      [Descriptions.name]: Descriptions,
      [Descriptions.Item.name]: Descriptions.Item,
      [Switch.name]: Switch,
    },
    setup() {
      // const table = useTableContext();
      const { t } = useI18n();
      const { currentRoute } = useRouter();
      const go = useGo();

      const content = ref('');
      const permissionList = ref<V1PermissionDefGroup[]>([]);
      const roleId = currentRoute.value.params.id as string;
      const { setTitle } = useTabs();
      const editable = ref(false);
      setTitle(t('role.role.editPermissions'));

      function goBack() {
        go('/system/roles');
      }
      async function handlModify() {
        const results: V1UpdateRolePermissionAcl[] = permissionList.value
          .flatMap((p) => p.def ?? [])
          .filter((p) => p.granted)
          .map((p) => {
            return { namespace: p.namespace!, resource: '*', action: p.action!, effect: 'GRANT' };
          });
        await putRolePermission({ id: roleId, body: { acl: results } });
      }
      onMounted(() => {
        if (roleId && roleId != '') {
          getdata();
        }
      });
      async function getdata() {
        const role = await getRoleDetail({ id: roleId });
        editable.value = !(role.data.isPreserved ?? false);
        const res = await getRolePermission({ id: roleId });
        permissionList.value = res.data.defGroups ?? [];

        if (permissionList.value.length == 0) {
          content.value = t('role.role.noPermissions');
        } else {
          content.value = '';
        }
      }
      function handleSuccess({ isUpdate, values }) {
        console.log(isUpdate, values);
      }
      async function onChange(def: V1PermissionDef, checked: boolean) {
        def.granted = checked;
        try {
          //TODO error handling
          await handlModify();
        } catch (e) {
          def.granted = !checked;
        }
      }

      return {
        t,
        permissionList,
        goBack,
        getdata,
        handlModify,
        handleSuccess,
        content,
        onChange,
        editable,
      };
    },
  });
</script>

<style scoped></style>
