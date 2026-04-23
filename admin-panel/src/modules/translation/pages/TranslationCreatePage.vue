<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import BaseCard from '@/shared/components/BaseCard.vue'
import TranslationForm from '../components/TranslationForm.vue'
import { useTranslationStore } from '../store/translationStore'

const router = useRouter()
const store = useTranslationStore()
const loading = ref(false)

const save = async (payload) => {
  try {
    loading.value = true
    await store.createTranslation(payload)

    await Swal.fire({
      icon: 'success',
      title: 'Tarjima yaratildi',
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
</script>

<template>
  <section class="page-section">
    <div class="page-title-row">
      <div>
        <h2 class="page-title">Yangi tarjima</h2>
        <p class="page-subtitle">Sayt matni uchun yangi kalit va tarjima qo'shish</p>
      </div>
    </div>

    <BaseCard>
      <TranslationForm :loading="loading" @submit="save" />
    </BaseCard>
  </section>
</template>
