<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ClientHeader from '../components/ClientHeader.vue'
import FilterSidebar from '../components/FilterSidebar.vue'
import VacancyCard from '../components/VacancyCard.vue'
import { getVacancies } from '../api/vacancyApi'
import { getResumes } from '../api/resumeApi'

const route = useRoute()
const router = useRouter()

const search = ref('')
const loading = ref(false)
const vacancies = ref([])
const totalVacancies = ref(0)
const totalMasters = ref(0)
const currentPage = ref(1)
const lastPage = ref(1)
const sortBy = ref('newest')
const mobileFilterOpen = ref(false)

const filters = reactive({
  category_ids: [],
  region_id: null,
  district_id: null,
  mahalla_id: null,
  min_price: '',
  max_price: '',
  is_active_filter: '',
  search: '',
})

const sortOptions = [
  { value: 'newest', label: 'Eng yangi' },
  { value: 'price_asc', label: "Narx ↑" },
  { value: 'price_desc', label: "Narx ↓" },
]

const activeFilterCount = computed(() => {
  let count = 0
  if (filters.region_id) count++
  if (filters.min_price || filters.max_price) count++
  if (filters.is_active_filter) count++
  return count
})

const syncToUrl = () => {
  const q = {}
  if (filters.search) q.search = filters.search
  if (filters.region_id) q.region_id = filters.region_id
  if (filters.district_id) q.district_id = filters.district_id
  if (filters.mahalla_id) q.mahalla_id = filters.mahalla_id
  if (filters.min_price) q.min_price = filters.min_price
  if (filters.max_price) q.max_price = filters.max_price
  if (filters.is_active_filter) q.is_active_filter = filters.is_active_filter
  if (sortBy.value !== 'newest') q.sort = sortBy.value
  if (currentPage.value > 1) q.page = currentPage.value
  router.replace({ query: q })
}

const loadFromUrl = () => {
  const q = route.query
  if (q.search) { filters.search = q.search; search.value = q.search }
  if (q.region_id) filters.region_id = q.region_id
  if (q.district_id) filters.district_id = q.district_id
  if (q.mahalla_id) filters.mahalla_id = q.mahalla_id
  if (q.min_price) filters.min_price = q.min_price
  if (q.max_price) filters.max_price = q.max_price
  if (q.is_active_filter) filters.is_active_filter = q.is_active_filter
  if (q.sort) sortBy.value = q.sort
  if (q.page) currentPage.value = Number(q.page) || 1
}

let searchTimer = null
watch(search, (val) => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    filters.search = val
    currentPage.value = 1
    syncToUrl()
    fetchVacancies()
  }, 450)
})

const buildParams = () => {
  const params = { page: currentPage.value, limit: 9 }
  if (filters.search) params.search = filters.search
  if (filters.region_id) params.region_id = filters.region_id
  if (filters.district_id) params.district_id = filters.district_id
  if (filters.mahalla_id) params.mahalla_id = filters.mahalla_id
  if (filters.min_price) params.min_price = filters.min_price
  if (filters.max_price) params.max_price = filters.max_price
  if (filters.is_active_filter === 'active') params.is_active = true
  else if (filters.is_active_filter === 'inactive') params.is_active = false
  if (sortBy.value === 'price_asc') { params.sort_by = 'price'; params.sort_order = 'asc' }
  else if (sortBy.value === 'price_desc') { params.sort_by = 'price'; params.sort_order = 'desc' }
  else { params.sort_by = 'id'; params.sort_order = 'desc' }
  return params
}

const fetchVacancies = async () => {
  loading.value = true
  try {
    const res = await getVacancies(buildParams())
    vacancies.value = res.data?.data || []
    const meta = res.data?.meta || {}
    totalVacancies.value = meta.total || 0
    lastPage.value = Math.ceil((meta.total || 0) / 9) || 1
  } catch {
    vacancies.value = []
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const res = await getResumes({ limit: 1 })
    totalMasters.value = res.data?.meta?.total || 0
  } catch {}
}

