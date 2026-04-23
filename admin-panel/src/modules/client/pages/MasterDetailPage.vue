<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ClientHeader from '../components/ClientHeader.vue'
import SharePanel from '@/shared/components/SharePanel.vue'
import { getResume } from '../api/resumeApi'
import { useI18n } from '@/shared/composables/useI18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const master = ref(null)
const loading = ref(true)
const search = ref('')

const initials = computed(() => {
  if (!master.value) return 'U'
  return (master.value.name || '').split(' ').slice(0, 2).map(w => w[0] || '').join('').toUpperCase() || 'U'
})

const avatarColor = computed(() => {
  const colors = ['#2563eb','#7c3aed','#059669','#d97706','#dc2626','#0891b2']
  return colors[(master.value?.id || 0) % colors.length]
})

const price = computed(() => {
  if (!master.value?.price) return null
  return Number(master.value.price).toLocaleString('uz-UZ')
})

const location = computed(() => {
  if (!master.value) return ''
  return [master.value.region_name, master.value.district_name, master.value.mahalla_name, master.value.adress].filter(Boolean).join(', ')
})

const skills = computed(() => {
  if (!master.value?.skills) return []
  return master.value.skills.split(',').map(s => s.trim()).filter(Boolean)
})

const pageUrl = computed(() => `${window.location.origin}/masters/${route.params.slug}`)

const goBack = () => {
  router.push({ name: 'home', query: { ...route.query } })
}

