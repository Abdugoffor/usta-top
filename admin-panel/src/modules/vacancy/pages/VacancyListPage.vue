<script setup>
import { reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { Icon } from '@iconify/vue'
import BaseCard from '@/shared/components/BaseCard.vue'
import VacancyTable from '../components/VacancyTable.vue'
import { useVacancyStore } from '../store/vacancyStore'

const router = useRouter()
const store = useVacancyStore()

const filters = reactive({
  search: '',
  is_active: '',
  sort_by: 'id',
  sort_order: 'desc',
})

const buildParams = () => ({
  search:     filters.search    || undefined,
  is_active:  filters.is_active !== '' ? filters.is_active : undefined,
  sort_by:    filters.sort_by   || undefined,
  sort_order: filters.sort_by   ? filters.sort_order : undefined,
})

const reload = async () => {
  store.resetCursor()
  await store.fetchVacancies(buildParams())
}

let searchTimer = null
watch(() => filters.search, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(reload, 400)
})

watch(() => filters.is_active, reload)

const handleSort = ({ field, order }) => {
  filters.sort_by = field
  filters.sort_order = order
  reload()
}

const goShow = (item) => router.push({ name: 'vacancy-admin-show', params: { id: item.slug || item.id } })

const deleteItem = async (item) => {
  const result = await Swal.fire({
    title: "O'chirishni xohlaysizmi?",
    html: `<span style="color:#94a3b8">"${item.name}"</span> vakansiyasi o'chiriladi`,
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

  await store.removeVacancy(item.id)
  await reload()

  Swal.fire({ icon: 'success', title: "O'chirildi", background: '#0f172a', color: '#f1f5f9', timer: 2000, showConfirmButton: false })
}

const goPrev = () => store.goPrev(buildParams())
const goNext = () => store.goNext(buildParams())

onMounted(reload)
</script>

<template>
  <section class="page-section">

    <div class="cl-header">
      <div class="cl-header__left">
        <div>
          <h2 class="cl-header__title">Vakansiyalar</h2>
          <p class="cl-header__sub">Barcha e'lon qilingan vakansiyalarni boshqaring</p>
        </div>
      </div>
    </div>

    <BaseCard>
      <div class="cl-filters-row">
        <div class="cl-search">
          <Icon icon="mdi:magnify" class="cl-search__icon" />
          <input v-model="filters.search" type="text" class="cl-search__input"
            placeholder="Lavozim, aloqa yoki manzil bo'yicha qidirish..." autocomplete="off" />
          <button v-if="filters.search" class="cl-search__clear" @click="filters.search = ''">
            <Icon icon="mdi:close-circle" />
          </button>
        </div>
        <select v-model="filters.is_active" class="cl-select">
          <option value="">Barcha holat</option>
          <option value="true">Faol</option>
          <option value="false">Nofaol</option>
        </select>
      </div>
      <div v-if="filters.search" class="cl-search__hint">
        <Icon icon="mdi:information-outline" />
        "{{ filters.search }}" bo'yicha qidirilmoqda
      </div>
    </BaseCard>

    <BaseCard>
      <VacancyTable
        :items="store.items"
        :loading="store.loading"
        :sort-by="filters.sort_by"
        :sort-order="filters.sort_order"
        @show="goShow"
        @delete="deleteItem"
        @sort="handleSort"
      />

      <div class="cl-pagination" v-if="store.hasPrev || store.hasNext">
        <span class="cl-pagination__info">
          Jami: <strong>{{ store.total }}</strong> ta
          <span v-if="store.total > 0"> · Sahifa {{ store.currentPage }}</span>
        </span>
        <div class="cl-pagination__btns">
          <button class="cl-pagination__btn" :disabled="!store.hasPrev || store.loading" @click="goPrev">
            <Icon icon="mdi:chevron-left" />
            Oldingi
          </button>
          <button class="cl-pagination__btn cl-pagination__btn--page" disabled>
            {{ store.currentPage }}
          </button>
          <button class="cl-pagination__btn" :disabled="!store.hasNext || store.loading" @click="goNext">
            Keyingi
            <Icon icon="mdi:chevron-right" />
          </button>
        </div>
      </div>
      <div v-else-if="!store.loading && store.total > 0" class="cl-pagination__info" style="margin-top:16px;padding-top:16px;border-top:1px solid var(--border)">
        Jami: <strong>{{ store.total }}</strong> ta vakansiya
      </div>
    </BaseCard>

  </section>
</template>

<style scoped>
.cl-header { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:16px; }
.cl-header__left { display:flex; align-items:center; gap:16px; }
.cl-header__title { font-size:26px; font-weight:800; background:linear-gradient(135deg,#f1f5f9,#94a3b8); -webkit-background-clip:text; -webkit-text-fill-color:transparent; background-clip:text; line-height:1.2; }
.cl-header__sub { font-size:13px; color:var(--muted); margin-top:3px; }

.cl-filters-row { display:flex; gap:12px; flex-wrap:wrap; }
.cl-search { position:relative; display:flex; align-items:center; flex:1; min-width:200px; }
.cl-search__icon { position:absolute; left:16px; font-size:20px; color:var(--muted); pointer-events:none; transition:color .2s; }
.cl-search:focus-within .cl-search__icon { color:#6366f1; }
.cl-search__input { width:100%; padding:12px 48px; background:var(--bg-elevated); border:1.5px solid var(--border); border-radius:14px; color:var(--text); font-size:15px; font-family:inherit; outline:none; transition:all .25s; }
.cl-search__input::placeholder { color:var(--muted); }
.cl-search__input:focus { border-color:#6366f1; box-shadow:0 0 0 4px rgba(99,102,241,.15); }
.cl-search__clear { position:absolute; right:14px; background:none; border:none; color:var(--muted); cursor:pointer; font-size:18px; display:flex; align-items:center; padding:4px; border-radius:6px; }
.cl-search__clear:hover { color:var(--text); }
.cl-search__hint { display:flex; align-items:center; gap:6px; font-size:12px; color:#6366f1; padding:8px 4px 0; }
.cl-select { padding:12px 16px; background:var(--bg-elevated); border:1.5px solid var(--border); border-radius:14px; color:var(--text); font-size:14px; font-family:inherit; outline:none; cursor:pointer; transition:border-color .2s; min-width:150px; }
.cl-select:focus { border-color:#6366f1; }
.cl-select option { background:var(--bg-elevated); }

.cl-pagination { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:12px; margin-top:20px; padding-top:18px; border-top:1px solid var(--border); }
.cl-pagination__info { font-size:13px; color:var(--muted); }
.cl-pagination__info strong { color:var(--text); }
.cl-pagination__btns { display:flex; align-items:center; gap:8px; }
.cl-pagination__btn { display:inline-flex; align-items:center; gap:6px; height:36px; padding:0 14px; background:var(--bg-elevated); color:var(--text-secondary); border:1px solid var(--border); border-radius:10px; cursor:pointer; font-size:13px; font-weight:600; font-family:inherit; transition:all .2s; }
.cl-pagination__btn:hover:not(:disabled) { border-color:#6366f1; color:#6366f1; background:rgba(99,102,241,.08); }
.cl-pagination__btn--page { background:linear-gradient(135deg,#6366f1,#3b82f6) !important; color:white !important; border-color:transparent !important; box-shadow:0 4px 12px rgba(99,102,241,.35); min-width:40px; justify-content:center; }
.cl-pagination__btn:disabled { opacity:.35; cursor:not-allowed; }
</style>
