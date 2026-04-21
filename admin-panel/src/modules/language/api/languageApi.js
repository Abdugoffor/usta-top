import axios from '@/app/providers/axios'

export const getLanguages = (params = {}) => axios.get('/languages', { params })
export const getLanguage  = (id)          => axios.get(`/languages/${id}`)
export const createLanguage = (data)      => axios.post('/languages', data)
export const updateLanguage = (id, data)  => axios.put(`/languages/${id}`, data)
export const deleteLanguage = (id)        => axios.delete(`/languages/${id}`)
