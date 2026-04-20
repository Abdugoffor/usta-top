<script setup>
import { Icon } from '@iconify/vue'
import BaseTableWrapper from '@/shared/components/BaseTableWrapper.vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  sortBy: { type: String, default: '' },
  sortOrder: { type: String, default: 'asc' },
})

const emit = defineEmits(['edit', 'delete', 'sort'])

const toggleSort = (field) => {
  if (props.sortBy === field) {
    emit('sort', { field, order: props.sortOrder === 'asc' ? 'desc' : 'asc' })
  } else {
    emit('sort', { field, order: 'asc' })
  }
}
</script>

<template>
  <BaseTableWrapper>

    <!-- Skeleton -->
    <template v-if="loading">
      <div class="ct-skeleton">
        <div v-for="i in 5" :key="i" class="ct-skeleton__row">
          <div class="ct-skeleton__cell ct-skeleton__cell--sm"></div>
          <div class="ct-skeleton__cell ct-skeleton__cell--lg"></div>
          <div class="ct-skeleton__cell ct-skeleton__cell--md"></div>
          <div class="ct-skeleton__cell ct-skeleton__cell--md"></div>
        </div>
      </div>
    </template>

    <!-- Empty state -->
    <template v-else-if="!items.length">
      <div class="ct-empty">
        <div class="ct-empty__icon">
          <Icon icon="mdi:folder-search-outline" />
        </div>
        <p class="ct-empty__title">Ma'lumot topilmadi</p>
        <p class="ct-empty__sub">Qidiruv so'rovingizni o'zgartiring yoki yangi kategoriya qo'shing</p>
      </div>
    </template>

    <!-- Table -->
    <table v-else class="ct-table">
      <thead>
        <tr>
          <th class="ct-th ct-th--num">#</th>
          <th class="ct-th ct-th--sortable" @click="toggleSort('name')">
            <div class="ct-th__inner">
              <span>Nomi</span>
              <span class="ct-sort-icon" :class="{ 'ct-sort-icon--active': sortBy === 'name' }">
                <Icon v-if="sortBy === 'name' && sortOrder === 'desc'" icon="mdi:arrow-down" />
                <Icon v-else-if="sortBy === 'name' && sortOrder === 'asc'" icon="mdi:arrow-up" />
                <Icon v-else icon="mdi:unfold-more-horizontal" />
              </span>
            </div>
          </th>
          <th class="ct-th ct-th--sortable" @click="toggleSort('is_active')">
            <div class="ct-th__inner">
              <span>Holati</span>
              <span class="ct-sort-icon" :class="{ 'ct-sort-icon--active': sortBy === 'is_active' }">
                <Icon v-if="sortBy === 'is_active' && sortOrder === 'desc'" icon="mdi:arrow-down" />
                <Icon v-else-if="sortBy === 'is_active' && sortOrder === 'asc'" icon="mdi:arrow-up" />
                <Icon v-else icon="mdi:unfold-more-horizontal" />
              </span>
            </div>
          </th>
          <th class="ct-th ct-th--actions">Amallar</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(item, index) in items"
          :key="item.id"
          class="ct-row"
        >
          <td class="ct-td ct-td--num">{{ index + 1 }}</td>

          <td class="ct-td">
            <div class="ct-name">
              <div class="ct-name__avatar">
                {{ item.name?.charAt(0)?.toUpperCase() }}
              </div>
              <span class="ct-name__text">{{ item.name }}</span>
            </div>
          </td>

          <td class="ct-td">
            <span :class="['ct-badge', item.is_active ? 'ct-badge--active' : 'ct-badge--inactive']">
              <span class="ct-badge__dot"></span>
              {{ item.is_active ? 'Faol' : 'Nofaol' }}
            </span>
          </td>

          <td class="ct-td ct-td--actions">
            <div class="ct-actions">
              <button class="ct-btn ct-btn--edit" @click="emit('edit', item)" title="Tahrirlash">
                <Icon icon="mdi:pencil-outline" />
              </button>
              <button class="ct-btn ct-btn--delete" @click="emit('delete', item)" title="O'chirish">
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
/* Skeleton */
.ct-skeleton {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.ct-skeleton__row {
  display: grid;
  grid-template-columns: 60px 1fr 140px 120px;
  gap: 16px;
  padding: 16px;
  border-bottom: 1px solid var(--border);
  align-items: center;
}

.ct-skeleton__cell {
  height: 14px;
  border-radius: 8px;
  background: linear-gradient(90deg, var(--bg-elevated) 25%, rgba(255,255,255,0.06) 50%, var(--bg-elevated) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}

.ct-skeleton__cell--sm { width: 32px; }
.ct-skeleton__cell--md { width: 80%; }
.ct-skeleton__cell--lg { width: 100%; }

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* Empty state */
.ct-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  gap: 12px;
  text-align: center;
}

.ct-empty__icon {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  color: var(--muted);
  margin-bottom: 4px;
}

.ct-empty__title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text);
}

