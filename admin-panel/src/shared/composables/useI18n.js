import { computed } from 'vue'
import { useLangStore } from '@/shared/stores/langStore'

/**
 * Laravel getTranslation() mantiqiga mos composable.
 * t('key') — cache bo'lsa qiymatni, bo'lmasa keyni qaytaradi
 * va fon rejimida APIdan tarjima yuklab cache'ni to'ldiradi.
 */
export function useI18n() {
  const store = useLangStore()

  const t = (key) => {
    const ck = `${store.currentLang}:${key}`
    if (store.cache[ck] !== undefined) return store.cache[ck]
    store.translate(key, store.currentLang)
    return key
  }

  return {
    t,
    lang: computed(() => store.currentLang),
    setLang: (code) => store.setLang(code),
    activeLanguages: computed(() => store.activeLanguages),
    langName: (nameMap) => store.langName(nameMap),
  }
}
