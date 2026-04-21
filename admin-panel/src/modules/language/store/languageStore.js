import { defineStore } from 'pinia'
import * as api from '../api/languageApi'

export const useLanguageStore = defineStore('language', {
  state: () => ({
    items: [],
    item: null,
    loading: false,
    total: 0,
    currentPage: 1,
    perPage: 10,
    lastPage: 1,
  }),

  actions: {
    async fetchLanguages(params = {}) {
      this.loading = true
      try {
        const { data } = await api.getLanguages(params)
        this.items       = data.data ?? []
        this.total       = data.meta?.total       ?? 0
        this.currentPage = data.meta?.page        ?? 1
        this.perPage     = data.meta?.limit       ?? 10
        this.lastPage    = data.meta?.total_pages ?? 1
      } finally {
        this.loading = false
      }
    },

    async fetchLanguage(id) {
      this.loading = true
      try {
        const { data } = await api.getLanguage(id)
        this.item = data
        return this.item
      } finally {
        this.loading = false
      }
    },

    async createLanguage(payload)      { return api.createLanguage(payload) },
    async updateLanguage(id, payload)  { return api.updateLanguage(id, payload) },
    async removeLanguage(id)           { return api.deleteLanguage(id) },
  },
})
