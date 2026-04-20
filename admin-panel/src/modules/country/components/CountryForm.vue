<script setup>
import { reactive, computed, watch, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSwitch from '@/shared/components/BaseSwitch.vue'
import { useCountryStore } from '../store/countryStore'

const props = defineProps({
  initialData: {
    type: Object,
    default: () => ({ name: '', parent_id: null, is_active: true }),
  },
  loading: { type: Boolean, default: false },
  editId: { type: [Number, String], default: null },
})

const emit = defineEmits(['submit'])
const store = useCountryStore()

const form = reactive({ name: '', parent_id: null, is_active: true })

watch(() => props.initialData, (val) => {
  form.name      = val?.name ?? ''
  form.parent_id = val?.parent_id ?? null
  form.is_active = val?.is_active !== undefined ? Boolean(val.is_active) : true
}, { immediate: true, deep: true })

onMounted(() => store.fetchFlatList())

const parentOptions = computed(() =>
  store.flatList.filter(n => String(n.id) !== String(props.editId))
)

const handleSubmit = () => {
  emit('submit', {
    name:      form.name,
    parent_id: form.parent_id ? Number(form.parent_id) : null,
    is_active: form.is_active,
  })
}
</script>

<template>
  <form class="cf-form" @submit.prevent="handleSubmit">

    <BaseInput
      v-model="form.name"
      label="Nomi"
      placeholder="Masalan: O'zbekiston"
    />

    <div class="cf-field">
      <label class="cf-label">
        <Icon icon="mdi:map-marker-path" />
        Parent (yuqori daraja)
      </label>
      <div class="cf-select-wrap">
        <Icon icon="mdi:chevron-down" class="cf-select-arrow" />
        <select v-model="form.parent_id" class="cf-select">
          <option :value="null">— Yuqori daraja (Root) —</option>
          <option
            v-for="opt in parentOptions"
            :key="opt.id"
            :value="opt.id"
          >
            {{ '—'.repeat(opt._depth) }} {{ opt.name }}
          </option>
        </select>
      </div>
    </div>

    <BaseSwitch v-model="form.is_active" label="Faol holati" />

    <div class="cf-actions">
      <button type="submit" class="cf-btn cf-btn--primary" :disabled="loading">
        <Icon icon="mdi:content-save-outline" />
        {{ loading ? 'Saqlanmoqda...' : 'Saqlash' }}
      </button>
    </div>

  </form>
</template>

<style scoped>
.cf-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 520px;
}

.cf-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.cf-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
}

.cf-select-wrap {
  position: relative;
}

.cf-select-arrow {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 18px;
  color: var(--muted);
  pointer-events: none;
}

.cf-select {
  width: 100%;
  padding: 11px 40px 11px 14px;
  background: var(--bg-elevated);
  border: 1.5px solid var(--border);
  border-radius: 12px;
  color: var(--text);
  font-size: 14px;
  font-family: inherit;
  outline: none;
  appearance: none;
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.cf-select:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.15);
}

.cf-select option {
  background: #1e293b;
  color: var(--text);
}

.cf-actions {
  display: flex;
  gap: 12px;
  padding-top: 4px;
}

.cf-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 11px 24px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 700;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cf-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.cf-btn--primary {
  background: linear-gradient(135deg, #6366f1, #3b82f6);
  color: white;
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.35);
}

.cf-btn--primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 28px rgba(99, 102, 241, 0.45);
}
</style>
