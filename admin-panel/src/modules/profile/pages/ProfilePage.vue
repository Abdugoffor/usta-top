<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store/authStore'
import { getResumes, deleteResume } from '@/modules/client/api/resumeApi'
import { getVacancies, deleteVacancy } from '@/modules/client/api/vacancyApi'
import ClientHeader from '@/modules/client/components/ClientHeader.vue'
import { formatNumber } from '@/shared/utils/formatNumber'

const router = useRouter()
const auth = useAuthStore()

const LIMIT = 20

const myResumes = ref([])
const resumesTotal = ref(0)
const resumesCursor = ref(null)
const resumesHasMore = ref(false)
const loadingResumes = ref(false)
const loadingMoreResumes = ref(false)

const myVacancies = ref([])
const vacanciesTotal = ref(0)
const vacanciesCursor = ref(null)
const vacanciesHasMore = ref(false)
const loadingVacancies = ref(false)
const loadingMoreVacancies = ref(false)

const activeTab = ref('resumes')

const fetchMyResumes = async () => {
  if (!auth.user?.id) return
  loadingResumes.value = true
  try {
    const res = await getResumes({ user_id: auth.user.id, limit: LIMIT })
    const { data, meta } = res.data
    myResumes.value = data || []
    resumesTotal.value = meta?.total ?? 0
    resumesCursor.value = meta?.next_cursor || null
    resumesHasMore.value = meta?.has_more || false
  } catch {
    myResumes.value = []
  } finally {
    loadingResumes.value = false
  }
}

const loadMoreResumes = async () => {
  if (!resumesCursor.value || loadingMoreResumes.value) return
  loadingMoreResumes.value = true
  try {
    const res = await getResumes({ user_id: auth.user.id, limit: LIMIT, cursor: resumesCursor.value })
    const { data, meta } = res.data
    myResumes.value.push(...(data || []))
    resumesCursor.value = meta?.next_cursor || null
    resumesHasMore.value = meta?.has_more || false
  } catch {
    // ignore
  } finally {
    loadingMoreResumes.value = false
  }
}

const fetchMyVacancies = async () => {
  if (!auth.user?.id) return
  loadingVacancies.value = true
  try {
    const res = await getVacancies({ user_id: auth.user.id, limit: LIMIT })
    const { data, meta } = res.data
    myVacancies.value = data || []
    vacanciesTotal.value = meta?.total ?? 0
    vacanciesCursor.value = meta?.next_cursor || null
    vacanciesHasMore.value = meta?.has_more || false
  } catch {
    myVacancies.value = []
  } finally {
    loadingVacancies.value = false
  }
}

const loadMoreVacancies = async () => {
  if (!vacanciesCursor.value || loadingMoreVacancies.value) return
  loadingMoreVacancies.value = true
  try {
    const res = await getVacancies({ user_id: auth.user.id, limit: LIMIT, cursor: vacanciesCursor.value })
    const { data, meta } = res.data
    myVacancies.value.push(...(data || []))
    vacanciesCursor.value = meta?.next_cursor || null
    vacanciesHasMore.value = meta?.has_more || false
  } catch {
    // ignore
  } finally {
    loadingMoreVacancies.value = false
  }
}

const doLogout = () => {
  auth.logout()
  router.push({ name: 'home' })
}

const doDeleteResume = async (id) => {
  if (!confirm("Resumeni o'chirishni tasdiqlaysizmi?")) return
  try {
    await deleteResume(id)
    myResumes.value = myResumes.value.filter(r => r.id !== id)
  } catch {
    alert("O'chirishda xatolik yuz berdi")
  }
}

