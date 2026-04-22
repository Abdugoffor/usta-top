import { defineStore } from 'pinia'
import * as api from '../api/vacancyApi'

export const useVacancyStore = defineStore('vacancy_admin', {
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
    async fetchVacancies(params = {}) {
      this.loading = true
      try {
        const cursor = this.cursorStack[this.stackIndex]
        const { data } = await api.getVacancies({
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
      await this.fetchVacancies(params)
    },

    async goPrev(params = {}) {
      if (this.stackIndex <= 0) return
      this.stackIndex--
      await this.fetchVacancies(params)
    },

    resetCursor() {
      this.cursorStack = ['']
      this.stackIndex = 0
    },

    async fetchVacancy(slug) {
      this.loading = true
      try {
        const { data } = await api.getVacancy(slug)
        this.item = data
        return this.item
      } finally {
        this.loading = false
      }
    },

    async updateVacancy(id, payload) { return api.updateVacancy(id, payload) },
    async removeVacancy(id)          { return api.deleteVacancy(id) },
  },
})
