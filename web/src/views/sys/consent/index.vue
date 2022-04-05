<template>
  <PageWrapper title="Consent">
    <template #headerContent>
      <div class="flex justify-between items-center">
        <span class="flex-1"> </span>
      </div>
      <div v-if="!hasChanllenge"> 参数错误，点击回到主页 </div>
    </template>
    <!-- TODO 显示当前用户的信息 -->
    <div v-if="consentInfoRef">
      {{ consentInfoRef.client?.id ?? '' }}
      {{ consentInfoRef.client?.name ?? '' }}
      {{ consentInfoRef.client?.logoUrl ?? '' }}
      请求权限
      <div v-for="item in consentInfoRef.requestedScope" :key="item"> {{ item }} </div>
      <Button type="primary" size="large" block @click="handleGrant" :loading="loading">
        同意
      </Button>
      <Button type="primary" size="large" block @click="handleReject" :loading="loading">
        拒绝
      </Button>
    </div>
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { PageWrapper } from '/@/components/Page';
  import { computed, onMounted, ref } from 'vue';
  import { useRoute } from 'vue-router';
  import { AuthWebApi, V1GetConsentResponse } from '/@/api-gen';
  import { useUserStore } from '/@/store/modules/user';
  import { router } from '/@/router';
  import { Button } from 'ant-design-vue';
  const { query } = useRoute();
  const { consent_challenge: consentChallenge } = query;
  const userStore = useUserStore();

  const hasLogin = computed(() => userStore.getIsLogin);
  console.log(hasLogin);
  //TODO 如果没有login，直接跳转到主页

  const loading = ref(false);

  const hasChanllenge = computed(() => consentChallenge?.toString() != null);

  const consentInfoRef = ref<Nullable<V1GetConsentResponse>>(null);
  onMounted(async () => {
    loading.value = true;
    try {
      const data = await new AuthWebApi().authWebGetConsent({
        consentChallenge: consentChallenge?.toString(),
      });
      consentInfoRef.value = data.data;
    } finally {
      loading.value = false;
    }
  });

  async function handleRedirect(redirect = '/') {
    if (redirect.indexOf('http://') === 0 || redirect.indexOf('https://') === 0) {
      location.replace(redirect);
    } else {
      await router.replace(redirect);
    }
  }

  async function handleReject() {
    const data = await new AuthWebApi().authWebGrantConsent({
      body: {
        challenge: consentChallenge?.toString(),
        reject: true,
        grantScope: consentInfoRef.value?.requestedScope,
      },
    });
    //check redirect
    await handleRedirect(data.data.redirect);
  }

  async function handleGrant() {
    const data = await new AuthWebApi().authWebGrantConsent({
      body: {
        challenge: consentChallenge?.toString(),
        reject: false,
        grantScope: consentInfoRef.value?.requestedScope,
      },
    });
    //check redirect
    await handleRedirect(data.data.redirect);
  }
</script>
