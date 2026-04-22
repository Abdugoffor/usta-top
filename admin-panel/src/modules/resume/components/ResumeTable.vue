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
      <div class="rt-skeleton">
        <div v-for="i in 5" :key="i" class="rt-skeleton__row">
          <div class="rt-skeleton__cell rt-skeleton__cell--sm" />
          <div class="rt-skeleton__cell rt-skeleton__cell--lg" />
          <div class="rt-skeleton__cell rt-skeleton__cell--md" />
          <div class="rt-skeleton__cell rt-skeleton__cell--md" />
          <div class="rt-skeleton__cell rt-skeleton__cell--sm" />
          <div class="rt-skeleton__cell rt-skeleton__cell--sm" />
        </div>
      </div>
    </template>

    <template v-else-if="!items.length">
      <div class="rt-empty">
        <div class="rt-empty__icon"><Icon icon="mdi:account-search-outline" /></div>
        <p class="rt-empty__title">Resume topilmadi</p>
        <p class="rt-empty__sub">Qidiruv so'rovingizni o'zgartiring</p>
      </div>
    </template>

    <table v-else class="rt-table">
      <thead>
        <tr>
          <th class="rt-th rt-th--sortable" @click="toggleSort('id')">
            <div class="rt-th__inner">
              <span>#</span>
              <Icon :icon="sortIcon('id')" class="rt-sort" :class="{ 'rt-sort--on': sortBy === 'id' }" />
            </div>
          </th>
          <th class="rt-th rt-th--sortable" @click="toggleSort('name')">
            <div class="rt-th__inner">
              <span>Ism / Kasb</span>
              <Icon :icon="sortIcon('name')" class="rt-sort" :class="{ 'rt-sort--on': sortBy === 'name' }" />
            </div>
          </th>
          <th class="rt-th">Aloqa</th>
          <th class="rt-th rt-th--sortable" @click="toggleSort('price')">
            <div class="rt-th__inner">
              <span>Narx</span>
              <Icon :icon="sortIcon('price')" class="rt-sort" :class="{ 'rt-sort--on': sortBy === 'price' }" />
            </div>
          </th>
          <th class="rt-th rt-th--sortable" @click="toggleSort('is_active')">
            <div class="rt-th__inner">
              <span>Holati</span>
              <Icon :icon="sortIcon('is_active')" class="rt-sort" :class="{ 'rt-sort--on': sortBy === 'is_active' }" />
            </div>
          </th>
          <th class="rt-th rt-th--actions">Amallar</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id" class="rt-row">
          <td class="rt-td rt-td--id">{{ item.id }}</td>

          <td class="rt-td">
            <div class="rt-person">
              <div class="rt-person__avatar">
                <img v-if="item.photo" :src="item.photo" :alt="item.name" class="rt-person__img" />
                <span v-else>{{ item.name?.charAt(0)?.toUpperCase() }}</span>
              </div>
              <div class="rt-person__info">
                <span class="rt-person__name">{{ item.name }}</span>
                <span class="rt-person__title">{{ item.title }}</span>
              </div>
            </div>
          </td>

          <td class="rt-td rt-td--contact">{{ item.contact || '—' }}</td>
          <td class="rt-td rt-td--price">{{ fmtPrice(item.price) }}</td>

          <td class="rt-td">
            <span :class="['rt-badge', item.is_active ? 'rt-badge--active' : 'rt-badge--inactive']">
              <span class="rt-badge__dot" />
              {{ item.is_active ? 'Faol' : 'Nofaol' }}
            </span>
          </td>

          <td class="rt-td rt-td--actions">
            <div class="rt-actions">
              <button class="rt-btn rt-btn--show"   @click="emit('show', item)"   title="Ko'rish">
                <Icon icon="mdi:eye-outline" />
              </button>
              <button class="rt-btn rt-btn--delete" @click="emit('delete', item)" title="O'chirish">
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
.rt-skeleton { display:flex; flex-direction:column; gap:1px; }
.rt-skeleton__row { display:grid; grid-template-columns:60px 1fr 120px 120px 100px 80px; gap:16px; padding:14px 16px; border-bottom:1px solid var(--border); align-items:center; }
.rt-skeleton__cell { height:14px; border-radius:8px; background:linear-gradient(90deg,var(--bg-elevated) 25%,rgba(255,255,255,.06) 50%,var(--bg-elevated) 75%); background-size:200% 100%; animation:shimmer 1.4s infinite; }
.rt-skeleton__cell--sm { width:40%; }
.rt-skeleton__cell--md { width:70%; }
.rt-skeleton__cell--lg { width:100%; }
@keyframes shimmer { 0%{background-position:200% 0} 100%{background-position:-200% 0} }

