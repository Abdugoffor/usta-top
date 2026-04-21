import axios from '@/app/providers/axios'

export const register = (data) => axios.post('/auth/register', data)
export const login = (data) => axios.post('/auth/login', data)
