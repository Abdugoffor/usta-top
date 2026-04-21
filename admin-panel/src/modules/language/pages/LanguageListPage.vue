<script setup>
import { reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { Icon } from '@iconify/vue'
import BaseCard from '@/shared/components/BaseCard.vue'
import LanguageTable from '../components/LanguageTable.vue'
import { useLanguageStore } from '../store/languageStore'

const router = useRouter()
const store = useLanguageStore()

const filters = reactive({
  name: '',
  is_active: '',
  page: 1,
  sort_by: '',
  sort_order: 'asc',
})

const loadData = () => store.fetchLanguages({
  name: filters.name || undefined,
  is_active: filters.is_active !== '' ? filters.is_active : undefined,
  page: filters.page,
  sort_by: filters.sort_by || undefined,
  sort_order: filters.sort_by ? filters.sort_order : undefined,
})

let searchTimer = null
watch(() => filters.name, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { filters.page = 1; loadData() }, 400)
})

watch(() => filters.is_active, () => { filters.page = 1; loadData() })

const handleSort = ({ field, order }) => {
  filters.sort_by = field
  filters.sort_order = order
  filters.page = 1
  loadData()
}

const goCreate = () => router.push({ name: 'language-create' })
const goEdit = (item) => router.push({ name: 'language-edit', params: { id: item.id } })

const deleteItem = async (item) => {
  const result = await Swal.fire({
    title: "O'chirishni xohlaysizmi?",
    html: `<span style="color:#94a3b8">"${item.name}"</span> o'chiriladi`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: "Ha, o'chir",
    cancelButtonText: 'Bekor qilish',
    background: '#0f172a',
    color: '#f1f5f9',
    confirmButtonColor: '#ef4444',
    cancelButtonColor: '#1e293b',
  })
  if (!result.isConfirmed) return

  await store.removeLanguage(item.id)
  await loadData()

  Swal.fire({ icon: 'success', title: "Muvaffaqiyatli o'chirildi", background: '#0f172a', color: '#f1f5f9', timer: 2000, showConfirmButton: false })
}

const changePage = (page) => {
  if (page < 1 || page > store.lastPage) return
  filters.page = page
  loadData()
}

const pageNumbers = computed(() => {
  const total = store.lastPage, cur = store.currentPage
  const pages = []
  for (let i = Math.max(1, cur - 2); i <= Math.min(total, cur + 2); i++) pages.push(i)
  return pages
})

onMounted(loadData)
</script>

<template>
  <section class="page-section">

    <!-- Header -->
    <div class="ll-header">
      <div class="ll-header__left">
        <div>
          <h2 class="ll-header__title">Tillar</h2>
          <p class="ll-header__sub">Barcha tillarni boshqaring</p>
        </div>
      </div>
      <button class="cl-add-btn" @click="goCreate">
        <span>Добавить</span>
      </button>
    </div>

    <!-- Filters -->
    <BaseCard>
      <div class="ll-filters">
        <div class="ll-search">
          <Icon icon="mdi:magnify" class="ll-search__icon" />
          <input v-model="filters.name" type="text" class="ll-search__input" placeholder="Til nomini qidirish..."
            autocomplete="off" />
          <button v-if="filters.name" class="ll-search__clear" @click="filters.name = ''">
            <Icon icon="mdi:close-circle" />
          </button>
        </div>

        <select v-model="filters.is_active" class="ll-select">
          <option value="">Barcha holat</option>
          <option value="true">Faol</option>
          <option value="false">Nofaol</option>
        </select>
      </div>
      <div v-if="filters.name" class="ll-search__hint">
        <Icon icon="mdi:information-outline" />
        "{{ filters.name }}" bo'yicha qidirilmoqda
      </div>
    </BaseCard>

    <!-- Table -->
    <BaseCard>
      <LanguageTable :items="store.items" :loading="store.loading" :sort-by="filters.sort_by"
        :sort-order="filters.sort_order" @edit="goEdit" @delete="deleteItem" @sort="handleSort" />

      <!-- Pagination -->
      <div class="ll-pagination" v-if="store.lastPage > 1">
        <span class="ll-pagination__info">
          {{ (store.currentPage - 1) * store.perPage + 1 }}–{{ Math.min(store.currentPage * store.perPage, store.total)
          }} / {{ store.total }}
        </span>
        <div class="ll-pagination__btns">
          <button class="ll-pagination__btn" :disabled="filters.page <= 1" @click="changePage(filters.page - 1)">
            <Icon icon="mdi:chevron-left" />
          </button>
          <button v-if="store.currentPage > 3" class="ll-pagination__btn" @click="changePage(1)">1</button>
          <span v-if="store.currentPage > 4" class="ll-pagination__dots">…</span>
          <button v-for="p in pageNumbers" :key="p" class="ll-pagination__btn"
            :class="{ 'll-pagination__btn--active': p === store.currentPage }" @click="changePage(p)">{{ p }}</button>
          <span v-if="store.currentPage < store.lastPage - 3" class="ll-pagination__dots">…</span>
          <button v-if="store.currentPage < store.lastPage - 2" class="ll-pagination__btn"
            @click="changePage(store.lastPage)">{{ store.lastPage }}</button>
          <button class="ll-pagination__btn" :disabled="filters.page >= store.lastPage"
            @click="changePage(filters.page + 1)">
            <Icon icon="mdi:chevron-right" />
          </button>
        </div>
      </div>
    </BaseCard>

  </section>
