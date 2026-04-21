import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '../../shared/layouts/AdminLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import categoryRoutes from '@/modules/category/routes'
import countryRoutes  from '@/modules/country/routes'
import languageRoutes from '@/modules/language/routes'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: AdminLayout,
      children: [
        {
          path: '',
          name: 'dashboard',
          component: DashboardView
        },
        ...categoryRoutes,
        ...countryRoutes,
        ...languageRoutes,
      ]
    }
  ]
})

export default router