const doDeleteVacancy = async (id) => {
  if (!confirm("Vakansiyani o'chirishni tasdiqlaysizmi?")) return
  try {
    await deleteVacancy(id)
    myVacancies.value = myVacancies.value.filter(v => v.id !== id)
  } catch {
    alert("O'chirishda xatolik yuz berdi")
  }
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
            Resumelarim ({{ formatNumber(resumesTotal) }})
          </button>
          <button class="profile-tab" :class="{ 'profile-tab--active': activeTab === 'vacancies' }" @click="activeTab = 'vacancies'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="7" width="20" height="14" rx="2" ry="2"/>
              <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>
            </svg>
            Vakansiyalarim ({{ formatNumber(vacanciesTotal) }})
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

          <template v-else>
            <div v-if="myResumes.length" class="profile-list">
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
                  <RouterLink :to="{ name: 'resume-edit', params: { id: r.slug } }" class="profile-item__edit">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                  </RouterLink>
                  <button class="profile-item__del" @click="doDeleteResume(r.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4h6v2"/></svg>
                  </button>
                </div>
              </div>
            </div>

            <div v-if="resumesHasMore" class="profile-load-more">
              <button class="load-more-btn" :disabled="loadingMoreResumes" @click="loadMoreResumes">
                <span v-if="loadingMoreResumes" class="load-more-btn__spinner"></span>
                <template v-else>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <path d="M12 5v14M5 12l7 7 7-7"/>
                  </svg>
                  Ko'proq yuklash
                </template>
              </button>
            </div>

            <div v-if="!myResumes.length" class="profile-empty">
              <div class="profile-empty__icon">📄</div>
              <h3>Resumelaringiz yo'q</h3>
              <p>Birinchi resumeingizni yarating</p>
              <RouterLink :to="{ name: 'resume-create' }" class="profile-create-btn profile-create-btn--lg">
                Resume yaratish
              </RouterLink>
            </div>
          </template>
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

          <template v-else>
            <div v-if="myVacancies.length" class="profile-list">
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
                  <RouterLink :to="{ name: 'vacancy-edit', params: { id: v.slug } }" class="profile-item__edit">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                  </RouterLink>
                  <button class="profile-item__del" @click="doDeleteVacancy(v.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4h6v2"/></svg>
                  </button>
                </div>
              </div>
            </div>

            <div v-if="vacanciesHasMore" class="profile-load-more">
              <button class="load-more-btn" :disabled="loadingMoreVacancies" @click="loadMoreVacancies">
                <span v-if="loadingMoreVacancies" class="load-more-btn__spinner"></span>
                <template v-else>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <path d="M12 5v14M5 12l7 7 7-7"/>
                  </svg>
                  Ko'proq yuklash
                </template>
              </button>
            </div>

            <div v-if="!myVacancies.length" class="profile-empty">
              <div class="profile-empty__icon">📋</div>
              <h3>Vakansiyalaringiz yo'q</h3>
              <p>Birinchi vakansiyangizni e'lon qiling</p>
              <RouterLink :to="{ name: 'vacancy-create' }" class="profile-create-btn profile-create-btn--lg">
                Vakansiya e'lon qilish
              </RouterLink>
            </div>
          </template>
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

.profile-item__edit {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: 1.5px solid #e5e7eb;
  border-radius: 8px;
  color: #6b7280;
  text-decoration: none;
  transition: all 0.2s;
  flex-shrink: 0;
}

.profile-item__edit svg { width: 15px; height: 15px; }
.profile-item__edit:hover { border-color: #2563eb; color: #2563eb; background: #eff6ff; }

.profile-item__del {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: 1.5px solid #e5e7eb;
  border-radius: 8px;
  color: #6b7280;
  background: none;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.profile-item__del svg { width: 15px; height: 15px; }
.profile-item__del:hover { border-color: #dc2626; color: #dc2626; background: #fef2f2; }

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

.profile-load-more {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.load-more-btn {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 13px 36px;
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  border: none;
  color: #fff;
  border-radius: 50px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.25s ease;
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.3);
}

.load-more-btn svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  transition: transform 0.25s ease;
}

.load-more-btn:hover:not(:disabled) {
  box-shadow: 0 6px 22px rgba(37, 99, 235, 0.45);
  transform: translateY(-2px);
}

.load-more-btn:hover:not(:disabled) svg { transform: translateY(3px); }

.load-more-btn:disabled { opacity: 0.65; cursor: not-allowed; }

.load-more-btn__spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;
}

@keyframes spin { to { transform: rotate(360deg); } }

.profile-empty__icon { font-size: 48px; margin-bottom: 16px; }

.profile-empty h3 { font-size: 18px; font-weight: 700; color: #111827; margin-bottom: 8px; }
.profile-empty p { color: #6b7280; font-size: 14px; }

@media (max-width: 640px) {
  .profile-tabs { flex-direction: column; }
  .profile-hero__logout { margin-left: 0; width: 100%; justify-content: center; }
}
</style>
