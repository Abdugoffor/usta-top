<script setup>
import { ref, onMounted, watch } from 'vue'
import { getCategories } from '@/modules/category/api/categoryApi'
import { getCountries } from '@/modules/country/api/countryApi'
import { useI18n } from '@/shared/composables/useI18n'

const props = defineProps({
  modelValue: { type: Object, required: true },
  showCategories: { type: Boolean, default: true },
  mobileOpen: { type: Boolean, default: false }
})
const emit = defineEmits(['update:modelValue', 'apply', 'close'])

const { t, langName } = useI18n()

const categories = ref([])
const regions = ref([])
const districts = ref([])
const mahallas = ref([])
const priceOpen = ref(!!(props.modelValue.min_price || props.modelValue.max_price))
const regionRequestVersion = ref(0)
const districtRequestVersion = ref(0)

const CATEGORY_ICONS = {
  Elektrik: '⚡', Santexnik: '🔧', Qurilish: '🏗️', Duradgor: '🪚',
  "Bo'yoqchi": '🖌️', Suvoqchi: '🧱', Temirchi: '⚙️', Haydovchi: '🚗',
  Konditsioner: '❄️', Kompyuter: '💻', Tikuvchi: '🧵', Oshpaz: '👨‍🍳',
}

const CATEGORY_COLORS = [
  { bg: '#fef3c7', color: '#b45309' }, { bg: '#dbeafe', color: '#1d4ed8' },
  { bg: '#dcfce7', color: '#15803d' }, { bg: '#fce7f3', color: '#be185d' },
  { bg: '#ede9fe', color: '#6d28d9' }, { bg: '#ffedd5', color: '#c2410c' },
  { bg: '#ecfdf5', color: '#065f46' }, { bg: '#f0f9ff', color: '#0369a1' },
  { bg: '#fdf4ff', color: '#7e22ce' }, { bg: '#fff7ed', color: '#9a3412' },
  { bg: '#f0fdf4', color: '#166534' }, { bg: '#fefce8', color: '#854d0e' },
]

const getCategoryStyle = (index) => CATEGORY_COLORS[index % CATEGORY_COLORS.length]

// Kategoriya belgisini default nom bo'yicha topadi
const getCategoryIcon = (category) => {
  const name = langName(category.name)
  return CATEGORY_ICONS[name] || '🔨'
}

const patchModel = (changes) => emit('update:modelValue', { ...props.modelValue, ...changes })
const update = (key, value) => patchModel({ [key]: value })

const toggleCategory = (id) => {
  const current = Array.isArray(props.modelValue.category_ids) ? [...props.modelValue.category_ids] : []
  const index = current.indexOf(id)
  if (index === -1) current.push(id)
  else current.splice(index, 1)
  patchModel({ category_ids: current })
  emit('apply')
}

const isCategorySelected = (id) =>
  Array.isArray(props.modelValue.category_ids) && props.modelValue.category_ids.includes(id)

const findLocationName = (items, id) => {
  if (!id) return ''
  const item = items.find((entry) => String(entry.id) === String(id))
  return item ? item.name : ''
}

const selectRegion   = (value) => { patchModel({ region_id: value || null, district_id: null, mahalla_id: null }); emit('apply') }
const selectDistrict = (value) => { patchModel({ district_id: value || null, mahalla_id: null }); emit('apply') }
const selectMahalla  = (value) => { update('mahalla_id', value || null); emit('apply') }

const resetAll = () => {
  emit('update:modelValue', { category_ids: [], region_id: null, district_id: null, mahalla_id: null, min_price: '', max_price: '', is_active_filter: '', search: '' })
  districts.value = []
  mahallas.value = []
  emit('apply')
}

const applyAndClose = () => { emit('apply'); emit('close') }

watch(
  () => !!(props.modelValue.min_price || props.modelValue.max_price),
  (hasPrice) => { if (hasPrice) priceOpen.value = true },
  { immediate: true }
)

watch(() => props.modelValue.region_id, async (regionId) => {
  const version = ++regionRequestVersion.value
  districts.value = []
  mahallas.value = []
  if (!regionId) return
  try {
    const response = await getCountries({ parent_id: regionId, limit: 100 })
    if (version !== regionRequestVersion.value) return
    districts.value = response.data?.data || []
  } catch {
    if (version !== regionRequestVersion.value) return
    districts.value = []
  }
}, { immediate: true })

watch(() => props.modelValue.district_id, async (districtId) => {
  const version = ++districtRequestVersion.value
  mahallas.value = []
  if (!districtId) return
  try {
    const response = await getCountries({ parent_id: districtId, limit: 100 })
    if (version !== districtRequestVersion.value) return
    mahallas.value = response.data?.data || []
  } catch {
    if (version !== districtRequestVersion.value) return
    mahallas.value = []
  }
}, { immediate: true })

onMounted(async () => {
  const requests = [getCountries({ limit: 100 })]
  if (props.showCategories) requests.unshift(getCategories({ limit: 50 }))

  const responses = await Promise.all(requests)
  const categoryResponse = props.showCategories ? responses[0] : null
  const regionResponse   = props.showCategories ? responses[1] : responses[0]

  categories.value = categoryResponse?.data?.data || []
  regions.value    = regionResponse?.data?.data   || []
})
</script>

