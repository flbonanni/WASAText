<template>
  <section class="p-4 max-w-sm mx-auto">
    <h2 class="text-2xl mb-4">Login</h2>
    <input v-model="username" placeholder="Username" class="input" />
    <button @click="doLogin" class="btn mt-4">Entra</button>
  </section>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

export default {
  setup() {
    const username = ref('')
    const router = useRouter()

    function doLogin() {
      console.log('cliccato login, username =', username.value)
      if (!username.value.trim()) return
      // Salva l’ID utente in localStorage
      localStorage.setItem('userId', username.value.trim())
      // Vai alla dashboard
      router.push({ name: 'dashboard' })
    }

    return { username, doLogin }
  },
}
</script>