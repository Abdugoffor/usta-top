<script setup>
import { Icon } from '@iconify/vue'
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar__brand">
      <div class="sidebar__logo">A</div>
      <div>
        <h2>Admin Panel</h2>
        <p>Vue 3 Dashboard</p>
      </div>
    </div>

    <nav class="sidebar__menu">
      <RouterLink to="/" class="sidebar__link">
        <Icon icon="mdi:view-dashboard-outline" />
        <span>Dashboard</span>
      </RouterLink>

      <RouterLink to="/categories" class="sidebar__link">
        <Icon icon="mdi:shape-outline" />
        <span>Categories</span>
      </RouterLink>
    </nav>
  </aside>
</template>