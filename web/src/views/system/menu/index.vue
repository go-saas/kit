<template>
  <BasicTable @register="registerTable" />
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { onMounted } from 'vue';
  import { BasicTable, useTable } from '/@/components/Table';
  import { getMenuColumns, getMenuData } from './data';

  export default defineComponent({
    components: { BasicTable },
    setup() {
      const [registerTable] = useTable({
        api: getMenuData,
        columns: getMenuColumns(),
      });
      onMounted(() => {});
      return {
        registerTable,
      };
    },
  });
</script>
