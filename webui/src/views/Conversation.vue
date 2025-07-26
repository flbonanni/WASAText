<template>
  <section class="p-4">
    <h2 class="text-xl font-bold mb-4">Chat</h2>
    <div v-for="msg in messages" :key="msg.id" class="flex items-center mb-2">
      <div class="flex-1">{{ msg.text }}</div>
      <MessageOptions
        :isOwner="msg.userId === user.id"
        :isCommented="!!msg.comment"
        @forward="() => doForward(msg)"
        @delete="() => doDelete(msg)"
        @comment="() => doComment(msg)"
        @uncomment="() => doUncomment(msg)"
      />
    </div>
    <div class="mt-4 flex">
      <input v-model="newText" class="flex-1 input" placeholder="Scrivi un messaggio..." />
      <button @click="send" class="btn">Invia</button>
    </div>
  </section>
</template>

<script>
import axios from 'axios'
import { inject, onMounted, ref } from 'vue'
import MessageOptions from '@/components/MessageOptions.vue'

export default {
  components: { MessageOptions },
  setup(props, { route }) {
    axios.defaults.withCredentials = true
    const user = inject('currentUser')
    const messages = ref([])
    const newText = ref('')
    const convId = route.params.id

    const load = async () => {
      const res = await axios.get(`http://localhost:8080/conversations/${convId}`, { withCredentials: true })
      messages.value = res.data.messages
    }

    const send = async () => {
      if (!newText.value) return
      await axios.post(`http://localhost:8080/conversations/${convId}/messages`, { text: newText.value }, { withCredentials: true })
      newText.value = ''
      load()
    }

    const doForward = msg     => {/* logic forward */}
    const doDelete  = msg     => axios.delete(`http://localhost:8080/conversations/${convId}/messages/${msg.id}`, { withCredentials: true }).then(load)
    const doComment = msg     => {/* logic comment */}
    const doUncomment = msg   => {/* logic uncomment */}

    onMounted(load)
    return { user, messages, newText, send, doForward, doDelete, doComment, doUncomment }
  }
}
</script>