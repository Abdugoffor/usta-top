import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './app/router'
import './shared/styles/main.scss'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.mount('#app')

// Aktive tillarni ilovani ishga tushirishda bir marta yuklash
import { useLangStore } from '@/shared/stores/langStore'
const langStore = useLangStore()
langStore.fetchActiveLanguages()
