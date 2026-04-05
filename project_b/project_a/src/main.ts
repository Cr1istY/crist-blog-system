import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// PrimeVue v4
import PrimeVue from 'primevue/config';
import ToastService from 'primevue/toastservice';
import Aura from '@primevue/themes/aura'; // ✅ v4 新方式
import 'primeicons/primeicons.css';

const app = createApp(App)

router.afterEach((to) => {
  // 如果 meta.title 存在且非空，直接使用
  if (to.meta.title) {
    document.title = to.meta.title as string
  } else {
    // 否则保留默认标题（或从其他地方获取，如组件内）
    document.title = 'foreveryang.cn'
  }
})

app.use(createPinia())
app.use(router)

app.use(PrimeVue, {
  theme: {
    preset: Aura, // ✅ 使用预设主题
    options: {
      darkModeSelector: '.my-app-dark', // 可选：暗色模式选择器
      cssLayer: {
        name: 'primevue',
        order: 'tailwind-base, primevue, tailwind-utilities'
      }
    }
  },
  ripple: true
});

app.use(ToastService);
app.mount('#app')
