<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import BaseCard from '@/shared/components/BaseCard.vue'
import LanguageForm from '../components/LanguageForm.vue'
import { useLanguageStore } from '../store/languageStore'

const router = useRouter()
const store  = useLanguageStore()
const loading = ref(false)

const save = async (payload) => {
  try {
    loading.value = true
    await store.createLanguage(payload)
    await Swal.fire({ icon: 'success', title: 'Til yaratildi', background: '#0f172a', color: '#f1f5f9', timer: 1800, showConfirmButton: false })
    router.push({ name: 'language-list' })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="page-section">
    <div class="page-title-row">
      <div>
        <h2 class="page-title">Yangi til</h2>
        <p class="page-subtitle">Yangi til qo'shish</p>
      </div>
    </div>
    <BaseCard>
      <LanguageForm :loading="loading" @submit="save" />
    </BaseCard>
  </section>
</template>
