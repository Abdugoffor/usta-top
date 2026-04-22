import VacancyListPage from './pages/VacancyListPage.vue'
import VacancyShowPage from './pages/VacancyShowPage.vue'

export default [
  { path: 'vacancies-admin', name: 'vacancy-admin-list', component: VacancyListPage },
  { path: 'vacancies-admin/:id', name: 'vacancy-admin-show', component: VacancyShowPage },
]
