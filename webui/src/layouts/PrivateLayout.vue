<template>
  <Header>
    <template #actions>
      <img
        :src="user.avatarUrl"
        class="w-8 h-8 rounded-full cursor-pointer"
        @click="goProfile"
      />
      <button @click="doLogout">Logout</button>
    </template>
  </Header>
  <router-view />
</template>

<script>
import Header from '@/components/Header.vue'
import axios from 'axios'
import { inject } from 'vue'

export default {
  components: { Header },
  setup() {
    const user = inject('currentUser')
    const goProfile = () => router.push({ name: 'my-profile' })
    const doLogout = async () => {
      await axios.post('http://localhost:8080/logout', {}, { withCredentials: true })
      window.location.href = '/'
    }
    return { user, goProfile, doLogout }
  }
}
</script>