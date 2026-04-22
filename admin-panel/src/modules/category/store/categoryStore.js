import { defineStore } from 'pinia'
import * as api from '../api/categoryApi'

export const useCategoryStore = defineStore('category', {
  state: () => ({
    items: [],
    item: null,
    loading: false,
    total: 0,
    hasMore: false,
    nextCursor: '',
    limit: 10,
    cursorStack: [''],
    stackIndex: 0,
  }),

  getters: {
    hasPrev: (s) => s.stackIndex > 0,
    hasNext: (s) => s.hasMore,
    currentPage: (s) => s.stackIndex + 1,
  },

  actions: {
    async fetchCategories(params = {}) {
      this.loading = true
      try {
        const cursor = this.cursorStack[this.stackIndex]
        const { data } = await api.getCategories({
          ...params,
          limit: this.limit,
          cursor: cursor || undefined,
        })
        this.items      = data.data ?? []
        this.total      = data.meta?.total ?? 0
        this.hasMore    = data.meta?.has_more ?? false
        this.nextCursor = data.meta?.next_cursor ?? ''
      } finally {
        this.loading = false
      }
    },

    async goNext(params = {}) {
      if (!this.hasMore || !this.nextCursor) return
      this.cursorStack = [...this.cursorStack.slice(0, this.stackIndex + 1), this.nextCursor]
      this.stackIndex++
      await this.fetchCategories(params)
    },

    async goPrev(params = {}) {
      if (this.stackIndex <= 0) return
      this.stackIndex--
      await this.fetchCategories(params)
    },

    resetCursor() {
      this.cursorStack = ['']
      this.stackIndex  = 0
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

    async createCategory(payload)     { return api.createCategory(payload) },
    async updateCategory(id, payload) { return api.updateCategory(id, payload) },
    async removeCategory(id)          { return api.deleteCategory(id) },
  },
})