<template>
  <div>
    <div v-if="mobileOpen" class="filter-backdrop" @click="emit('close')"></div>

    <aside class="filter-sidebar" :class="{ 'filter-sidebar--mobile-open': mobileOpen }">
      <div class="filter-sidebar__head">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" class="filter-sidebar__icon">
          <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
        </svg>
        <span>{{ t('filter_title') }}</span>
        <button class="filter-sidebar__reset" @click="resetAll">{{ t('filter_reset') }}</button>
        <button class="filter-sidebar__close-btn" @click="emit('close')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M18 6 6 18M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Kategoriya -->
      <div v-if="showCategories" class="filter-sidebar__section">
        <div class="filter-sidebar__section-title">
          <span>{{ t('filter_category') }}</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="m18 15-6-6-6 6" />
          </svg>
        </div>

        <div class="filter-sidebar__chips">
          <button
            v-for="(category, index) in categories"
            :key="category.id"
            class="filter-chip"
            :class="{ 'filter-chip--active': isCategorySelected(category.id) }"
            :style="isCategorySelected(category.id) ? {} : { background: getCategoryStyle(index).bg, color: getCategoryStyle(index).color }"
            @click="toggleCategory(category.id)"
          >
            <span>{{ getCategoryIcon(category) }}</span>
            {{ langName(category.name) }}
          </button>
        </div>

        <div v-if="modelValue.category_ids?.length" class="filter-sidebar__selected-tags">
          <span v-for="categoryId in modelValue.category_ids" :key="categoryId" class="selected-tag">
            {{ langName(categories.find((c) => c.id === categoryId)?.name) }}
            <button @click="toggleCategory(categoryId)">×</button>
          </span>
        </div>
      </div>

      <!-- Joylashuv -->
      <div class="filter-sidebar__section">
        <div class="filter-sidebar__section-title">
          <span>{{ t('filter_location') }}</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="m18 15-6-6-6 6" />
          </svg>
        </div>

        <div class="filter-sidebar__selects">
          <div class="filter-sidebar__select-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z" />
              <circle cx="12" cy="9" r="2.5" />
            </svg>
            <select :value="modelValue.region_id || ''" class="filter-sidebar__select" @change="selectRegion($event.target.value || null)">
              <option value="">{{ t('filter_select_region') }}</option>
              <option v-for="region in regions" :key="region.id" :value="region.id">{{ region.name }}</option>
            </select>
          </div>

          <div v-if="modelValue.region_id" class="filter-sidebar__location-pill">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z" />
              <circle cx="12" cy="9" r="2.5" />
            </svg>
            <span>{{ findLocationName(regions, modelValue.region_id) }}</span>
            <button @click="selectRegion(null)">×</button>
          </div>

          <div v-if="districts.length" class="filter-sidebar__select-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
              <polyline points="9 22 9 12 15 12 15 22" />
            </svg>
            <select :value="modelValue.district_id || ''" class="filter-sidebar__select" @change="selectDistrict($event.target.value || null)">
              <option value="">{{ t('filter_select_district') }}</option>
              <option v-for="district in districts" :key="district.id" :value="district.id">{{ district.name }}</option>
            </select>
          </div>

          <div v-if="modelValue.district_id" class="filter-sidebar__location-pill filter-sidebar__location-pill--sm">
            <span>{{ findLocationName(districts, modelValue.district_id) }}</span>
            <button @click="selectDistrict(null)">×</button>
          </div>

          <div v-if="mahallas.length" class="filter-sidebar__select-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
              <polyline points="9 22 9 12 15 12 15 22" />
            </svg>
            <select :value="modelValue.mahalla_id || ''" class="filter-sidebar__select" @change="selectMahalla($event.target.value || null)">
              <option value="">{{ t('filter_select_mahalla') }}</option>
              <option v-for="mahalla in mahallas" :key="mahalla.id" :value="mahalla.id">{{ mahalla.name }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Narx -->
      <div class="filter-sidebar__section">
        <button class="filter-sidebar__section-title filter-sidebar__section-title--btn" @click="priceOpen = !priceOpen">
          <span>{{ t('filter_price') }}</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" :style="{ transform: priceOpen ? 'rotate(0)' : 'rotate(180deg)', transition: 'transform 0.2s' }">
            <path d="m18 15-6-6-6 6" />
          </svg>
        </button>

        <div v-if="priceOpen" class="filter-sidebar__price">
          <div class="filter-sidebar__price-row">
            <input class="filter-sidebar__price-input" type="number" :placeholder="t('filter_price_min')" :value="modelValue.min_price" @input="update('min_price', $event.target.value)" />
            <span class="filter-sidebar__price-sep">—</span>
            <input class="filter-sidebar__price-input" type="number" :placeholder="t('filter_price_max')" :value="modelValue.max_price" @input="update('max_price', $event.target.value)" />
          </div>
          <button class="filter-sidebar__price-apply" @click="emit('apply')">{{ t('filter_apply') }}</button>
        </div>
      </div>

      <!-- Holat -->
      <div class="filter-sidebar__section">
        <div class="filter-sidebar__section-title">
          <span>{{ t('filter_status') }}</span>
        </div>

        <div class="filter-sidebar__status-btns">
          <button class="status-btn" :class="{ 'status-btn--active': !modelValue.is_active_filter }"
            @click="update('is_active_filter', ''); emit('apply')">
            {{ t('filter_all') }}
          </button>
          <button class="status-btn status-btn--green" :class="{ 'status-btn--active': modelValue.is_active_filter === 'active' }"
            @click="update('is_active_filter', 'active'); emit('apply')">
            {{ t('filter_active') }}
          </button>
          <button class="status-btn status-btn--gray" :class="{ 'status-btn--active': modelValue.is_active_filter === 'inactive' }"
            @click="update('is_active_filter', 'inactive'); emit('apply')">
            {{ t('filter_inactive') }}
          </button>
        </div>
      </div>

      <button class="filter-sidebar__apply-btn" @click="applyAndClose">
        {{ t('filter_submit') }}
      </button>
    </aside>
  </div>
</template>
