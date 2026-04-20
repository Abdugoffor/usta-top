<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import BaseCard from '@/shared/components/BaseCard.vue'
import CategoryForm from '../components/CategoryForm.vue'
import { useCategoryStore } from '../store/categoryStore'

const route = useRoute()
const router = useRouter()
const categoryStore = useCategoryStore()
const loading = ref(false)
const formData = ref({
  name: '',
  is_active: true
})

const loadItem = async () => {
  const item = await categoryStore.fetchCategory(route.params.id)
  formData.value = {
    name: item.name,
    is_active: Boolean(item.is_active)
  }
}

const updateItem = async (payload) => {
  try {
    loading.value = true
    await categoryStore.updateCategory(route.params.id, payload)

    await Swal.fire({
      icon: 'success',
      title: 'Category yangilandi'
    })

    router.push({ name: 'category-list' })
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
        <h2 class="page-title">Edit category</h2>
        <p class="page-subtitle">Category ma'lumotini yangilash</p>
      </div>
    </div>

    <BaseCard>
      <CategoryForm
        :initial-data="formData"
        :loading="loading"
        @submit="updateItem"
      />
    </BaseCard>
  </section>
</template>