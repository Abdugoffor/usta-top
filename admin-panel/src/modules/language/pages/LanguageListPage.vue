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

const filters = reactive({ name: '', is_active: '', sort_by: '', sort_order: 'asc' })

const buildParams = () => ({
  name:      filters.name      || undefined,
  is_active: filters.is_active !== '' ? filters.is_active : undefined,
})

const reload = async () => {
  store.resetCursor()
  await store.fetchLanguages(buildParams())
}

let searchTimer = null
watch(() => filters.name, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(reload, 400)
})
watch(() => filters.is_active, () => reload())

const handleSort = ({ field, order }) => {
  filters.sort_by    = field
  filters.sort_order = order
}

const sortedItems = computed(() => {
  if (!filters.sort_by) return store.items
  return [...store.items].sort((a, b) => {
    let av = a[filters.sort_by]
    let bv = b[filters.sort_by]
    if (typeof av === 'string') av = av.toLowerCase()
    if (typeof bv === 'string') bv = bv.toLowerCase()
    if (av < bv) return filters.sort_order === 'asc' ? -1 : 1
    if (av > bv) return filters.sort_order === 'asc' ? 1 : -1
    return 0
  })
})

const goCreate = () => router.push({ name: 'language-create' })
const goShow   = (item) => router.push({ name: 'language-show', params: { id: item.id } })
const goEdit   = (item) => router.push({ name: 'language-edit', params: { id: item.id } })

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
  await reload()

  Swal.fire({ icon: 'success', title: "Muvaffaqiyatli o'chirildi", background: '#0f172a', color: '#f1f5f9', timer: 2000, showConfirmButton: false })
}

const goPrev = () => store.goPrev(buildParams())
const goNext = () => store.goNext(buildParams())

onMounted(reload)
</script>

<template>
  <section class="page-section">

    <div class="ll-header">
      <div class="ll-header__left">
        <div>
          <h2 class="ll-header__title">Tillar</h2>
          <p class="ll-header__sub">Barcha tillarni boshqaring</p>
        </div>
      </div>
      <button class="cl-add-btn" @click="goCreate">
        <span>Qo'shish</span>
      </button>
    </div>

    <BaseCard>
      <div class="ll-filters">
        <div class="ll-search">
          <Icon icon="mdi:magnify" class="ll-search__icon" />
          <input v-model="filters.name" type="text" class="ll-search__input"
            placeholder="Til nomini qidirish..." autocomplete="off" />
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

    <BaseCard>
      <LanguageTable
        :items="sortedItems"
        :loading="store.loading"
        :sort-by="filters.sort_by"
        :sort-order="filters.sort_order"
        @show="goShow"
        @edit="goEdit"
        @delete="deleteItem"
        @sort="handleSort"
      />

      <div class="ll-pagination" v-if="store.hasPrev || store.hasNext">
        <span class="ll-pagination__info">Sahifa {{ store.currentPage }}</span>
        <div class="ll-pagination__btns">
          <button class="ll-pagination__btn" :disabled="!store.hasPrev || store.loading" @click="goPrev">
            <Icon icon="mdi:chevron-left" /> Oldingi
          </button>
          <button class="ll-pagination__btn ll-pagination__btn--active" disabled>{{ store.currentPage }}</button>
          <button class="ll-pagination__btn" :disabled="!store.hasNext || store.loading" @click="goNext">
            Keyingi <Icon icon="mdi:chevron-right" />
          </button>
        </div>
      </div>
    </BaseCard>

  </section>
</template>

<style scoped>
.ll-header { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:16px; }
.ll-header__left { display:flex; align-items:center; gap:16px; }
.ll-header__title { font-size:26px; font-weight:800; background:linear-gradient(135deg,#f1f5f9,#94a3b8); -webkit-background-clip:text; -webkit-text-fill-color:transparent; background-clip:text; }
.ll-header__sub { font-size:13px; color:var(--muted); margin-top:3px; }

.cl-add-btn { display:inline-flex; align-items:center; gap:8px; padding:11px 22px; background:linear-gradient(135deg,#6366f1,#3b82f6); color:white; border:none; border-radius:14px; font-size:14px; font-weight:700; cursor:pointer; transition:all .25s; box-shadow:0 6px 20px rgba(99,102,241,.35); white-space:nowrap; }
.cl-add-btn:hover { transform:translateY(-2px); box-shadow:0 10px 28px rgba(99,102,241,.5); }
.cl-add-btn:active { transform:translateY(0); }

.ll-filters { display:flex; gap:12px; flex-wrap:wrap; }

.ll-search { position:relative; display:flex; align-items:center; flex:1; min-width:200px; }
.ll-search__icon { position:absolute; left:16px; font-size:20px; color:var(--muted); pointer-events:none; transition:color .2s; }
.ll-search:focus-within .ll-search__icon { color:#6366f1; }
.ll-search__input { width:100%; padding:12px 48px; background:var(--bg-elevated); border:1.5px solid var(--border); border-radius:14px; color:var(--text); font-size:15px; font-family:inherit; outline:none; transition:all .25s; }
.ll-search__input::placeholder { color:var(--muted); }
.ll-search__input:focus { border-color:#6366f1; box-shadow:0 0 0 4px rgba(99,102,241,.15); }
.ll-search__clear { position:absolute; right:14px; background:none; border:none; color:var(--muted); cursor:pointer; font-size:18px; display:flex; align-items:center; padding:4px; border-radius:6px; }
.ll-search__clear:hover { color:var(--text); }
.ll-search__hint { display:flex; align-items:center; gap:6px; font-size:12px; color:#6366f1; padding:8px 4px 0; animation:fadeIn .2s ease; }
@keyframes fadeIn { from{opacity:0;transform:translateY(-4px)} to{opacity:1;transform:translateY(0)} }

.ll-select { padding:12px 16px; background:var(--bg-elevated); border:1.5px solid var(--border); border-radius:14px; color:var(--text); font-size:14px; font-family:inherit; outline:none; cursor:pointer; transition:border-color .2s; min-width:150px; }
.ll-select:focus { border-color:#6366f1; }
.ll-select option { background:var(--bg-elevated); }

.ll-pagination { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:12px; margin-top:20px; padding-top:18px; border-top:1px solid var(--border); }
.ll-pagination__info { font-size:13px; color:var(--muted); }
.ll-pagination__btns { display:flex; align-items:center; gap:8px; }
.ll-pagination__btn { display:inline-flex; align-items:center; gap:6px; min-width:36px; height:36px; padding:0 14px; background:var(--bg-elevated); color:var(--text-secondary); border:1px solid var(--border); border-radius:10px; cursor:pointer; font-size:13px; font-weight:600; font-family:inherit; transition:all .2s; }
.ll-pagination__btn:hover:not(:disabled) { border-color:#6366f1; color:#6366f1; background:rgba(99,102,241,.08); }
.ll-pagination__btn--active { background:linear-gradient(135deg,#6366f1,#3b82f6) !important; color:white !important; border-color:transparent !important; box-shadow:0 4px 12px rgba(99,102,241,.35); min-width:40px; justify-content:center; }
.ll-pagination__btn:disabled { opacity:.35; cursor:not-allowed; }

@media (max-width:480px) { .cl-add-btn span { display:none; } .cl-add-btn { padding:11px 14px; } }
</style>
