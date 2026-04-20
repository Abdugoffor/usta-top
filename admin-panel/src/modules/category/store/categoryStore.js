import { defineStore } from 'pinia'
import * as api from '../api/categoryApi'

export const useCategoryStore = defineStore('category', {
  state: () => ({
    items: [],
    item: null,
    loading: false,
    total: 0,
    currentPage: 1,
    perPage: 10,
    lastPage: 1
  }),

  actions: {
    async fetchCategories(params = {}) {
      this.loading = true
      try {
        const { data } = await api.getCategories(params)
        this.items = data.data ?? []
        this.total = data.meta?.total ?? 0
        this.currentPage = data.meta?.page ?? 1
        this.perPage = data.meta?.limit ?? 10
        this.lastPage = data.meta?.total_pages ?? 1
      } finally {
        this.loading = false
      }
    },

    async fetchCategory(id) {
      this.loading = true
      try {
        const { data } = await api.getCategory(id)
        this.item = data
        return this.item
      } finally {
        this.loading = false
      }
    },

    async createCategory(payload) {
      return await api.createCategory(payload)
    },

    async updateCategory(id, payload) {
      return await api.updateCategory(id, payload)
    },

    async removeCategory(id) {
      return await api.deleteCategory(id)
    }
  }
})