.ct-empty__sub {
  font-size: 13px;
  color: var(--muted);
  max-width: 320px;
  line-height: 1.5;
}

/* Table */
.ct-table {
  width: 100%;
  border-collapse: collapse;
}

/* Header */
.ct-th {
  padding: 12px 16px;
  text-align: left;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.07em;
  color: var(--muted);
  border-bottom: 1px solid var(--border);
  white-space: nowrap;
  background: rgba(255, 255, 255, 0.02);
}

.ct-th--num {
  width: 60px;
  text-align: center;
}

.ct-th--sortable {
  cursor: pointer;
  user-select: none;
  transition: color 0.15s;
}

.ct-th--sortable:hover {
  color: var(--text-secondary);
}

.ct-th__inner {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.ct-sort-icon {
  font-size: 14px;
  color: var(--border-hover);
  display: flex;
  align-items: center;
  transition: color 0.15s;
}

.ct-sort-icon--active {
  color: #6366f1;
}

.ct-th--actions {
  width: 120px;
  text-align: center;
}

/* Row */
.ct-row {
  transition: background 0.15s ease;
  animation: rowIn 0.25s ease both;
}

.ct-row:nth-child(n) { animation-delay: calc(var(--i, 0) * 40ms); }

@keyframes rowIn {
  from { opacity: 0; transform: translateX(-6px); }
  to { opacity: 1; transform: translateX(0); }
}

.ct-row:hover {
  background: rgba(99, 102, 241, 0.04);
}

.ct-row:not(:last-child) .ct-td {
  border-bottom: 1px solid var(--border);
}

/* Cell */
.ct-td {
  padding: 14px 16px;
  color: var(--text);
  font-size: 14px;
  vertical-align: middle;
}

.ct-td--num {
  text-align: center;
  color: var(--muted);
  font-size: 12px;
  font-weight: 600;
}

.ct-td--actions {
  text-align: center;
}

/* Name with avatar */
.ct-name {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ct-name__avatar {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.25), rgba(59, 130, 246, 0.25));
  border: 1px solid rgba(99, 102, 241, 0.2);
  color: #818cf8;
  font-size: 13px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ct-name__text {
  font-weight: 500;
  color: var(--text);
}

/* Badge */
.ct-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}

.ct-badge__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.ct-badge--active {
  background: rgba(16, 185, 129, 0.12);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.22);
}

.ct-badge--active .ct-badge__dot {
  background: #34d399;
  box-shadow: 0 0 6px rgba(52, 211, 153, 0.6);
  animation: pulse 2s infinite;
}

.ct-badge--inactive {
  background: rgba(239, 68, 68, 0.1);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.ct-badge--inactive .ct-badge__dot {
  background: #f87171;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(0.85); }
}

/* Action buttons */
.ct-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.ct-btn {
  width: 34px;
  height: 34px;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  transition: all 0.2s ease;
  color: white;
}

.ct-btn:hover {
  transform: translateY(-2px);
}

.ct-btn:active {
  transform: translateY(0);
}

.ct-btn--edit {
  background: linear-gradient(135deg, #3b82f6, #6366f1);
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
}

.ct-btn--edit:hover {
  box-shadow: 0 6px 16px rgba(99, 102, 241, 0.45);
}

.ct-btn--delete {
  background: linear-gradient(135deg, #ef4444, #f97316);
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.3);
}

.ct-btn--delete:hover {
  box-shadow: 0 6px 16px rgba(239, 68, 68, 0.45);
}

/* Responsive */
@media (max-width: 640px) {
  .ct-skeleton__row {
    grid-template-columns: 40px 1fr 100px 80px;
  }

  .ct-td, .ct-th {
    padding: 12px 10px;
  }

  .ct-name__avatar {
    display: none;
  }
}
</style>
