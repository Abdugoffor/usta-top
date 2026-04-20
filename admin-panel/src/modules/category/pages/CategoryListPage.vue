<script setup>
import { reactive, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { Icon } from '@iconify/vue'
import BaseCard from '@/shared/components/BaseCard.vue'
import CategoryTable from '../components/CategoryTable.vue'
import { useCategoryStore } from '../store/categoryStore'

const router = useRouter()
const categoryStore = useCategoryStore()

const filters = reactive({ name: '', page: 1, sort_by: '', sort_order: 'asc' })

const loadData = async () => {
  await categoryStore.fetchCategories({
    name: filters.name,
    page: filters.page,
    sort_by: filters.sort_by || undefined,
    sort_order: filters.sort_by ? filters.sort_order : undefined,
  })
}

let searchTimer = null
watch(() => filters.name, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    filters.page = 1
    loadData()
  }, 400)
})

const handleSort = ({ field, order }) => {
  filters.sort_by = field
  filters.sort_order = order
  filters.page = 1
  loadData()
}

const goCreate = () => router.push({ name: 'category-create' })

const goEdit = (item) => router.push({ name: 'category-edit', params: { id: item.id } })

const deleteItem = async (item) => {
  const result = await Swal.fire({
    title: "Ochirishni xohlaysizmi?",
    html: `<span style="color:#94a3b8">"${item.name}"</span> o'chiriladi`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Ha, o\'chir',
    cancelButtonText: "Bekor qilish",
    background: '#0f172a',
    color: '#f1f5f9',
    confirmButtonColor: '#ef4444',
    cancelButtonColor: '#1e293b',
    borderRadius: '16px',
  })

  if (!result.isConfirmed) return

  await categoryStore.removeCategory(item.id)
  await loadData()

  Swal.fire({
    icon: 'success',
    title: 'Muvaffaqiyatli o\'chirildi',
    background: '#0f172a',
    color: '#f1f5f9',
    timer: 2000,
    showConfirmButton: false,
  })
}

const changePage = async (page) => {
  if (page < 1 || page > categoryStore.lastPage) return
  filters.page = page
  await loadData()
}

const activeCount   = computed(() => categoryStore.items.filter(i => i.is_active).length)
const inactiveCount = computed(() => categoryStore.items.length - activeCount.value)

const pageNumbers = computed(() => {
  const total = categoryStore.lastPage
  const cur = categoryStore.currentPage
  const pages = []
  const delta = 2
  for (let i = Math.max(1, cur - delta); i <= Math.min(total, cur + delta); i++) {
    pages.push(i)
  }
  return pages
})

onMounted(loadData)
</script>

<template>
  <section class="page-section">

    <!-- Header -->
    <div class="cl-header">
      <div class="cl-header__left">
        <div class="cl-header__icon">
          <Icon icon="mdi:shape-plus" />
        </div>
        <div>
          <h2 class="cl-header__title">Kategoriyalar</h2>
          <p class="cl-header__sub">Barcha kategoriyalarni boshqaring</p>
        </div>
      </div>
      <button class="cl-add-btn" @click="goCreate">
        <Icon icon="mdi:plus" />
        <span>Yangi kategoriya</span>
      </button>
    </div>

    <!-- Stats row -->
    <div class="cl-stats">
      <div class="cl-stat">
        <Icon icon="mdi:layers-triple" class="cl-stat__icon cl-stat__icon--blue" />
        <div>
          <span class="cl-stat__val">{{ categoryStore.total }}</span>
          <span class="cl-stat__label">Jami</span>
        </div>
      </div>
      <div class="cl-stat">
        <Icon icon="mdi:check-circle" class="cl-stat__icon cl-stat__icon--green" />
        <div>
          <span class="cl-stat__val">{{ activeCount }}</span>
          <span class="cl-stat__label">Faol</span>
        </div>
      </div>
      <div class="cl-stat">
        <Icon icon="mdi:close-circle" class="cl-stat__icon cl-stat__icon--red" />
        <div>
          <span class="cl-stat__val">{{ inactiveCount }}</span>
          <span class="cl-stat__label">Nofaol</span>
        </div>
      </div>
    </div>

    <!-- Search card -->
    <BaseCard>
      <div class="cl-search-wrap">
        <div class="cl-search">
          <Icon icon="mdi:magnify" class="cl-search__icon" />
          <input
            v-model="filters.name"
            type="text"
            class="cl-search__input"
            placeholder="Kategoriya nomini qidirish..."
            autocomplete="off"
          />
          <button v-if="filters.name" class="cl-search__clear" @click="filters.name = ''">
            <Icon icon="mdi:close-circle" />
          </button>
        </div>
        <div v-if="filters.name" class="cl-search__hint">
          <Icon icon="mdi:information-outline" />
          "{{ filters.name }}" bo'yicha qidirilmoqda
        </div>
      </div>
    </BaseCard>

    <!-- Table card -->
    <BaseCard>
      <CategoryTable
        :items="categoryStore.items"
        :loading="categoryStore.loading"
        :sort-by="filters.sort_by"
        :sort-order="filters.sort_order"
        @edit="goEdit"
        @delete="deleteItem"
        @sort="handleSort"
      />

      <!-- Pagination -->
      <div class="cl-pagination" v-if="categoryStore.lastPage > 1">
        <span class="cl-pagination__info">
          {{ (categoryStore.currentPage - 1) * categoryStore.perPage + 1 }}–{{ Math.min(categoryStore.currentPage * categoryStore.perPage, categoryStore.total) }} / {{ categoryStore.total }}
        </span>

        <div class="cl-pagination__btns">
          <button
            class="cl-pagination__btn"
            :disabled="filters.page <= 1"
            @click="changePage(filters.page - 1)"
          >
            <Icon icon="mdi:chevron-left" />
          </button>

          <button
            v-if="categoryStore.currentPage > 3"
            class="cl-pagination__btn"
            @click="changePage(1)"
          >1</button>
          <span v-if="categoryStore.currentPage > 4" class="cl-pagination__dots">…</span>

          <button
            v-for="p in pageNumbers"
            :key="p"
            class="cl-pagination__btn"
            :class="{ 'cl-pagination__btn--active': p === categoryStore.currentPage }"
            @click="changePage(p)"
          >{{ p }}</button>

          <span v-if="categoryStore.currentPage < categoryStore.lastPage - 3" class="cl-pagination__dots">…</span>
          <button
            v-if="categoryStore.currentPage < categoryStore.lastPage - 2"
            class="cl-pagination__btn"
            @click="changePage(categoryStore.lastPage)"
          >{{ categoryStore.lastPage }}</button>

          <button
            class="cl-pagination__btn"
            :disabled="filters.page >= categoryStore.lastPage"
            @click="changePage(filters.page + 1)"
          >
            <Icon icon="mdi:chevron-right" />
          </button>
        </div>
      </div>
    </BaseCard>

  </section>
