<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ClientHeader from '../components/ClientHeader.vue'
import FilterSidebar from '../components/FilterSidebar.vue'
import MasterCard from '../components/MasterCard.vue'
import { getResumes } from '../api/resumeApi'
import { getVacancies } from '../api/vacancyApi'

const route = useRoute()
const router = useRouter()

const search = ref('')
const loading = ref(false)
const loadingMore = ref(false)
const masters = ref([])
const nextCursor = ref('')
const hasMore = ref(false)
const totalMasters = ref(0)
const totalVacancies = ref(0)
const sortBy = ref('newest')
const mobileFilterOpen = ref(false)
const sentinel = ref(null)

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
  { value: 'experience', label: 'Tajriba' },
]

const activeFilterCount = computed(() => {
  let count = 0
  if (filters.category_ids?.length) count++
  if (filters.region_id) count++
  if (filters.min_price || filters.max_price) count++
  if (filters.is_active_filter) count++
  return count
})


const syncToUrl = () => {
  const q = {}
  if (filters.search) q.search = filters.search
  if (filters.category_ids?.length) q.category_ids = filters.category_ids.join(',')
  if (filters.region_id) q.region_id = filters.region_id
  if (filters.district_id) q.district_id = filters.district_id
  if (filters.mahalla_id) q.mahalla_id = filters.mahalla_id
  if (filters.min_price) q.min_price = filters.min_price
  if (filters.max_price) q.max_price = filters.max_price
  if (filters.is_active_filter) q.is_active_filter = filters.is_active_filter
  if (sortBy.value !== 'newest') q.sort = sortBy.value
  router.replace({ query: q })
}

const loadFromUrl = () => {
  const q = route.query
  if (q.search) { filters.search = q.search; search.value = q.search }
  if (q.category_ids) filters.category_ids = q.category_ids.split(',').map(Number).filter(Boolean)
  if (q.region_id) filters.region_id = q.region_id
  if (q.district_id) filters.district_id = q.district_id
  if (q.mahalla_id) filters.mahalla_id = q.mahalla_id
  if (q.min_price) filters.min_price = q.min_price
  if (q.max_price) filters.max_price = q.max_price
  if (q.is_active_filter) filters.is_active_filter = q.is_active_filter
  if (q.sort) sortBy.value = q.sort
}

let searchTimer = null
watch(search, (val) => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    filters.search = val
    syncToUrl()
    resetAndFetch()
  }, 450)
})

const buildParams = (cursor = '') => {
  const params = { limit: 9 }
  if (cursor) params.cursor = cursor
  if (filters.search) params.search = filters.search
  if (filters.category_ids?.length) params.category_ids = filters.category_ids.join(',')
  if (filters.region_id) params.region_id = filters.region_id
  if (filters.district_id) params.district_id = filters.district_id
  if (filters.mahalla_id) params.mahalla_id = filters.mahalla_id
  if (filters.min_price) params.min_price = filters.min_price
  if (filters.max_price) params.max_price = filters.max_price
  if (filters.is_active_filter === 'active') params.is_active = true
  else if (filters.is_active_filter === 'inactive') params.is_active = false
  if (sortBy.value === 'price_asc') { params.sort_by = 'price'; params.sort_order = 'asc' }
  else if (sortBy.value === 'price_desc') { params.sort_by = 'price'; params.sort_order = 'desc' }
  else if (sortBy.value === 'experience') { params.sort_by = 'experience_year'; params.sort_order = 'desc' }
  return params
}

const fetchMasters = async () => {
  loading.value = true
  try {
    const res = await getResumes(buildParams(''))
    masters.value = res.data?.data || []
    const meta = res.data?.meta || {}
    hasMore.value = meta.has_more || false
    nextCursor.value = meta.next_cursor || ''
    totalMasters.value = meta.total || 0
  } catch {
    masters.value = []
    hasMore.value = false
    nextCursor.value = ''
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  if (!hasMore.value || loadingMore.value || !nextCursor.value) return
  loadingMore.value = true
  const cursorToUse = nextCursor.value
  try {
    const res = await getResumes(buildParams(cursorToUse))
    const newItems = res.data?.data || []
    masters.value = [...masters.value, ...newItems]
    const meta = res.data?.meta || {}
    hasMore.value = meta.has_more || false
    nextCursor.value = meta.next_cursor || ''
  } catch {
    hasMore.value = false
  } finally {
    loadingMore.value = false
  }
}

const resetAndFetch = () => {
  masters.value = []
  nextCursor.value = ''
  hasMore.value = false
  fetchMasters()
}

const fetchVacancyCount = async () => {
  try {
    const res = await getVacancies({ limit: 1 })
    totalVacancies.value = res.data?.meta?.total || 0
  } catch { }
}

const applyFilters = () => {
  syncToUrl()
  resetAndFetch()
}

const onSortChange = () => {
  syncToUrl()
  resetAndFetch()
}

let observer = null

onMounted(async () => {
  loadFromUrl()
  await Promise.all([fetchMasters(), fetchVacancyCount()])

  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) loadMore()
  }, { rootMargin: '300px' })

  if (sentinel.value) observer.observe(sentinel.value)
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})
</script>

