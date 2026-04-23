<script setup>
import { Icon } from '@iconify/vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from '@/shared/composables/useI18n'

const router = useRouter()
const route  = useRoute()
const { lang, setLang, activeLanguages } = useI18n()

const changeLang = (e) => {
  const code = e.target.value
  setLang(code)
  router.push({ ...route, params: { ...route.params, lang: code } })
}
</script>

<template>
  <header class="app-header">
    <div>
      <h1>Admin Dashboard</h1>
      <p>Usta Top boshqaruv paneli</p>
    </div>

    <div class="app-header__actions">
      <div v-if="activeLanguages.length > 0" class="lang-switcher">
        <Icon icon="mdi:translate" class="lang-switcher__icon" />
        <select class="lang-switcher__select" :value="lang" @change="changeLang">
          <option
            v-for="l in activeLanguages"
            :key="l.name"
            :value="l.name.toLowerCase()"
          >{{ l.name.toUpperCase() }}</option>
        </select>
        <Icon icon="mdi:chevron-down" class="lang-switcher__arrow" />
      </div>

      <button class="icon-btn">
        <Icon icon="mdi:bell-outline" />
      </button>

      <div class="user-chip">QA</div>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 28px;
  border-bottom: 1px solid var(--border);
  background: var(--bg);
  gap: 16px;
}

.app-header h1 {
  font-size: 18px;
  font-weight: 700;
  color: var(--text);
  line-height: 1.2;
}

.app-header p {
  font-size: 12px;
  color: var(--muted);
  margin-top: 2px;
}

.app-header__actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* Language switcher */
.lang-switcher {
  position: relative;
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-elevated);
  border: 1.5px solid var(--border);
  border-radius: 12px;
  padding: 0 10px 0 10px;
  height: 38px;
  transition: border-color 0.2s;
  cursor: pointer;
}

.lang-switcher:focus-within {
  border-color: #6366f1;
}

.lang-switcher__icon {
  font-size: 16px;
  color: var(--muted);
  flex-shrink: 0;
}

.lang-switcher__arrow {
  font-size: 14px;
  color: var(--muted);
  flex-shrink: 0;
  pointer-events: none;
}

.lang-switcher__select {
  background: transparent;
  border: none;
  outline: none;
  color: var(--text);
  font-size: 13px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
  padding: 0 2px;
  min-width: 36px;
}

.lang-switcher__select option {
  background: var(--bg-elevated);
  color: var(--text);
}

/* Icon button */
.icon-btn {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  background: var(--bg-elevated);
  border: 1.5px solid var(--border);
  color: var(--muted);
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-btn:hover {
  border-color: #6366f1;
  color: #6366f1;
}

/* User chip */
.user-chip {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366f1, #3b82f6);
  color: white;
  font-size: 13px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.35);
  flex-shrink: 0;
}
</style>
