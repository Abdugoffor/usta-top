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
      <RouterLink to="/admin" class="sidebar__link" exact-active-class="router-link-active">
        <Icon icon="mdi:view-dashboard-outline" />
        <span>Dashboard</span>
      </RouterLink>

      <RouterLink to="/admin/categories" class="sidebar__link">
        <Icon icon="mdi:shape-outline" />
        <span>Categories</span>
      </RouterLink>

      <RouterLink to="/admin/countries" class="sidebar__link">
        <Icon icon="mdi:earth" />
        <span>Countries</span>
      </RouterLink>

      <RouterLink to="/admin/languages" class="sidebar__link">
        <Icon icon="mdi:translate" />
        <span>Languages</span>
      </RouterLink>

      <div class="sidebar__divider"></div>

      <RouterLink to="/" class="sidebar__link sidebar__link--site">
        <Icon icon="mdi:web" />
        <span>Saytga o'tish</span>
      </RouterLink>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar__divider {
  height: 1px;
  background: var(--border);
  margin: 8px 0;
}

.sidebar__link--site {
  color: #10b981 !important;
}

.sidebar__link--site:hover {
  background: rgba(16, 185, 129, 0.1) !important;
  color: #10b981 !important;
}
</style>
