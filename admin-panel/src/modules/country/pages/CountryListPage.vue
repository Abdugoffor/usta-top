<script setup>
import { reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { Icon } from '@iconify/vue'
import BaseCard from '@/shared/components/BaseCard.vue'
import CountryTreeRow from '../components/CountryTreeRow.vue'
import { useCountryStore } from '../store/countryStore'

const router = useRouter()
const store  = useCountryStore()

const filters = reactive({ name: '', page: 1 })

const loadData = () =>
  store.fetchTree({ name: filters.name, page: filters.page })

let timer = null
watch(() => filters.name, () => {
  clearTimeout(timer)
  timer = setTimeout(() => { filters.page = 1; loadData() }, 400)
})

const goCreate = (parentId = null) =>
  router.push({ name: 'country-create', query: parentId ? { parent_id: parentId } : {} })

const goEdit = (item) =>
  router.push({ name: 'country-edit', params: { id: item.id } })

const deleteItem = async (item) => {
  const result = await Swal.fire({
    title: 'O\'chirishni tasdiqlang',
    html: `<span style="color:#94a3b8">"${item.name}"</span>${item.children?.length ? `<br><span style="color:#f87171;font-size:13px">Ichki ${item.children.length} ta ham o'chiriladi!</span>` : ''}`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Ha, o\'chir',
    cancelButtonText: 'Bekor',
    background: '#0f172a',
    color: '#f1f5f9',
    confirmButtonColor: '#ef4444',
    cancelButtonColor: '#1e293b',
  })
  if (!result.isConfirmed) return

  await store.removeCountry(item.id)
  await loadData()

  Swal.fire({
    icon: 'success',
    title: 'O\'chirildi',
    background: '#0f172a',
    color: '#f1f5f9',
    timer: 2000,
    showConfirmButton: false,
  })
}

const changePage = (p) => {
  if (p < 1 || p > store.lastPage) return
  filters.page = p
  loadData()
}

const pageNumbers = computed(() => {
  const total = store.lastPage, cur = store.currentPage
  const pages = []
  for (let i = Math.max(1, cur - 2); i <= Math.min(total, cur + 2); i++) pages.push(i)
  return pages
})

const countNodes = computed(() => {
  const walk = (nodes) => {
    let n = 0
    for (const node of nodes ?? []) { n++; n += walk(node.children) }
    return n
  }
  return walk(store.tree)
})

onMounted(loadData)
</script>

<template>
  <section class="page-section">

    <!-- Header -->
    <div class="cl-header">
      <div class="cl-header__left">
        <div class="cl-header__icon" style="background: linear-gradient(135deg,#10b981,#3b82f6);">
          <Icon icon="mdi:map-marker-multiple" />
        </div>
        <div>
          <h2 class="cl-header__title">Mamlakatlar</h2>
          <p class="cl-header__sub">Rekursiv daraxht ko'rinishida boshqarish</p>
        </div>
      </div>
      <button class="cl-add-btn" @click="goCreate()">
        <Icon icon="mdi:plus" />
        <span>Yangi mamlakat</span>
      </button>
    </div>

    <!-- Stats -->
    <div class="cl-stats">
      <div class="cl-stat">
        <Icon icon="mdi:earth" class="cl-stat__icon" style="color:#3b82f6" />
        <div>
          <span class="cl-stat__val">{{ store.total }}</span>
          <span class="cl-stat__label">Jami (root)</span>
        </div>
      </div>
      <div class="cl-stat">
        <Icon icon="mdi:sitemap" class="cl-stat__icon" style="color:#6366f1" />
        <div>
          <span class="cl-stat__val">{{ countNodes }}</span>
          <span class="cl-stat__label">Ko'rsatilayotgan</span>
        </div>
      </div>
      <div class="cl-stat">
        <Icon icon="mdi:file-tree" class="cl-stat__icon" style="color:#10b981" />
        <div>
          <span class="cl-stat__val">{{ store.tree.length }}</span>
          <span class="cl-stat__label">Root darajasi</span>
        </div>
      </div>
    </div>

    <!-- Search -->
    <BaseCard>
      <div class="cl-search">
        <Icon icon="mdi:magnify" class="cl-search__icon" />
        <input
          v-model="filters.name"
          type="text"
          class="cl-search__input"
          placeholder="Mamlakat nomini qidirish..."
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
    </BaseCard>

    <!-- Tree card -->
    <BaseCard>

      <!-- Table header -->
      <div class="cnt-header">
        <span class="cnt-header__col">Nomi</span>
        <span class="cnt-header__col cnt-header__col--status">Holat</span>
        <span class="cnt-header__col cnt-header__col--actions">Amallar</span>
      </div>

      <!-- Skeleton -->
      <template v-if="store.loading">
        <div v-for="i in 5" :key="i" class="cnt-skeleton">
          <div class="cnt-skeleton__cell cnt-skeleton__cell--lg" />
          <div class="cnt-skeleton__cell cnt-skeleton__cell--sm" />
          <div class="cnt-skeleton__cell cnt-skeleton__cell--sm" />
        </div>
      </template>

      <!-- Empty -->
      <div v-else-if="!store.tree.length" class="cnt-empty">
        <Icon icon="mdi:earth-off" />
        <p>Mamlakat topilmadi</p>
        <span>Qidiruvni o'zgartiring yoki yangi qo'shing</span>
      </div>

      <!-- Tree rows -->
      <template v-else>
        <CountryTreeRow
          v-for="node in store.tree"
          :key="node.id"
          :node="node"
          :depth="0"
          @edit="goEdit"
          @delete="deleteItem"
          @add-child="n => goCreate(n.id)"
        />
      </template>

      <!-- Pagination -->
      <div class="cl-pagination" v-if="store.lastPage > 1">
        <span class="cl-pagination__info">
          {{ (store.currentPage - 1) * store.perPage + 1 }}–{{ Math.min(store.currentPage * store.perPage, store.total) }} / {{ store.total }}
        </span>
        <div class="cl-pagination__btns">
          <button class="cl-pagination__btn" :disabled="filters.page <= 1" @click="changePage(filters.page - 1)">
            <Icon icon="mdi:chevron-left" />
          </button>
          <button v-if="store.currentPage > 3" class="cl-pagination__btn" @click="changePage(1)">1</button>
          <span v-if="store.currentPage > 4" class="cl-pagination__dots">…</span>
          <button
            v-for="p in pageNumbers" :key="p"
            class="cl-pagination__btn"
            :class="{ 'cl-pagination__btn--active': p === store.currentPage }"
            @click="changePage(p)"
          >{{ p }}</button>
          <span v-if="store.currentPage < store.lastPage - 3" class="cl-pagination__dots">…</span>
          <button v-if="store.currentPage < store.lastPage - 2" class="cl-pagination__btn" @click="changePage(store.lastPage)">{{ store.lastPage }}</button>
          <button class="cl-pagination__btn" :disabled="filters.page >= store.lastPage" @click="changePage(filters.page + 1)">
            <Icon icon="mdi:chevron-right" />
          </button>
        </div>
      </div>

    </BaseCard>
  </section>
</template>

<style scoped>
/* Reuse category page styles via global, add page-specific below */
.cl-header { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:16px; }
.cl-header__left { display:flex; align-items:center; gap:16px; }
.cl-header__icon { width:52px; height:52px; border-radius:16px; display:flex; align-items:center; justify-content:center; font-size:24px; color:white; box-shadow:0 8px 24px rgba(16,185,129,.3); flex-shrink:0; }
.cl-header__title { font-size:26px; font-weight:800; background:linear-gradient(135deg,#f1f5f9,#94a3b8); -webkit-background-clip:text; -webkit-text-fill-color:transparent; background-clip:text; }
.cl-header__sub { font-size:13px; color:var(--muted); margin-top:3px; }

.cl-add-btn { display:inline-flex; align-items:center; gap:8px; padding:11px 22px; background:linear-gradient(135deg,#10b981,#3b82f6); color:white; border:none; border-radius:14px; font-size:14px; font-weight:700; cursor:pointer; transition:all .25s; box-shadow:0 6px 20px rgba(16,185,129,.3); white-space:nowrap; font-family:inherit; }
.cl-add-btn:hover { transform:translateY(-2px); box-shadow:0 10px 28px rgba(16,185,129,.45); }

.cl-stats { display:grid; grid-template-columns:repeat(3,1fr); gap:14px; }
.cl-stat { display:flex; align-items:center; gap:14px; background:var(--bg-card); border:1px solid var(--border); border-radius:16px; padding:16px 20px; transition:border-color .2s, transform .2s; }
.cl-stat:hover { border-color:var(--border-hover); transform:translateY(-2px); }
.cl-stat__icon { font-size:28px; flex-shrink:0; }
.cl-stat__val { display:block; font-size:22px; font-weight:800; color:var(--text); line-height:1; }
.cl-stat__label { display:block; font-size:12px; color:var(--muted); margin-top:3px; }

.cl-search { position:relative; display:flex; align-items:center; }
.cl-search__icon { position:absolute; left:16px; font-size:20px; color:var(--muted); pointer-events:none; transition:color .2s; }
.cl-search:focus-within .cl-search__icon { color:#6366f1; }
.cl-search__input { width:100%; padding:13px 48px; background:var(--bg-elevated); border:1.5px solid var(--border); border-radius:14px; color:var(--text); font-size:15px; font-family:inherit; outline:none; transition:all .25s; }
.cl-search__input::placeholder { color:var(--muted); }
.cl-search__input:focus { border-color:#6366f1; box-shadow:0 0 0 4px rgba(99,102,241,.15); }
.cl-search__clear { position:absolute; right:14px; background:none; border:none; color:var(--muted); cursor:pointer; font-size:18px; display:flex; align-items:center; padding:4px; border-radius:6px; }
.cl-search__clear:hover { color:var(--text); }
.cl-search__hint { display:flex; align-items:center; gap:6px; font-size:12px; color:#6366f1; padding:8px 4px 0; animation:fadeIn .2s ease; }
@keyframes fadeIn { from{opacity:0;transform:translateY(-4px)} to{opacity:1;transform:translateY(0)} }

/* Tree table header */
.cnt-header { display:grid; grid-template-columns:1fr auto auto; align-items:center; gap:16px; padding:10px 16px; border-bottom:1px solid var(--border); background:rgba(255,255,255,.02); }
.cnt-header__col { font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:.07em; color:var(--muted); }
.cnt-header__col--status { width:90px; text-align:center; }
.cnt-header__col--actions { width:110px; text-align:center; }

/* Skeleton */
.cnt-skeleton { display:grid; grid-template-columns:1fr 90px 110px; gap:16px; padding:14px 16px; border-bottom:1px solid var(--border); align-items:center; }
.cnt-skeleton__cell { height:14px; border-radius:8px; background:linear-gradient(90deg,var(--bg-elevated) 25%,rgba(255,255,255,.06) 50%,var(--bg-elevated) 75%); background-size:200% 100%; animation:shimmer 1.4s infinite; }
.cnt-skeleton__cell--lg { width:100%; }
.cnt-skeleton__cell--sm { width:70%; }
@keyframes shimmer { 0%{background-position:200% 0} 100%{background-position:-200% 0} }

/* Empty */
.cnt-empty { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:60px 24px; gap:10px; text-align:center; color:var(--muted); font-size:32px; }
.cnt-empty p { font-size:16px; font-weight:700; color:var(--text); }
.cnt-empty span { font-size:13px; color:var(--muted); }

/* Pagination */
.cl-pagination { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:12px; margin-top:20px; padding-top:18px; border-top:1px solid var(--border); }
.cl-pagination__info { font-size:13px; color:var(--muted); }
.cl-pagination__btns { display:flex; align-items:center; gap:6px; }
.cl-pagination__btn { min-width:36px; height:36px; padding:0 10px; background:var(--bg-elevated); color:var(--text-secondary); border:1px solid var(--border); border-radius:10px; cursor:pointer; font-size:13px; font-weight:600; font-family:inherit; display:flex; align-items:center; justify-content:center; transition:all .2s; }
.cl-pagination__btn:hover:not(:disabled) { border-color:#6366f1; color:#6366f1; background:rgba(99,102,241,.08); }
.cl-pagination__btn--active { background:linear-gradient(135deg,#6366f1,#3b82f6) !important; color:white !important; border-color:transparent !important; box-shadow:0 4px 12px rgba(99,102,241,.35); }
.cl-pagination__btn:disabled { opacity:.35; cursor:not-allowed; }
.cl-pagination__dots { color:var(--muted); font-size:13px; padding:0 4px; }

@media (max-width:768px) {
  .cl-stats { grid-template-columns:1fr; }
  .cl-header__title { font-size:20px; }
  .cl-add-btn span { display:none; }
  .cl-add-btn { padding:11px 14px; }
}
</style>
