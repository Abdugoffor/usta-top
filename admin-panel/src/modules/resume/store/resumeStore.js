import { defineStore } from 'pinia'
import * as api from '../api/resumeApi'

export const useResumeStore = defineStore('resume', {
  state: () => ({
    items: [],
    item: null,
    loading: false,
    total: 0,
    hasMore: false,
    nextCursor: '',
    limit: 10,
    // cursor stack for prev/next navigation
    cursorStack: [''],  // [''] means first page (empty cursor)
    stackIndex: 0,
  }),

  getters: {
    hasPrev: (s) => s.stackIndex > 0,
    hasNext: (s) => s.hasMore,
    currentPage: (s) => s.stackIndex + 1,
    totalPages: (s) => s.hasMore
      ? s.stackIndex + 2
      : s.stackIndex + 1,
  },

  actions: {
    async fetchResumes(params = {}) {
      this.loading = true
      try {
        const cursor = this.cursorStack[this.stackIndex]
        const { data } = await api.getResumes({
          ...params,
          limit: this.limit,
          cursor: cursor || undefined,
        })
        this.items     = data.data ?? []
        this.total     = data.meta?.total ?? 0
        this.hasMore   = data.meta?.has_more ?? false
        this.nextCursor = data.meta?.next_cursor ?? ''
      } finally {
        this.loading = false
      }
    },

    async goNext(params = {}) {
      if (!this.hasMore || !this.nextCursor) return
      // push next cursor onto stack
      this.cursorStack = [...this.cursorStack.slice(0, this.stackIndex + 1), this.nextCursor]
      this.stackIndex++
      await this.fetchResumes(params)
    },

    async goPrev(params = {}) {
      if (this.stackIndex <= 0) return
      this.stackIndex--
      await this.fetchResumes(params)
    },

    resetCursor() {
      this.cursorStack = ['']
      this.stackIndex = 0
    },

    async fetchResume(slug) {
      this.loading = true
      try {
        const { data } = await api.getResume(slug)
        this.item = data
        return this.item
      } finally {
        this.loading = false
      }
    },

    async updateResume(id, payload) { return api.updateResume(id, payload) },
    async removeResume(id)          { return api.deleteResume(id) },
  },
})
