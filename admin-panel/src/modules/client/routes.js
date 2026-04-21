import ClientLayout from '@/shared/layouts/ClientLayout.vue'
import HomePage from './pages/HomePage.vue'
import VacanciesPage from './pages/VacanciesPage.vue'
import MasterDetailPage from './pages/MasterDetailPage.vue'
import VacancyDetailPage from './pages/VacancyDetailPage.vue'

export default [
  {
    path: '/',
    component: ClientLayout,
    children: [
      {
        path: '',
        name: 'home',
        component: HomePage,
      },
      {
        path: 'vacancies',
        name: 'vacancies',
        component: VacanciesPage,
      },
      {
        path: 'masters/:slug',
        name: 'master-detail',
        component: MasterDetailPage,
      },
      {
        path: 'vacancies/:slug',
        name: 'vacancy-detail',
        component: VacancyDetailPage,
      },
    ],
  },
]
