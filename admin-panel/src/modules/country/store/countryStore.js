import { defineStore } from 'pinia'
import * as api from '../api/countryApi'

export const useCountryStore = defineStore('country', {
  state: () => ({
    tree: [],
    flatList: [],
    item: null,
    loading: false,
    total: 0,
    currentPage: 1,
    perPage: 20,
    lastPage: 1,
  }),

  actions: {
    async fetchTree(params = {}) {
      this.loading = true
      try {
        const { data } = await api.getCountries({ limit: 20, ...params })
        this.tree        = data.data ?? []
        this.total       = data.meta?.total ?? 0
        this.currentPage = data.meta?.page ?? 1
        this.perPage     = data.meta?.limit ?? 20
        this.lastPage    = data.meta?.total_pages ?? 1
      } finally {
        this.loading = false
      }
    },

    async fetchFlatList() {
      try {
        const { data } = await api.getCountries({ limit: 500 })
        const flatten = (nodes, depth = 0) => {
          const result = []
          for (const n of nodes ?? []) {
            result.push({ ...n, _depth: depth })
            result.push(...flatten(n.children, depth + 1))
          }
          return result
        }
        this.flatList = flatten(data.data ?? [])
      } catch (err) { console.error('fetchFlatList failed:', err) }
    },

    async fetchCountry(id) {
      this.loading = true
      try {
        const { data } = await api.getCountry(id)
        this.item = data
        return data
      } finally {
        this.loading = false
      }
    },

    async createCountry(payload) {
      return await api.createCountry(payload)
    },

    async updateCountry(id, payload) {
      return await api.updateCountry(id, payload)
    },

    async removeCountry(id) {
      return await api.deleteCountry(id)
    },
  },
})
