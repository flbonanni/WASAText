<template>
  <section class="container mx-auto mt-20 max-w-xs p-6 bg-white bg-opacity-10 rounded-lg">
    <h2 class="text-2xl mb-6 gradient-title text-center">Login</h2>

    <input
      v-model="username"
      placeholder="Username"
      class="input w-full mb-4"
    />

    <button
      @click="doLogin"
      class="btn w-full py-2"
    >
      Entra
    </button>
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
      if (!username.value.trim()) return
      localStorage.setItem('userId', username.value.trim())
      router.push({ name: 'dashboard' })
    }

    return { username, doLogin }
  },
}
</script>

<style scoped>
/* Se vuoi un background semitrasparente per la box */
section {
  backdrop-filter: blur(10px);
}

/* Se vuoi override specifici per input & btn solo qui */
.input {
  margin-bottom: 1rem;
}

.btn {
  margin-top: 0.5rem;
}
</style>