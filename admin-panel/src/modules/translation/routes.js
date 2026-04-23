import TranslationListPage   from './pages/TranslationListPage.vue'
import TranslationCreatePage from './pages/TranslationCreatePage.vue'
import TranslationEditPage   from './pages/TranslationEditPage.vue'
import TranslationShowPage   from './pages/TranslationShowPage.vue'

export default [
  { path: 'translations',              name: 'translation-list',   component: TranslationListPage },
  { path: 'translations/create',       name: 'translation-create', component: TranslationCreatePage },
  { path: 'translations/:id',          name: 'translation-show',   component: TranslationShowPage },
  { path: 'translations/:id/edit',     name: 'translation-edit',   component: TranslationEditPage },
]
