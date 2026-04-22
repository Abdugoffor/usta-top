<script setup>
import { Icon } from '@iconify/vue'
import BaseTableWrapper from '@/shared/components/BaseTableWrapper.vue'

const props = defineProps({
  items:     { type: Array,   default: () => [] },
  loading:   { type: Boolean, default: false },
  sortBy:    { type: String,  default: '' },
  sortOrder: { type: String,  default: 'desc' },
})

const emit = defineEmits(['show', 'delete', 'sort'])

const toggleSort = (field) => {
  if (props.sortBy === field) {
    emit('sort', { field, order: props.sortOrder === 'asc' ? 'desc' : 'asc' })
  } else {
    emit('sort', { field, order: 'desc' })
  }
}

const sortIcon = (field) => {
  if (props.sortBy !== field) return 'mdi:unfold-more-horizontal'
  return props.sortOrder === 'desc' ? 'mdi:arrow-down' : 'mdi:arrow-up'
}

const fmtPrice = (v) => v ? Number(v).toLocaleString('uz-UZ') + ' so\'m' : '—'
</script>

<template>
  <BaseTableWrapper>

    <template v-if="loading">
      <div class="vt-skeleton">
        <div v-for="i in 5" :key="i" class="vt-skeleton__row">
          <div class="vt-skeleton__cell vt-skeleton__cell--sm" />
          <div class="vt-skeleton__cell vt-skeleton__cell--lg" />
          <div class="vt-skeleton__cell vt-skeleton__cell--md" />
          <div class="vt-skeleton__cell vt-skeleton__cell--md" />
          <div class="vt-skeleton__cell vt-skeleton__cell--sm" />
          <div class="vt-skeleton__cell vt-skeleton__cell--sm" />
        </div>
      </div>
    </template>

    <template v-else-if="!items.length">
      <div class="vt-empty">
        <div class="vt-empty__icon"><Icon icon="mdi:briefcase-search-outline" /></div>
        <p class="vt-empty__title">Vakansiya topilmadi</p>
        <p class="vt-empty__sub">Qidiruv so'rovingizni o'zgartiring</p>
      </div>
    </template>

    <table v-else class="vt-table">
      <thead>
        <tr>
          <th class="vt-th vt-th--sortable" @click="toggleSort('id')">
            <div class="vt-th__inner">
              <span>#</span>
              <Icon :icon="sortIcon('id')" class="vt-sort" :class="{ 'vt-sort--on': sortBy === 'id' }" />
            </div>
          </th>
          <th class="vt-th vt-th--sortable" @click="toggleSort('name')">
            <div class="vt-th__inner">
              <span>Lavozim / Sarlavha</span>
              <Icon :icon="sortIcon('name')" class="vt-sort" :class="{ 'vt-sort--on': sortBy === 'name' }" />
            </div>
          </th>
          <th class="vt-th">Aloqa</th>
          <th class="vt-th vt-th--sortable" @click="toggleSort('price')">
            <div class="vt-th__inner">
              <span>Maosh</span>
              <Icon :icon="sortIcon('price')" class="vt-sort" :class="{ 'vt-sort--on': sortBy === 'price' }" />
            </div>
          </th>
          <th class="vt-th vt-th--sortable" @click="toggleSort('is_active')">
            <div class="vt-th__inner">
              <span>Holati</span>
              <Icon :icon="sortIcon('is_active')" class="vt-sort" :class="{ 'vt-sort--on': sortBy === 'is_active' }" />
            </div>
          </th>
          <th class="vt-th vt-th--actions">Amallar</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id" class="vt-row">
          <td class="vt-td vt-td--id">{{ item.id }}</td>

          <td class="vt-td">
            <div class="vt-job">
              <div class="vt-job__icon">
                <Icon icon="mdi:briefcase-outline" />
              </div>
              <div class="vt-job__info">
                <span class="vt-job__name">{{ item.name }}</span>
                <span class="vt-job__title">{{ item.title }}</span>
              </div>
            </div>
          </td>

          <td class="vt-td vt-td--contact">{{ item.contact || '—' }}</td>
          <td class="vt-td vt-td--price">{{ fmtPrice(item.price) }}</td>

          <td class="vt-td">
            <span :class="['vt-badge', item.is_active ? 'vt-badge--active' : 'vt-badge--inactive']">
              <span class="vt-badge__dot" />
              {{ item.is_active ? 'Faol' : 'Nofaol' }}
            </span>
          </td>

          <td class="vt-td vt-td--actions">
            <div class="vt-actions">
              <button class="vt-btn vt-btn--show"   @click="emit('show', item)"   title="Ko'rish">
                <Icon icon="mdi:eye-outline" />
              </button>
              <button class="vt-btn vt-btn--delete" @click="emit('delete', item)" title="O'chirish">
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
.vt-skeleton { display:flex; flex-direction:column; gap:1px; }
.vt-skeleton__row { display:grid; grid-template-columns:60px 1fr 120px 120px 100px 80px; gap:16px; padding:14px 16px; border-bottom:1px solid var(--border); align-items:center; }
.vt-skeleton__cell { height:14px; border-radius:8px; background:linear-gradient(90deg,var(--bg-elevated) 25%,rgba(255,255,255,.06) 50%,var(--bg-elevated) 75%); background-size:200% 100%; animation:shimmer 1.4s infinite; }
.vt-skeleton__cell--sm { width:40%; }
.vt-skeleton__cell--md { width:70%; }
.vt-skeleton__cell--lg { width:100%; }
@keyframes shimmer { 0%{background-position:200% 0} 100%{background-position:-200% 0} }

