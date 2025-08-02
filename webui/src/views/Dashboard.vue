<template>
  <section class="container mx-auto mt-12 max-w-lg p-4">
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
      <button @click="searchUser" class="btn w-full mb-4">Cerca</button>
      <button @click="goToConversations" class="btn w-full">Le mie conversazioni</button>
    </div>

    <!-- Risultati ricerca -->
    <div v-if="loading" class="text-center mb-4">Ricerca in corso…</div>
    <div v-else-if="error" class="text-red-500 mb-4">{{ error }}</div>
    <ul v-else class="space-y-2 mb-8 max-w-md mx-auto">
      <!-- ... list items ... -->
    </ul>
  </section>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { getCurrentUser } from '@/services/userApi'

export default {
  setup() {
    const router = useRouter()
    const user = ref({})
    const searchTerm = ref('')
    const results = ref([])
    const loading = ref(false)
    const error = ref('')
    const defaultAvatar = '/default-avatar.png'

    onMounted(async () => {
      try {
        const res = await getCurrentUser()
        user.value = res.data
      } catch {
        router.push({ name: 'login' })
      }
    })

    const searchUser = async () => {
      /* ... */
    }
    const goToConversations = () => { /* ... */ }

    return {
      user, searchTerm, results, loading, error,
      defaultAvatar, searchUser, goToConversations
    }
  }
}
</script>
