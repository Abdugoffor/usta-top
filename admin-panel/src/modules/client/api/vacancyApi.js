import axios from '@/app/providers/axios'

export const getVacancies = (params = {}) => axios.get('/vacancies', { params })
export const getVacancy = (slug) => axios.get(`/vacancies/${slug}`)
