<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ClientHeader from '../components/ClientHeader.vue'
import { getVacancy } from '../api/vacancyApi'

const route = useRoute()
const router = useRouter()
const vacancy = ref(null)
const loading = ref(true)
const search = ref('')

const iconColor = computed(() => {
  const colors = ['#f59e0b','#10b981','#8b5cf6','#ef4444','#06b6d4','#f97316']
  return colors[(vacancy.value?.id || 0) % colors.length]
})

const price = computed(() => {
  if (!vacancy.value?.price) return null
  return Number(vacancy.value.price).toLocaleString('uz-UZ')
})

const location = computed(() => {
  if (!vacancy.value) return ''
  return [vacancy.value.region_name, vacancy.value.district_name, vacancy.value.mahalla_name, vacancy.value.adress].filter(Boolean).join(', ')
})

const timeAgo = computed(() => {
  if (!vacancy.value?.created_at) return ''
  const date = new Date(vacancy.value.created_at)
  const now = new Date()
  const diff = Math.floor((now - date) / 1000)
  if (diff < 3600) return `${Math.floor(diff/60)} daqiqa oldin`
  if (diff < 86400) return `${Math.floor(diff/3600)} soat oldin`
  return `${Math.floor(diff/86400)} kun oldin`
})

onMounted(async () => {
  try {
    const res = await getVacancy(route.params.slug)
    vacancy.value = res.data
  } catch {
    router.push('/vacancies')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="vac-detail-page">
    <ClientHeader v-model="search" />

    <div v-if="loading" class="vac-detail-loading">
      <div class="spinner"></div>
    </div>

    <template v-else-if="vacancy">
      <div class="vac-detail-hero">
        <div class="vac-detail-hero__inner">
          <button class="detail-back" @click="router.back()">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m15 18-6-6 6-6"/></svg>
            Orqaga
          </button>
          <div class="vac-detail-hero__card">
            <div class="vac-detail-hero__icon" :style="{ background: iconColor + '33', color: iconColor }">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M20 6h-2.18c.07-.44.18-.88.18-1.36C18 2.96 16.04 1 13.64 1c-1.3 0-2.55.57-3.42 1.56L9 4 7.78 2.56C6.91 1.57 5.66 1 4.36 1 1.96 1 0 2.96 0 5.36c0 .48.11.92.18 1.36H0v2h20V6zm-7.37-2.21c.43-.51 1.06-.79 1.73-.79.97 0 1.64.67 1.64 1.64 0 .96-.77 1.72-1.64 1.72h-2.55l.82-.97.82-.97.18-.63zm-9.26.64C3.37 4.43 4.04 3.76 5.01 3.76c.67 0 1.3.28 1.73.79l.82.97.82.97-.18.63H5.64C4.77 7.12 4 6.36 4 5.4c0-.34.1-.65.24-.94-.14-.02-.24.01-.37.03zM0 19c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V9H0v10z"/>
              </svg>
            </div>
            <div class="vac-detail-hero__info">
              <div class="vac-detail-hero__badges">
                <span class="detail-badge" :class="vacancy.is_active ? 'detail-badge--active' : 'detail-badge--inactive'">
                  {{ vacancy.is_active ? 'Faol' : 'Yopiq' }}
                </span>
              </div>
              <h1 class="vac-detail-hero__name">{{ vacancy.name || vacancy.title }}</h1>
              <p v-if="vacancy.title && vacancy.title !== vacancy.name" class="vac-detail-hero__sub">{{ vacancy.title }}</p>
              <div class="vac-detail-hero__meta">
                <span v-if="location">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
                    <circle cx="12" cy="9" r="2.5"/>
                  </svg>
                  {{ location }}
                </span>
                <span v-if="timeAgo">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
                  </svg>
                  {{ timeAgo }}
                </span>
                <span v-if="vacancy.views_count">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                  {{ vacancy.views_count }} ko'rish
                </span>
              </div>
            </div>
            <div class="vac-detail-hero__price-box">
              <div class="vac-detail-hero__price-label">Maosh</div>
              <div class="vac-detail-hero__price">{{ price ? price + ' so\'m' : 'Kelishiladi' }}</div>
              <button v-if="vacancy.contact" class="vac-detail-contact-btn">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.69 12 19.79 19.79 0 0 1 1.61 3.4 2 2 0 0 1 3.6 1.22h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L7.91 8.78A16 16 0 0 0 15 15.87l.85-.85a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/>
                </svg>
                {{ vacancy.contact }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="vac-detail-body">
        <div class="vac-detail-body__inner">
          <div class="vac-detail-main">
            <div v-if="vacancy.text" class="detail-section">
              <h2 class="detail-section__title">Vakansiya haqida</h2>
              <p class="detail-section__text">{{ vacancy.text }}</p>
            </div>
          </div>

          <aside class="vac-detail-aside">
            <div class="detail-aside__card">
              <h3 class="detail-aside__title">Aloqa ma'lumotlari</h3>
              <div v-if="vacancy.contact" class="detail-aside__row">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.69 12 19.79 19.79 0 0 1 1.61 3.4 2 2 0 0 1 3.6 1.22h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L7.91 8.78A16 16 0 0 0 15 15.87l.85-.85a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/>
                </svg>
                <span>{{ vacancy.contact }}</span>
              </div>
              <div v-if="location" class="detail-aside__row">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
                  <circle cx="12" cy="9" r="2.5"/>
                </svg>
                <span>{{ location }}</span>
              </div>
            </div>

            <div class="detail-aside__card">
              <div class="vac-detail-stat-row">
                <div class="detail-aside__stat">
                  <div class="detail-aside__stat-val">{{ vacancy.views_count || 0 }}</div>
                  <div class="detail-aside__stat-label">Ko'rish</div>
                </div>
                <div class="detail-aside__stat">
                  <div class="detail-aside__stat-val" :class="vacancy.is_active ? 'text--green' : 'text--gray'">
                    {{ vacancy.is_active ? 'Faol' : 'Yopiq' }}
                  </div>
                  <div class="detail-aside__stat-label">Holat</div>
                </div>
              </div>
            </div>

            <RouterLink to="/vacancies" class="back-link">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
              Barcha vakansiyalar
            </RouterLink>
          </aside>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.vac-detail-page { min-height: 100vh; background: #f3f4f6; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; }

.vac-detail-loading { display: flex; justify-content: center; align-items: center; height: 50vh; }
.spinner { width: 44px; height: 44px; border: 3px solid #e5e7eb; border-top-color: #2563eb; border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.vac-detail-hero { background: linear-gradient(135deg, #1e3a8a, #1d4ed8, #2563eb); padding: 32px 20px 40px; }
.vac-detail-hero__inner { max-width: 1100px; margin: 0 auto; }

.detail-back { display: inline-flex; align-items: center; gap: 6px; color: rgba(255,255,255,0.8); background: rgba(255,255,255,0.1); border: 1px solid rgba(255,255,255,0.2); border-radius: 10px; padding: 8px 16px; font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; margin-bottom: 20px; transition: all 0.2s; }
.detail-back:hover { background: rgba(255,255,255,0.2); color: #fff; }
.detail-back svg { width: 16px; height: 16px; }

.vac-detail-hero__card { display: flex; align-items: flex-start; gap: 24px; background: rgba(255,255,255,0.1); backdrop-filter: blur(10px); border: 1px solid rgba(255,255,255,0.2); border-radius: 20px; padding: 28px; }

.vac-detail-hero__icon { width: 72px; height: 72px; border-radius: 18px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.vac-detail-hero__icon svg { width: 32px; height: 32px; }

.vac-detail-hero__info { flex: 1; }
.vac-detail-hero__badges { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }

.detail-badge { font-size: 12px; font-weight: 600; padding: 4px 12px; border-radius: 20px; }
.detail-badge--active { background: #d1fae5; color: #065f46; }
.detail-badge--inactive { background: rgba(255,255,255,0.2); color: rgba(255,255,255,0.8); }

.vac-detail-hero__name { font-size: 26px; font-weight: 800; color: #fff; margin-bottom: 4px; }
.vac-detail-hero__sub { font-size: 14px; color: rgba(255,255,255,0.7); margin-bottom: 14px; }

.vac-detail-hero__meta { display: flex; flex-wrap: wrap; gap: 16px; }
.vac-detail-hero__meta span { display: flex; align-items: center; gap: 6px; font-size: 13px; color: rgba(255,255,255,0.8); }
.vac-detail-hero__meta svg { width: 14px; height: 14px; flex-shrink: 0; }

.vac-detail-hero__price-box { text-align: right; flex-shrink: 0; }
.vac-detail-hero__price-label { font-size: 12px; color: rgba(255,255,255,0.6); margin-bottom: 4px; }
.vac-detail-hero__price { font-size: 24px; font-weight: 800; color: #fff; margin-bottom: 12px; }

.vac-detail-contact-btn { display: flex; align-items: center; gap: 8px; background: rgba(255,255,255,0.15); border: 1px solid rgba(255,255,255,0.3); color: #fff; border-radius: 10px; padding: 10px 16px; font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; transition: all 0.2s; }
.vac-detail-contact-btn:hover { background: rgba(255,255,255,0.25); }
.vac-detail-contact-btn svg { width: 14px; height: 14px; }

.vac-detail-body { padding: 28px 20px 48px; }
.vac-detail-body__inner { max-width: 1100px; margin: 0 auto; display: grid; grid-template-columns: 1fr 300px; gap: 24px; align-items: start; }

.vac-detail-main { display: flex; flex-direction: column; gap: 16px; }

.detail-section { background: #fff; border: 1px solid #e5e7eb; border-radius: 16px; padding: 24px; }
.detail-section__title { font-size: 16px; font-weight: 700; color: #111827; margin-bottom: 12px; }
.detail-section__text { font-size: 14px; color: #4b5563; line-height: 1.7; white-space: pre-wrap; }

.vac-detail-aside { display: flex; flex-direction: column; gap: 16px; position: sticky; top: 80px; }

.detail-aside__card { background: #fff; border: 1px solid #e5e7eb; border-radius: 16px; padding: 20px; }
.detail-aside__title { font-size: 14px; font-weight: 700; color: #111827; margin-bottom: 14px; }
.detail-aside__row { display: flex; align-items: center; gap: 10px; font-size: 13px; color: #4b5563; margin-bottom: 10px; }
.detail-aside__row svg { width: 16px; height: 16px; color: #9ca3af; flex-shrink: 0; }

.vac-detail-stat-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; text-align: center; }
.detail-aside__stat-val { font-size: 20px; font-weight: 800; color: #111827; }
.detail-aside__stat-label { font-size: 11px; color: #6b7280; margin-top: 2px; }
.text--green { color: #059669; }
.text--gray { color: #6b7280; }

.back-link { display: flex; align-items: center; gap: 6px; padding: 12px 16px; background: #fff; border: 1.5px solid #e5e7eb; border-radius: 12px; font-size: 13px; font-weight: 600; color: #374151; text-decoration: none; transition: all 0.2s; }
.back-link:hover { border-color: #2563eb; color: #2563eb; background: #eff6ff; }
.back-link svg { width: 16px; height: 16px; }

@media (max-width: 900px) {
  .vac-detail-hero__card { flex-direction: column; }
  .vac-detail-hero__price-box { text-align: left; }
  .vac-detail-body__inner { grid-template-columns: 1fr; }
}
@media (max-width: 640px) {
  .vac-detail-hero { padding: 24px 16px 32px; }
  .vac-detail-body { padding: 20px 16px 40px; }
  .vac-detail-hero__name { font-size: 20px; }
}
</style>