<template>
  <div class="home-page">
    <ClientHeader v-model="search" />

    <section class="hero">
      <div class="hero__inner">
        <h1 class="hero__title">O'zbekistondagi eng yaxshi<br />ustalarni va ishlarni toping</h1>
        <p class="hero__sub">Elektrik, santexnik, qurilish va yuzlab mutaxassislar bir joyda</p>
      </div>
    </section>

    <div class="stats-row">
      <div class="stats-row__inner">
        <div class="stats-card stats-card--blue">
          <div class="stats-card__icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          </div>
          <div class="stats-card__body">
            <div class="stats-card__number">{{ totalMasters }}</div>
            <div class="stats-card__label">Usta / Mutaxassis</div>
          </div>
        </div>

        <RouterLink to="/vacancies" class="stats-card stats-card--white">
          <div class="stats-card__icon stats-card__icon--amber">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="7" width="20" height="14" rx="2" ry="2" />
              <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16" />
            </svg>
          </div>
          <div class="stats-card__body">
            <div class="stats-card__number stats-card__number--amber">{{ totalVacancies }}</div>
            <div class="stats-card__label">Ish e'loni</div>
          </div>
          <svg class="stats-card__arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="m9 18 6-6-6-6" />
          </svg>
        </RouterLink>
      </div>
    </div>

    <div class="main-content">
      <div class="main-content__inner">
        <div class="sidebar-sticky">
          <FilterSidebar :model-value="filters" @update:model-value="Object.assign(filters, $event)"
            :mobile-open="mobileFilterOpen" @apply="applyFilters" @close="mobileFilterOpen = false" />
        </div>
        <div class="masters-section">
          <div class="masters-section__toolbar">
            <span class="masters-section__count">
              <strong>{{ totalMasters }}</strong> ta usta
            </span>
            <div class="toolbar-right">
              <button class="filter-toggle-btn" @click="mobileFilterOpen = true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
                </svg>
                Filtr
                <span v-if="activeFilterCount > 0" class="filter-toggle-btn__badge">{{ activeFilterCount }}</span>
              </button>
              <select v-model="sortBy" class="masters-section__sort" @change="onSortChange">
                <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>
            </div>
          </div>

          <div v-if="loading" class="masters-grid masters-grid--loading">
            <div v-for="i in 9" :key="i" class="skeleton-card">
              <div class="skeleton skeleton--avatar"></div>
              <div class="skeleton-lines">
                <div class="skeleton skeleton--line skeleton--w80"></div>
                <div class="skeleton skeleton--line skeleton--w60"></div>
              </div>
            </div>
          </div>

          <template v-else>
            <div v-if="masters.length" class="masters-grid">
              <MasterCard v-for="m in masters" :key="m.id" :item="m" />
            </div>

            <div v-else class="empty-state">
              <div class="empty-state__icon">🔍</div>
              <h3>Usta topilmadi</h3>
              <p>Filtrlarni o'zgartirib ko'ring</p>
            </div>
          </template>

          <div ref="sentinel" class="sentinel"></div>

          <div v-if="loadingMore" class="load-more-spinner">
            <div class="spinner"></div>
          </div>

          <div v-if="!hasMore && masters.length > 0 && !loading" class="end-of-list">
            Barcha ustalar ko'rsatildi
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-page {
  min-height: 100vh;
  background: #f3f4f6;
  color: #111827;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.hero {
  background: linear-gradient(135deg, #1a3fbd 0%, #1d56db 40%, #2563eb 100%);
  padding: 36px 20px 40px;
}

.hero__inner {
  max-width: 680px;
  margin: 0 auto;
  text-align: center;
}

.hero__title {
  font-size: 34px;
  font-weight: 800;
  color: #fff;
  line-height: 1.2;
  margin-bottom: 10px;
}

.hero__sub {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.8);
}

.stats-row {
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  padding: 16px 20px;
}

.stats-row__inner {
  max-width: 1280px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.stats-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 20px;
  border-radius: 14px;
  transition: all 0.2s;
  text-decoration: none;
  cursor: pointer;
}

.stats-card--blue {
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  color: #fff;
}

.stats-card--white {
  background: #fff;
  border: 1.5px solid #e5e7eb;
  color: #111827;
}

