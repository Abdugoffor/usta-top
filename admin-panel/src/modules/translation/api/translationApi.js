import axios from '@/app/providers/axios'

export const getTranslations = (params = {}) => axios.get('/translations', { params })
export const getTranslation  = (id)          => axios.get(`/translations/${id}`)
export const createTranslation = (data)      => axios.post('/translations', data)
export const updateTranslation = (id, data)  => axios.put(`/translations/${id}`, data)
export const deleteTranslation = (id)        => axios.delete(`/translations/${id}`)
