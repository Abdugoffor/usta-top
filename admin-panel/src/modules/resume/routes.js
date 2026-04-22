import ResumeListPage from './pages/ResumeListPage.vue'
import ResumeShowPage from './pages/ResumeShowPage.vue'

export default [
  { path: 'resumes', name: 'resume-admin-list', component: ResumeListPage },
  { path: 'resumes/:id', name: 'resume-admin-show', component: ResumeShowPage },
]