.stats-card--white:hover {
  border-color: #93c5fd;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.1);
}

.stats-card__icon {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stats-card__icon svg {
  width: 22px;
  height: 22px;
}

.stats-card__icon--amber {
  background: #fef3c7;
  color: #d97706;
}

.stats-card__body {
  flex: 1;
  min-width: 0;
}

.stats-card__number {
  font-size: 28px;
  font-weight: 800;
  line-height: 1;
  color: #fff;
}

.stats-card__number--amber {
  color: #d97706;
}

.stats-card--white .stats-card__number {
  color: #111827;
}

.stats-card__label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.75);
  margin-top: 3px;
}

.stats-card--white .stats-card__label {
  color: #6b7280;
}

.stats-card__arrow {
  width: 18px;
  height: 18px;
  color: rgba(255, 255, 255, 0.7);
  flex-shrink: 0;
}

.main-content {
  padding: 24px 16px 48px;
}

.main-content__inner {
  max-width: 1280px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 272px 1fr;
  gap: 20px;
  align-items: start;
}

.masters-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.masters-section__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  flex-wrap: wrap;
}

.masters-section__count {
  font-size: 14px;
  color: #374151;
}

.masters-section__count strong {
  color: #111827;
  font-weight: 700;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-toggle-btn {
  display: none;
  align-items: center;
  gap: 6px;
  padding: 9px 14px;
  background: #fff;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.2s;
  position: relative;
}

.filter-toggle-btn svg {
  width: 15px;
  height: 15px;
}

.filter-toggle-btn:hover {
  border-color: #2563eb;
  color: #2563eb;
  background: #eff6ff;
}

.filter-toggle-btn__badge {
  position: absolute;
  top: -6px;
  right: -6px;
  background: #2563eb;
  color: #fff;
  font-size: 10px;
  font-weight: 700;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.masters-section__sort {
  padding: 9px 32px 9px 12px;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  font-size: 13px;
  color: #374151;
  background: #fff url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%236b7280' stroke-width='2'%3E%3Cpath d='m6 9 6 6 6-6'/%3E%3C/svg%3E") no-repeat right 8px center / 14px;
  outline: none;
  cursor: pointer;
  appearance: none;
  font-family: inherit;
  transition: border-color 0.2s;
}

.masters-section__sort:focus {
  border-color: #2563eb;
}

.masters-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
}

.masters-grid--loading {
  pointer-events: none;
}

.skeleton-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 18px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.skeleton {
  background: linear-gradient(90deg, #f3f4f6 25%, #e5e7eb 50%, #f3f4f6 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 8px;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }

  100% {
    background-position: -200% 0;
  }
}

.skeleton--avatar {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  flex-shrink: 0;
}

.skeleton-lines {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-top: 4px;
}

.skeleton--line {
  height: 14px;
}

.skeleton--w80 {
  width: 80%;
}

.skeleton--w60 {
  width: 60%;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e5e7eb;
}

.empty-state__icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 8px;
}

.empty-state p {
  color: #6b7280;
  font-size: 14px;
}

.sentinel {
  height: 1px;
}

.load-more-spinner {
  display: flex;
  justify-content: center;
  padding: 24px 0;
}

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid #e5e7eb;
  border-top-color: #2563eb;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.end-of-list {
  text-align: center;
  font-size: 13px;
  color: #9ca3af;
  padding: 16px 0;
}

@media (max-width: 1100px) {
  .masters-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 900px) {
  .main-content__inner {
    grid-template-columns: 1fr;
  }

  .filter-toggle-btn {
    display: flex;
  }

  .hero__title {
    font-size: 26px;
  }
}

@media (max-width: 640px) {
  .hero {
    padding: 28px 16px 32px;
  }

  .hero__title {
    font-size: 21px;
  }

  .hero__sub {
    font-size: 13px;
  }

  .stats-row {
    padding: 12px 16px;
  }

  .main-content {
    padding: 16px 12px 40px;
  }

  .masters-section__toolbar {
    flex-wrap: nowrap;
  }
}

@media (max-width: 520px) {
  .masters-grid {
    grid-template-columns: 1fr;
  }

  .stats-row__inner {
    grid-template-columns: 1fr 1fr;
    gap: 8px;
  }

  .stats-card {
    padding: 12px 14px;
    gap: 10px;
  }

  .stats-card__icon {
    width: 38px;
    height: 38px;
    border-radius: 10px;
  }

  .stats-card__number {
    font-size: 22px;
  }

  .stats-card__label {
    font-size: 11px;
  }
}

.sidebar-sticky {
  position: sticky;
  top: 88px;
  align-self: start;
}
</style>
