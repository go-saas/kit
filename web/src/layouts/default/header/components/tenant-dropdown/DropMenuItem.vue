<template>
  <MenuItem :key="itemKey">
    <span class="flex items-center">
      <img :src="logoUrl" class="mr-1" />
      <span>{{ displayName }}</span>
    </span>
  </MenuItem>
</template>
<script lang="ts">
  import { Menu } from 'ant-design-vue';

  import { computed, defineComponent } from 'vue';
  import { propTypes } from '/@/utils/propTypes';

  export default defineComponent({
    name: 'TenantDropdownMenuItem',
    components: { MenuItem: Menu.Item },
    props: {
      id: propTypes.string,
      name: propTypes.string,
      displayName: propTypes.string,
      logoUrl: propTypes.string,
    },
    setup(props) {
      const itemKey = computed(() => props.id);
      return { itemKey };
    },
  });
</script>

<style lang="less" scoped>
  img {
    width: 24px;
    height: 24px;
    margin-right: 12px;
  }
</style>
