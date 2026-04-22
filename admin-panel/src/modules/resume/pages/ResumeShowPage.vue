<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import { Icon } from '@iconify/vue'
import BaseCard from '@/shared/components/BaseCard.vue'
import { useResumeStore } from '../store/resumeStore'

const route = useRoute()
const router = useRouter()
const store = useResumeStore()
const item = ref(null)
const toggling = ref(false)

const fmt = (iso) => iso ? new Date(iso).toLocaleString('uz-UZ') : '—'
const fmtPrice = (v) => v ? Number(v).toLocaleString('uz-UZ') + ' so\'m' : '—'

const toggleActive = async () => {
  if (toggling.value) return
  toggling.value = true
  try {
    await store.updateResume(item.value.id, { is_active: !item.value.is_active })
    item.value.is_active = !item.value.is_active
  } finally {
    toggling.value = false
  }
}

const deleteItem = async () => {
  const result = await Swal.fire({
    title: "O'chirishni xohlaysizmi?",
    html: `<span style="color:#94a3b8">"${item.value.name}"</span> resumesi o'chiriladi`,
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
  await store.removeResume(item.value.id)
  router.push({ name: 'resume-admin-list' })
}

onMounted(async () => {
  item.value = await store.fetchResume(route.params.id)
})
</script>

<template>
  <section class="page-section">
    <div class="cl-header">
      <div class="cl-header__left">
        <button class="cl-back-btn" @click="router.push({ name: 'resume-admin-list' })">
          <Icon icon="mdi:arrow-left" />
        </button>
        <div>
          <h2 class="cl-header__title">Resume ma'lumotlari</h2>
          <p class="cl-header__sub">Batafsil ko'rish</p>
        </div>
      </div>
      <div class="cl-header__actions" v-if="item">
        <button
          class="cl-toggle-btn"
          :class="item.is_active ? 'cl-toggle-btn--deactivate' : 'cl-toggle-btn--activate'"
          :disabled="toggling"
          @click="toggleActive"
        >
          <Icon :icon="item.is_active ? 'mdi:eye-off-outline' : 'mdi:eye-outline'" />
          {{ item.is_active ? 'Nofaol qilish' : 'Faol qilish' }}
        </button>
        <button class="cl-delete-btn" @click="deleteItem">
          <Icon icon="mdi:trash-can-outline" />
          O'chirish
        </button>
      </div>
    </div>

    <BaseCard v-if="!item">
      <div class="show-loading">
        <div class="cl-spinner"></div>
        <p>Yuklanmoqda...</p>
      </div>
    </BaseCard>

    <template v-else>
      <BaseCard>
        <div class="show-top">
          <div class="show-avatar">
            <img v-if="item.photo" :src="item.photo" :alt="item.name" class="show-avatar__img" />
            <span v-else class="show-avatar__letter">{{ item.name?.charAt(0)?.toUpperCase() }}</span>
          </div>
          <div class="show-intro">
            <h3 class="show-intro__name">{{ item.name }}</h3>
            <p class="show-intro__title">{{ item.title }}</p>
            <span :class="['show-badge', item.is_active ? 'show-badge--active' : 'show-badge--inactive']">
              <span class="show-badge__dot" />
              {{ item.is_active ? 'Faol' : 'Nofaol' }}
            </span>
          </div>
        </div>
      </BaseCard>

      <BaseCard>
        <div class="show-grid">
          <div class="show-field">
            <span class="show-field__label">ID</span>
            <span class="show-field__value show-field__value--mono">{{ item.id }}</span>
          </div>
          <div class="show-field" v-if="item.slug">
            <span class="show-field__label">Slug</span>
            <span class="show-field__value show-field__value--mono">{{ item.slug }}</span>
          </div>
          <div class="show-field">
            <span class="show-field__label">Aloqa</span>
            <span class="show-field__value">{{ item.contact || '—' }}</span>
          </div>
          <div class="show-field">
            <span class="show-field__label">Narx</span>
            <span class="show-field__value">{{ fmtPrice(item.price) }}</span>
          </div>
          <div class="show-field" v-if="item.experience_year">
            <span class="show-field__label">Tajriba</span>
            <span class="show-field__value">{{ item.experience_year }} yil</span>
          </div>
          <div class="show-field" v-if="item.adress">
            <span class="show-field__label">Manzil</span>
            <span class="show-field__value">{{ item.adress }}</span>
          </div>
          <div class="show-field" v-if="item.skills">
            <span class="show-field__label">Ko'nikmalar</span>
            <div class="show-chips">
              <span v-for="s in item.skills.split(',')" :key="s" class="show-chip">{{ s.trim() }}</span>
            </div>
          </div>
          <div class="show-field" v-if="item.text">
            <span class="show-field__label">Tavsif</span>
            <p class="show-field__value show-field__value--text">{{ item.text }}</p>
          </div>
          <div class="show-field">
            <span class="show-field__label">Yaratilgan vaqt</span>
            <span class="show-field__value show-field__value--date">
              <Icon icon="mdi:calendar-outline" />
              {{ fmt(item.created_at) }}
            </span>
          </div>
          <div class="show-field">
            <span class="show-field__label">Yangilangan vaqt</span>
            <span class="show-field__value show-field__value--date">
              <Icon icon="mdi:calendar-clock" />
              {{ fmt(item.updated_at) }}
            </span>
          </div>
        </div>
      </BaseCard>
    </template>
  </section>
</template>

<style scoped>
.cl-header { display:flex; align-items:center; justify-content:space-between; flex-wrap:wrap; gap:16px; }
.cl-header__left { display:flex; align-items:center; gap:14px; }
.cl-header__actions { display:flex; gap:10px; }
.cl-back-btn { width:42px; height:42px; border-radius:12px; background:var(--bg-elevated); border:1px solid var(--border); color:var(--text-secondary); font-size:20px; display:flex; align-items:center; justify-content:center; cursor:pointer; transition:all .2s; }
.cl-back-btn:hover { border-color:#6366f1; color:#6366f1; }
.cl-header__title { font-size:24px; font-weight:800; background:linear-gradient(135deg,#f1f5f9,#94a3b8); -webkit-background-clip:text; -webkit-text-fill-color:transparent; background-clip:text; line-height:1.2; }
.cl-header__sub { font-size:13px; color:var(--muted); margin-top:3px; }

.cl-toggle-btn { display:inline-flex; align-items:center; gap:7px; padding:10px 18px; border:none; border-radius:12px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; transition:all .2s; }
.cl-toggle-btn--activate { background:rgba(16,185,129,.15); color:#34d399; }
.cl-toggle-btn--activate:hover { background:rgba(16,185,129,.25); }
.cl-toggle-btn--deactivate { background:rgba(251,191,36,.12); color:#fbbf24; }
.cl-toggle-btn--deactivate:hover { background:rgba(251,191,36,.22); }
.cl-toggle-btn:disabled { opacity:.6; cursor:not-allowed; }

.cl-delete-btn { display:inline-flex; align-items:center; gap:7px; padding:10px 18px; background:rgba(239,68,68,.12); color:#f87171; border:none; border-radius:12px; font-size:13px; font-weight:600; font-family:inherit; cursor:pointer; transition:all .2s; }
.cl-delete-btn:hover { background:rgba(239,68,68,.22); }

.show-loading { display:flex; flex-direction:column; align-items:center; gap:12px; padding:40px; color:var(--muted); font-size:14px; }
.cl-spinner { width:32px; height:32px; border:3px solid var(--border); border-top-color:#6366f1; border-radius:50%; animation:spin .8s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }

.show-top { display:flex; align-items:center; gap:20px; }
.show-avatar { width:72px; height:72px; border-radius:20px; flex-shrink:0; overflow:hidden; background:linear-gradient(135deg,rgba(99,102,241,.3),rgba(59,130,246,.3)); border:2px solid rgba(99,102,241,.2); display:flex; align-items:center; justify-content:center; }
.show-avatar__img { width:100%; height:100%; object-fit:cover; }
.show-avatar__letter { font-size:28px; font-weight:800; color:#818cf8; }
.show-intro { display:flex; flex-direction:column; gap:6px; }
.show-intro__name { font-size:20px; font-weight:700; color:var(--text); }
.show-intro__title { font-size:14px; color:var(--muted); }

.show-badge { display:inline-flex; align-items:center; gap:6px; padding:4px 12px; border-radius:20px; font-size:12px; font-weight:600; }
.show-badge__dot { width:6px; height:6px; border-radius:50%; }
.show-badge--active { background:rgba(16,185,129,.12); color:#34d399; border:1px solid rgba(16,185,129,.22); }
.show-badge--active .show-badge__dot { background:#34d399; box-shadow:0 0 6px rgba(52,211,153,.6); animation:pulse 2s infinite; }
.show-badge--inactive { background:rgba(239,68,68,.1); color:#f87171; border:1px solid rgba(239,68,68,.2); }
.show-badge--inactive .show-badge__dot { background:#f87171; }
@keyframes pulse { 0%,100%{opacity:1;transform:scale(1)} 50%{opacity:.6;transform:scale(.85)} }

.show-grid { display:flex; flex-direction:column; }
.show-field { display:flex; align-items:flex-start; gap:16px; padding:16px 0; border-bottom:1px solid var(--border); }
.show-field:last-child { border-bottom:none; }
.show-field__label { font-size:12px; font-weight:700; text-transform:uppercase; letter-spacing:.07em; color:var(--muted); min-width:160px; flex-shrink:0; padding-top:2px; }
.show-field__value { font-size:14px; color:var(--text); }
.show-field__value--mono { font-family:monospace; font-size:13px; color:var(--muted); }
.show-field__value--date { display:flex; align-items:center; gap:6px; color:var(--muted); font-size:13px; }
.show-field__value--text { line-height:1.6; color:var(--text-secondary); white-space:pre-wrap; }

.show-chips { display:flex; flex-wrap:wrap; gap:6px; }
.show-chip { padding:4px 12px; background:rgba(99,102,241,.12); color:#818cf8; border:1px solid rgba(99,102,241,.2); border-radius:20px; font-size:12px; font-weight:600; }
</style>
