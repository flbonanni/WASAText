<template>
  <section class="p-4 max-w-md mx-auto">
    <h2 class="text-2xl mb-4">My Profile</h2>
    <img :src="user.avatarUrl" class="w-24 h-24 rounded-full mb-4" />
    <p class="font-semibold">{{ user.name }}</p>
    <button @click="changePhoto" class="btn mt-2">Cambia foto</button>
    <button @click="changeName"  class="btn mt-2">Cambia nome</button>
  </section>
</template>

<script>
import axios from 'axios'
import { inject, ref } from 'vue'

export default {
  setup() {
    const user = inject('currentUser')
    const changePhoto = async () => {
      const url = prompt('Nuova URL foto:')
      if (url) {
        await axios.put('http://localhost:8080/me', { avatarUrl: url }, { withCredentials: true })
        user.avatarUrl = url
      }
    }
    const changeName = async () => {
      const name = prompt('Nuovo nome:')
      if (name) {
        await axios.put('http://localhost:8080/me', { name }, { withCredentials: true })
        user.name = name
      }
    }
    return { user, changePhoto, changeName }
  }
}
</script>