.vt-empty { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:64px 24px; gap:12px; text-align:center; }
.vt-empty__icon { width:72px; height:72px; border-radius:20px; background:var(--bg-elevated); border:1px solid var(--border); display:flex; align-items:center; justify-content:center; font-size:32px; color:var(--muted); margin-bottom:4px; }
.vt-empty__title { font-size:16px; font-weight:700; color:var(--text); }
.vt-empty__sub { font-size:13px; color:var(--muted); max-width:320px; line-height:1.5; }

.vt-table { width:100%; border-collapse:collapse; }
.vt-th { padding:12px 14px; text-align:left; font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:.07em; color:var(--muted); border-bottom:1px solid var(--border); white-space:nowrap; background:rgba(255,255,255,.02); }
.vt-th--sortable { cursor:pointer; user-select:none; transition:color .15s; }
.vt-th--sortable:hover { color:var(--text-secondary); }
.vt-th--actions { text-align:center; width:90px; }
.vt-th__inner { display:inline-flex; align-items:center; gap:5px; }
.vt-sort { font-size:14px; color:var(--border-hover); transition:color .15s; }
.vt-sort--on { color:#6366f1; }

.vt-row { transition:background .15s; }
.vt-row:hover { background:rgba(99,102,241,.04); }
.vt-row:not(:last-child) .vt-td { border-bottom:1px solid var(--border); }

.vt-td { padding:13px 14px; color:var(--text); font-size:14px; vertical-align:middle; }
.vt-td--id { color:var(--muted); font-size:12px; font-weight:600; width:60px; }
.vt-td--contact { color:var(--muted); font-size:13px; white-space:nowrap; }
.vt-td--price { font-size:13px; font-weight:600; color:var(--text); white-space:nowrap; }
.vt-td--actions { text-align:center; }

.vt-job { display:flex; align-items:center; gap:12px; }
.vt-job__icon { width:36px; height:36px; border-radius:10px; background:linear-gradient(135deg,rgba(99,102,241,.2),rgba(59,130,246,.2)); border:1px solid rgba(99,102,241,.2); color:#818cf8; font-size:18px; display:flex; align-items:center; justify-content:center; flex-shrink:0; }
.vt-job__info { display:flex; flex-direction:column; min-width:0; }
.vt-job__name { font-size:14px; font-weight:600; color:var(--text); white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.vt-job__title { font-size:12px; color:var(--muted); white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }

.vt-badge { display:inline-flex; align-items:center; gap:6px; padding:5px 12px; border-radius:20px; font-size:12px; font-weight:600; white-space:nowrap; }
.vt-badge__dot { width:6px; height:6px; border-radius:50%; flex-shrink:0; }
.vt-badge--active { background:rgba(16,185,129,.12); color:#34d399; border:1px solid rgba(16,185,129,.22); }
.vt-badge--active .vt-badge__dot { background:#34d399; box-shadow:0 0 6px rgba(52,211,153,.6); animation:pulse 2s infinite; }
.vt-badge--inactive { background:rgba(239,68,68,.1); color:#f87171; border:1px solid rgba(239,68,68,.2); }
.vt-badge--inactive .vt-badge__dot { background:#f87171; }
@keyframes pulse { 0%,100%{opacity:1;transform:scale(1)} 50%{opacity:.6;transform:scale(.85)} }

.vt-actions { display:flex; align-items:center; justify-content:center; gap:6px; }
.vt-btn { width:34px; height:34px; border:none; border-radius:10px; cursor:pointer; display:flex; align-items:center; justify-content:center; font-size:16px; transition:all .2s; color:white; }
.vt-btn:hover { transform:translateY(-2px); }
.vt-btn--show { background:linear-gradient(135deg,#0ea5e9,#06b6d4); box-shadow:0 2px 8px rgba(14,165,233,.3); }
.vt-btn--show:hover { box-shadow:0 6px 16px rgba(14,165,233,.45); }
.vt-btn--delete { background:linear-gradient(135deg,#ef4444,#f97316); box-shadow:0 2px 8px rgba(239,68,68,.3); }
.vt-btn--delete:hover { box-shadow:0 6px 16px rgba(239,68,68,.45); }

@media (max-width:768px) {
  .vt-td--contact, .vt-td--price { display:none; }
}
</style>
