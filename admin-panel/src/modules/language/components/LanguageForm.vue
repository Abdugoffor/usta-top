<script setup>
import { reactive, watch } from 'vue'
import BaseInput  from '@/shared/components/BaseInput.vue'
import BaseSwitch from '@/shared/components/BaseSwitch.vue'
import BaseButton from '@/shared/components/BaseButton.vue'

const props = defineProps({
  initialData: { type: Object,  default: () => ({ name: '', description: '', is_active: true }) },
  loading:     { type: Boolean, default: false },
})

const emit = defineEmits(['submit'])

const form = reactive({ name: '', description: '', is_active: true })

watch(() => props.initialData, (v) => {
  form.name        = v?.name        ?? ''
  form.description = v?.description ?? ''
  form.is_active   = v?.is_active   ?? true
}, { immediate: true, deep: true })

const handleSubmit = () => emit('submit', {
  name:        form.name,
  description: form.description || undefined,
  is_active:   form.is_active,
})
</script>

<template>
  <form class="lf-form" @submit.prevent="handleSubmit">
    <BaseInput v-model="form.name" label="Til nomi" placeholder="Masalan: O'zbek tili" />
    <BaseInput v-model="form.description" label="Tavsif" placeholder="Til haqida qisqacha ma'lumot" />
    <BaseSwitch v-model="form.is_active" label="Faol holati" />
    <div class="lf-actions">
      <BaseButton type="submit">{{ loading ? 'Saqlanmoqda...' : 'Saqlash' }}</BaseButton>
    </div>
  </form>
</template>

<style scoped>
.lf-form { display: flex; flex-direction: column; gap: 20px; }
.lf-actions { display: flex; justify-content: flex-end; margin-top: 8px; }
</style>
