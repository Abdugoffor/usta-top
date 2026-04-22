import CategoryListPage from './pages/CategoryListPage.vue'
import CategoryCreatePage from './pages/CategoryCreatePage.vue'
import CategoryEditPage from './pages/CategoryEditPage.vue'
import CategoryShowPage from './pages/CategoryShowPage.vue'

export default [
  {
    path: 'categories',
    name: 'category-list',
    component: CategoryListPage
  },
  {
    path: 'categories/create',
    name: 'category-create',
    component: CategoryCreatePage
  },
  {
    path: 'categories/:id',
    name: 'category-show',
    component: CategoryShowPage
  },
  {
    path: 'categories/:id/edit',
    name: 'category-edit',
    component: CategoryEditPage
  }
]