.rt-empty { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:64px 24px; gap:12px; text-align:center; }
.rt-empty__icon { width:72px; height:72px; border-radius:20px; background:var(--bg-elevated); border:1px solid var(--border); display:flex; align-items:center; justify-content:center; font-size:32px; color:var(--muted); margin-bottom:4px; }
.rt-empty__title { font-size:16px; font-weight:700; color:var(--text); }
.rt-empty__sub { font-size:13px; color:var(--muted); max-width:320px; line-height:1.5; }

.rt-table { width:100%; border-collapse:collapse; }
.rt-th { padding:12px 14px; text-align:left; font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:.07em; color:var(--muted); border-bottom:1px solid var(--border); white-space:nowrap; background:rgba(255,255,255,.02); }
.rt-th--sortable { cursor:pointer; user-select:none; transition:color .15s; }
.rt-th--sortable:hover { color:var(--text-secondary); }
.rt-th--actions { text-align:center; width:90px; }
.rt-th__inner { display:inline-flex; align-items:center; gap:5px; }
.rt-sort { font-size:14px; color:var(--border-hover); transition:color .15s; }
.rt-sort--on { color:#6366f1; }

.rt-row { transition:background .15s; }
.rt-row:hover { background:rgba(99,102,241,.04); }
.rt-row:not(:last-child) .rt-td { border-bottom:1px solid var(--border); }

.rt-td { padding:13px 14px; color:var(--text); font-size:14px; vertical-align:middle; }
.rt-td--id { color:var(--muted); font-size:12px; font-weight:600; width:60px; }
.rt-td--contact { color:var(--muted); font-size:13px; white-space:nowrap; }
.rt-td--price { font-size:13px; font-weight:600; color:var(--text); white-space:nowrap; }
.rt-td--actions { text-align:center; }

.rt-person { display:flex; align-items:center; gap:12px; }
.rt-person__avatar { width:38px; height:38px; border-radius:12px; background:linear-gradient(135deg,rgba(99,102,241,.25),rgba(59,130,246,.25)); border:1px solid rgba(99,102,241,.2); color:#818cf8; font-size:14px; font-weight:800; display:flex; align-items:center; justify-content:center; flex-shrink:0; overflow:hidden; }
.rt-person__img { width:100%; height:100%; object-fit:cover; }
.rt-person__info { display:flex; flex-direction:column; min-width:0; }
.rt-person__name { font-size:14px; font-weight:600; color:var(--text); white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.rt-person__title { font-size:12px; color:var(--muted); white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }

.rt-badge { display:inline-flex; align-items:center; gap:6px; padding:5px 12px; border-radius:20px; font-size:12px; font-weight:600; white-space:nowrap; }
.rt-badge__dot { width:6px; height:6px; border-radius:50%; flex-shrink:0; }
.rt-badge--active { background:rgba(16,185,129,.12); color:#34d399; border:1px solid rgba(16,185,129,.22); }
.rt-badge--active .rt-badge__dot { background:#34d399; box-shadow:0 0 6px rgba(52,211,153,.6); animation:pulse 2s infinite; }
.rt-badge--inactive { background:rgba(239,68,68,.1); color:#f87171; border:1px solid rgba(239,68,68,.2); }
.rt-badge--inactive .rt-badge__dot { background:#f87171; }
@keyframes pulse { 0%,100%{opacity:1;transform:scale(1)} 50%{opacity:.6;transform:scale(.85)} }

.rt-actions { display:flex; align-items:center; justify-content:center; gap:6px; }
.rt-btn { width:34px; height:34px; border:none; border-radius:10px; cursor:pointer; display:flex; align-items:center; justify-content:center; font-size:16px; transition:all .2s; color:white; }
.rt-btn:hover { transform:translateY(-2px); }
.rt-btn--show { background:linear-gradient(135deg,#0ea5e9,#06b6d4); box-shadow:0 2px 8px rgba(14,165,233,.3); }
.rt-btn--show:hover { box-shadow:0 6px 16px rgba(14,165,233,.45); }
.rt-btn--delete { background:linear-gradient(135deg,#ef4444,#f97316); box-shadow:0 2px 8px rgba(239,68,68,.3); }
.rt-btn--delete:hover { box-shadow:0 6px 16px rgba(239,68,68,.45); }

@media (max-width:768px) {
  .rt-td--contact, .rt-td--price { display:none; }
  .rt-person__avatar { display:none; }
}
</style>
