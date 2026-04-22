<script setup>
import { Icon } from '@iconify/vue'
import BaseTableWrapper from '@/shared/components/BaseTableWrapper.vue'

const props = defineProps({
  items:     { type: Array,   default: () => [] },
  loading:   { type: Boolean, default: false },
  sortBy:    { type: String,  default: '' },
  sortOrder: { type: String,  default: 'asc' },
})

const emit = defineEmits(['show', 'edit', 'delete', 'sort'])

const toggleSort = (field) => {
  if (props.sortBy === field) {
    emit('sort', { field, order: props.sortOrder === 'asc' ? 'desc' : 'asc' })
  } else {
    emit('sort', { field, order: 'asc' })
  }
}

const sortIcon = (field) => {
  if (props.sortBy !== field) return 'mdi:unfold-more-horizontal'
  return props.sortOrder === 'desc' ? 'mdi:arrow-down' : 'mdi:arrow-up'
}
</script>

<template>
  <BaseTableWrapper>

    <template v-if="loading">
      <div class="lt-skeleton">
        <div v-for="i in 5" :key="i" class="lt-skeleton__row">
          <div class="lt-skeleton__cell lt-skeleton__cell--sm" />
          <div class="lt-skeleton__cell lt-skeleton__cell--md" />
          <div class="lt-skeleton__cell lt-skeleton__cell--lg" />
          <div class="lt-skeleton__cell lt-skeleton__cell--sm" />
          <div class="lt-skeleton__cell lt-skeleton__cell--sm" />
        </div>
      </div>
    </template>

    <template v-else-if="!items.length">
      <div class="lt-empty">
        <div class="lt-empty__icon"><Icon icon="mdi:translate-off" /></div>
        <p class="lt-empty__title">Til topilmadi</p>
        <p class="lt-empty__sub">Qidiruv so'rovingizni o'zgartiring yoki yangi til qo'shing</p>
      </div>
    </template>

    <table v-else class="lt-table">
      <thead>
        <tr>
          <th class="lt-th lt-th--sortable" @click="toggleSort('id')">
            <div class="lt-th__inner">
              <span>#</span>
              <Icon :icon="sortIcon('id')" class="lt-sort" :class="{ 'lt-sort--on': sortBy === 'id' }" />
            </div>
          </th>
          <th class="lt-th lt-th--sortable" @click="toggleSort('name')">
            <div class="lt-th__inner">
              <span>Nomi</span>
              <Icon :icon="sortIcon('name')" class="lt-sort" :class="{ 'lt-sort--on': sortBy === 'name' }" />
            </div>
          </th>
          <th class="lt-th">Tavsif</th>
          <th class="lt-th lt-th--sortable" @click="toggleSort('is_active')">
            <div class="lt-th__inner">
              <span>Holati</span>
              <Icon :icon="sortIcon('is_active')" class="lt-sort" :class="{ 'lt-sort--on': sortBy === 'is_active' }" />
            </div>
          </th>
          <th class="lt-th lt-th--actions">Amallar</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id" class="lt-row">
          <td class="lt-td lt-td--id">{{ item.id }}</td>

          <td class="lt-td">
            <div class="lt-name">
              <div class="lt-name__avatar">{{ item.name?.charAt(0)?.toUpperCase() }}</div>
              <span>{{ item.name }}</span>
            </div>
          </td>

          <td class="lt-td lt-td--desc">{{ item.description || '—' }}</td>

          <td class="lt-td">
            <span :class="['lt-badge', item.is_active ? 'lt-badge--active' : 'lt-badge--inactive']">
              <span class="lt-badge__dot" />
              {{ item.is_active ? 'Faol' : 'Nofaol' }}
            </span>
          </td>

          <td class="lt-td lt-td--actions">
            <div class="lt-actions">
              <button class="lt-btn lt-btn--show"   @click="emit('show', item)"   title="Ko'rish">
                <Icon icon="mdi:eye-outline" />
              </button>
              <button class="lt-btn lt-btn--edit"   @click="emit('edit', item)"   title="Tahrirlash">
                <Icon icon="mdi:pencil-outline" />
              </button>
              <button class="lt-btn lt-btn--delete" @click="emit('delete', item)" title="O'chirish">
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
.lt-skeleton { display:flex; flex-direction:column; gap:1px; }
.lt-skeleton__row { display:grid; grid-template-columns:60px 1fr 1fr 100px 120px; gap:16px; padding:14px 16px; border-bottom:1px solid var(--border); align-items:center; }
.lt-skeleton__cell { height:14px; border-radius:8px; background:linear-gradient(90deg,var(--bg-elevated) 25%,rgba(255,255,255,.06) 50%,var(--bg-elevated) 75%); background-size:200% 100%; animation:shimmer 1.4s infinite; }
.lt-skeleton__cell--sm { width:50%; }
.lt-skeleton__cell--md { width:75%; }
.lt-skeleton__cell--lg { width:100%; }
@keyframes shimmer { 0%{background-position:200% 0} 100%{background-position:-200% 0} }

