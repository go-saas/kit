<template>
  <BasicTable @register="registerTable" />
</template>
<script lang="ts">
  import { defineComponent } from 'vue';
  import { onMounted } from 'vue';
  import { BasicTable, useTable } from '/@/components/Table';
  import { getUserColumns, getUserData } from './data';

  export default defineComponent({
    components: { BasicTable },
    setup() {
      const [registerTable] = useTable({
        api: getUserData,
        columns: getUserColumns(),
      });
      onMounted(() => {});
      return {
        registerTable,
      };
    },
  });
</script>
<style lang=""></style>
