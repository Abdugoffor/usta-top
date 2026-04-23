<script setup>
import { Icon } from '@iconify/vue'
import BaseTableWrapper from '@/shared/components/BaseTableWrapper.vue'

const props = defineProps({
  items:     { type: Array,   default: () => [] },
  loading:   { type: Boolean, default: false },
})

const emit = defineEmits(['show', 'edit', 'delete'])

const defaultName = (item) => item.name?.default || Object.values(item.name ?? {})[0] || '—'
</script>

<template>
  <BaseTableWrapper>

    <template v-if="loading">
      <div class="tt-skeleton">
        <div v-for="i in 6" :key="i" class="tt-skeleton__row">
          <div class="tt-skeleton__cell tt-skeleton__cell--sm" />
          <div class="tt-skeleton__cell tt-skeleton__cell--md" />
          <div class="tt-skeleton__cell tt-skeleton__cell--lg" />
          <div class="tt-skeleton__cell tt-skeleton__cell--sm" />
          <div class="tt-skeleton__cell tt-skeleton__cell--sm" />
        </div>
      </div>
    </template>

    <template v-else-if="!items.length">
      <div class="tt-empty">
        <div class="tt-empty__icon"><Icon icon="mdi:translate-off" /></div>
        <p class="tt-empty__title">Tarjima topilmadi</p>
        <p class="tt-empty__sub">Qidiruv so'rovingizni o'zgartiring yoki yangi tarjima qo'shing</p>
      </div>
    </template>

    <table v-else class="tt-table">
      <thead>
        <tr>
          <th class="tt-th">#</th>
          <th class="tt-th">Kalit (slug)</th>
          <th class="tt-th">Standart matn</th>
          <th class="tt-th">Tillar</th>
          <th class="tt-th">Holati</th>
          <th class="tt-th tt-th--actions">Amallar</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id" class="tt-row">
          <td class="tt-td tt-td--id">{{ item.id }}</td>

          <td class="tt-td">
            <code class="tt-slug">{{ item.slug }}</code>
          </td>

          <td class="tt-td tt-td--text">{{ defaultName(item) }}</td>

          <td class="tt-td">
            <div class="tt-langs">
              <span
                v-for="(val, key) in item.name"
                :key="key"
                :class="['tt-lang-chip', key === 'default' ? 'tt-lang-chip--default' : '']"
                :title="val"
              >{{ key }}</span>
            </div>
          </td>

          <td class="tt-td">
            <span :class="['tt-badge', item.is_active ? 'tt-badge--active' : 'tt-badge--inactive']">
              <span class="tt-badge__dot" />
              {{ item.is_active ? 'Faol' : 'Nofaol' }}
            </span>
          </td>

          <td class="tt-td tt-td--actions">
            <div class="tt-actions">
              <button class="tt-btn tt-btn--show"   @click="emit('show', item)"   title="Ko'rish">
                <Icon icon="mdi:eye-outline" />
              </button>
              <button class="tt-btn tt-btn--edit"   @click="emit('edit', item)"   title="Tahrirlash">
                <Icon icon="mdi:pencil-outline" />
              </button>
              <button class="tt-btn tt-btn--delete" @click="emit('delete', item)" title="O'chirish">
                <Icon icon="mdi:trash-can-outline" />
              </button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

  </BaseTableWrapper>
</template>

<style scoped>
.tt-skeleton { display:flex; flex-direction:column; gap:1px; }
.tt-skeleton__row { display:grid; grid-template-columns:60px 1fr 1.5fr 120px 100px 120px; gap:16px; padding:14px 16px; border-bottom:1px solid var(--border); align-items:center; }
.tt-skeleton__cell { height:14px; border-radius:8px; background:linear-gradient(90deg,var(--bg-elevated) 25%,rgba(255,255,255,.06) 50%,var(--bg-elevated) 75%); background-size:200% 100%; animation:shimmer 1.4s infinite; }
.tt-skeleton__cell--sm { width:50%; }
.tt-skeleton__cell--md { width:75%; }
.tt-skeleton__cell--lg { width:100%; }
@keyframes shimmer { 0%{background-position:200% 0} 100%{background-position:-200% 0} }

