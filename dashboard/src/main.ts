import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import { addCollection } from '@iconify/vue'
import lucideIcons from '@iconify-json/lucide/icons.json'

// Register lucide icons locally for offline use
addCollection(lucideIcons)

createApp(App).use(router).mount('#app')
