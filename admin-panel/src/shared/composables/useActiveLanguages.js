import { ref, onMounted } from 'vue'
import { getLanguages } from '@/modules/language/api/languageApi'

export function useActiveLanguages() {
  const languages = ref([])
  const loading = ref(false)

  onMounted(async () => {
    loading.value = true
    try {
      const { data } = await getLanguages({ is_active: true, limit: 100 })
      languages.value = data.data ?? []
    } finally {
      loading.value = false
    }
  })

  return { languages, loading }
}
