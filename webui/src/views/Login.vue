<template>
  <!-- Aggiungi mt-20 per distanziare dal top, riduci il max-width -->
  <section class="mt-20 mx-auto w-full max-w-xs p-6 bg-white bg-opacity-10 rounded-lg">
    <!-- Mantieni il titolo, erediterà gradient-title dal tuo Header se lo usi qui -->
    <h2 class="text-2xl mb-6 text-center gradient-title">Login</h2>

    <!-- Input centrato e non appiccicato al bordo -->
    <input
      v-model="username"
      placeholder="Username"
      class="input mb-4"
    />

    <!-- Maggior padding sul bottone -->
    <button
      @click="doLogin"
      class="btn px-6 py-2 block mx-auto"
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