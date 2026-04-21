<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import BaseCard from '@/shared/components/BaseCard.vue'
import LanguageForm from '../components/LanguageForm.vue'
import { useLanguageStore } from '../store/languageStore'

const route   = useRoute()
const router  = useRouter()
const store   = useLanguageStore()
const loading = ref(false)
const formData = ref({ name: '', description: '', is_active: true })

onMounted(async () => {
  const item = await store.fetchLanguage(route.params.id)
  formData.value = { name: item.name, description: item.description ?? '', is_active: Boolean(item.is_active) }
})

const update = async (payload) => {
  try {
    loading.value = true
    await store.updateLanguage(route.params.id, payload)
    await Swal.fire({ icon: 'success', title: 'Til yangilandi', background: '#0f172a', color: '#f1f5f9', timer: 1800, showConfirmButton: false })
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
        <h2 class="page-title">Tilni tahrirlash</h2>
        <p class="page-subtitle">Til ma'lumotini yangilash</p>
      </div>
    </div>
    <BaseCard>
      <LanguageForm :initial-data="formData" :loading="loading" @submit="update" />
    </BaseCard>
  </section>
</template>
