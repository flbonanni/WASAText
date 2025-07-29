import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import './index.css'

axios.defaults.baseURL = '/api'

// Dì a axios di inviare i cookie se userai anche sessioni, ma per header non serve
axios.defaults.withCredentials = true

// Interceptor per aggiungere sempre l’Authorization header
axios.interceptors.request.use(config => {
  const userId = localStorage.getItem('userId')
  if (userId) {
    config.headers.Authorization = `Bearer ${userId}`
  }
  return config
})

const app = createApp(App)
app.use(router)
app.mount('#app')
