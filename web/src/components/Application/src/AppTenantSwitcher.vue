<!--
 * @Author: Vben
 * @Description: logo component
-->
<template>
  <div class="wrap">
    <div class="anticon" v-if="!loading">
      {{ getDisplayName(currentTenant!) }}
    </div>
    <Input.Search
      :placeholder="t('saas.tenant.name')"
      size="large"
      allowClear
      v-model:value="searchValue"
      @search="handleChangeTenant"
    >
      <template #enterButton>
        <Button> {{ t('saas.tenant.change') }} </Button>
      </template>
    </Input.Search>
  </div>
</template>
<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { Input, Button } from 'ant-design-vue';
  import { useUserStore } from '/@/store/modules/user';
  import { TenantServiceApi, V1GetCurrentTenantReply } from '/@/api-gen';
  import { useI18n } from '/@/hooks/web/useI18n';

  const { t } = useI18n();
  // const props = defineProps({});
  const loading = ref(true);
  const userStore = useUserStore();
  const currentTenant = ref<Nullable<V1GetCurrentTenantReply>>(null);
  const searchValue = ref('');
  onMounted(async () => {
    loading.value = true;
    const current = await new TenantServiceApi().tenantServiceGetCurrentTenant();
    currentTenant.value = current.data;
    searchValue.value = current.data?.tenant?.name ?? '';
    loading.value = false;
  });

  function handleChangeTenant() {
    const key = searchValue.value;
    userStore.changeTenant(key);
  }
  const getDisplayName = (tenanReply: V1GetCurrentTenantReply) => {
    if (tenanReply.isHost) {
      return t('saas.hostSide');
    }
    return tenanReply?.tenant?.displayName ?? '';
  };
</script>
<style lang="less" scoped>
  .wrap {
    margin-bottom: 80px;

    .anticon {
      font-weight: bold;
      font-size: 20px;
      display: inline-block;
      margin-bottom: 25px;
    }
  }
</style>
