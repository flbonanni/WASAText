<!-- webui/src/views/MyProfile.vue -->
<template>
  <section class="p-4 max-w-md mx-auto">
    <h2 class="text-2xl mb-4">My Profile</h2>
    <img :src="user.avatarUrl" class="w-24 h-24 rounded-full mb-4" />
    <p class="font-semibold">{{ user.name }}</p>
    <!-- bottoni per cambiare -->
  </section>
</template>

<script>
import { ref, onMounted } from 'vue'
import { getCurrentUser } from '@/services/userApi'

export default {
  setup() {
    const user = ref({})
    onMounted(async () => {
      const res = await getCurrentUser()
      user.value = res.data
    })
    return { user }
  }
}
</script>