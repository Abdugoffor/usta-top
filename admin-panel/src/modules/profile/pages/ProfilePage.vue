<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store/authStore'
import { getResumes } from '@/modules/client/api/resumeApi'
import { getVacancies } from '@/modules/client/api/vacancyApi'
import ClientHeader from '@/modules/client/components/ClientHeader.vue'

const router = useRouter()
const auth = useAuthStore()

const myResumes = ref([])
const myVacancies = ref([])
const loadingResumes = ref(false)
const loadingVacancies = ref(false)
const activeTab = ref('resumes')

const fetchMyResumes = async () => {
  if (!auth.user?.id) return
  loadingResumes.value = true
  try {
    const res = await getResumes({ user_id: auth.user.id, limit: 50 })
    myResumes.value = res.data?.data || []
  } catch {
    myResumes.value = []
  } finally {
    loadingResumes.value = false
  }
}

const fetchMyVacancies = async () => {
  if (!auth.user?.id) return
  loadingVacancies.value = true
  try {
    const res = await getVacancies({ user_id: auth.user.id, limit: 50 })
    myVacancies.value = res.data?.data || []
  } catch {
    myVacancies.value = []
  } finally {
    loadingVacancies.value = false
  }
}

const doLogout = () => {
  auth.logout()
  router.push({ name: 'home' })
}

const formatDate = (dt) => {
  if (!dt) return ''
  return new Date(dt).toLocaleDateString('uz-UZ', { day: '2-digit', month: 'long', year: 'numeric' })
}

const formatPrice = (p) => p ? Number(p).toLocaleString('uz-UZ') + " so'm" : '—'

onMounted(() => {
  if (!auth.isLoggedIn) { router.push({ name: 'login' }); return }
  fetchMyResumes()
  fetchMyVacancies()
})
</script>

<template>
  <div class="profile-page">
    <ClientHeader />

    <div class="profile-hero">
      <div class="profile-hero__inner">
        <div class="profile-hero__avatar">
          {{ auth.user?.full_name?.charAt(0)?.toUpperCase() || 'U' }}
        </div>
        <div class="profile-hero__info">
          <h1 class="profile-hero__name">{{ auth.user?.full_name }}</h1>
          <p class="profile-hero__phone">{{ auth.user?.phone }}</p>
          <span class="profile-hero__role">
            {{ auth.user?.role === 'employer' ? '🏢 Ish beruvchi' : '👷 Usta / Ishchi' }}
          </span>
        </div>
        <button class="profile-hero__logout" @click="doLogout">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
          Chiqish
        </button>
      </div>
    </div>

    <div class="profile-main">
      <div class="profile-main__inner">
        <!-- Tabs -->
        <div class="profile-tabs">
          <button class="profile-tab" :class="{ 'profile-tab--active': activeTab === 'resumes' }" @click="activeTab = 'resumes'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
            </svg>
            Resumelarim ({{ myResumes.length }})
          </button>
          <button class="profile-tab" :class="{ 'profile-tab--active': activeTab === 'vacancies' }" @click="activeTab = 'vacancies'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="7" width="20" height="14" rx="2" ry="2"/>
              <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>
            </svg>
            Vakansiyalarim ({{ myVacancies.length }})
          </button>
        </div>

        <!-- Resumes Tab -->
        <div v-if="activeTab === 'resumes'">
          <div class="profile-section__header">
            <h2>Mening resumelarim</h2>
            <RouterLink :to="{ name: 'resume-create' }" class="profile-create-btn">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              Yangi resume
            </RouterLink>
          </div>

          <div v-if="loadingResumes" class="profile-loading">Yuklanmoqda...</div>

          <div v-else-if="myResumes.length" class="profile-list">
            <div v-for="r in myResumes" :key="r.id" class="profile-item">
              <div class="profile-item__left">
                <div class="profile-item__avatar" :style="{ background: '#2563eb' }">
                  {{ (r.name || 'U').charAt(0).toUpperCase() }}
                </div>
                <div class="profile-item__info">
                  <div class="profile-item__name">{{ r.name }}</div>
                  <div class="profile-item__sub">{{ r.title }}</div>
                  <div class="profile-item__meta">
                    <span>{{ formatPrice(r.price) }}</span>
                    <span v-if="r.experience_year">• {{ r.experience_year }} yil tajriba</span>
                    <span>• {{ formatDate(r.created_at) }}</span>
                  </div>
                </div>
              </div>
              <div class="profile-item__right">
                <span class="profile-item__badge" :class="r.is_active ? 'badge--active' : 'badge--inactive'">
                  {{ r.is_active ? 'Faol' : 'Yopiq' }}
                </span>
                <RouterLink :to="{ name: 'master-detail', params: { slug: r.slug } }" class="profile-item__view">Ko'rish</RouterLink>
              </div>
            </div>
          </div>

          <div v-else class="profile-empty">
            <div class="profile-empty__icon">📄</div>
            <h3>Resumelaringiz yo'q</h3>
            <p>Birinchi resumeingizni yarating</p>
            <RouterLink :to="{ name: 'resume-create' }" class="profile-create-btn profile-create-btn--lg">
              Resume yaratish
            </RouterLink>
          </div>
        </div>

        <!-- Vacancies Tab -->
        <div v-if="activeTab === 'vacancies'">
          <div class="profile-section__header">
            <h2>Mening vakansiyalarim</h2>
            <RouterLink :to="{ name: 'vacancy-create' }" class="profile-create-btn">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              Yangi vakansiya
            </RouterLink>
          </div>

          <div v-if="loadingVacancies" class="profile-loading">Yuklanmoqda...</div>

          <div v-else-if="myVacancies.length" class="profile-list">
            <div v-for="v in myVacancies" :key="v.id" class="profile-item">
              <div class="profile-item__left">
                <div class="profile-item__avatar" :style="{ background: '#f59e0b', color: '#92400e' }">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:20px;height:20px;">
                    <rect x="2" y="7" width="20" height="14" rx="2" ry="2"/>
                    <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>
                  </svg>
                </div>
                <div class="profile-item__info">
                  <div class="profile-item__name">{{ v.name }}</div>
                  <div class="profile-item__sub">{{ v.title }}</div>
                  <div class="profile-item__meta">
                    <span>{{ formatPrice(v.price) }}</span>
                    <span>• {{ formatDate(v.created_at) }}</span>
                  </div>
                </div>
              </div>
              <div class="profile-item__right">
                <span class="profile-item__badge" :class="v.is_active ? 'badge--active' : 'badge--inactive'">
                  {{ v.is_active ? 'Faol' : 'Yopiq' }}
                </span>
                <RouterLink :to="{ name: 'vacancy-detail', params: { slug: v.slug } }" class="profile-item__view">Ko'rish</RouterLink>
              </div>
            </div>
          </div>

          <div v-else class="profile-empty">
            <div class="profile-empty__icon">📋</div>
            <h3>Vakansiyalaringiz yo'q</h3>
            <p>Birinchi vakansiyangizni e'lon qiling</p>
            <RouterLink :to="{ name: 'vacancy-create' }" class="profile-create-btn profile-create-btn--lg">
              Vakansiya e'lon qilish
            </RouterLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #f3f4f6;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.profile-hero {
  background: linear-gradient(135deg, #1a3fbd 0%, #1d56db 40%, #2563eb 100%);
  padding: 40px 20px;
}

