<script setup>
import { ref, onMounted } from 'vue'
import { getCategories } from '@/modules/category/api/categoryApi'
import { getCountries } from '@/modules/country/api/countryApi'

const props = defineProps({
  modelValue: { type: Object, required: true },
  showCategories: { type: Boolean, default: true },
  mobileOpen: { type: Boolean, default: false }
})
const emit = defineEmits(['update:modelValue', 'apply', 'close'])

const categories = ref([])
const regions = ref([])
const districts = ref([])
const mahallas = ref([])
const priceOpen = ref(false)

const CATEGORY_ICONS = {
  'Elektrik': '⚡', 'Santexnik': '🔧', 'Qurilish': '🏗️',
  'Duradgor': '🪚', "Bo'yoqchi": '🖌️', 'Suvoqchi': '🧱',
  'Temirchi': '⚙️', 'Haydovchi': '🚗', 'Konditsioner': '❄️',
  'Kompyuter': '💻', 'Tikuvchi': '🧵', 'Oshpaz': '👨‍🍳',
}

const CATEGORY_COLORS = [
  { bg: '#fef3c7', color: '#b45309' }, { bg: '#dbeafe', color: '#1d4ed8' },
  { bg: '#dcfce7', color: '#15803d' }, { bg: '#fce7f3', color: '#be185d' },
  { bg: '#ede9fe', color: '#6d28d9' }, { bg: '#ffedd5', color: '#c2410c' },
  { bg: '#ecfdf5', color: '#065f46' }, { bg: '#f0f9ff', color: '#0369a1' },
  { bg: '#fdf4ff', color: '#7e22ce' }, { bg: '#fff7ed', color: '#9a3412' },
  { bg: '#f0fdf4', color: '#166534' }, { bg: '#fefce8', color: '#854d0e' },
]

const getCategoryStyle = (idx) => CATEGORY_COLORS[idx % CATEGORY_COLORS.length]

const update = (key, val) => {
  emit('update:modelValue', { ...props.modelValue, [key]: val })
}

const toggleCategory = (id) => {
  const current = Array.isArray(props.modelValue.category_ids) ? [...props.modelValue.category_ids] : []
  const idx = current.indexOf(id)
  if (idx === -1) current.push(id)
  else current.splice(idx, 1)
  update('category_ids', current)
  emit('apply')
}

const isCategorySelected = (id) => {
  return Array.isArray(props.modelValue.category_ids) && props.modelValue.category_ids.includes(id)
}

const selectedRegionName = (id) => {
  if (!id) return ''
  const r = regions.value.find(x => String(x.id) === String(id))
  return r ? r.name : ''
}

const selectedDistrictName = (id) => {
  if (!id) return ''
  const d = districts.value.find(x => String(x.id) === String(id))
  return d ? d.name : ''
}

const selectedMahallaName = (id) => {
  if (!id) return ''
  const m = mahallas.value.find(x => String(x.id) === String(id))
  return m ? m.name : ''
}

const selectRegion = async (val) => {
  update('region_id', val || null)
  update('district_id', null)
  update('mahalla_id', null)
  districts.value = []
  mahallas.value = []
  if (val) {
    const res = await getCountries({ parent_id: val, limit: 100 })
    districts.value = res.data?.data || []
  }
  emit('apply')
}

const selectDistrict = async (val) => {
  update('district_id', val || null)
  update('mahalla_id', null)
  mahallas.value = []
  if (val) {
    const res = await getCountries({ parent_id: val, limit: 100 })
    mahallas.value = res.data?.data || []
  }
  emit('apply')
}

const selectMahalla = (val) => {
  update('mahalla_id', val || null)
  emit('apply')
}

const applyPrice = () => emit('apply')

const resetAll = () => {
  emit('update:modelValue', {
    category_ids: [], region_id: null, district_id: null,
    mahalla_id: null, min_price: '', max_price: '', is_active_filter: '', search: ''
  })
  districts.value = []
  mahallas.value = []
  emit('apply')
}

const applyAndClose = () => {
  emit('apply')
  emit('close')
}

onMounted(async () => {
  const [catRes, regRes] = await Promise.all([
    getCategories({ limit: 50 }),
    getCountries({ limit: 100 })
  ])
  categories.value = catRes.data?.data || []
  regions.value = regRes.data?.data?.filter(r => !r.parent_id) || []

  if (props.modelValue.region_id) {
    const dr = await getCountries({ parent_id: props.modelValue.region_id, limit: 100 })
    districts.value = dr.data?.data || []
  }
  if (props.modelValue.district_id) {
    const mr = await getCountries({ parent_id: props.modelValue.district_id, limit: 100 })
    mahallas.value = mr.data?.data || []
  }
})
</script>

