<template>
  <ConfigProvider :locale="getAntdLocale">
    <AppProvider>
      <RouterView />
    </AppProvider>
  </ConfigProvider>
</template>
<script lang="ts" setup>
  import { ConfigProvider } from 'ant-design-vue';
  import { AppProvider } from '/@/components/Application';
  import { useTitle } from '/@/hooks/web/useTitle';
  import { useLocale } from '/@/locales/useLocale';
  import { onMounted } from 'vue';
  import { AuthApi } from '/@/api-gen';
  import { defHttp } from '/@/utils/http/axios';
  import { setDefaultAxiosInstance } from '@kit/api';
  // support Multi-language
  const { getAntdLocale } = useLocale();
  setDefaultAxiosInstance(defHttp.getAxios());
  // Listening to page changes and dynamically changing site titles
  useTitle();

  onMounted(() => {
    new AuthApi().authGetCsrfToken();
  });
</script>