.tt-empty { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:64px 24px; gap:12px; text-align:center; }
.tt-empty__icon { width:72px; height:72px; border-radius:20px; background:var(--bg-elevated); border:1px solid var(--border); display:flex; align-items:center; justify-content:center; font-size:32px; color:var(--muted); margin-bottom:4px; }
.tt-empty__title { font-size:16px; font-weight:700; color:var(--text); }
.tt-empty__sub { font-size:13px; color:var(--muted); max-width:320px; line-height:1.5; }

.tt-table { width:100%; border-collapse:collapse; }
.tt-th { padding:12px 14px; text-align:left; font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:.07em; color:var(--muted); border-bottom:1px solid var(--border); white-space:nowrap; background:rgba(255,255,255,.02); }
.tt-th--actions { text-align:center; width:120px; }

.tt-row { transition:background .15s; }
.tt-row:hover { background:rgba(99,102,241,.04); }
.tt-row:not(:last-child) .tt-td { border-bottom:1px solid var(--border); }

.tt-td { padding:13px 14px; color:var(--text); font-size:14px; vertical-align:middle; }
.tt-td--id { color:var(--muted); font-size:12px; font-weight:600; width:60px; }
.tt-td--text { color:var(--text-secondary); font-size:13px; max-width:240px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.tt-td--actions { text-align:center; }

.tt-slug { background:rgba(99,102,241,.12); color:#818cf8; border-radius:6px; padding:3px 8px; font-size:12px; font-family:monospace; }

.tt-langs { display:flex; flex-wrap:wrap; gap:4px; }
.tt-lang-chip { padding:2px 7px; border-radius:6px; font-size:10px; font-weight:700; text-transform:uppercase; background:rgba(99,102,241,.12); color:#818cf8; border:1px solid rgba(99,102,241,.2); }
.tt-lang-chip--default { background:rgba(16,185,129,.1); color:#34d399; border-color:rgba(16,185,129,.2); }

.tt-badge { display:inline-flex; align-items:center; gap:6px; padding:4px 11px; border-radius:20px; font-size:12px; font-weight:600; white-space:nowrap; }
.tt-badge__dot { width:6px; height:6px; border-radius:50%; flex-shrink:0; }
.tt-badge--active { background:rgba(16,185,129,.12); color:#34d399; border:1px solid rgba(16,185,129,.22); }
.tt-badge--active .tt-badge__dot { background:#34d399; box-shadow:0 0 6px rgba(52,211,153,.6); animation:pulse 2s infinite; }
.tt-badge--inactive { background:rgba(239,68,68,.1); color:#f87171; border:1px solid rgba(239,68,68,.2); }
.tt-badge--inactive .tt-badge__dot { background:#f87171; }
@keyframes pulse { 0%,100%{opacity:1;transform:scale(1)} 50%{opacity:.6;transform:scale(.85)} }

.tt-actions { display:flex; align-items:center; justify-content:center; gap:6px; }
.tt-btn { width:34px; height:34px; border:none; border-radius:10px; cursor:pointer; display:flex; align-items:center; justify-content:center; font-size:16px; transition:all .2s; color:white; }
.tt-btn:hover { transform:translateY(-2px); }
.tt-btn--show { background:linear-gradient(135deg,#0ea5e9,#06b6d4); box-shadow:0 2px 8px rgba(14,165,233,.3); }
.tt-btn--show:hover { box-shadow:0 6px 16px rgba(14,165,233,.45); }
.tt-btn--edit { background:linear-gradient(135deg,#3b82f6,#6366f1); box-shadow:0 2px 8px rgba(99,102,241,.3); }
.tt-btn--edit:hover { box-shadow:0 6px 16px rgba(99,102,241,.45); }
.tt-btn--delete { background:linear-gradient(135deg,#ef4444,#f97316); box-shadow:0 2px 8px rgba(239,68,68,.3); }
.tt-btn--delete:hover { box-shadow:0 6px 16px rgba(239,68,68,.45); }
</style>