.profile-hero__inner {
  max-width: 1000px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.profile-hero__avatar {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: rgba(255,255,255,0.2);
  border: 3px solid rgba(255,255,255,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 32px;
  font-weight: 800;
  flex-shrink: 0;
}

.profile-hero__info { flex: 1; }

.profile-hero__name {
  font-size: 24px;
  font-weight: 800;
  color: #fff;
  margin-bottom: 4px;
}

.profile-hero__phone { font-size: 14px; color: rgba(255,255,255,0.75); margin-bottom: 8px; }

.profile-hero__role {
  display: inline-block;
  background: rgba(255,255,255,0.2);
  color: #fff;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}

.profile-hero__logout {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255,255,255,0.15);
  border: 1.5px solid rgba(255,255,255,0.3);
  color: #fff;
  padding: 10px 18px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.2s;
  margin-left: auto;
}

.profile-hero__logout svg { width: 16px; height: 16px; }
.profile-hero__logout:hover { background: rgba(255,255,255,0.25); }

.profile-main { padding: 28px 20px 48px; }

.profile-main__inner { max-width: 1000px; margin: 0 auto; }

.profile-tabs {
  display: flex;
  gap: 4px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 4px;
  margin-bottom: 24px;
}

.profile-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: #6b7280;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.2s;
}

.profile-tab svg { width: 16px; height: 16px; }

.profile-tab:hover { background: #f3f4f6; color: #374151; }

.profile-tab--active {
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  color: #fff;
}

.profile-section__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.profile-section__header h2 {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.profile-create-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  color: #fff;
  padding: 10px 18px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
  transition: opacity 0.2s;
}

.profile-create-btn svg { width: 14px; height: 14px; }
.profile-create-btn:hover { opacity: 0.9; }

.profile-create-btn--lg { padding: 12px 24px; font-size: 15px; margin-top: 16px; }

.profile-loading {
  text-align: center;
  padding: 40px;
  color: #6b7280;
  font-size: 14px;
}

.profile-list { display: flex; flex-direction: column; gap: 12px; }

.profile-item {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.2s;
}

.profile-item:hover { border-color: #93c5fd; box-shadow: 0 2px 8px rgba(37,99,235,0.08); }

.profile-item__left { display: flex; align-items: center; gap: 14px; flex: 1; min-width: 0; }

.profile-item__avatar {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  flex-shrink: 0;
}

.profile-item__info { flex: 1; min-width: 0; }

.profile-item__name {
  font-size: 15px;
  font-weight: 700;
  color: #111827;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.profile-item__sub { font-size: 13px; color: #6b7280; margin-top: 2px; }

.profile-item__meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: #9ca3af;
  margin-top: 4px;
  flex-wrap: wrap;
}

.profile-item__right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.profile-item__badge {
  font-size: 11px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 20px;
}

.badge--active { background: #d1fae5; color: #065f46; }
.badge--inactive { background: #f3f4f6; color: #6b7280; }

.profile-item__view {
  font-size: 13px;
  font-weight: 600;
  color: #2563eb;
  text-decoration: none;
  padding: 6px 14px;
  border: 1.5px solid #bfdbfe;
  border-radius: 8px;
  transition: all 0.2s;
}

.profile-item__view:hover { background: #eff6ff; }

.profile-empty {
  text-align: center;
  padding: 60px 20px;
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.profile-empty__icon { font-size: 48px; margin-bottom: 16px; }

.profile-empty h3 { font-size: 18px; font-weight: 700; color: #111827; margin-bottom: 8px; }
.profile-empty p { color: #6b7280; font-size: 14px; }

@media (max-width: 640px) {
  .profile-tabs { flex-direction: column; }
  .profile-hero__logout { margin-left: 0; width: 100%; justify-content: center; }
}
</style>