</template>

<style scoped>
/* Header */
.cl-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.cl-header__left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.cl-header__icon {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  background: linear-gradient(135deg, #6366f1, #3b82f6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35);
  flex-shrink: 0;
}

.cl-header__title {
  font-size: 26px;
  font-weight: 800;
  background: linear-gradient(135deg, #f1f5f9, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1.2;
}

.cl-header__sub {
  font-size: 13px;
  color: var(--muted);
  margin-top: 3px;
}

.cl-add-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 11px 22px;
  background: linear-gradient(135deg, #6366f1, #3b82f6);
  color: white;
  border: none;
  border-radius: 14px;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.25s ease;
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.35);
  white-space: nowrap;
}

.cl-add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 28px rgba(99, 102, 241, 0.5);
}

.cl-add-btn:active {
  transform: translateY(0);
}

.cl-add-btn .iconify {
  font-size: 18px;
}

/* Stats */
.cl-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
}

.cl-stat {
  display: flex;
  align-items: center;
  gap: 14px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 16px 20px;
  transition: border-color 0.2s, transform 0.2s;
}

.cl-stat:hover {
  border-color: var(--border-hover);
  transform: translateY(-2px);
}

.cl-stat__icon {
  font-size: 28px;
  flex-shrink: 0;
}

.cl-stat__icon--blue { color: #6366f1; }
.cl-stat__icon--green { color: #10b981; }
.cl-stat__icon--red { color: #ef4444; }

.cl-stat__val {
  display: block;
  font-size: 22px;
  font-weight: 800;
  color: var(--text);
  line-height: 1;
}

.cl-stat__label {
  display: block;
  font-size: 12px;
  color: var(--muted);
  margin-top: 3px;
}

/* Search */
.cl-search-wrap {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.cl-search {
  position: relative;
  display: flex;
  align-items: center;
}

.cl-search__icon {
  position: absolute;
  left: 16px;
  font-size: 20px;
  color: var(--muted);
  pointer-events: none;
  transition: color 0.2s;
}

.cl-search__input {
  width: 100%;
  padding: 13px 48px 13px 48px;
  background: var(--bg-elevated);
  border: 1.5px solid var(--border);
  border-radius: 14px;
  color: var(--text);
  font-size: 15px;
  font-family: inherit;
  outline: none;
  transition: all 0.25s ease;
}

.cl-search__input::placeholder {
  color: var(--muted);
}

.cl-search__input:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.15);
  background: rgba(30, 41, 59, 0.9);
}

.cl-search__input:focus ~ .cl-search__icon,
.cl-search:focus-within .cl-search__icon {
  color: #6366f1;
}

.cl-search__clear {
  position: absolute;
  right: 14px;
  background: none;
  border: none;
  color: var(--muted);
  cursor: pointer;
  font-size: 18px;
  display: flex;
  align-items: center;
  padding: 4px;
  border-radius: 6px;
  transition: color 0.2s;
}

.cl-search__clear:hover {
  color: var(--text);
}

.cl-search__hint {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #6366f1;
  padding-left: 4px;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Pagination */
.cl-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 20px;
  padding-top: 18px;
  border-top: 1px solid var(--border);
}

.cl-pagination__info {
  font-size: 13px;
  color: var(--muted);
}

.cl-pagination__btns {
  display: flex;
  align-items: center;
  gap: 6px;
}

.cl-pagination__btn {
  min-width: 36px;
  height: 36px;
  padding: 0 10px;
  background: var(--bg-elevated);
  color: var(--text-secondary);
  border: 1px solid var(--border);
  border-radius: 10px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  font-family: inherit;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.cl-pagination__btn:hover:not(:disabled) {
  border-color: #6366f1;
  color: #6366f1;
  background: rgba(99, 102, 241, 0.08);
}

.cl-pagination__btn--active {
  background: linear-gradient(135deg, #6366f1, #3b82f6) !important;
  color: white !important;
  border-color: transparent !important;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.35);
}

.cl-pagination__btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.cl-pagination__dots {
  color: var(--muted);
  font-size: 13px;
  padding: 0 4px;
}

/* Responsive */
@media (max-width: 768px) {
  .cl-stats {
    grid-template-columns: 1fr;
  }

  .cl-header__title {
    font-size: 20px;
  }

  .cl-pagination {
    justify-content: center;
  }

  .cl-pagination__info {
    width: 100%;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .cl-add-btn span {
    display: none;
  }

  .cl-add-btn {
    padding: 11px 14px;
  }
}
</style>
