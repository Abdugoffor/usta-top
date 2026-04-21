import LanguageListPage   from './pages/LanguageListPage.vue'
import LanguageCreatePage from './pages/LanguageCreatePage.vue'
import LanguageEditPage   from './pages/LanguageEditPage.vue'

export default [
  { path: '/languages',            name: 'language-list',   component: LanguageListPage   },
  { path: '/languages/create',     name: 'language-create', component: LanguageCreatePage },
  { path: '/languages/:id/edit',   name: 'language-edit',   component: LanguageEditPage   },
]
