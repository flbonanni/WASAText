<template>
  <div class="relative inline-block">
    <button @click="open = !open">⋯</button>
    <ul v-if="open" class="absolute bg-white shadow p-2 rounded right-0">
      <li><button @click="$emit('forward')">Forward</button></li>
      <li v-if="isOwner"><button @click="$emit('delete')">Delete</button></li>
      <li v-if="!isOwner && !isCommented"><button @click="$emit('comment')">Comment</button></li>
      <li v-if="!isOwner && isCommented"><button @click="$emit('uncomment')">Uncomment</button></li>
    </ul>
  </div>
</template>

<script>
export default {
  props: { isOwner: Boolean, isCommented: Boolean },
  data() { return { open: false } }
}
</script>