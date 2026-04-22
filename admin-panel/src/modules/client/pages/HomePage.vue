<script setup>
import { useRoute, useRouter } from 'vue-router'
import ClientHeader from '../components/ClientHeader.vue'
import FilterSidebar from '../components/FilterSidebar.vue'
import MasterCard from '../components/MasterCard.vue'
import { getResumes } from '../api/resumeApi'
import { getVacancies } from '../api/vacancyApi'
import { useClientListingPage } from '../composables/useClientListingPage'

const route = useRoute()
const router = useRouter()

const sortOptions = [
  { value: 'newest', label: 'Eng yangi' },
  { value: 'price_asc', label: 'Narx ↑' },
  { value: 'price_desc', label: 'Narx ↓' },
  { value: 'experience', label: 'Tajriba' },
]

const syncToUrl = ({ filters, sortBy, router }) => {
  const q = {}
  if (filters.search) q.search = filters.search
  if (filters.category_ids?.length) q.category_ids = filters.category_ids.join(',')
  if (filters.region_id) q.region_id = filters.region_id
  if (filters.district_id) q.district_id = filters.district_id
  if (filters.mahalla_id) q.mahalla_id = filters.mahalla_id
  if (filters.min_price) q.min_price = filters.min_price
  if (filters.max_price) q.max_price = filters.max_price
  if (filters.is_active_filter) q.is_active_filter = filters.is_active_filter
  if (sortBy !== 'newest') q.sort = sortBy
  router.replace({ query: q })
}

const loadFromUrl = ({ query, filters, search, sortBy }) => {
  if (query.search) {
    filters.search = query.search
    search.value = query.search
  }
  if (query.category_ids) filters.category_ids = query.category_ids.split(',').map(Number).filter(Boolean)
  if (query.region_id) filters.region_id = query.region_id
  if (query.district_id) filters.district_id = query.district_id
  if (query.mahalla_id) filters.mahalla_id = query.mahalla_id
  if (query.min_price) filters.min_price = query.min_price
  if (query.max_price) filters.max_price = query.max_price
  if (query.is_active_filter) filters.is_active_filter = query.is_active_filter
  if (query.sort) sortBy.value = query.sort
}

const buildParams = ({ filters, sortBy, cursor }) => {
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
  if (sortBy === 'price_asc') {
    params.sort_by = 'price'
    params.sort_order = 'asc'
  } else if (sortBy === 'price_desc') {
    params.sort_by = 'price'
    params.sort_order = 'desc'
  } else if (sortBy === 'experience') {
    params.sort_by = 'experience_year'
    params.sort_order = 'desc'
  }
  return params
}

const {
  activeFilterCount,
  applyFilters,
  filters,
  hasMore,
  items: masters,
  loading,
  loadingMore,
  mobileFilterOpen,
  onSortChange,
  search,
  secondaryTotal: totalVacancies,
  sentinel,
  sortBy,
  total: totalMasters,
} = useClientListingPage({
  route,
  router,
  createFilters: () => ({
    category_ids: [],
    region_id: null,
    district_id: null,
    mahalla_id: null,
    min_price: '',
    max_price: '',
    is_active_filter: '',
    search: '',
  }),
  fetchList: (params) => getResumes(params),
  fetchSecondaryCount: () => getVacancies({ limit: 1 }),
  buildParams,
  syncToUrl,
  loadFromUrl,
  countActiveFilters: (filters) => {
    let count = 0
    if (filters.category_ids?.length) count++
    if (filters.region_id) count++
    if (filters.min_price || filters.max_price) count++
    if (filters.is_active_filter) count++
    return count
  },
})
</script>

<template>
  <div class="client-page client-page--home">
    <ClientHeader v-model="search" @search="applyFilters" />

    <section class="client-hero">
      <div class="client-hero__inner">
        <h1 class="client-hero__title">O'zbekistondagi eng yaxshi<br />ustalarni va ishlarni toping</h1>
        <p class="client-hero__sub">Elektrik, santexnik, qurilish va yuzlab mutaxassislar bir joyda</p>
      </div>
    </section>

    <div class="client-stats">
      <div class="client-stats__inner">
        <div class="client-stats-card client-stats-card--primary">
          <div class="client-stats-card__icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          </div>
          <div class="client-stats-card__body">
            <div class="client-stats-card__number">{{ totalMasters }}</div>
            <div class="client-stats-card__label">Usta / Mutaxassis</div>
          </div>
        </div>

        <RouterLink to="/vacancies" class="client-stats-card client-stats-card--surface">
          <div class="client-stats-card__icon client-stats-card__icon--accent">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="7" width="20" height="14" rx="2" ry="2" />
              <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16" />
            </svg>
          </div>
          <div class="client-stats-card__body">
            <div class="client-stats-card__number client-stats-card__number--accent">{{ totalVacancies }}</div>
            <div class="client-stats-card__label">Ish e'loni</div>
          </div>
          <svg class="client-stats-card__arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="m9 18 6-6-6-6" />
          </svg>
        </RouterLink>
      </div>
    </div>

    <div class="client-content">
      <div class="client-content__inner">
        <div class="client-sidebar">
          <FilterSidebar
            :model-value="filters"
            :mobile-open="mobileFilterOpen"
            @update:model-value="Object.assign(filters, $event)"
            @apply="applyFilters"
            @close="mobileFilterOpen = false"
          />
        </div>

        <div class="client-results">
          <div class="client-toolbar">
            <span class="client-toolbar__count">
              <strong>{{ totalMasters }}</strong> ta usta
            </span>

            <div class="client-toolbar__actions">
              <button class="client-filter-toggle" @click="mobileFilterOpen = true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
                </svg>
                Filtr
                <span v-if="activeFilterCount > 0" class="client-filter-toggle__badge">{{ activeFilterCount }}</span>
              </button>

              <select v-model="sortBy" class="client-sort" @change="onSortChange">
                <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>
            </div>
          </div>

          <div v-if="loading" class="client-grid client-grid--loading">
            <div v-for="i in 9" :key="i" class="client-skeleton-card">
              <div class="client-skeleton client-skeleton--avatar"></div>
              <div class="client-skeleton-lines">
                <div class="client-skeleton client-skeleton--line client-skeleton--w80"></div>
                <div class="client-skeleton client-skeleton--line client-skeleton--w60"></div>
              </div>
            </div>
          </div>

          <template v-else>
            <div v-if="masters.length" class="client-grid">
              <MasterCard v-for="master in masters" :key="master.id" :item="master" />
            </div>

            <div v-else class="client-empty">
              <div class="client-empty__icon">🔍</div>
              <h3>Usta topilmadi</h3>
              <p>Filtrlarni o'zgartirib ko'ring</p>
            </div>
          </template>

          <div ref="sentinel" class="client-sentinel"></div>

          <div v-if="loadingMore" class="client-load-more">
            <div class="client-spinner"></div>
          </div>

          <div v-if="!hasMore && masters.length > 0 && !loading" class="client-end">
            Barcha ustalar ko'rsatildi
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
