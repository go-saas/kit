"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[7132],{2648:(e,n,o)=>{o.r(n),o.d(n,{assets:()=>c,contentTitle:()=>a,default:()=>u,frontMatter:()=>s,metadata:()=>i,toc:()=>d});var t=o(1527),r=o(4744);const s={sidebar_label:"Background Job",title:"Background Job"},a=void 0,i={id:"learn/fundamentals/background-job",title:"Background Job",description:"A background job, also known as a background task or asynchronous task, refers to a process or piece of work that is executed independently and concurrently with the main execution flow of a software application or system.",source:"@site/docs/02-learn/01-fundamentals/09-background-job.mdx",sourceDirName:"02-learn/01-fundamentals",slug:"/learn/fundamentals/background-job",permalink:"/kit/zh-Hans/docs/learn/fundamentals/background-job",draft:!1,unlisted:!1,editUrl:"https://github.com/go-saas/kit/tree/main/docs/docs/02-learn/01-fundamentals/09-background-job.mdx",tags:[],version:"current",sidebarPosition:9,frontMatter:{sidebar_label:"Background Job",title:"Background Job"},sidebar:"tutorialSidebar",previous:{title:"Events",permalink:"/kit/zh-Hans/docs/learn/fundamentals/events"},next:{title:"Registry & Service Discovery",permalink:"/kit/zh-Hans/docs/learn/fundamentals/registry"}},c={},d=[{value:"Server",id:"server",level:2},{value:"UI",id:"ui",level:2}];function l(e){const n={a:"a",code:"code",h2:"h2",p:"p",pre:"pre",...(0,r.a)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(n.p,{children:"A background job, also known as a background task or asynchronous task, refers to a process or piece of work that is executed independently and concurrently with the main execution flow of a software application or system.\nYou can use background job to do tasks like sending emails, long-run data analysis."}),"\n",(0,t.jsxs)(n.p,{children:["The implementation of background job is based on ",(0,t.jsx)(n.a,{href:"https://github.com/hibiken/asynq",children:"asynq"})]}),"\n",(0,t.jsx)(n.h2,{id:"server",children:"Server"}),"\n",(0,t.jsxs)(n.p,{children:[(0,t.jsx)(n.a,{href:"https://github.com/go-saas/kit/blob/main/pkg/job/server.go",children:(0,t.jsx)(n.code,{children:"job.Server"})})," implements kratos ",(0,t.jsx)(n.code,{children:"transport.Server"}),"\nYou can apply middleware pattern to ",(0,t.jsx)(n.code,{children:"job.Server"})]}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-go",children:'import (\n\tklog "github.com/go-kratos/kratos/v2/log"\n\t"github.com/go-saas/kit/pkg/job"\n\t"github.com/go-saas/kit/user/private/biz"\n\t"github.com/hibiken/asynq"\n)\n\nfunc NewJobServer(opt asynq.RedisConnOpt, log klog.Logger, handlers []*job.Handler) *job.Server {\n    // set queue\n\tsrv := job.NewServer(opt, job.WithQueues(map[string]int{\n\t\tstring(biz.ConnName): 1,\n\t}))\n\tsrv.Use(job.TracingServer(), job.Logging(log))\n\tjob.RegisterHandlers(srv, handlers...)\n\treturn srv\n}\n\n'})}),"\n",(0,t.jsx)(n.h2,{id:"ui",children:"UI"}),"\n",(0,t.jsxs)(n.p,{children:["The default ui is exposed in ",(0,t.jsx)(n.a,{href:"../../modules/sys",children:"sys module"})]}),"\n",(0,t.jsx)("img",{src:"/kit/img/background-job-ui.png",alt:"background-job-ui"})]})}function u(e={}){const{wrapper:n}={...(0,r.a)(),...e.components};return n?(0,t.jsx)(n,{...e,children:(0,t.jsx)(l,{...e})}):l(e)}},4744:(e,n,o)=>{o.d(n,{Z:()=>i,a:()=>a});var t=o(959);const r={},s=t.createContext(r);function a(e){const n=t.useContext(s);return t.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function i(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:a(e.components),t.createElement(s.Provider,{value:n},e.children)}}}]);