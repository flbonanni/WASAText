<template>
  <section class="p-4">
    <h2 class="text-xl font-bold mb-4">Chat</h2>

    <div v-for="msg in messages" :key="msg.id" class="flex items-center mb-2">
      <div class="flex-1">{{ msg.text }}</div>
      <MessageOptions
        :isOwner="msg.userId === currentUser.id"
        :isCommented="!!msg.comment"
        @forward="() => doForward(msg)"
        @delete="() => doDelete(msg)"
        @comment="() => doComment(msg)"
        @uncomment="() => doUncomment(msg)"
      />
    </div>

    <div class="mt-4 flex">
      <input
        v-model="newText"
        class="flex-1 input"
        placeholder="Scrivi un messaggio..."
      />
      <button @click="send" class="btn">Invia</button>
    </div>
  </section>
</template>

<script>
import axios from 'axios'
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MessageOptions from '@/components/MessageOptions.vue'
import { getCurrentUser } from '@/services/userApi'

export default {
  components: { MessageOptions },
  setup() {
    const route        = useRoute()
    const router       = useRouter()
    const convId       = route.params.id
    const currentUser  = ref(null)  // <-- utente loggato
    const messages     = ref([])
    const newText      = ref('')

    // 1) Carica il profilo e poi i messaggi
    onMounted(async () => {
      try {
        const me = await getCurrentUser()
        currentUser.value = me.data
      } catch {
        return router.push({ name: 'login' })
      }
      await load()
    })

    // 2) Load conversazione
    const load = async () => {
      const res = await axios.get(
        `http://localhost:8080/conversations/${convId}`,
        { withCredentials: true }
      )
      messages.value = res.data.messages
    }

    // 3) Invia messaggio
    const send = async () => {
      if (!newText.value) return
      await axios.post(
        `http://localhost:8080/conversations/${convId}/messages`,
        { text: newText.value },
        { withCredentials: true }
      )
      newText.value = ''
      await load()
    }

    // 4) Opzioni messaggio
    const doForward = async (msg) => { /* ... */ }
    const doDelete  = async (msg) => {
      await axios.delete(
        `http://localhost:8080/conversations/${convId}/messages/${msg.id}`,
        { withCredentials: true }
      )
      await load()
    }
    const doComment   = async (msg) => { /* ... */ }
    const doUncomment = async (msg) => { /* ... */ }

    return {
      currentUser,
      messages,
      newText,
      send,
      doForward,
      doDelete,
      doComment,
      doUncomment
    }
  }
}
</script>

<style scoped>
/* ... gli stessi stili di prima ... */
</style>