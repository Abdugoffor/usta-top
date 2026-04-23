<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store/authStore'
import { useI18n } from '@/shared/composables/useI18n'

const props = defineProps({
  modelValue: { type: String, default: '' }
})
const emit = defineEmits(['update:modelValue', 'search', 'brand-click'])

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const { t, lang, setLang, activeLanguages } = useI18n()

const mobileMenuOpen = ref(false)
const userMenuOpen = ref(false)

const onInput   = (e) => emit('update:modelValue', e.target.value)
const onKeydown = (e) => { if (e.key === 'Enter') emit('search') }
const clearSearch = () => { emit('update:modelValue', '') }
const onBrandClick = () => { emit('brand-click'); mobileMenuOpen.value = false }

const doLogout = () => {
  auth.logout()
  userMenuOpen.value = false
  mobileMenuOpen.value = false
  router.push({ name: 'home' })
}

const changeLang = (e) => {
  const code = e.target.value
  setLang(code)
  router.push({ ...route, params: { ...route.params, lang: code } })
}
</script>

<template>
  <header class="site-header">
    <div class="site-header__inner">
      <RouterLink :to="`/${lang}`" class="site-header__brand" @click="onBrandClick">
        <div class="site-header__logo">UT</div>
        <span class="site-header__name">Usta<strong>Top</strong></span>
      </RouterLink>

      <div class="site-header__search">
        <svg class="site-header__search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
        </svg>
        <input
          class="site-header__search-input"
          type="text"
          :placeholder="t('search_placeholder')"
          :value="modelValue"
          @input="onInput"
          @keydown="onKeydown"
        />
        <button v-if="modelValue" type="button" class="site-header__search-clear" @click="clearSearch" :aria-label="t('filter_reset')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M18 6 6 18M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <nav class="site-header__nav">
        <RouterLink :to="`/${lang}`" class="site-header__nav-link" exact-active-class="active">{{ t('nav_masters') }}</RouterLink>
        <RouterLink :to="`/${lang}/vacancies`" class="site-header__nav-link" active-class="active">{{ t('nav_vacancies') }}</RouterLink>
      </nav>

      <!-- Til tanlash -->
      <div v-if="activeLanguages.length > 0" class="lang-switcher">
        <svg class="lang-switcher__icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M2 12h20M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
        </svg>
        <select class="lang-switcher__select" :value="lang" @change="changeLang">
          <option
            v-for="l in activeLanguages"
            :key="l.name"
            :value="l.name.toLowerCase()"
          >{{ l.name.toUpperCase() }}</option>
        </select>
      </div>

      <div class="site-header__actions">
        <template v-if="auth.isLoggedIn">
          <div class="user-menu-wrap">
            <button class="user-menu-btn" @click="userMenuOpen = !userMenuOpen">
              <div class="user-menu-btn__avatar">
                {{ auth.user?.full_name?.charAt(0)?.toUpperCase() || 'U' }}
              </div>
              <span class="user-menu-btn__name">{{ auth.user?.full_name?.split(' ')[0] }}</span>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"
                :style="{ transform: userMenuOpen ? 'rotate(180deg)' : '', transition: 'transform 0.2s' }">
                <path d="m6 9 6 6 6-6"/>
              </svg>
            </button>
            <div v-if="userMenuOpen" class="user-dropdown">
              <RouterLink :to="`/${lang}/profile`" class="user-dropdown__item" @click="userMenuOpen = false">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
                {{ t('profile_my') }}
              </RouterLink>
              <RouterLink :to="`/${lang}/profile/resumes/create`" class="user-dropdown__item" @click="userMenuOpen = false">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14 2 14 8 20 8"/>
                  <line x1="12" y1="18" x2="12" y2="12"/>
                  <line x1="9" y1="15" x2="15" y2="15"/>
                </svg>
                {{ t('resume_add') }}
              </RouterLink>
              <RouterLink :to="`/${lang}/profile/vacancies/create`" class="user-dropdown__item" @click="userMenuOpen = false">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="2" y="7" width="20" height="14" rx="2" ry="2"/>
                  <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>
                </svg>
                {{ t('vacancy_add') }}
              </RouterLink>
              <div class="user-dropdown__divider"></div>
              <button class="user-dropdown__item user-dropdown__item--danger" @click="doLogout">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                  <polyline points="16 17 21 12 16 7"/>
                  <line x1="21" y1="12" x2="9" y2="12"/>
                </svg>
                {{ t('auth_logout') }}
              </button>
            </div>
          </div>
        </template>
        <template v-else>
          <RouterLink :to="`/${lang}/login`"    class="site-header__btn site-header__btn--ghost">{{ t('auth_login') }}</RouterLink>
          <RouterLink :to="`/${lang}/register`" class="site-header__btn site-header__btn--primary">{{ t('auth_register') }}</RouterLink>
        </template>
      </div>

      <button class="site-header__burger" @click="mobileMenuOpen = !mobileMenuOpen">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path v-if="!mobileMenuOpen" d="M4 6h16M4 12h16M4 18h16"/>
          <path v-else d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </div>

    <div v-if="mobileMenuOpen" class="site-header__mobile-menu">
      <div class="site-header__mobile-search">
        <svg class="site-header__search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
        </svg>
        <input
          class="site-header__search-input"
          type="text"
          :placeholder="t('search_placeholder')"
          :value="modelValue"
          @input="onInput"
          @keydown="onKeydown"
        />
        <button v-if="modelValue" type="button" class="site-header__search-clear" @click="clearSearch" :aria-label="t('filter_reset')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M18 6 6 18M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <!-- Mobil til tanlash -->
      <div v-if="activeLanguages.length > 0" class="lang-switcher lang-switcher--mobile">
        <svg class="lang-switcher__icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M2 12h20M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
        </svg>
        <select class="lang-switcher__select" :value="lang" @change="changeLang">
          <option v-for="l in activeLanguages" :key="l.name" :value="l.name.toLowerCase()">
            {{ l.name.toUpperCase() }} — {{ l.description || l.name }}
          </option>
        </select>
      </div>

      <RouterLink :to="`/${lang}`" class="site-header__mobile-link" @click="mobileMenuOpen = false">{{ t('nav_masters') }}</RouterLink>
      <RouterLink :to="`/${lang}/vacancies`" class="site-header__mobile-link" @click="mobileMenuOpen = false">{{ t('nav_vacancies') }}</RouterLink>

      <template v-if="auth.isLoggedIn">
        <RouterLink :to="`/${lang}/profile`" class="site-header__mobile-link" @click="mobileMenuOpen = false">
          👤 {{ auth.user?.full_name }}
        </RouterLink>
        <RouterLink :to="`/${lang}/profile/resumes/create`" class="site-header__mobile-link" @click="mobileMenuOpen = false">
          📄 {{ t('resume_add') }}
        </RouterLink>
        <RouterLink :to="`/${lang}/profile/vacancies/create`" class="site-header__mobile-link" @click="mobileMenuOpen = false">
          💼 {{ t('vacancy_add') }}
        </RouterLink>
        <button class="site-header__mobile-link site-header__mobile-link--danger" @click="doLogout"
          style="border:none;text-align:left;cursor:pointer;font-family:inherit;font-size:15px;">
          🚪 {{ t('auth_logout') }}
        </button>
      </template>
      <template v-else>
        <div class="site-header__mobile-actions">
          <RouterLink :to="`/${lang}/login`"    class="site-header__btn site-header__btn--ghost"   style="flex:1;text-align:center;text-decoration:none;" @click="mobileMenuOpen = false">{{ t('auth_login') }}</RouterLink>
          <RouterLink :to="`/${lang}/register`" class="site-header__btn site-header__btn--primary" style="flex:1;text-align:center;text-decoration:none;" @click="mobileMenuOpen = false">{{ t('auth_register') }}</RouterLink>
        </div>
      </template>
    </div>
  </header>
