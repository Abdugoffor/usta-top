<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Icon } from '@iconify/vue'
import BaseCard from '@/shared/components/BaseCard.vue'
import { useTranslationStore } from '../store/translationStore'

const route = useRoute()
const router = useRouter()
const store = useTranslationStore()
const item = ref(null)

const fmt = (iso) => iso ? new Date(iso).toLocaleString('uz-UZ') : '—'

onMounted(async () => {
  item.value = await store.fetchTranslation(route.params.id)
})
</script>

<template>
  <section class="page-section">
    <div class="cl-header">
      <div class="cl-header__left">
        <button class="cl-back-btn" @click="router.push({ name: 'translation-list' })">
          <Icon icon="mdi:arrow-left" />
        </button>
        <div>
          <h2 class="cl-header__title">Tarjima ma'lumotlari</h2>
          <p class="cl-header__sub">Batafsil ko'rish</p>
        </div>
      </div>
      <button class="cl-edit-btn" @click="router.push({ name: 'translation-edit', params: { id: route.params.id } })">
        <Icon icon="mdi:pencil-outline" />
        <span>Tahrirlash</span>
      </button>
    </div>

    <BaseCard v-if="!item">
      <div class="ct-empty">
        <div class="cl-spinner"></div>
        <p style="color:var(--muted);font-size:14px">Yuklanmoqda...</p>
      </div>
    </BaseCard>

    <BaseCard v-else>
      <div class="show-grid">
        <div class="show-field">
          <span class="show-field__label">ID</span>
          <span class="show-field__value show-field__value--mono">{{ item.id }}</span>
        </div>

        <div class="show-field">
          <span class="show-field__label">Kalit (slug)</span>
          <code class="show-slug">{{ item.slug }}</code>
        </div>

        <div class="show-field show-field--col">
          <span class="show-field__label">Tarjimalar</span>
          <div class="show-translations">
            <div
              v-for="(val, key) in item.name"
              :key="key"
              class="show-translation-row"
            >
              <span :class="['show-lang-chip', key === 'default' ? 'show-lang-chip--default' : '']">
                {{ key }}
              </span>
              <span class="show-translation-val">{{ val || '—' }}</span>
            </div>
          </div>
        </div>

        <div class="show-field">
          <span class="show-field__label">Holati</span>
          <span :class="['lt-badge', item.is_active ? 'lt-badge--active' : 'lt-badge--inactive']">
            <span class="lt-badge__dot" />
            {{ item.is_active ? 'Faol' : 'Nofaol' }}
          </span>
        </div>

        <div class="show-field">
          <span class="show-field__label">Yaratilgan</span>
          <span class="show-field__value show-field__value--date">
            <Icon icon="mdi:calendar-outline" />
            {{ fmt(item.created_at) }}
          </span>
        </div>

        <div class="show-field">
          <span class="show-field__label">Yangilangan</span>
          <span class="show-field__value show-field__value--date">
            <Icon icon="mdi:calendar-clock" />
            {{ fmt(item.updated_at) }}
          </span>
        </div>
      </div>
    </BaseCard>
  </section>
</template>

<style scoped>
.cl-header { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:16px; }
.cl-header__left { display:flex; align-items:center; gap:14px; }
.cl-back-btn { width:42px; height:42px; border-radius:12px; background:var(--bg-elevated); border:1px solid var(--border); color:var(--text-secondary); font-size:20px; display:flex; align-items:center; justify-content:center; cursor:pointer; transition:all .2s; }
.cl-back-btn:hover { border-color:#6366f1; color:#6366f1; }
.cl-header__title { font-size:24px; font-weight:800; background:linear-gradient(135deg,#f1f5f9,#94a3b8); -webkit-background-clip:text; -webkit-text-fill-color:transparent; background-clip:text; }
.cl-header__sub { font-size:13px; color:var(--muted); margin-top:3px; }
.cl-edit-btn { display:inline-flex; align-items:center; gap:8px; padding:11px 22px; background:linear-gradient(135deg,#6366f1,#3b82f6); color:white; border:none; border-radius:14px; font-size:14px; font-weight:700; cursor:pointer; transition:all .25s; box-shadow:0 6px 20px rgba(99,102,241,.35); }
.cl-edit-btn:hover { transform:translateY(-2px); box-shadow:0 10px 28px rgba(99,102,241,.5); }

.ct-empty { display:flex; flex-direction:column; align-items:center; gap:12px; padding:40px; }
.cl-spinner { width:32px; height:32px; border:3px solid var(--border); border-top-color:#6366f1; border-radius:50%; animation:spin .8s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }

.show-grid { display:flex; flex-direction:column; }
.show-field { display:flex; align-items:center; gap:16px; padding:16px 0; border-bottom:1px solid var(--border); }
.show-field--col { flex-direction:column; align-items:flex-start; }
.show-field:last-child { border-bottom:none; }
.show-field__label { font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:.07em; color:var(--muted); min-width:140px; flex-shrink:0; }
.show-field__value { font-size:14px; color:var(--text); }
.show-field__value--mono { font-family:monospace; font-size:13px; color:var(--muted); }
.show-field__value--date { display:flex; align-items:center; gap:6px; color:var(--muted); font-size:13px; }

.show-slug { background:rgba(99,102,241,.12); color:#818cf8; border-radius:8px; padding:4px 12px; font-size:13px; font-family:monospace; }

.show-translations { display:flex; flex-direction:column; gap:8px; width:100%; margin-top:4px; }
.show-translation-row { display:flex; align-items:center; gap:12px; padding:10px 14px; background:var(--bg-elevated); border:1px solid var(--border); border-radius:10px; }
.show-lang-chip { padding:2px 9px; border-radius:7px; font-size:11px; font-weight:700; text-transform:uppercase; background:rgba(99,102,241,.15); color:#818cf8; border:1px solid rgba(99,102,241,.25); flex-shrink:0; min-width:52px; text-align:center; }
.show-lang-chip--default { background:rgba(16,185,129,.12); color:#34d399; border-color:rgba(16,185,129,.22); }
.show-translation-val { font-size:14px; color:var(--text); }

.lt-badge { display:inline-flex; align-items:center; gap:6px; padding:5px 12px; border-radius:20px; font-size:12px; font-weight:600; }
.lt-badge__dot { width:6px; height:6px; border-radius:50%; flex-shrink:0; }
.lt-badge--active { background:rgba(16,185,129,.12); color:#34d399; border:1px solid rgba(16,185,129,.22); }
.lt-badge--active .lt-badge__dot { background:#34d399; box-shadow:0 0 6px rgba(52,211,153,.6); animation:pulse 2s infinite; }
.lt-badge--inactive { background:rgba(239,68,68,.1); color:#f87171; border:1px solid rgba(239,68,68,.2); }
.lt-badge--inactive .lt-badge__dot { background:#f87171; }
@keyframes pulse { 0%,100%{opacity:1;transform:scale(1)} 50%{opacity:.6;transform:scale(.85)} }
</style>