.lt-empty { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:64px 24px; gap:12px; text-align:center; }
.lt-empty__icon { width:72px; height:72px; border-radius:20px; background:var(--bg-elevated); border:1px solid var(--border); display:flex; align-items:center; justify-content:center; font-size:32px; color:var(--muted); margin-bottom:4px; }
.lt-empty__title { font-size:16px; font-weight:700; color:var(--text); }
.lt-empty__sub { font-size:13px; color:var(--muted); max-width:320px; line-height:1.5; }

.lt-table { width:100%; border-collapse:collapse; }
.lt-th { padding:12px 14px; text-align:left; font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:.07em; color:var(--muted); border-bottom:1px solid var(--border); white-space:nowrap; background:rgba(255,255,255,.02); }
.lt-th--sortable { cursor:pointer; user-select:none; transition:color .15s; }
.lt-th--sortable:hover { color:var(--text-secondary); }
.lt-th--actions { text-align:center; width:120px; }
.lt-th__inner { display:inline-flex; align-items:center; gap:5px; }
.lt-sort { font-size:14px; color:var(--border-hover); transition:color .15s; }
.lt-sort--on { color:#6366f1; }

.lt-row { transition:background .15s; }
.lt-row:hover { background:rgba(99,102,241,.04); }
.lt-row:not(:last-child) .lt-td { border-bottom:1px solid var(--border); }

.lt-td { padding:13px 14px; color:var(--text); font-size:14px; vertical-align:middle; }
.lt-td--id { color:var(--muted); font-size:12px; font-weight:600; width:60px; }
.lt-td--desc { color:var(--muted); font-size:13px; max-width:220px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.lt-td--actions { text-align:center; }

.lt-name { display:flex; align-items:center; gap:10px; }
.lt-name__avatar { width:32px; height:32px; border-radius:10px; background:linear-gradient(135deg,rgba(99,102,241,.25),rgba(59,130,246,.25)); border:1px solid rgba(99,102,241,.2); color:#818cf8; font-size:13px; font-weight:800; display:flex; align-items:center; justify-content:center; flex-shrink:0; }

.lt-badge { display:inline-flex; align-items:center; gap:6px; padding:4px 11px; border-radius:20px; font-size:12px; font-weight:600; white-space:nowrap; }
.lt-badge__dot { width:6px; height:6px; border-radius:50%; flex-shrink:0; }
.lt-badge--active { background:rgba(16,185,129,.12); color:#34d399; border:1px solid rgba(16,185,129,.22); }
.lt-badge--active .lt-badge__dot { background:#34d399; box-shadow:0 0 6px rgba(52,211,153,.6); animation:pulse 2s infinite; }
.lt-badge--inactive { background:rgba(239,68,68,.1); color:#f87171; border:1px solid rgba(239,68,68,.2); }
.lt-badge--inactive .lt-badge__dot { background:#f87171; }
@keyframes pulse { 0%,100%{opacity:1;transform:scale(1)} 50%{opacity:.6;transform:scale(.85)} }

.lt-actions { display:flex; align-items:center; justify-content:center; gap:6px; }
.lt-btn { width:34px; height:34px; border:none; border-radius:10px; cursor:pointer; display:flex; align-items:center; justify-content:center; font-size:16px; transition:all .2s; color:white; }
.lt-btn:hover { transform:translateY(-2px); }
.lt-btn--show { background:linear-gradient(135deg,#0ea5e9,#06b6d4); box-shadow:0 2px 8px rgba(14,165,233,.3); }
.lt-btn--show:hover { box-shadow:0 6px 16px rgba(14,165,233,.45); }
.lt-btn--edit { background:linear-gradient(135deg,#3b82f6,#6366f1); box-shadow:0 2px 8px rgba(99,102,241,.3); }
.lt-btn--edit:hover { box-shadow:0 6px 16px rgba(99,102,241,.45); }
.lt-btn--delete { background:linear-gradient(135deg,#ef4444,#f97316); box-shadow:0 2px 8px rgba(239,68,68,.3); }
.lt-btn--delete:hover { box-shadow:0 6px 16px rgba(239,68,68,.45); }
</style>
