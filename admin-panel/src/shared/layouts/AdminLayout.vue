<script setup>
import AppSidebar from '@/widgets/sidebar/AppSidebar.vue'
import AppHeader from '@/widgets/header/AppHeader.vue'
</script>

<template>
  <div class="admin-layout">
    <AppSidebar />

    <div class="admin-layout__main">
      <AppHeader />

      <main class="admin-layout__content">
        <RouterView />
      </main>
    </div>
  </div>
</template>