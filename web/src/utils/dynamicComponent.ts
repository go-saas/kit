import { Component, h } from 'vue';
import { AppRouteRecordRaw } from '../router/types';
import { useUserStore } from '../store/modules/user';

export function getDynamicComponent(param: AppRouteRecordRaw): Component {
  return {
    render(_ctx) {
      return h('div', {}, [
        h(
          'micro-app',
          {
            name: param.name,
            url: _ctx.url,
            baseroute: '',
            data: _ctx.microAppData,
            onCreated: _ctx.handleCreate,
            onBeforemount: _ctx.handleBeforeMount,
            onMounted: _ctx.handleMount,
            onUnmount: _ctx.handleUnmount,
            onError: _ctx.handleError,
            onDatachange: _ctx.handleDataChange,
          },
          [],
        ),
      ]);
    },
    data() {
      const userStore = useUserStore();
      return {
        url: param.microApp,
        microAppData: { msg: userStore.getUserInfo },
      };
    },
    methods: {
      handleCreate(): void {
        console.log(param.title, '创建了');
      },
      handleBeforeMount(): void {
        console.log(param.title, '即将被渲染');
      },
      handleMount(): void {
        console.log(param.title, '已经渲染完成');
        setTimeout(() => {
          // @ts-ignore
          this.microAppData = { msg: '来自基座的新数据' };
        }, 2000);
      },
      handleUnmount(): void {
        console.log(param.title, '卸载了');
      },
      handleError(): void {
        console.log(param.title, '加载出错了');
      },
      handleDataChange(e: CustomEvent): void {
        console.log(`来自子应用${param.title}的数据:`, e.detail.data);
      },
    },
  };
}
