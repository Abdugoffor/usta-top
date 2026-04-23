import { defineStore } from 'pinia'
import * as api from '../api/translationApi'

export const useTranslationStore = defineStore('translation', {
  state: () => ({
    items: [],
    item: null,
    loading: false,
    total: 0,
    hasMore: false,
    nextCursor: '',
    limit: 15,
    cursorStack: [''],
    stackIndex: 0,
  }),

  getters: {
    hasPrev: (s) => s.stackIndex > 0,
    hasNext: (s) => s.hasMore,
    currentPage: (s) => s.stackIndex + 1,
  },

  actions: {
    async fetchTranslations(params = {}) {
      this.loading = true
      try {
        const cursor = this.cursorStack[this.stackIndex]
        const { data } = await api.getTranslations({
          ...params,
          limit: this.limit,
          cursor: cursor || undefined,
        })
        this.items      = data.data ?? []
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
      await this.fetchTranslations(params)
    },

    async goPrev(params = {}) {
      if (this.stackIndex <= 0) return
      this.stackIndex--
      await this.fetchTranslations(params)
    },

    resetCursor() {
      this.cursorStack = ['']
      this.stackIndex  = 0
    },

    async fetchTranslation(id) {
      this.loading = true
      try {
        const { data } = await api.getTranslation(id)
        this.item = data
        return this.item
      } finally {
        this.loading = false
      }
    },

    async createTranslation(payload)     { return api.createTranslation(payload) },
    async updateTranslation(id, payload) { return api.updateTranslation(id, payload) },
    async removeTranslation(id)          { return api.deleteTranslation(id) },
  },
})
