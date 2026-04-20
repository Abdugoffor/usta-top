<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import BaseCard from '@/shared/components/BaseCard.vue'
import CategoryForm from '../components/CategoryForm.vue'
import { useCategoryStore } from '../store/categoryStore'

const router = useRouter()
const categoryStore = useCategoryStore()
const loading = ref(false)

const saveCategory = async (payload) => {
  try {
    loading.value = true
    await categoryStore.createCategory(payload)

    await Swal.fire({
      icon: 'success',
      title: 'Category yaratildi'
    })

    router.push({ name: 'category-list' })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="page-section">
    <div class="page-title-row">
      <div>
        <h2 class="page-title">Create category</h2>
        <p class="page-subtitle">Yangi category qo‘shish</p>
      </div>
    </div>

    <BaseCard>
      <CategoryForm :loading="loading" @submit="saveCategory" />
    </BaseCard>
  </section>
</template>