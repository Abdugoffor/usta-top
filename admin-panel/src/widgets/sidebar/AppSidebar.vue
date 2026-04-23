<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Icon } from '@iconify/vue'

const route = useRoute()
const lang = computed(() => route.params.lang || 'uz')
const a = (path) => `/${lang.value}/admin${path}`
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
      <RouterLink :to="a('')" class="sidebar__link" exact-active-class="router-link-active">
        <Icon icon="mdi:view-dashboard-outline" />
        <span>Dashboard</span>
      </RouterLink>

      <div class="sidebar__section-label">Katalog</div>

      <RouterLink :to="a('/categories')" class="sidebar__link">
        <Icon icon="mdi:shape-outline" />
        <span>Kategoriyalar</span>
      </RouterLink>

      <RouterLink :to="a('/countries')" class="sidebar__link">
        <Icon icon="mdi:earth" />
        <span>Hududlar</span>
      </RouterLink>

      <RouterLink :to="a('/languages')" class="sidebar__link">
        <Icon icon="mdi:translate" />
        <span>Tillar</span>
      </RouterLink>

      <RouterLink :to="a('/translations')" class="sidebar__link">
        <Icon icon="mdi:text-box-multiple-outline" />
        <span>Tarjimalar</span>
      </RouterLink>

      <div class="sidebar__section-label">Foydalanuvchilar</div>

      <RouterLink :to="a('/resumes')" class="sidebar__link">
        <Icon icon="mdi:account-box-outline" />
        <span>Resumelar</span>
      </RouterLink>

      <RouterLink :to="a('/vacancies-admin')" class="sidebar__link">
        <Icon icon="mdi:briefcase-outline" />
        <span>Vakansiyalar</span>
      </RouterLink>

      <div class="sidebar__divider"></div>

      <RouterLink :to="`/${lang}`" class="sidebar__link sidebar__link--site">
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
