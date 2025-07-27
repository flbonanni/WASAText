<template>
  <section class="p-4 max-w-sm mx-auto">
    <h2 class="text-2xl mb-4">Login</h2>
    <input v-model="email" placeholder="Email" class="input" />
    <input
      v-model="password"
      type="password"
      placeholder="Password"
      class="input mt-2"
    />
    <button @click="doLogin" class="btn mt-4">Login</button>
    <p v-if="error" class="text-red-500 mt-2">{{ error }}</p>
  </section>
</template>

<script>
import { ref } from 'vue'
import axios from 'axios'

export default {
  setup() {
    const email = ref('')
    const password = ref('')
    const error = ref('')

    async function doLogin() {
      try {
        await axios.post(
          'http://localhost:8080/login',
          { email: email.value, password: password.value },
          { withCredentials: true }
        )
        window.location.href = '/app'
      } catch {
        error.value = 'Credenziali non valide'
      }
    }

    return { email, password, error, doLogin }
  },
}
</script>

<style scoped>
.input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 0.25rem;
}
.btn {
  padding: 0.5rem 1rem;
  background-color: #2563eb;
  color: white;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
}
.btn:hover {
  background-color: #1e40af;
}
</style>