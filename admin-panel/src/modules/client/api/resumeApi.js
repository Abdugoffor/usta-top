import axios from '@/app/providers/axios'

export const getResumes = (params = {}) => axios.get('/resumes', { params })
export const getResume = (slug) => axios.get(`/resumes/${slug}`)
