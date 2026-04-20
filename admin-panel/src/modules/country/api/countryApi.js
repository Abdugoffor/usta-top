import axios from '@/app/providers/axios'

export const getCountries = (params = {}) => axios.get('/countries', { params })
export const getCountry   = (id)          => axios.get(`/countries/${id}`)
export const createCountry = (data)       => axios.post('/countries', data)
export const updateCountry = (id, data)   => axios.put(`/countries/${id}`, data)
export const deleteCountry = (id)         => axios.delete(`/countries/${id}`)
