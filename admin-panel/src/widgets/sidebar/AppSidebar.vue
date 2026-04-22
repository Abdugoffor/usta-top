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

      <div class="sidebar__section-label">Katalog</div>

      <RouterLink to="/admin/categories" class="sidebar__link">
        <Icon icon="mdi:shape-outline" />
        <span>Kategoriyalar</span>
      </RouterLink>

      <RouterLink to="/admin/countries" class="sidebar__link">
        <Icon icon="mdi:earth" />
        <span>Hududlar</span>
      </RouterLink>

      <RouterLink to="/admin/languages" class="sidebar__link">
        <Icon icon="mdi:translate" />
        <span>Tillar</span>
      </RouterLink>

      <div class="sidebar__section-label">Foydalanuvchilar</div>

      <RouterLink to="/admin/resumes" class="sidebar__link">
        <Icon icon="mdi:account-box-outline" />
        <span>Resumelar</span>
      </RouterLink>

      <RouterLink to="/admin/vacancies-admin" class="sidebar__link">
        <Icon icon="mdi:briefcase-outline" />
        <span>Vakansiyalar</span>
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
.sidebar__section-label {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--muted);
  padding: 12px 16px 4px;
  opacity: 0.6;
}

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
