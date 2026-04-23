<script setup>
import { reactive, watch } from 'vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSwitch from '@/shared/components/BaseSwitch.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import { useActiveLanguages } from '@/shared/composables/useActiveLanguages'

const props = defineProps({
  initialData: {
    type: Object,
    default: () => ({ slug: '', name: {}, is_active: true })
  },
  loading:   { type: Boolean, default: false },
  editMode:  { type: Boolean, default: false },
})

const emit = defineEmits(['submit'])
const { languages } = useActiveLanguages()

const form = reactive({
  slug: '',
  name: { default: '' },
  is_active: true
})

watch(
  () => props.initialData,
  (value) => {
    form.slug = value?.slug ?? ''
    const incoming = value?.name
    form.name = (incoming && typeof incoming === 'object') ? { ...incoming } : { default: '' }
    form.is_active = Boolean(value?.is_active ?? true)
  },
  { immediate: true, deep: true }
)

const handleSubmit = () => {
  emit('submit', {
    slug: form.slug.trim(),
    name: { ...form.name },
    is_active: form.is_active
  })
}
</script>

<template>
  <form class="tf-form" @submit.prevent="handleSubmit">

    <div class="tf-slug-wrap">
      <BaseInput
        v-model="form.slug"
        label="Kalit (slug) *"
        placeholder="Masalan: home_title, about_description"
        :readonly="editMode"
      />
      <p class="tf-slug-hint">
        Frontend tarjimada ishlatiladi: <code>key="{{ form.slug || 'home_title' }}"</code>
      </p>
    </div>

    <div class="tf-divider">
      <span>Tarjimalar</span>
    </div>

    <div class="tf-lang-section">
      <div class="tf-lang-label">
        <span class="tf-lang-badge tf-lang-badge--default">default</span>
        Standart matn <span class="tf-required">*</span>
      </div>
      <BaseInput
        v-model="form.name.default"
        placeholder="Standart (fallback) matn"
      />
    </div>

    <div
      v-for="lang in languages"
      :key="lang.name"
      class="tf-lang-section"
    >
      <div class="tf-lang-label">
        <span class="tf-lang-badge">{{ lang.name }}</span>
        {{ lang.description || lang.name }}
      </div>
      <BaseInput
        v-model="form.name[lang.name]"
        :placeholder="`${lang.description || lang.name} tilidagi tarjima`"
      />
    </div>

    <BaseSwitch v-model="form.is_active" label="Faol holati" />

    <div class="tf-actions">
      <BaseButton type="submit">
        {{ loading ? 'Saqlanmoqda...' : 'Saqlash' }}
      </BaseButton>
    </div>
  </form>
</template>

<style scoped>
.tf-form { display: flex; flex-direction: column; gap: 20px; }

.tf-slug-wrap { display: flex; flex-direction: column; gap: 6px; }
.tf-slug-hint { font-size: 12px; color: var(--muted); padding: 0 2px; }
.tf-slug-hint code { background: rgba(99,102,241,.15); color: #818cf8; border-radius: 5px; padding: 1px 6px; font-size: 11px; }

.tf-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--muted);
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: .08em;
}
.tf-divider::before,
.tf-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border);
}

.tf-lang-section { display: flex; flex-direction: column; gap: 8px; }

.tf-lang-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
}

.tf-lang-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 9px;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: .05em;
  background: rgba(99,102,241,.15);
  color: #818cf8;
  border: 1px solid rgba(99,102,241,.25);
}
.tf-lang-badge--default {
  background: rgba(16,185,129,.12);
  color: #34d399;
  border-color: rgba(16,185,129,.22);
}

.tf-required { color: #f87171; }
.tf-actions { display: flex; justify-content: flex-end; margin-top: 8px; }
</style>