const applyFilters = () => {
  currentPage.value = 1
  syncToUrl()
  fetchVacancies()
}

const onSortChange = () => {
  currentPage.value = 1
  syncToUrl()
  fetchVacancies()
}

const changePage = (p) => {
  if (p < 1 || p > lastPage.value) return
  currentPage.value = p
  syncToUrl()
  fetchVacancies()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const pageNumbers = computed(() => {
  const pages = []
  const delta = 2
  for (let i = Math.max(1, currentPage.value - delta); i <= Math.min(lastPage.value, currentPage.value + delta); i++) {
    pages.push(i)
  }
  return pages
})

onMounted(async () => {
  loadFromUrl()
  await Promise.all([fetchVacancies(), fetchStats()])
})
</script>

<template>
  <div class="vacancies-page">
    <ClientHeader v-model="search" />

    <section class="hero">
      <div class="hero__inner">
        <h1 class="hero__title">O'zbekistondagi eng yaxshi<br />ish e'lonlarini toping</h1>
        <p class="hero__sub">Kompaniyalar va tadbirkorlar tomonidan e'lon qilingan vakansiyalar</p>
      </div>
    </section>

    <div class="stats-row">
      <div class="stats-row__inner">
        <RouterLink to="/" class="stats-card stats-card--white">
          <div class="stats-card__icon stats-card__icon--blue">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
          </div>
          <div class="stats-card__body">
            <div class="stats-card__number stats-card__number--blue">{{ totalMasters }}</div>
            <div class="stats-card__label">Usta / Mutaxassis</div>
          </div>
          <svg class="stats-card__arrow--dark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="m9 18 6-6-6-6"/>
          </svg>
        </RouterLink>

        <div class="stats-card stats-card--blue">
          <div class="stats-card__icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="7" width="20" height="14" rx="2" ry="2"/>
              <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>
            </svg>
          </div>
          <div class="stats-card__body">
            <div class="stats-card__number">{{ totalVacancies }}</div>
            <div class="stats-card__label">Ish e'loni</div>
          </div>
        </div>
      </div>
    </div>

    <div class="main-content">
      <div class="main-content__inner">
        <FilterSidebar
          :model-value="filters"
          @update:model-value="Object.assign(filters, $event)"
          :show-categories="false"
          :mobile-open="mobileFilterOpen"
          @apply="applyFilters"
          @close="mobileFilterOpen = false"
        />

        <div class="vacancies-section">
          <div class="vacancies-section__toolbar">
            <span class="vacancies-section__count">
              <strong>{{ totalVacancies }}</strong> ta vakansiya
            </span>
            <div class="toolbar-right">
              <button class="filter-toggle-btn" @click="mobileFilterOpen = true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
                </svg>
                Filtr
                <span v-if="activeFilterCount > 0" class="filter-toggle-btn__badge">{{ activeFilterCount }}</span>
              </button>
              <select v-model="sortBy" class="vacancies-section__sort" @change="onSortChange">
                <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>
            </div>
          </div>

          <div v-if="loading" class="vacancies-grid vacancies-grid--loading">
            <div v-for="i in 9" :key="i" class="skeleton-card">
              <div class="skeleton skeleton--avatar"></div>
              <div class="skeleton-lines">
                <div class="skeleton skeleton--line skeleton--w80"></div>
                <div class="skeleton skeleton--line skeleton--w60"></div>
              </div>
            </div>
          </div>

          <div v-else-if="vacancies.length" class="vacancies-grid">
            <VacancyCard v-for="v in vacancies" :key="v.id" :item="v" />
          </div>

          <div v-else class="empty-state">
            <div class="empty-state__icon">📋</div>
            <h3>Vakansiya topilmadi</h3>
            <p>Filtrlarni o'zgartirib ko'ring</p>
          </div>

          <div v-if="lastPage > 1" class="pagination">
            <button class="pagination__btn" :disabled="currentPage <= 1" @click="changePage(currentPage - 1)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m15 18-6-6 6-6"/></svg>
            </button>
            <button v-if="currentPage > 3" class="pagination__btn" @click="changePage(1)">1</button>
            <span v-if="currentPage > 4" class="pagination__dots">…</span>
            <button v-for="p in pageNumbers" :key="p" class="pagination__btn"
              :class="{ 'pagination__btn--active': p === currentPage }" @click="changePage(p)">{{ p }}</button>
            <span v-if="currentPage < lastPage - 3" class="pagination__dots">…</span>
            <button v-if="currentPage < lastPage - 2" class="pagination__btn" @click="changePage(lastPage)">{{ lastPage }}</button>
            <button class="pagination__btn" :disabled="currentPage >= lastPage" @click="changePage(currentPage + 1)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m9 18 6-6-6-6"/></svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vacancies-page { min-height: 100vh; background: #f3f4f6; color: #111827; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; }

.hero { background: linear-gradient(135deg, #1e3a8a 0%, #1d4ed8 50%, #2563eb 100%); padding: 36px 20px 40px; }
.hero__inner { max-width: 680px; margin: 0 auto; text-align: center; }
.hero__title { font-size: 34px; font-weight: 800; color: #fff; line-height: 1.2; margin-bottom: 10px; }
.hero__sub { font-size: 15px; color: rgba(255,255,255,0.8); }

.stats-row { background: #fff; border-bottom: 1px solid #e5e7eb; padding: 16px 20px; }
.stats-row__inner { max-width: 1280px; margin: 0 auto; display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.stats-card { display: flex; align-items: center; gap: 14px; padding: 16px 20px; border-radius: 14px; transition: all 0.2s; text-decoration: none; cursor: pointer; }
.stats-card--blue { background: linear-gradient(135deg, #1d4ed8, #2563eb); color: #fff; }
.stats-card--white { background: #fff; border: 1.5px solid #e5e7eb; color: #111827; }
.stats-card--white:hover { border-color: #93c5fd; box-shadow: 0 4px 12px rgba(37,99,235,0.1); }
.stats-card__icon { width: 46px; height: 46px; border-radius: 12px; background: rgba(255,255,255,0.2); display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.stats-card__icon svg { width: 22px; height: 22px; }
.stats-card__icon--blue { background: #dbeafe; color: #1d4ed8; }
.stats-card__body { flex: 1; min-width: 0; }
.stats-card__number { font-size: 28px; font-weight: 800; line-height: 1; color: #fff; }
.stats-card__number--blue { color: #1d4ed8; }
.stats-card--white .stats-card__number { color: #111827; }
.stats-card__label { font-size: 12px; color: rgba(255,255,255,0.75); margin-top: 3px; }
.stats-card--white .stats-card__label { color: #6b7280; }
.stats-card__arrow--dark { width: 18px; height: 18px; color: #9ca3af; flex-shrink: 0; }

.main-content { padding: 24px 16px 48px; }
.main-content__inner { max-width: 1280px; margin: 0 auto; display: grid; grid-template-columns: 272px 1fr; gap: 20px; align-items: start; }

.vacancies-section { display: flex; flex-direction: column; gap: 16px; }
.vacancies-section__toolbar { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; }
.vacancies-section__count { font-size: 14px; color: #374151; }
.vacancies-section__count strong { color: #111827; font-weight: 700; }
.toolbar-right { display: flex; align-items: center; gap: 8px; }

.filter-toggle-btn { display: none; align-items: center; gap: 6px; padding: 9px 14px; background: #fff; border: 1.5px solid #e5e7eb; border-radius: 10px; font-size: 13px; font-weight: 600; color: #374151; cursor: pointer; font-family: inherit; transition: all 0.2s; position: relative; }
.filter-toggle-btn svg { width: 15px; height: 15px; }
.filter-toggle-btn:hover { border-color: #2563eb; color: #2563eb; background: #eff6ff; }
.filter-toggle-btn__badge { position: absolute; top: -6px; right: -6px; background: #2563eb; color: #fff; font-size: 10px; font-weight: 700; width: 18px; height: 18px; border-radius: 50%; display: flex; align-items: center; justify-content: center; }

.vacancies-section__sort { padding: 9px 32px 9px 12px; border: 1.5px solid #e5e7eb; border-radius: 10px; font-size: 13px; color: #374151; background: #fff url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%236b7280' stroke-width='2'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E") no-repeat right 8px center / 14px; outline: none; cursor: pointer; appearance: none; font-family: inherit; transition: border-color 0.2s; }
.vacancies-section__sort:focus { border-color: #2563eb; }

.vacancies-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
.vacancies-grid--loading { pointer-events: none; }

.skeleton-card { background: #fff; border: 1px solid #e5e7eb; border-radius: 16px; padding: 18px; display: flex; align-items: flex-start; gap: 12px; }
.skeleton { background: linear-gradient(90deg, #f3f4f6 25%, #e5e7eb 50%, #f3f4f6 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; border-radius: 8px; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
.skeleton--avatar { width: 48px; height: 48px; border-radius: 14px; flex-shrink: 0; }
.skeleton-lines { flex: 1; display: flex; flex-direction: column; gap: 8px; padding-top: 4px; }
.skeleton--line { height: 14px; }
.skeleton--w80 { width: 80%; }
.skeleton--w60 { width: 60%; }

.empty-state { text-align: center; padding: 60px 20px; background: #fff; border-radius: 16px; border: 1px solid #e5e7eb; }
.empty-state__icon { font-size: 48px; margin-bottom: 16px; }
.empty-state h3 { font-size: 18px; font-weight: 700; color: #111827; margin-bottom: 8px; }
.empty-state p { color: #6b7280; font-size: 14px; }

.pagination { display: flex; justify-content: center; align-items: center; gap: 5px; flex-wrap: wrap; }
.pagination__btn { min-width: 38px; height: 38px; padding: 0 10px; background: #fff; color: #374151; border: 1.5px solid #e5e7eb; border-radius: 10px; cursor: pointer; font-size: 14px; font-weight: 600; font-family: inherit; display: flex; align-items: center; justify-content: center; transition: all 0.2s; }
.pagination__btn svg { width: 16px; height: 16px; }
.pagination__btn:hover:not(:disabled) { border-color: #2563eb; color: #2563eb; background: #eff6ff; }
.pagination__btn--active { background: linear-gradient(135deg, #1d4ed8, #2563eb) !important; color: #fff !important; border-color: transparent !important; box-shadow: 0 2px 8px rgba(37,99,235,0.3); }
.pagination__btn:disabled { opacity: 0.4; cursor: not-allowed; }
.pagination__dots { color: #9ca3af; font-size: 14px; padding: 0 4px; }

@media (max-width: 1100px) { .vacancies-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 900px) { .main-content__inner { grid-template-columns: 1fr; } .filter-toggle-btn { display: flex; } .hero__title { font-size: 26px; } }
@media (max-width: 640px) { .hero { padding: 28px 16px 32px; } .hero__title { font-size: 21px; } .hero__sub { font-size: 13px; } .stats-row { padding: 12px 16px; } .main-content { padding: 16px 12px 40px; } }
@media (max-width: 520px) { .vacancies-grid { grid-template-columns: 1fr; } .stats-row__inner { grid-template-columns: 1fr 1fr; gap: 8px; } .stats-card { padding: 12px 14px; gap: 10px; } .stats-card__icon { width: 38px; height: 38px; border-radius: 10px; } .stats-card__number { font-size: 22px; } .stats-card__label { font-size: 11px; } }
</style>
