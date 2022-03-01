import { Component, h } from 'vue';

export function getDynamicComponent(param: string): Component {
  return {
    render(_ctx) {
      return h('div', {}, [
        h(
          'micro-app',
          {
            name: param,
            url: _ctx.url,
            baseroute: '/',
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
      return {
        url: 'http://localhost:4009/testpage',
        microAppData: { msg: 'base' },
      };
    },
    methods: {
      handleCreate(): void {
        console.log('child-vue3 创建了');
      },
      handleBeforeMount(): void {
        console.log('child-vue3 即将被渲染');
      },
      handleMount(): void {
        console.log('child-vue3 已经渲染完成');
        setTimeout(() => {
          // @ts-ignore
          this.microAppData = { msg: '来自基座的新数据' };
        }, 2000);
      },
      handleUnmount(): void {
        console.log('child-vue3 卸载了');
      },
      handleError(): void {
        console.log('child-vue3 加载出错了');
      },
      handleDataChange(e: CustomEvent): void {
        console.log('来自子应用 child-vue3 的数据:', e.detail.data);
      },
    },
  };
}
