<script setup>
import { ref, computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps({
  node:  { type: Object,  required: true },
  depth: { type: Number,  default: 0 },
})

const emit = defineEmits(['edit', 'delete', 'add-child'])

const open = ref(props.depth === 0)

const hasChildren = computed(() => props.node.children?.length > 0)
</script>

<template>
  <div class="ctr-wrap">

    <!-- Row -->
    <div
      class="ctr-row"
      :style="{ '--depth': depth }"
      :class="{ 'ctr-row--child': depth > 0 }"
    >
      <!-- Indent + toggle -->
      <div class="ctr-indent">
        <span v-for="d in depth" :key="d" class="ctr-line" />
        <button
          v-if="hasChildren"
          class="ctr-toggle"
          :class="{ 'ctr-toggle--open': open }"
          @click="open = !open"
        >
          <Icon icon="mdi:chevron-right" />
        </button>
        <span v-else class="ctr-leaf">
          <Icon icon="mdi:circle-small" />
        </span>
      </div>

      <!-- Avatar + name -->
      <div class="ctr-name">
        <div
          class="ctr-avatar"
          :style="depth === 0
            ? 'background: linear-gradient(135deg,rgba(99,102,241,.25),rgba(59,130,246,.25)); color:#818cf8'
            : 'background: rgba(255,255,255,0.04); color: var(--muted)'"
        >
          {{ node.name?.charAt(0)?.toUpperCase() }}
        </div>
        <div>
          <span class="ctr-name__text">{{ node.name }}</span>
          <span v-if="node.children?.length" class="ctr-children-count">
            {{ node.children.length }} ta ichki
          </span>
        </div>
      </div>

      <!-- Badge -->
      <div class="ctr-status">
        <span :class="['ctr-badge', node.is_active ? 'ctr-badge--active' : 'ctr-badge--inactive']">
          <span class="ctr-badge__dot" />
          {{ node.is_active ? 'Faol' : 'Nofaol' }}
        </span>
      </div>

      <!-- Actions -->
      <div class="ctr-actions">
        <button class="ctr-btn ctr-btn--child" @click="emit('add-child', node)" title="Ichki qo'shish">
          <Icon icon="mdi:plus" />
        </button>
        <button class="ctr-btn ctr-btn--edit" @click="emit('edit', node)" title="Tahrirlash">
          <Icon icon="mdi:pencil-outline" />
        </button>
        <button class="ctr-btn ctr-btn--delete" @click="emit('delete', node)" title="O'chirish">
          <Icon icon="mdi:trash-can-outline" />
        </button>
      </div>
    </div>

    <!-- Children (recursive) -->
    <Transition name="tree-expand">
      <div v-if="open && hasChildren" class="ctr-children">
        <CountryTreeRow
          v-for="child in node.children"
          :key="child.id"
          :node="child"
          :depth="depth + 1"
          @edit="emit('edit', $event)"
          @delete="emit('delete', $event)"
          @add-child="emit('add-child', $event)"
        />
      </div>
    </Transition>

  </div>
</template>

<style scoped>
.ctr-wrap {
  display: flex;
  flex-direction: column;
}

/* Row */
.ctr-row {
  display: grid;
  grid-template-columns: auto 1fr auto auto;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
  transition: background 0.15s;
  animation: rowIn 0.2s ease both;
}

.ctr-row:hover {
  background: rgba(99, 102, 241, 0.04);
}

.ctr-row--child {
  background: rgba(255, 255, 255, 0.01);
}

@keyframes rowIn {
  from { opacity: 0; transform: translateX(-4px); }
  to   { opacity: 1; transform: translateX(0); }
}

/* Indent line + toggle */
.ctr-indent {
  display: flex;
  align-items: center;
  gap: 0;
  flex-shrink: 0;
}

.ctr-line {
  display: inline-block;
  width: 20px;
  border-left: 1.5px solid var(--border);
  height: 36px;
  margin-left: 10px;
  flex-shrink: 0;
}

.ctr-toggle {
  width: 26px;
  height: 26px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--bg-elevated);
  color: var(--muted);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 16px;
  flex-shrink: 0;
}

.ctr-toggle:hover {
  border-color: #6366f1;
  color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.ctr-toggle--open .iconify {
  transform: rotate(90deg);
}

.ctr-toggle .iconify {
  transition: transform 0.2s;
}

.ctr-leaf {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--border-hover);
  font-size: 20px;
  flex-shrink: 0;
}

/* Name */
.ctr-name {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.ctr-avatar {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.ctr-name__text {
  font-size: 14px;
  font-weight: 500;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ctr-children-count {
  display: block;
  font-size: 11px;
  color: var(--muted);
  margin-top: 1px;
}

/* Badge */
.ctr-status {
  flex-shrink: 0;
}

.ctr-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}

.ctr-badge__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.ctr-badge--active {
  background: rgba(16, 185, 129, 0.12);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.22);
}

.ctr-badge--active .ctr-badge__dot {
  background: #34d399;
  box-shadow: 0 0 5px rgba(52, 211, 153, 0.5);
  animation: pulse 2s infinite;
}

.ctr-badge--inactive {
  background: rgba(239, 68, 68, 0.1);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.ctr-badge--inactive .ctr-badge__dot {
  background: #f87171;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50%       { opacity: 0.5; transform: scale(0.8); }
}

/* Actions */
.ctr-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.ctr-btn {
  width: 30px;
  height: 30px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: white;
  transition: all 0.2s;
}

.ctr-btn:hover { transform: translateY(-2px); }
.ctr-btn:active { transform: translateY(0); }

.ctr-btn--child {
  background: linear-gradient(135deg, #10b981, #059669);
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.3);
}

.ctr-btn--child:hover { box-shadow: 0 4px 14px rgba(16, 185, 129, 0.45); }

.ctr-btn--edit {
  background: linear-gradient(135deg, #3b82f6, #6366f1);
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
}

.ctr-btn--edit:hover { box-shadow: 0 4px 14px rgba(99, 102, 241, 0.45); }

.ctr-btn--delete {
  background: linear-gradient(135deg, #ef4444, #f97316);
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.3);
}

.ctr-btn--delete:hover { box-shadow: 0 4px 14px rgba(239, 68, 68, 0.45); }

/* Children container */
.ctr-children {
  display: flex;
  flex-direction: column;
}

/* Transition */
.tree-expand-enter-active,
.tree-expand-leave-active {
  transition: all 0.22s ease;
  overflow: hidden;
}

.tree-expand-enter-from,
.tree-expand-leave-to {
  opacity: 0;
  max-height: 0;
}

.tree-expand-enter-to,
.tree-expand-leave-from {
  opacity: 1;
  max-height: 2000px;
}

/* Responsive */
@media (max-width: 640px) {
  .ctr-avatar { display: none; }
  .ctr-children-count { display: none; }
  .ctr-row { gap: 8px; padding: 10px 12px; }
}
</style>