</template>

<style scoped>
.ll-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.ll-header__left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.ll-header__icon {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
  box-shadow: 0 8px 24px rgba(139, 92, 246, .35);
  flex-shrink: 0;
}

.ll-header__title {
  font-size: 26px;
  font-weight: 800;
  background: linear-gradient(135deg, #f1f5f9, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.ll-header__sub {
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

.ll-filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.ll-search {
  position: relative;
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 200px;
}

.ll-search__icon {
  position: absolute;
  left: 16px;
  font-size: 20px;
  color: var(--muted);
  pointer-events: none;
  transition: color .2s;
}

.ll-search:focus-within .ll-search__icon {
  color: #6366f1;
}

.ll-search__input {
  width: 100%;
  padding: 12px 48px;
  background: var(--bg-elevated);
  border: 1.5px solid var(--border);
  border-radius: 14px;
  color: var(--text);
  font-size: 15px;
  font-family: inherit;
  outline: none;
  transition: all .25s;
}

.ll-search__input::placeholder {
  color: var(--muted);
}

.ll-search__input:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, .15);
}

.ll-search__clear {
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
}

.ll-search__clear:hover {
  color: var(--text);
}

.ll-search__hint {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #6366f1;
  padding: 8px 4px 0;
  animation: fadeIn .2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-4px)
  }

  to {
    opacity: 1;
    transform: translateY(0)
  }
}

.ll-select {
  padding: 12px 16px;
  background: var(--bg-elevated);
  border: 1.5px solid var(--border);
  border-radius: 14px;
  color: var(--text);
  font-size: 14px;
  font-family: inherit;
  outline: none;
  cursor: pointer;
  transition: border-color .2s;
  min-width: 150px;
}

.ll-select:focus {
  border-color: #6366f1;
}

.ll-select option {
  background: var(--bg-elevated);
}

.ll-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 20px;
  padding-top: 18px;
  border-top: 1px solid var(--border);
}

.ll-pagination__info {
  font-size: 13px;
  color: var(--muted);
}

.ll-pagination__btns {
  display: flex;
  align-items: center;
  gap: 6px;
}

.ll-pagination__btn {
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
  transition: all .2s;
}

.ll-pagination__btn:hover:not(:disabled) {
  border-color: #6366f1;
  color: #6366f1;
  background: rgba(99, 102, 241, .08);
}

.ll-pagination__btn--active {
  background: linear-gradient(135deg, #6366f1, #3b82f6) !important;
  color: white !important;
  border-color: transparent !important;
  box-shadow: 0 4px 12px rgba(99, 102, 241, .35);
}

.ll-pagination__btn:disabled {
  opacity: .35;
  cursor: not-allowed;
}

.ll-pagination__dots {
  color: var(--muted);
  font-size: 13px;
  padding: 0 4px;
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
