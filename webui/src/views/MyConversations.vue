<template>
  <section class="p-4 max-w-lg mx-auto">
    <h2 class="text-2xl font-bold mb-4">Le mie conversazioni</h2>

    <div v-if="loading" class="text-center">Caricamento...</div>
    <div v-else-if="error" class="text-red-500">{{ error }}</div>
    <ul v-else>
      <li
        v-for="conv in conversations"
        :key="conv.id"
        class="p-3 mb-2 border rounded hover:bg-gray-100 cursor-pointer flex justify-between items-center"
        @click="openConversation(conv.id)"
      >
        <div>
          <span class="font-semibold">
            <!-- Se è un gruppo, mostra il nome; altrimenti il nome dell'altro utente -->
            {{ conv.isGroup ? conv.name : otherParticipantName(conv.participants) }}
          </span>
          <div class="text-sm text-gray-600">
            <!-- Anteprima del messaggio più recente -->
            {{ conv.lastMessage.text }}
          </div>
        </div>
        <div class="text-xs text-gray-500">
          {{ formatDate(conv.lastMessage.timestamp) }}
        </div>
      </li>
    </ul>
  </section>
</template>

<script>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

export default {
  setup() {
    const conversations = ref([])
    const loading = ref(false)
    const error = ref('')
    const router = useRouter()

    const loadConversations = async () => {
      loading.value = true
      error.value = ''
      try {
        const res = await axios.get('http://localhost:8080/conversations', {
          withCredentials: true
        })
        conversations.value = res.data
      } catch (e) {
        error.value = 'Impossibile caricare le conversazioni'
      } finally {
        loading.value = false
      }
    }

    const openConversation = (id) => {
      router.push({ name: 'conversation', params: { id } })
    }

    // Se non è un gruppo, trova il partecipante che non è l'utente corrente
    const otherParticipantName = (participants) => {
      const meId = participants.meId
      const other = participants.list.find(u => u.id !== meId)
      return other ? other.name : 'Chat privata'
    }

    // Formatta timestamp in modo semplice
    const formatDate = (iso) => {
      const d = new Date(iso)
      return d.toLocaleString('it-IT', {
        day: '2-digit',
        month: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    onMounted(loadConversations)

    return {
      conversations,
      loading,
      error,
      openConversation,
      otherParticipantName,
      formatDate
    }
  }
}
</script>

<style scoped>
ul {
  list-style: none;
  padding: 0;
}
</style>