<template>
  <div>
    <div v-if="mobileOpen" class="filter-backdrop" @click="emit('close')"></div>

    <aside class="filter-sidebar" :class="{ 'filter-sidebar--mobile-open': mobileOpen }">
      <div class="filter-sidebar__head">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" class="filter-sidebar__icon">
          <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/>
        </svg>
        <span>Filtr</span>
        <button class="filter-sidebar__reset" @click="resetAll">Tozalash</button>
        <button class="filter-sidebar__close-btn" @click="emit('close')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M18 6 6 18M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <!-- Kategoriyalar - ko'p tanlash -->
      <div v-if="showCategories" class="filter-sidebar__section">
        <div class="filter-sidebar__section-title">
          <span>Kategoriya</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="m18 15-6-6-6 6"/>
          </svg>
        </div>
        <div class="filter-sidebar__chips">
          <button
            v-for="(cat, idx) in categories"
            :key="cat.id"
            class="filter-chip"
            :class="{ 'filter-chip--active': isCategorySelected(cat.id) }"
            :style="isCategorySelected(cat.id) ? {} : { background: getCategoryStyle(idx).bg, color: getCategoryStyle(idx).color }"
            @click="toggleCategory(cat.id)"
          >
            <span>{{ CATEGORY_ICONS[cat.name] || '🔨' }}</span>
            {{ cat.name }}
          </button>
        </div>
        <div v-if="modelValue.category_ids?.length" class="filter-sidebar__selected-tags">
          <span
            v-for="catId in modelValue.category_ids"
            :key="catId"
            class="selected-tag"
          >
            {{ categories.find(c => c.id === catId)?.name }}
            <button @click="toggleCategory(catId)">×</button>
          </span>
        </div>
      </div>

      <!-- Joylashuv -->
      <div class="filter-sidebar__section">
        <div class="filter-sidebar__section-title">
          <span>Joylashuv</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="m18 15-6-6-6 6"/>
          </svg>
        </div>
        <div class="filter-sidebar__selects">
          <div class="filter-sidebar__select-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
              <circle cx="12" cy="9" r="2.5"/>
            </svg>
            <select :value="modelValue.region_id || ''" @change="selectRegion($event.target.value || null)" class="filter-sidebar__select">
              <option value="">Viloyat tanlang</option>
              <option v-for="r in regions" :key="r.id" :value="r.id">{{ r.name }}</option>
            </select>
          </div>

          <div v-if="modelValue.region_id" class="filter-sidebar__location-pill">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
              <circle cx="12" cy="9" r="2.5"/>
            </svg>
            <span>{{ selectedRegionName(modelValue.region_id) }}</span>
            <button @click="selectRegion(null)">×</button>
          </div>

          <div class="filter-sidebar__select-wrap" v-if="districts.length">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
              <polyline points="9 22 9 12 15 12 15 22"/>
            </svg>
            <select :value="modelValue.district_id || ''" @change="selectDistrict($event.target.value || null)" class="filter-sidebar__select">
              <option value="">Tuman</option>
              <option v-for="d in districts" :key="d.id" :value="d.id">{{ d.name }}</option>
            </select>
          </div>

          <div v-if="modelValue.district_id" class="filter-sidebar__location-pill filter-sidebar__location-pill--sm">
            <span>{{ selectedDistrictName(modelValue.district_id) }}</span>
            <button @click="selectDistrict(null)">×</button>
          </div>

          <div class="filter-sidebar__select-wrap" v-if="mahallas.length">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
              <polyline points="9 22 9 12 15 12 15 22"/>
            </svg>
            <select :value="modelValue.mahalla_id || ''" @change="selectMahalla($event.target.value || null)" class="filter-sidebar__select">
              <option value="">Mahalla</option>
              <option v-for="m in mahallas" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Narx -->
      <div class="filter-sidebar__section">
        <button class="filter-sidebar__section-title filter-sidebar__section-title--btn" @click="priceOpen = !priceOpen">
          <span>Narx (so'm)</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"
            :style="{ transform: priceOpen ? 'rotate(0)' : 'rotate(180deg)', transition: 'transform 0.2s' }">
            <path d="m18 15-6-6-6 6"/>
          </svg>
        </button>
        <div v-if="priceOpen" class="filter-sidebar__price">
          <div class="filter-sidebar__price-row">
            <input type="number" class="filter-sidebar__price-input" placeholder="Min"
              :value="modelValue.min_price" @input="update('min_price', $event.target.value)" />
            <span class="filter-sidebar__price-sep">—</span>
            <input type="number" class="filter-sidebar__price-input" placeholder="Max"
              :value="modelValue.max_price" @input="update('max_price', $event.target.value)" />
          </div>
          <button class="filter-sidebar__price-apply" @click="applyPrice">Qo'llash</button>
        </div>
      </div>

      <!-- Holat filtri - 3 holat -->
      <div class="filter-sidebar__section">
        <div class="filter-sidebar__section-title">
          <span>Holat</span>
        </div>
        <div class="filter-sidebar__status-btns">
          <button class="status-btn" :class="{ 'status-btn--active': !modelValue.is_active_filter }"
            @click="update('is_active_filter', ''); $emit('apply')">
            Hammasi
          </button>
          <button class="status-btn status-btn--green" :class="{ 'status-btn--active': modelValue.is_active_filter === 'active' }"
            @click="update('is_active_filter', 'active'); $emit('apply')">
            Faol
          </button>
          <button class="status-btn status-btn--gray" :class="{ 'status-btn--active': modelValue.is_active_filter === 'inactive' }"
            @click="update('is_active_filter', 'inactive'); $emit('apply')">
            Yopiq
          </button>
        </div>
      </div>

      <button class="filter-sidebar__apply-btn" @click="applyAndClose">
        Filtrlash
      </button>
    </aside>
  </div>
</template>

<style scoped>
.filter-backdrop { display: none; }

.filter-sidebar {
  background: #fff; border: 1px solid #e5e7eb; border-radius: 16px; padding: 20px;
  display: flex; flex-direction: column; gap: 20px; height: fit-content;
  position: sticky; top: 80px;
}

.filter-sidebar__head { display: flex; align-items: center; gap: 8px; font-size: 15px; font-weight: 700; color: #111827; }
.filter-sidebar__icon { width: 18px; height: 18px; color: #374151; }

.filter-sidebar__reset {
  margin-left: auto; background: none; border: none; font-size: 12px; color: #2563eb;
  cursor: pointer; font-weight: 600; font-family: inherit; padding: 4px 8px; border-radius: 6px; transition: background 0.2s;
}
.filter-sidebar__reset:hover { background: #eff6ff; }

.filter-sidebar__close-btn {
  display: none; background: none; border: none; cursor: pointer; padding: 4px; color: #6b7280;
  border-radius: 8px; align-items: center; justify-content: center; flex-shrink: 0;
}
.filter-sidebar__close-btn svg { width: 20px; height: 20px; }
.filter-sidebar__close-btn:hover { background: #f3f4f6; }

.filter-sidebar__section {
  display: flex; flex-direction: column; gap: 10px; padding-top: 16px; border-top: 1px solid #f3f4f6;
}

.filter-sidebar__section-title {
  display: flex; align-items: center; justify-content: space-between;
  font-size: 13px; font-weight: 700; color: #374151; text-transform: uppercase; letter-spacing: 0.05em;
}
.filter-sidebar__section-title svg { width: 16px; height: 16px; }

.filter-sidebar__section-title--btn {
  background: none; border: none; cursor: pointer; font-family: inherit; width: 100%; text-align: left; padding: 0;
}

.filter-sidebar__chips { display: flex; flex-wrap: wrap; gap: 6px; }

.filter-chip {
  display: inline-flex; align-items: center; gap: 4px; padding: 6px 12px; border-radius: 20px;
  font-size: 12px; font-weight: 600; cursor: pointer; border: 1.5px solid transparent;
  transition: all 0.2s; font-family: inherit;
}
.filter-chip:hover { opacity: 0.85; transform: translateY(-1px); }
.filter-chip--active {
  background: linear-gradient(135deg, #1d4ed8, #2563eb) !important;
  color: #fff !important; border-color: transparent;
  box-shadow: 0 2px 8px rgba(37,99,235,0.3);
}

.filter-sidebar__selected-tags {
  display: flex; flex-wrap: wrap; gap: 6px; margin-top: 4px;
}

.selected-tag {
  display: inline-flex; align-items: center; gap: 4px;
  background: #eff6ff; color: #1d4ed8; font-size: 12px; font-weight: 600;
  padding: 4px 8px 4px 10px; border-radius: 20px;
}

.selected-tag button {
  background: none; border: none; cursor: pointer; color: #1d4ed8;
  font-size: 14px; line-height: 1; padding: 0; display: flex; align-items: center;
}

.filter-sidebar__selects { display: flex; flex-direction: column; gap: 8px; }

.filter-sidebar__select-wrap {
  position: relative; display: flex; align-items: center;
}
.filter-sidebar__select-wrap svg { position: absolute; left: 12px; width: 16px; height: 16px; color: #9ca3af; pointer-events: none; }

.filter-sidebar__select {
  width: 100%; padding: 10px 12px 10px 36px; border: 1.5px solid #e5e7eb; border-radius: 10px;
  font-size: 13px; color: #374151; background: #f9fafb; outline: none; cursor: pointer;
  transition: border-color 0.2s; appearance: none; font-family: inherit;
}
.filter-sidebar__select:focus { border-color: #2563eb; background: #fff; }

.filter-sidebar__location-pill {
  display: flex; align-items: center; gap: 6px; background: #eff6ff; color: #1d4ed8;
  border-radius: 10px; padding: 8px 12px; font-size: 13px; font-weight: 600;
}
.filter-sidebar__location-pill svg { width: 14px; height: 14px; flex-shrink: 0; }
.filter-sidebar__location-pill span { flex: 1; }
.filter-sidebar__location-pill button { background: none; border: none; cursor: pointer; color: #1d4ed8; font-size: 16px; line-height: 1; padding: 0; }

.filter-sidebar__location-pill--sm { background: #f0fdf4; color: #166534; }
.filter-sidebar__location-pill--sm button { color: #166534; }

.filter-sidebar__price { display: flex; flex-direction: column; gap: 8px; }
.filter-sidebar__price-row { display: flex; align-items: center; gap: 8px; }
.filter-sidebar__price-input { flex: 1; padding: 9px 12px; border: 1.5px solid #e5e7eb; border-radius: 10px; font-size: 13px; color: #374151; background: #f9fafb; outline: none; transition: border-color 0.2s; font-family: inherit; min-width: 0; }
.filter-sidebar__price-input:focus { border-color: #2563eb; background: #fff; }
.filter-sidebar__price-input::placeholder { color: #9ca3af; }
.filter-sidebar__price-sep { color: #9ca3af; font-size: 14px; flex-shrink: 0; }
.filter-sidebar__price-apply { background: linear-gradient(135deg, #1d4ed8, #2563eb); color: #fff; border: none; border-radius: 10px; padding: 9px 16px; font-size: 13px; font-weight: 600; cursor: pointer; font-family: inherit; transition: opacity 0.2s; }
.filter-sidebar__price-apply:hover { opacity: 0.9; }

/* Holat 3-state tugmalari */
.filter-sidebar__status-btns { display: flex; gap: 6px; }

.status-btn {
  flex: 1; padding: 8px 6px; border: 1.5px solid #e5e7eb; border-radius: 10px;
  font-size: 12px; font-weight: 600; color: #374151; background: #f9fafb;
  cursor: pointer; font-family: inherit; transition: all 0.2s; text-align: center;
}
.status-btn:hover { border-color: #93c5fd; }

.status-btn--active {
  background: linear-gradient(135deg, #1d4ed8, #2563eb) !important;
  color: #fff !important; border-color: transparent !important;
  box-shadow: 0 2px 8px rgba(37,99,235,0.25);
}

.status-btn--green.status-btn--active {
  background: linear-gradient(135deg, #059669, #10b981) !important;
  box-shadow: 0 2px 8px rgba(5,150,105,0.25) !important;
}

.status-btn--gray.status-btn--active {
  background: linear-gradient(135deg, #4b5563, #6b7280) !important;
  box-shadow: 0 2px 8px rgba(75,85,99,0.25) !important;
}

.filter-sidebar__apply-btn {
  display: none;
  background: linear-gradient(135deg, #1d4ed8, #2563eb); color: #fff; border: none;
  border-radius: 12px; padding: 14px; font-size: 15px; font-weight: 700;
  cursor: pointer; font-family: inherit; width: 100%; text-align: center; margin-top: 4px;
}

/* Mobile */
@media (max-width: 900px) {
  .filter-backdrop { display: block; position: fixed; inset: 0; background: rgba(0,0,0,0.45); z-index: 499; backdrop-filter: blur(2px); }

  .filter-sidebar {
    display: none; position: fixed; bottom: 0; left: 0; right: 0; top: auto;
    max-height: 88vh; overflow-y: auto; border-radius: 20px 20px 0 0; z-index: 500;
    box-shadow: 0 -8px 32px rgba(0,0,0,0.18); padding: 20px 20px 32px; border: none;
  }

  .filter-sidebar--mobile-open { display: flex; }
  .filter-sidebar__close-btn { display: flex; }
  .filter-sidebar__apply-btn { display: block; }
}
</style>
