<template>
  <header class="relative bg-[var(--background)] text-[var(--text)] py-4">
    <!-- Titolo assoluto, centrato -->
    <h1 class="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 text-3xl font-bold gradient-title">
      WASATEXT
    </h1>
    <!-- Slot azioni, posizionato a destra -->
    <div class="absolute right-4 top-1/2 transform -translate-y-1/2 flex items-center space-x-2">
      <slot name="actions" />
    </div>
  </header>
</template>

<script>
export default {
  name: 'Header',
}
</script>

<style scoped>
/* niente scoped in più qui, stilizzazioni già in index.css */
</style>