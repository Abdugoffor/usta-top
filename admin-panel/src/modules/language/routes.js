import LanguageListPage from './pages/LanguageListPage.vue'
import LanguageCreatePage from './pages/LanguageCreatePage.vue'
import LanguageEditPage from './pages/LanguageEditPage.vue'
import LanguageShowPage from './pages/LanguageShowPage.vue'

export default [
  { path: 'languages', name: 'language-list', component: LanguageListPage },
  { path: 'languages/create', name: 'language-create', component: LanguageCreatePage },
  { path: 'languages/:id', name: 'language-show', component: LanguageShowPage },
  { path: 'languages/:id/edit', name: 'language-edit', component: LanguageEditPage },
]
