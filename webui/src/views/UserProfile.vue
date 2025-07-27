<template>
  <section class="p-4 max-w-md mx-auto">
    <h2 class="text-2xl font-bold mb-4">User Profile</h2>
    <div v-if="loading" class="text-center">Caricamento...</div>
    <div v-else>
      <img
        :src="user.avatarUrl || defaultAvatar"
        class="w-24 h-24 rounded-full mb-4"
        alt="Avatar utente"
      />
      <p class="text-lg font-semibold">{{ user.name }}</p>
      <p class="text-gray-600">{{ user.email }}</p>
    </div>
    <p v-if="error" class="text-red-500 mt-4">{{ error }}</p>
  </section>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

export default {
  setup() {
    const route = useRoute()
    const user = ref({})
    const loading = ref(true)
    const error = ref('')
    const defaultAvatar = '/default-avatar.png' // metti in public/

    const loa
