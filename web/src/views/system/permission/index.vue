<template>
  <BasicTable @register="registerTable" />
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { onMounted } from 'vue';
  import { BasicTable, useTable } from '/@/components/Table';
  import { getPremissionColumns, getPremissionData } from './data';

  export default defineComponent({
    components: { BasicTable },
    setup() {
      const [registerTable] = useTable({
        api: getPremissionData,
        columns: getPremissionColumns(),
      });
      onMounted(() => {});
      return {
        registerTable,
      };
    },
  });
</script>
<style lang=""></style>