onMounted(async () => {
  try {
    const res = await getResume(route.params.slug)
    master.value = res.data
  } catch {
    router.push({ name: 'home' })
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="detail-page">
    <ClientHeader v-model="search" />

    <div v-if="loading" class="detail-loading">
      <div class="spinner"></div>
    </div>

    <template v-else-if="master">
      <div class="detail-hero">
        <div class="detail-hero__inner">
          <div class="detail-hero__nav">
            <button class="detail-back" @click="goBack">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m15 18-6-6 6-6"/></svg>
              {{ t('btn_back') }}
            </button>
            <SharePanel :url="pageUrl" :title="master.name" :phone="master.contact" :dark="true" />
          </div>
          <div class="detail-hero__card">
            <div class="detail-hero__avatar" :style="{ background: avatarColor }">{{ initials }}</div>
            <div class="detail-hero__info">
              <div class="detail-hero__badges">
                <span class="detail-badge" :class="master.is_active ? 'detail-badge--active' : 'detail-badge--inactive'">
                  {{ master.is_active ? t('status_active') : t('status_inactive') }}
                </span>
                <span v-if="master.category_name" class="detail-badge detail-badge--blue">{{ master.category_name }}</span>
              </div>
              <h1 class="detail-hero__name">{{ master.name }}</h1>
              <p class="detail-hero__title">{{ master.title }}</p>
              <div class="detail-hero__meta">
                <span v-if="location">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
                    <circle cx="12" cy="9" r="2.5"/>
                  </svg>
                  {{ location }}
                </span>
                <span v-if="master.experience_year">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
                  </svg>
                  {{ master.experience_year }} {{ t('label_years_experience') }}
                </span>
                <span v-if="master.views_count">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  {{ master.views_count }} {{ t('label_views_count') }}
                </span>
              </div>
            </div>
            <div class="detail-hero__price-box">
              <div class="detail-hero__price-label">{{ t('label_daily_price') }}</div>
              <div class="detail-hero__price">{{ price ? price + ' ' + t('price_sum') : t('price_negotiable') }}</div>
              <a v-if="master.contact" :href="'tel:' + master.contact" class="detail-hero__contact-btn">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.69 12 19.79 19.79 0 0 1 1.61 3.4 2 2 0 0 1 3.6 1.22h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L7.91 8.78A16 16 0 0 0 15 15.87l.85-.85a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/>
                </svg>
                {{ master.contact }}
              </a>
            </div>
          </div>
        </div>
      </div>

      <div class="detail-body">
        <div class="detail-body__inner">
          <div class="detail-main">
            <div v-if="master.text" class="detail-section">
              <h2 class="detail-section__title">{{ t('section_about_me') }}</h2>
              <p class="detail-section__text">{{ master.text }}</p>
            </div>

            <div v-if="skills.length" class="detail-section">
              <h2 class="detail-section__title">{{ t('section_skills') }}</h2>
              <div class="detail-skills">
                <span v-for="s in skills" :key="s" class="detail-skill">{{ s }}</span>
              </div>
            </div>

            <div v-if="master.adress" class="detail-section">
              <h2 class="detail-section__title">{{ t('section_address') }}</h2>
              <p class="detail-section__text">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:16px;height:16px;display:inline;vertical-align:middle;margin-right:4px">
                  <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
                  <circle cx="12" cy="9" r="2.5"/>
                </svg>
                {{ location || master.adress }}
              </p>
            </div>
          </div>

          <aside class="detail-aside">
            <div class="detail-aside__card">
              <h3 class="detail-aside__title">{{ t('section_contact_info') }}</h3>
              <a v-if="master.contact" :href="'tel:' + master.contact" class="detail-aside__row detail-aside__phone">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.69 12 19.79 19.79 0 0 1 1.61 3.4 2 2 0 0 1 3.6 1.22h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L7.91 8.78A16 16 0 0 0 15 15.87l.85-.85a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/>
                </svg>
                <span>{{ master.contact }}</span>
              </a>
              <div v-if="location" class="detail-aside__row">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
                  <circle cx="12" cy="9" r="2.5"/>
                </svg>
                <span>{{ location }}</span>
              </div>
            </div>

            <div class="detail-aside__card detail-aside__stats">
              <div class="detail-aside__stat">
                <div class="detail-aside__stat-val">{{ master.experience_year || 0 }}</div>
                <div class="detail-aside__stat-label">{{ t('stat_years_experience') }}</div>
              </div>
              <div class="detail-aside__stat">
                <div class="detail-aside__stat-val">{{ master.views_count || 0 }}</div>
                <div class="detail-aside__stat-label">{{ t('stat_views') }}</div>
              </div>
              <div class="detail-aside__stat">
                <div class="detail-aside__stat-val" :class="master.is_active ? 'text--green' : 'text--gray'">
                  {{ master.is_active ? t('status_active') : t('status_inactive') }}
                </div>
                <div class="detail-aside__stat-label">{{ t('stat_status') }}</div>
              </div>
            </div>
          </aside>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.detail-page {
  min-height: 100vh;
  background: #f3f4f6;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.detail-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50vh;
}

.spinner {
  width: 44px;
  height: 44px;
  border: 3px solid #e5e7eb;
  border-top-color: #2563eb;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.detail-hero {
  background: linear-gradient(135deg, #1a3fbd, #2563eb);
  padding: 32px 20px 40px;
}

.detail-hero__inner { max-width: 1100px; margin: 0 auto; }

.detail-hero__nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.detail-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: rgba(255,255,255,0.8);
  background: rgba(255,255,255,0.1);
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 10px;
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.2s;
}

.detail-back:hover { background: rgba(255,255,255,0.2); color: #fff; }
.detail-back svg { width: 16px; height: 16px; }

.detail-hero__card {
  display: flex;
  align-items: flex-start;
  gap: 24px;
  background: rgba(255,255,255,0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 20px;
  padding: 28px;
}

.detail-hero__avatar {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 800;
  color: #fff;
  flex-shrink: 0;
}

.detail-hero__info { flex: 1; }

.detail-hero__badges { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }

.detail-badge {
  font-size: 12px;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 20px;
}

.detail-badge--active { background: #d1fae5; color: #065f46; }
.detail-badge--inactive { background: rgba(255,255,255,0.2); color: rgba(255,255,255,0.8); }
.detail-badge--blue { background: rgba(255,255,255,0.2); color: #fff; }

.detail-hero__name {
  font-size: 28px;
  font-weight: 800;
  color: #fff;
  margin-bottom: 4px;
}

.detail-hero__title { font-size: 15px; color: rgba(255,255,255,0.75); margin-bottom: 16px; }

.detail-hero__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.detail-hero__meta span {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: rgba(255,255,255,0.8);
}

.detail-hero__meta svg { width: 14px; height: 14px; flex-shrink: 0; }

.detail-hero__price-box {
  text-align: right;
  flex-shrink: 0;
}

.detail-hero__price-label { font-size: 12px; color: rgba(255,255,255,0.6); margin-bottom: 4px; }

.detail-hero__price {
  font-size: 24px;
  font-weight: 800;
  color: #fff;
  margin-bottom: 12px;
}

.detail-hero__contact-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255,255,255,0.15);
  border: 1px solid rgba(255,255,255,0.3);
  color: #fff;
  border-radius: 10px;
  padding: 10px 16px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  text-decoration: none;
  transition: all 0.2s;
}

.detail-hero__contact-btn:hover { background: rgba(255,255,255,0.25); }
.detail-hero__contact-btn svg { width: 14px; height: 14px; }

.detail-aside__phone {
  text-decoration: none;
  color: #2563eb;
  font-weight: 600;
  border-radius: 8px;
  transition: background 0.15s;
}
.detail-aside__phone:hover { background: #eff6ff; }
.detail-aside__phone svg { color: #2563eb; }

.detail-body { padding: 28px 20px 48px; }

.detail-body__inner {
  max-width: 1100px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 24px;
  align-items: start;
}

.detail-main { display: flex; flex-direction: column; gap: 16px; }

.detail-section {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 24px;
}

.detail-section__title {
  font-size: 16px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 12px;
}

.detail-section__text { font-size: 14px; color: #4b5563; line-height: 1.7; }

.detail-skills { display: flex; flex-wrap: wrap; gap: 8px; }

.detail-skill {
  background: #eff6ff;
  color: #1d4ed8;
  font-size: 13px;
  font-weight: 600;
  padding: 6px 14px;
  border-radius: 20px;
}

.detail-aside { display: flex; flex-direction: column; gap: 16px; position: sticky; top: 80px; }

.detail-aside__card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 20px;
}

.detail-aside__title {
  font-size: 14px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 14px;
}

.detail-aside__row {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: #4b5563;
  margin-bottom: 10px;
}

.detail-aside__row svg { width: 16px; height: 16px; color: #9ca3af; flex-shrink: 0; }

.detail-aside__stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  text-align: center;
}

.detail-aside__stat-val {
  font-size: 20px;
  font-weight: 800;
  color: #111827;
}

.detail-aside__stat-label { font-size: 11px; color: #6b7280; margin-top: 2px; }

.text--green { color: #059669; }
.text--gray { color: #6b7280; }

@media (max-width: 900px) {
  .detail-hero__card { flex-direction: column; }
  .detail-hero__price-box { text-align: left; }
  .detail-body__inner { grid-template-columns: 1fr; }
  .detail-hero__name { font-size: 22px; }
}

@media (max-width: 640px) {
  .detail-hero { padding: 24px 16px 32px; }
  .detail-body { padding: 20px 16px 40px; }
}
</style>
