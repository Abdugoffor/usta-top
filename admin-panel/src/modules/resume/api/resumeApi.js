import axios from '@/app/providers/axios'

export const getResumes = (params = {}) => axios.get('/resumes', { params })
export const getResume = (id) => axios.get(`/resumes/${id}`)
export const updateResume = (id, data) => axios.put(`/resumes/${id}`, data)
export const deleteResume = (id) => axios.delete(`/resumes/${id}`)
