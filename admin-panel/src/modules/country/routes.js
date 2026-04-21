import CountryListPage from './pages/CountryListPage.vue'
import CountryCreatePage from './pages/CountryCreatePage.vue'
import CountryEditPage from './pages/CountryEditPage.vue'

export default [
  {
    path: 'countries',
    name: 'country-list',
    component: CountryListPage,
  },
  {
    path: 'countries/create',
    name: 'country-create',
    component: CountryCreatePage,
  },
  {
    path: 'countries/:id/edit',
    name: 'country-edit',
    component: CountryEditPage,
  },
]
