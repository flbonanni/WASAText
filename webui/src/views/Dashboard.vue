<template>
  <section class="container mx-auto mt-12 max-w-lg p-4">
    <!-- Titolo centrato con gradiente -->
    <h2 class="text-3xl font-bold mb-6 gradient-title text-center">Dashboard</h2>

    <!-- Saluto -->
    <p v-if="user.name" class="text-center mb-6">Benvenuto, {{ user.name }}!</p>

    <!-- Ricerca utente -->
    <div class="max-w-md mx-auto mb-8">
      <input
        v-model="searchTerm"
        @keyup.enter="searchUser"
        placeholder="Cerca utente..."
        class="input w-full mb-4"
      />
      <button @click="searchUser" class="btn w-full mb-2">Cerca</button>
      <button @click="goToConversations" class="btn w-full">Le mie conversazioni</button>
    </div>

    <!-- Risultati ricerca -->
    <div v-if="loading" class="text-center mb-4">Ricerca in corso…</div>
    <div v-else-if="error" class="text-red-500 mb-4">{{ error }}</div>
    <ul v-else class="space-y-2 mb-8 max-w-md mx-auto">
      <li
        v-for="u in results"
        :key="u.id"
        class="p-3 border rounded hover:bg-gray-100 cursor-pointer flex items-center"
        @click="goToUser(u.id)"
      >
        <img
          :src="u.avatarUrl || defaultAvatar"
          class="w-8 h-8 rounded-full mr-3"
          alt="Avatar"
        />
        <div>
          <p class="font-semibold">{{ u.name }}</p>
          <p class="text-sm text-gray-600">{{ u.email }}</p>
        </div>
      </li>
    </ul>
  </section>
</template>

<script>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { getCurrentUser } from '@/services/userApi'

export default {
  setup() {
    const router = useRouter()
    const user       = ref({})       // <-- utente loggato
    const searchTerm = ref('')
    const results    = ref([])
    const loading    = ref(false)
    const error      = ref('')
    const defaultAvatar = '/default-avatar.png'

    // 1) Carica il profilo corrente
    onMounted(async () => {
      try {
        const res = await getCurrentUser()
        user.value = res.data
      } catch {
        // se non sei autenticato, torna al login
        return router.push({ name: 'login' })
      }
    })

    // 2) Funzioni originali
    const searchUser = async () => {
      if (!searchTerm.value.trim()) return
      loading.value = true
      error.value   = ''
      try {
        const res = await axios.get('/users?search=' + encodeURIComponent(searchTerm.value))
        results.value = res.data
      } catch {
        error.value = 'Errore durante la ricerca'
      } finally {
        loading.value = false
      }
    }

    const goToUser = (id) => {
      router.push({ name: 'user-profile', params: { id } })
    }

    const goToConversations = () => {
      router.push({ name: 'my-conversations' })
    }

    return {
      user,
      searchTerm,
      results,
      loading,
      error,
      defaultAvatar,
      searchUser,
      goToUser,
      goToConversations
    }
  },
}
</script>

<style scoped>
/* ... gli stessi stili di prima ... */
</style>