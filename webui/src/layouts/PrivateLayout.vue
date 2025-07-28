<!-- webui/src/layouts/PrivateLayout.vue -->
<template>
  <Header>
    <template #actions>
      <span class="text-white mr-4">👤 {{ user.name }}</span>
      <button @click="doLogout" class="text-white">Logout</button>
    </template>
  </Header>
  <router-view />
</template>

<script>
import { ref, onMounted } from 'vue'
import Header from '@/components/Header.vue'
import { getCurrentUser } from '@/services/userApi'
import { useRouter } from 'vue-router'

export default {
  components: { Header },
  setup() {
    const user = ref({ name: '' })
    const router = useRouter()

    onMounted(async () => {
      try {
        const res = await getCurrentUser()
        user.value = res.data
      } catch {
        router.push({ name: 'login' })
      }
    })

    const doLogout = () => {
      localStorage.removeItem('userId')
      router.push({ name: 'login' })
    }

    return { user, doLogout }
  }
}
</script>