<script setup>
import { reactive, watch } from 'vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseSwitch from '@/shared/components/BaseSwitch.vue'
import BaseButton from '@/shared/components/BaseButton.vue'

const props = defineProps({
  initialData: {
    type: Object,
    default: () => ({
      name: '',
      is_active: true
    })
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit'])

const form = reactive({
  name: '',
  is_active: true
})

watch(
  () => props.initialData,
  (value) => {
    form.name = value?.name ?? ''
    form.is_active = Boolean(value?.is_active)
  },
  { immediate: true, deep: true }
)

const handleSubmit = () => {
  emit('submit', {
    name: form.name,
    is_active: form.is_active
  })
}
</script>

<template>
  <form class="category-form" @submit.prevent="handleSubmit">
    <BaseInput
      v-model="form.name"
      label="Category name"
      placeholder="Masalan: Elektronika"
    />

    <BaseSwitch
      v-model="form.is_active"
      label="Active holati"
    />

    <div class="form-actions">
      <BaseButton type="submit">
        {{ loading ? 'Saqlanmoqda...' : 'Saqlash' }}
      </BaseButton>
    </div>
  </form>
</template>