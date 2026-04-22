import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'

export function useClientListingPage({
  route,
  router,
  createFilters,
  defaultSort = 'newest',
  fetchList,
  fetchSecondaryCount,
  buildParams,
  syncToUrl,
  loadFromUrl,
  countActiveFilters,
  searchDelay = 450,
}) {
  const search = ref('')
  const loading = ref(false)
  const loadingMore = ref(false)
  const items = ref([])
  const nextCursor = ref('')
  const hasMore = ref(false)
  const total = ref(0)
  const secondaryTotal = ref(0)
  const sortBy = ref(defaultSort)
  const mobileFilterOpen = ref(false)
  const requestVersion = ref(0)
  const filters = reactive(createFilters())
  const activeFilterCount = computed(() => countActiveFilters(filters))
  const remaining = computed(() => Math.max(0, total.value - items.value.length))

  let searchTimer = null

  const fetchItems = async () => {
    const version = ++requestVersion.value
    loading.value = true

    try {
      const res = await fetchList(buildParams({ filters, sortBy: sortBy.value, cursor: '' }))
      if (version !== requestVersion.value) return

      items.value = res.data?.data || []
      const meta = res.data?.meta || {}
      hasMore.value = meta.has_more || false
      nextCursor.value = meta.next_cursor || ''
      total.value = meta.total || 0
    } catch {
      if (version !== requestVersion.value) return

      items.value = []
      hasMore.value = false
      nextCursor.value = ''
      total.value = 0
    } finally {
      if (version !== requestVersion.value) return
      loading.value = false
    }
  }

  const loadMore = async () => {
    if (!hasMore.value || loadingMore.value || !nextCursor.value) return

    const version = requestVersion.value
    loadingMore.value = true

    try {
      const res = await fetchList(buildParams({ filters, sortBy: sortBy.value, cursor: nextCursor.value }))
      if (version !== requestVersion.value) return

      const newItems = res.data?.data || []
      items.value = [...items.value, ...newItems]

      const meta = res.data?.meta || {}
      hasMore.value = meta.has_more || false
      nextCursor.value = meta.next_cursor || ''
    } catch {
      if (version !== requestVersion.value) return
      hasMore.value = false
    } finally {
      if (version !== requestVersion.value) return
      loadingMore.value = false
    }
  }

  const resetAndFetch = () => {
    items.value = []
    nextCursor.value = ''
    hasMore.value = false
    fetchItems()
  }

  const updateUrl = () => {
    syncToUrl({ filters, sortBy: sortBy.value, router })
  }

  const applyFilters = () => {
    clearTimeout(searchTimer)
    updateUrl()
    resetAndFetch()
  }

  const onSortChange = () => {
    updateUrl()
    resetAndFetch()
  }

  const fetchSecondaryTotal = async () => {
    if (!fetchSecondaryCount) return

    try {
      const res = await fetchSecondaryCount()
      secondaryTotal.value = res.data?.meta?.total || 0
    } catch {
      secondaryTotal.value = 0
    }
  }

  watch(search, (value) => {
    clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
      filters.search = value
      updateUrl()
      resetAndFetch()
    }, searchDelay)
  })

  // filters.search tashqaridan o'zgarganda (masalan resetAll) search inputni ham yangilash
  watch(() => filters.search, (newVal) => {
    if (search.value !== newVal) {
      clearTimeout(searchTimer)
      search.value = newVal
    }
  })

  onMounted(async () => {
    loadFromUrl({ query: route.query, filters, search, sortBy })
    clearTimeout(searchTimer)
    await Promise.all([fetchItems(), fetchSecondaryTotal()])
  })

  onUnmounted(() => {
    clearTimeout(searchTimer)
  })

  return {
    activeFilterCount,
    applyFilters,
    filters,
    hasMore,
    items,
    loading,
    loadMore,
    loadingMore,
    mobileFilterOpen,
    onSortChange,
    remaining,
    search,
    secondaryTotal,
    sortBy,
    total,
  }
}
