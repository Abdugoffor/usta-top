import { ref, reactive } from 'vue'
import { defineStore } from 'pinia'
import { getLanguages } from '@/modules/language/api/languageApi'
import axios from '@/app/providers/axios'

export const useLangStore = defineStore('lang', () => {
  const currentLang = ref((localStorage.getItem('site_lang') || 'uz').toLowerCase())
  const activeLanguages = ref([])
  const cache = reactive({})
  const pending = new Set()

  const setLang = (code) => {
    const lc = (code || 'uz').toLowerCase()
    currentLang.value = lc
    localStorage.setItem('site_lang', lc)
    Object.keys(cache).forEach((k) => delete cache[k])
  }

  const fetchActiveLanguages = async () => {
    try {
      const { data } = await getLanguages({ is_active: true, limit: 100 })
      activeLanguages.value = data.data ?? []
    } catch {}
  }

  // Bir kalit bo'yicha tarjima olish (background, cache-first)
  const translate = async (key, lang) => {
    const ck = `${lang}:${key}`
    if (cache[ck] !== undefined || pending.has(ck)) return
    pending.add(ck)
    try {
      const { data } = await axios.get('/t', { params: { key, lang } })
      cache[ck] = data.value ?? key
    } catch {
      cache[ck] = key
    } finally {
      pending.delete(ck)
    }
  }

  // Kategoriya nomini joriy tildan olish (JSONB map)
  const langName = (nameMap) => {
    if (!nameMap || typeof nameMap !== 'object') return nameMap || ''
    return (
      nameMap[currentLang.value] ||
      nameMap['default'] ||
      Object.values(nameMap)[0] ||
      ''
    )
  }

  return { currentLang, activeLanguages, cache, setLang, fetchActiveLanguages, translate, langName }
})
