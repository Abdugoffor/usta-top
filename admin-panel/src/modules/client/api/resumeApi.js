import axios from '@/app/providers/axios'

export const getResumes = (params = {}) => axios.get('/resumes', { params })
export const getResume = (slug) => axios.get(`/resumes/${slug}`)
export const updateResume = (id, payload) => axios.put(`/resumes/${id}`, payload)
export const deleteResume = (id) => axios.delete(`/resumes/${id}`)
