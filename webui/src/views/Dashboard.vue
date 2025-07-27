<template>
  <section class="p-4 max-w-lg mx-auto">
    <h2 class="text-2xl font-bold mb-4">Dashboard</h2>

    <!-- Ricerca utente -->
    <div class="flex mb-6">
      <input
        v-model="searchTerm"
        @keyup.enter="searchUser"
        placeholder="Cerca utente..."
        class="flex-1 input"
      />
      <button @click="searchUser" class="btn ml-2">Cerca</button>
    </div>

    <!-- Risultati ricerca -->
    <div v-if="loading" class="text-center mb-4">Ricerca in corso…</div>
    <div v-else-if="error" class="text-red-500 mb-4">{{ error }}</div>
    <ul v-else class="space-y-2 mb-8">
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

    <!-- Pulsante My Conversations -->
    <button @click="goToConversations" class="btn">
      Le mie conversazioni
    </button>
  </section>
</template>

<script>
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

export default {
  setup() {
    const searchTerm = ref('')
    const results = ref([])
    const loading = ref(false)
    const error = ref('')
    const router = useRouter()
    const defaultAvatar = '/default-avatar.png'

    const searchUser = async () => {
      if (!searchTerm.value.trim()) return
      loading.value = true
      error.value = ''
      try {
        const res = await axios.get(
          `http://localhost:8080/users?search=${encodeURIComponent(searchTerm.value)}`,
          { withCredentials: true }
        )
        results.value = res.data  // array di utenti trovati
      } catch (e) {
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
.input {
  flex: 1;
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
