<script setup>
import { reactive, watch } from 'vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSwitch from '@/shared/components/BaseSwitch.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import { useActiveLanguages } from '@/shared/composables/useActiveLanguages'

const props = defineProps({
  initialData: {
    type: Object,
    default: () => ({ name: {}, is_active: true })
  },
  loading: { type: Boolean, default: false }
})

const emit = defineEmits(['submit'])
const { languages } = useActiveLanguages()

const form = reactive({
  name: { default: '' },
  is_active: true
})

watch(
  () => props.initialData,
  (value) => {
    const incoming = value?.name
    if (incoming && typeof incoming === 'object') {
      form.name = { ...incoming }
    } else {
      form.name = { default: '' }
    }
    form.is_active = Boolean(value?.is_active ?? true)
  },
  { immediate: true, deep: true }
)

const handleSubmit = () => {
  emit('submit', {
    name: { ...form.name },
    is_active: form.is_active
  })
}
</script>

<template>
  <form class="cf-form" @submit.prevent="handleSubmit">

    <div class="cf-lang-section">
      <div class="cf-lang-label">
        <span class="cf-lang-badge cf-lang-badge--default">default</span>
        Standart nom <span class="cf-required">*</span>
      </div>
      <BaseInput
        v-model="form.name.default"
        placeholder="Masalan: Elektronika"
      />
    </div>

    <div
      v-for="lang in languages"
      :key="lang.name"
      class="cf-lang-section"
    >
      <div class="cf-lang-label">
        <span class="cf-lang-badge">{{ lang.name }}</span>
        {{ lang.description || lang.name }}
      </div>
      <BaseInput
        v-model="form.name[lang.name]"
        :placeholder="`${lang.description || lang.name} tilidagi tarjima`"
      />
    </div>

    <BaseSwitch v-model="form.is_active" label="Faol holati" />

    <div class="cf-actions">
      <BaseButton type="submit">
        {{ loading ? 'Saqlanmoqda...' : 'Saqlash' }}
      </BaseButton>
    </div>
  </form>
</template>

<style scoped>
.cf-form { display: flex; flex-direction: column; gap: 20px; }

.cf-lang-section { display: flex; flex-direction: column; gap: 8px; }

.cf-lang-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
}

.cf-lang-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 9px;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: .05em;
  background: rgba(99, 102, 241, .15);
  color: #818cf8;
  border: 1px solid rgba(99, 102, 241, .25);
}

.cf-lang-badge--default {
  background: rgba(16, 185, 129, .12);
  color: #34d399;
  border-color: rgba(16, 185, 129, .22);
}

.cf-required { color: #f87171; }

.cf-actions { display: flex; justify-content: flex-end; margin-top: 8px; }
</style>
