<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Swal from 'sweetalert2'
import { Icon } from '@iconify/vue'
import BaseCard from '@/shared/components/BaseCard.vue'
import CountryForm from '../components/CountryForm.vue'
import { useCountryStore } from '../store/countryStore'

const router = useRouter()
const route  = useRoute()
const store  = useCountryStore()
const loading = ref(false)

const parentId = route.query.parent_id ? Number(route.query.parent_id) : null

const save = async (payload) => {
  try {
    loading.value = true
    await store.createCountry(payload)
    await Swal.fire({
      icon: 'success',
      title: 'Muvaffaqiyatli yaratildi',
      background: '#0f172a',
      color: '#f1f5f9',
      timer: 2000,
      showConfirmButton: false,
    })
    router.push({ name: 'country-list' })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="page-section">
    <div class="cp-header">
      <button class="cp-back" @click="router.push({ name: 'country-list' })">
        <Icon icon="mdi:arrow-left" />
      </button>
      <div>
        <h2 class="page-title">Yangi mamlakat</h2>
        <p class="page-subtitle">Mamlakat yoki hudud qo'shish</p>
      </div>
    </div>

    <BaseCard>
      <CountryForm
        :initial-data="{ name: '', parent_id: parentId, is_active: true }"
        :loading="loading"
        @submit="save"
      />
    </BaseCard>
  </section>
</template>

<style scoped>
.cp-header { display:flex; align-items:center; gap:16px; }
.cp-back { width:40px; height:40px; border-radius:12px; border:1px solid var(--border); background:var(--bg-elevated); color:var(--text-secondary); display:flex; align-items:center; justify-content:center; cursor:pointer; font-size:20px; transition:all .2s; flex-shrink:0; }
.cp-back:hover { border-color:#6366f1; color:#6366f1; }
</style>
