import axios from '@/app/providers/axios'

export const getVacancies = (params = {}) => axios.get('/vacancies', { params })
export const getVacancy = (id) => axios.get(`/vacancies/${id}`)
export const updateVacancy = (id, data) => axios.put(`/vacancies/${id}`, data)
export const deleteVacancy = (id) => axios.delete(`/vacancies/${id}`)
