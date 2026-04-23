<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import BaseCard from '@/shared/components/BaseCard.vue'
import TranslationForm from '../components/TranslationForm.vue'
import { useTranslationStore } from '../store/translationStore'

const route = useRoute()
const router = useRouter()
const store = useTranslationStore()
const loading = ref(false)
const formData = ref({ slug: '', name: { default: '' }, is_active: true })

const loadItem = async () => {
  const item = await store.fetchTranslation(route.params.id)
  formData.value = {
    slug:      item.slug,
    name:      item.name && typeof item.name === 'object' ? { ...item.name } : { default: '' },
    is_active: Boolean(item.is_active)
  }
}

const save = async (payload) => {
  try {
    loading.value = true
    await store.updateTranslation(route.params.id, payload)

    await Swal.fire({
      icon: 'success',
      title: 'Tarjima yangilandi',
      background: '#0f172a',
      color: '#f1f5f9',
      timer: 2000,
      showConfirmButton: false
    })

    router.push({ name: 'translation-list' })
  } catch (e) {
    const msg = e?.response?.data?.error || 'Xatolik yuz berdi'
    Swal.fire({ icon: 'error', title: 'Xatolik', text: msg, background: '#0f172a', color: '#f1f5f9' })
  } finally {
    loading.value = false
  }
}

onMounted(loadItem)
</script>

<template>
  <section class="page-section">
    <div class="page-title-row">
      <div>
        <h2 class="page-title">Tarjimani tahrirlash</h2>
        <p class="page-subtitle">Tarjima ma'lumotlarini yangilash</p>
      </div>
    </div>

    <BaseCard>
      <TranslationForm
        :initial-data="formData"
        :loading="loading"
        :edit-mode="true"
        @submit="save"
      />
    </BaseCard>
  </section>
</template>
