import { createRouter, createWebHistory } from 'vue-router'
import { useLangStore } from '@/shared/stores/langStore'
import AdminLayout from '../../shared/layouts/AdminLayout.vue'
import ClientLayout from '../../shared/layouts/ClientLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import categoryRoutes    from '@/modules/category/routes'
import countryRoutes     from '@/modules/country/routes'
import languageRoutes    from '@/modules/language/routes'
import translationRoutes from '@/modules/translation/routes'
import resumeAdminRoutes from '@/modules/resume/routes'
import vacancyAdminRoutes from '@/modules/vacancy/routes'

import HomePage from '@/modules/client/pages/HomePage.vue'
import VacanciesPage from '@/modules/client/pages/VacanciesPage.vue'
import MasterDetailPage from '@/modules/client/pages/MasterDetailPage.vue'
import VacancyDetailPage from '@/modules/client/pages/VacancyDetailPage.vue'

import LoginPage from '@/modules/auth/pages/LoginPage.vue'
import RegisterPage from '@/modules/auth/pages/RegisterPage.vue'
import ProfilePage from '@/modules/profile/pages/ProfilePage.vue'
import ResumeCreatePage from '@/modules/profile/pages/ResumeCreatePage.vue'
import ResumeEditPage from '@/modules/profile/pages/ResumeEditPage.vue'
import VacancyCreatePage from '@/modules/profile/pages/VacancyCreatePage.vue'
import VacancyEditPage from '@/modules/profile/pages/VacancyEditPage.vue'

const defaultLang = () => (localStorage.getItem('site_lang') || 'uz').toLowerCase()

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: () => `/${defaultLang()}` },
    { path: '/admin', redirect: () => `/${defaultLang()}/admin` },
    {
      path: '/:lang',
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
        { path: 'profile/resumes/:id/edit', name: 'resume-edit', component: ResumeEditPage },
        { path: 'profile/vacancies/create', name: 'vacancy-create', component: VacancyCreatePage },
        { path: 'profile/vacancies/:id/edit', name: 'vacancy-edit', component: VacancyEditPage },
      ],
    },
    {
      path: '/:lang/admin',
      component: AdminLayout,
      children: [
        { path: '', name: 'dashboard', component: DashboardView },
        ...categoryRoutes,
        ...countryRoutes,
        ...languageRoutes,
        ...translationRoutes,
        ...resumeAdminRoutes,
        ...vacancyAdminRoutes,
      ],
    },
  ],
})

// URL dagi lang parametrini langStore bilan sinxronlashtirish.
// Named route push'lar uchun lang yo'q bo'lsa — joriy routedan olinadi.
router.beforeEach((to, from, next) => {
  const needsLang = to.matched.some(r => r.path.startsWith('/:lang'))
  if (needsLang && !to.params.lang) {
    const lang = (from.params.lang || defaultLang()).toLowerCase()
    return next({ ...to, params: { ...to.params, lang } })
  }
  if (to.params.lang) {
    const lc = to.params.lang.toLowerCase()
    if (lc !== to.params.lang) {
      return next({ ...to, params: { ...to.params, lang: lc } })
    }
    const store = useLangStore()
    if (store.currentLang !== lc) {
      store.setLang(lc)
    }
  }
  next()
})

export default router