</template>

<style scoped>
.site-header {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
}

.site-header__inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 20px;
  height: 64px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.site-header__brand { display: flex; align-items: center; gap: 10px; text-decoration: none; flex-shrink: 0; }

.site-header__logo {
  width: 38px; height: 38px; border-radius: 10px;
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  color: #fff; font-size: 13px; font-weight: 800;
  display: flex; align-items: center; justify-content: center; letter-spacing: 0.5px;
}

.site-header__name { font-size: 18px; font-weight: 500; color: #111827; }
.site-header__name strong { font-weight: 800; }

.site-header__search { flex: 1; max-width: 480px; position: relative; display: flex; align-items: center; }
.site-header__search-icon { position: absolute; left: 14px; width: 18px; height: 18px; color: #9ca3af; pointer-events: none; }
.site-header__search-input { width: 100%; padding: 10px 40px 10px 42px; border: 1.5px solid #e5e7eb; border-radius: 10px; font-size: 14px; color: #111827; background: #f9fafb; outline: none; transition: all 0.2s; font-family: inherit; }
.site-header__search-input:focus { border-color: #2563eb; background: #fff; box-shadow: 0 0 0 3px rgba(37,99,235,0.1); }
.site-header__search-input::placeholder { color: #9ca3af; }
.site-header__search-clear {
  position: absolute;
  right: 8px;
  width: 26px;
  height: 26px;
  border: none;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.15s;
  padding: 0;
}
.site-header__search-clear:hover { background: #f3f4f6; color: #111827; }
.site-header__search-clear svg { width: 16px; height: 16px; }

.site-header__nav { display: flex; align-items: center; gap: 4px; flex-shrink: 0; }
.site-header__nav-link { padding: 8px 14px; border-radius: 8px; font-size: 14px; font-weight: 500; color: #4b5563; text-decoration: none; transition: all 0.2s; }
.site-header__nav-link:hover, .site-header__nav-link.active { background: #eff6ff; color: #2563eb; }

/* Til tanlash */
.lang-switcher {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}
.lang-switcher__icon { width: 16px; height: 16px; color: #6b7280; flex-shrink: 0; }
.lang-switcher__select {
  padding: 6px 10px;
  border: 1.5px solid #e5e7eb;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  background: #f9fafb;
  outline: none;
  cursor: pointer;
  font-family: inherit;
  transition: border-color 0.2s;
}
.lang-switcher__select:focus { border-color: #2563eb; }
.lang-switcher--mobile { padding: 8px 0; }
.lang-switcher--mobile .lang-switcher__select { width: 100%; }

.site-header__actions { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }

.site-header__btn { padding: 9px 18px; border-radius: 10px; font-size: 14px; font-weight: 600; cursor: pointer; transition: all 0.2s; font-family: inherit; white-space: nowrap; display: inline-block; }
.site-header__btn--ghost { background: transparent; border: 1.5px solid #d1d5db; color: #374151; }
.site-header__btn--ghost:hover { border-color: #2563eb; color: #2563eb; background: #eff6ff; }
.site-header__btn--primary { background: linear-gradient(135deg, #1d4ed8, #2563eb); border: none; color: #fff; box-shadow: 0 2px 8px rgba(37,99,235,0.3); }
.site-header__btn--primary:hover { transform: translateY(-1px); box-shadow: 0 4px 12px rgba(37,99,235,0.4); }

.user-menu-wrap { position: relative; }
.user-menu-btn { display: flex; align-items: center; gap: 8px; padding: 6px 12px 6px 6px; border: 1.5px solid #e5e7eb; border-radius: 10px; background: #fff; cursor: pointer; font-family: inherit; transition: all 0.2s; }
.user-menu-btn:hover { border-color: #93c5fd; background: #eff6ff; }
.user-menu-btn__avatar { width: 30px; height: 30px; border-radius: 8px; background: linear-gradient(135deg, #1d4ed8, #2563eb); color: #fff; font-size: 13px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.user-menu-btn__name { font-size: 14px; font-weight: 600; color: #111827; }
.user-menu-btn svg { width: 14px; height: 14px; color: #6b7280; }

.user-dropdown { position: absolute; top: calc(100% + 8px); right: 0; background: #fff; border: 1.5px solid #e5e7eb; border-radius: 14px; box-shadow: 0 8px 24px rgba(0,0,0,0.12); min-width: 200px; padding: 6px; z-index: 200; }
.user-dropdown__item { display: flex; align-items: center; gap: 10px; padding: 10px 12px; border-radius: 8px; font-size: 14px; font-weight: 500; color: #374151; text-decoration: none; transition: all 0.15s; cursor: pointer; background: none; border: none; width: 100%; text-align: left; font-family: inherit; }
.user-dropdown__item svg { width: 16px; height: 16px; color: #6b7280; flex-shrink: 0; }
.user-dropdown__item:hover { background: #f3f4f6; color: #111827; }
.user-dropdown__item--danger { color: #dc2626; }
.user-dropdown__item--danger svg { color: #dc2626; }
.user-dropdown__item--danger:hover { background: #fef2f2; }
.user-dropdown__divider { height: 1px; background: #f3f4f6; margin: 4px 0; }

.site-header__burger { display: none; width: 40px; height: 40px; border: none; background: transparent; cursor: pointer; align-items: center; justify-content: center; border-radius: 8px; color: #374151; flex-shrink: 0; }
.site-header__burger svg { width: 22px; height: 22px; }
.site-header__burger:hover { background: #f3f4f6; }

.site-header__mobile-menu { border-top: 1px solid #e5e7eb; padding: 16px 20px; display: flex; flex-direction: column; gap: 8px; background: #fff; }
.site-header__mobile-search { position: relative; display: flex; align-items: center; margin-bottom: 4px; }
.site-header__mobile-search .site-header__search-input { width: 100%; }
.site-header__mobile-link { padding: 12px 16px; border-radius: 10px; font-size: 15px; font-weight: 500; color: #374151; text-decoration: none; background: #f9fafb; }
.site-header__mobile-link:hover { background: #eff6ff; color: #2563eb; }
.site-header__mobile-link--danger { color: #dc2626; }
.site-header__mobile-link--danger:hover { background: #fef2f2 !important; }
.site-header__mobile-actions { display: flex; gap: 8px; margin-top: 4px; }

@media (max-width: 768px) {
  .site-header__nav { display: none; }
  .site-header__actions { display: none; }
  .site-header__search { max-width: none; }
  .site-header__burger { display: flex; }
  .lang-switcher:not(.lang-switcher--mobile) { display: none; }
}

@media (max-width: 480px) {
  .site-header__search { display: none; }
}
</style>
