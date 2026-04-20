<script setup>
import { Icon } from '@iconify/vue'
</script>

<template>
  <header class="app-header">
    <div>
      <h1>Admin Dashboard</h1>
      <p>Category CRUD management</p>
    </div>

    <div class="app-header__actions">
      <button class="icon-btn">
        <Icon icon="mdi:bell-outline" />
      </button>
      <div class="user-chip">QA</div>
    </div>
  </header>
</template>