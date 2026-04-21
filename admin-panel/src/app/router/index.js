import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '../../shared/layouts/AdminLayout.vue'
import ClientLayout from '../../shared/layouts/ClientLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import categoryRoutes from '@/modules/category/routes'
import countryRoutes from '@/modules/country/routes'
import languageRoutes from '@/modules/language/routes'

import HomePage from '@/modules/client/pages/HomePage.vue'
import VacanciesPage from '@/modules/client/pages/VacanciesPage.vue'
import MasterDetailPage from '@/modules/client/pages/MasterDetailPage.vue'
import VacancyDetailPage from '@/modules/client/pages/VacancyDetailPage.vue'

import LoginPage from '@/modules/auth/pages/LoginPage.vue'
import RegisterPage from '@/modules/auth/pages/RegisterPage.vue'
import ProfilePage from '@/modules/profile/pages/ProfilePage.vue'
import ResumeCreatePage from '@/modules/profile/pages/ResumeCreatePage.vue'
import VacancyCreatePage from '@/modules/profile/pages/VacancyCreatePage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: ClientLayout,
      children: [
        { path: '', name: 'home', component: HomePage },
        { path: 'vacancies', name: 'vacancies', component: VacanciesPage },
        { path: 'masters/:slug', name: 'master-detail', component: MasterDetailPage },
        { path: 'vacancies/:slug', name: 'vacancy-detail', component: VacancyDetailPage },
        { path: 'login', name: 'login', component: LoginPage },
        { path: 'register', name: 'register', component: RegisterPage },
        { path: 'profile', name: 'profile', component: ProfilePage },
        { path: 'profile/resumes/create', name: 'resume-create', component: ResumeCreatePage },
        { path: 'profile/vacancies/create', name: 'vacancy-create', component: VacancyCreatePage },
      ],
    },
    {
      path: '/admin',
      component: AdminLayout,
      children: [
        { path: '', name: 'dashboard', component: DashboardView },
        ...categoryRoutes,
        ...countryRoutes,
        ...languageRoutes,
      ],
    },
  ],
})

export default router
