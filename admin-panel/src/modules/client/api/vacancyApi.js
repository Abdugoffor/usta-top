import axios from '@/app/providers/axios'

export const getVacancies = (params = {}) => axios.get('/vacancies', { params })
export const getVacancy = (slug) => axios.get(`/vacancies/${slug}`)
export const updateVacancy = (id, payload) => axios.put(`/vacancies/${id}`, payload)
export const deleteVacancy = (id) => axios.delete(`/vacancies/${id}`)
