<script setup>
import { Icon } from '@iconify/vue'
import BaseTableWrapper from '@/shared/components/BaseTableWrapper.vue'

defineProps({
  items: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['edit', 'delete'])
</script>

<template>
  <BaseTableWrapper>
    <table class="data-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>Status</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="!items.length">
          <td colspan="4" class="empty-row">Ma'lumot topilmadi</td>
        </tr>

        <tr v-for="item in items" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>
            <span :class="['badge', item.is_active ? 'badge--success' : 'badge--danger']">
              {{ item.is_active ? 'Active' : 'Inactive' }}
            </span>
          </td>
          <td>
            <div class="table-actions">
              <button class="action-btn action-btn--edit" @click="emit('edit', item)">
                <Icon icon="mdi:pencil-outline" />
              </button>
              <button class="action-btn action-btn--delete" @click="emit('delete', item)">
                <Icon icon="mdi:trash-can-outline" />
              </button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </BaseTableWrapper>
</template>