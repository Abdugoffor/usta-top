<script setup>
import { reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { Icon } from '@iconify/vue'
import BaseCard from '@/shared/components/BaseCard.vue'
import TranslationTable from '../components/TranslationTable.vue'
import { useTranslationStore } from '../store/translationStore'

const router = useRouter()
const store = useTranslationStore()

const filters = reactive({ slug: '', name: '', is_active: '' })

const buildParams = () => ({
  slug:      filters.slug      || undefined,
  name:      filters.name      || undefined,
  is_active: filters.is_active !== '' ? filters.is_active : undefined,
})

const reload = async () => {
  store.resetCursor()
  await store.fetchTranslations(buildParams())
}

let timer = null
watch(() => [filters.slug, filters.name], () => {
  clearTimeout(timer)
  timer = setTimeout(reload, 400)
})
watch(() => filters.is_active, () => reload())

const goCreate = () => router.push({ name: 'translation-create' })
const goShow   = (item) => router.push({ name: 'translation-show', params: { id: item.id } })
const goEdit   = (item) => router.push({ name: 'translation-edit', params: { id: item.id } })

const deleteItem = async (item) => {
  const result = await Swal.fire({
    title: "O'chirishni xohlaysizmi?",
    html: `<code style="color:#818cf8">"${item.slug}"</code> tarjimasi o'chiriladi`,
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

  await store.removeTranslation(item.id)
  await reload()

  Swal.fire({ icon: 'success', title: "O'chirildi", background: '#0f172a', color: '#f1f5f9', timer: 1800, showConfirmButton: false })
}

const goPrev = () => store.goPrev(buildParams())
const goNext = () => store.goNext(buildParams())

onMounted(reload)
</script>

<template>
  <section class="page-section">

    <div class="tl-header">
      <div>
        <h2 class="tl-header__title">Tarjimalar</h2>
        <p class="tl-header__sub">Sayt matnlari va labellarini boshqaring</p>
      </div>
      <button class="tl-add-btn" @click="goCreate">
        <Icon icon="mdi:plus" />
        <span>Qo'shish</span>
      </button>
    </div>

    <BaseCard>
      <div class="tl-filters">
        <div class="tl-search">
          <Icon icon="mdi:key-outline" class="tl-search__icon" />
          <input v-model="filters.slug" type="text" class="tl-search__input"
            placeholder="Kalit (slug) bo'yicha..." autocomplete="off" />
          <button v-if="filters.slug" class="tl-search__clear" @click="filters.slug = ''">
            <Icon icon="mdi:close-circle" />
          </button>
        </div>
        <div class="tl-search">
          <Icon icon="mdi:text-search" class="tl-search__icon" />
          <input v-model="filters.name" type="text" class="tl-search__input"
            placeholder="Matn bo'yicha..." autocomplete="off" />
          <button v-if="filters.name" class="tl-search__clear" @click="filters.name = ''">
            <Icon icon="mdi:close-circle" />
          </button>
        </div>
        <select v-model="filters.is_active" class="tl-select">
          <option value="">Barcha holat</option>
          <option value="true">Faol</option>
          <option value="false">Nofaol</option>
        </select>
      </div>
    </BaseCard>

    <BaseCard>
      <TranslationTable
        :items="store.items"
        :loading="store.loading"
        @show="goShow"
        @edit="goEdit"
        @delete="deleteItem"
      />

      <div class="tl-pagination" v-if="store.hasPrev || store.hasNext">
        <span class="tl-pagination__info">Sahifa {{ store.currentPage }}</span>
        <div class="tl-pagination__btns">
          <button class="tl-pagination__btn" :disabled="!store.hasPrev || store.loading" @click="goPrev">
            <Icon icon="mdi:chevron-left" /> Oldingi
          </button>
          <button class="tl-pagination__btn tl-pagination__btn--active" disabled>{{ store.currentPage }}</button>
          <button class="tl-pagination__btn" :disabled="!store.hasNext || store.loading" @click="goNext">
            Keyingi <Icon icon="mdi:chevron-right" />
          </button>
        </div>
      </div>
    </BaseCard>

  </section>
</template>

<style scoped>
.tl-header { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:16px; }
.tl-header__title { font-size:26px; font-weight:800; background:linear-gradient(135deg,#f1f5f9,#94a3b8); -webkit-background-clip:text; -webkit-text-fill-color:transparent; background-clip:text; }
.tl-header__sub { font-size:13px; color:var(--muted); margin-top:3px; }

.tl-add-btn { display:inline-flex; align-items:center; gap:8px; padding:11px 22px; background:linear-gradient(135deg,#6366f1,#3b82f6); color:white; border:none; border-radius:14px; font-size:14px; font-weight:700; cursor:pointer; transition:all .25s; box-shadow:0 6px 20px rgba(99,102,241,.35); }
.tl-add-btn:hover { transform:translateY(-2px); box-shadow:0 10px 28px rgba(99,102,241,.5); }

.tl-filters { display:flex; gap:12px; flex-wrap:wrap; }

.tl-search { position:relative; display:flex; align-items:center; flex:1; min-width:180px; }
.tl-search__icon { position:absolute; left:14px; font-size:18px; color:var(--muted); pointer-events:none; }
.tl-search:focus-within .tl-search__icon { color:#6366f1; }
.tl-search__input { width:100%; padding:11px 40px 11px 40px; background:var(--bg-elevated); border:1.5px solid var(--border); border-radius:12px; color:var(--text); font-size:14px; font-family:inherit; outline:none; transition:all .25s; }
.tl-search__input:focus { border-color:#6366f1; box-shadow:0 0 0 4px rgba(99,102,241,.12); }
.tl-search__clear { position:absolute; right:12px; background:none; border:none; color:var(--muted); cursor:pointer; font-size:17px; display:flex; align-items:center; }
.tl-search__clear:hover { color:var(--text); }

.tl-select { padding:11px 14px; background:var(--bg-elevated); border:1.5px solid var(--border); border-radius:12px; color:var(--text); font-size:14px; font-family:inherit; outline:none; cursor:pointer; transition:border-color .2s; min-width:140px; }
.tl-select:focus { border-color:#6366f1; }

.tl-pagination { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:12px; margin-top:20px; padding-top:18px; border-top:1px solid var(--border); }
.tl-pagination__info { font-size:13px; color:var(--muted); }
.tl-pagination__btns { display:flex; align-items:center; gap:8px; }
.tl-pagination__btn { display:inline-flex; align-items:center; gap:6px; min-width:36px; height:36px; padding:0 14px; background:var(--bg-elevated); color:var(--text-secondary); border:1px solid var(--border); border-radius:10px; cursor:pointer; font-size:13px; font-weight:600; font-family:inherit; transition:all .2s; }
.tl-pagination__btn:hover:not(:disabled) { border-color:#6366f1; color:#6366f1; }
.tl-pagination__btn--active { background:linear-gradient(135deg,#6366f1,#3b82f6) !important; color:white !important; border-color:transparent !important; box-shadow:0 4px 12px rgba(99,102,241,.35); min-width:40px; justify-content:center; }
.tl-pagination__btn:disabled { opacity:.35; cursor:not-allowed; }
</style>
