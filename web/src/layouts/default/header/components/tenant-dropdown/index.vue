<template>
  <Dropdown placement="bottomLeft" :overlayClassName="`${prefixCls}-dropdown-overlay`">
    <span :class="[prefixCls, `${prefixCls}--${theme}`]" class="flex">
      <img :class="`${prefixCls}__header`" :src="getCurrentTenantInfo.logo" />
      <span :class="`${prefixCls}__info hidden md:block`">
        <span :class="`${prefixCls}__name  `" class="truncate">
          {{ getCurrentTenantInfo.name }}
        </span>
      </span>
    </span>

    <template #overlay v-if="showDropdown">
      <Menu @click="handleChangeTenant">
        <div style="padding: 10px">{{ t('saas.tenant.change') }}</div>
        <template v-for="(item, index) in getAvailableTenants" :key="item.tenant?.id ?? ''">
          <template v-if="index > 0"><MenuDivider /></template>
          <MenuItem
            :id="item.tenant?.id ?? ''"
            :name="item.tenant?.name ?? ''"
            :displayName="getDisplayName(item)"
            :logoUrl="item.tenant?.logo?.url ?? headerImg"
          />
        </template>
      </Menu>
    </template>
  </Dropdown>
</template>
<script lang="ts">
  // components
  import { Dropdown, Menu } from 'ant-design-vue';
  import type { UserTenantInfo } from '/#/store';
  import { defineComponent, computed } from 'vue';

  import { useUserStore } from '/@/store/modules/user';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useDesign } from '/@/hooks/web/useDesign';
  import { propTypes } from '/@/utils/propTypes';
  import headerImg from '/@/assets/images/header.jpg';
  import { createAsyncComponent } from '/@/utils/factory/createAsyncComponent';

  export default defineComponent({
    name: 'TenantDropdown',
    components: {
      Dropdown,
      Menu,
      MenuItem: createAsyncComponent(() => import('./DropMenuItem.vue')),
      MenuDivider: Menu.Divider,
    },
    props: {
      theme: propTypes.oneOf(['dark', 'light']),
    },
    setup() {
      const { prefixCls } = useDesign('header-tenant-dropdown');
      const { t } = useI18n();
      const userStore = useUserStore();

      const getCurrentTenantInfo = computed(() => {
        let name = '';
        let logo: string | undefined = undefined;
        if (userStore.getCurrentIsHost) {
          name = t('saas.hostSide');
          logo = headerImg;
        } else {
          name = userStore.getCurrentTenant?.tenant?.displayName ?? '';
          logo = userStore.getCurrentTenant?.tenant?.logo?.url;
          logo = logo || headerImg;
        }
        return { name, logo };
      });

      const getAvailableTenants = computed(() => {
        let ret = userStore.getUserInfo?.tenants ?? [];
        ret = ret.filter((x) => x.tenant?.id != userStore.getCurrentTenant.tenant?.id);
        return ret;
      });

      const showDropdown = computed(() => getAvailableTenants.value.length > 0);

      function handleChangeTenant(e: { key?: string }) {
        const { key } = e;
        userStore.changeTenant(key);
      }
      const getDisplayName = (userTenant: UserTenantInfo) => {
        if (userTenant.isHost) {
          return t('saas.hostSide');
        }
        return userTenant?.tenant?.displayName ?? '';
      };
      return {
        prefixCls,
        getCurrentTenantInfo,
        getAvailableTenants,
        handleChangeTenant,
        headerImg,
        showDropdown,
        getDisplayName,
        t,
      };
    },
  });
</script>
<style lang="less">
  @prefix-cls: ~'@{namespace}-header-tenant-dropdown';

  .@{prefix-cls} {
    height: @header-height;
    padding: 0 0 0 10px;
    padding-right: 10px;
    overflow: hidden;
    font-size: 12px;
    cursor: pointer;
    align-items: center;

    img {
      width: 24px;
      height: 24px;
      margin-right: 12px;
    }

    &__header {
      border-radius: 50%;
    }

    &__name {
      font-size: 14px;
    }

    &--dark {
      &:hover {
        background-color: @header-dark-bg-hover-color;
      }
    }

    &--light {
      &:hover {
        background-color: @header-light-bg-hover-color;
      }

      .@{prefix-cls}__name {
        color: @text-color-base;
      }

      .@{prefix-cls}__desc {
        color: @header-light-desc-color;
      }
    }

    &-dropdown-overlay {
      .ant-dropdown-menu-item {
        min-width: 160px;
      }
    }
  }
</style>
