<script setup>
import { reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import BaseCard from '@/shared/components/BaseCard.vue'
import BaseButton from '@/shared/components/BaseButton.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import CategoryTable from '../components/CategoryTable.vue'
import { useCategoryStore } from '../store/categoryStore'

const router = useRouter()
const categoryStore = useCategoryStore()

const filters = reactive({ name: '', page: 1 })

const loadData = async () => {
  await categoryStore.fetchCategories({ name: filters.name, page: filters.page })
}

const goCreate = () => {
  router.push({ name: 'category-create' })
}

const goEdit = (item) => {
  router.push({ name: 'category-edit', params: { id: item.id } })
}

const deleteItem = async (item) => {
  const result = await Swal.fire({
    title: "Ochirishni xohlaysizmi?",
    text: item.name,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Ha',
    cancelButtonText: "Yo'q",
    background: '#0f172a',
    color: '#f1f5f9',
    confirmButtonColor: '#ef4444'
  })

  if (!result.isConfirmed) return

  await categoryStore.removeCategory(item.id)
  await loadData()

  Swal.fire({
    icon: 'success',
    title: 'Ochirildi',
    background: '#0f172a',
    color: '#f1f5f9'
  })
}

const changePage = async (page) => {
  if (page < 1 || page > categoryStore.lastPage) return
  filters.page = page
  await loadData()
}

onMounted(loadData)
</script>

<template>
  <section class="page-section">
    <div class="page-title-row">
      <div>
        <h2 class="page-title">Category CRUD</h2>
        <p class="page-subtitle">Category larni boshqarish bo'limi</p>
      </div>
      <BaseButton @click="goCreate">+ Add category</BaseButton>
    </div>

    <BaseCard>
      <div class="filters-grid">
        <BaseInput
          v-model="filters.name"
          label="Search"
          placeholder="Category nomi bo'yicha qidirish"
        />
        <div class="filter-btn-wrap">
          <BaseButton @click="loadData">Qidirish</BaseButton>
        </div>
      </div>
    </BaseCard>

    <BaseCard>
      <CategoryTable
        :items="categoryStore.items"
        @edit="goEdit"
        @delete="deleteItem"
      />

      <div class="pagination" v-if="categoryStore.lastPage > 1">
        <button class="pagination__btn" @click="changePage(filters.page - 1)">Prev</button>
        <span class="pagination__info">{{ categoryStore.currentPage }} / {{ categoryStore.lastPage }}</span>
        <button class="pagination__btn" @click="changePage(filters.page + 1)">Next</button>
      </div>
    </BaseCard>
  </section>
</template>
