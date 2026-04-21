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
  'Elektrik': '⚡',   'Santexnik': '🔧',  'Qurilish': '🏗️',
  'Duradgor': '🪚',   "Bo'yoqchi": '🖌️',  'Suvoqchi': '🧱',
  'Temirchi': '⚙️',   'Haydovchi': '🚗',  'Konditsioner': '❄️',
  'Kompyuter': '💻',  'Tikuvchi': '🧵',   'Oshpaz': '👨‍🍳',
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
  const current = props.modelValue.category_id
  update('category_id', current === id ? null : id)
  emit('apply')
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
    category_id: null, region_id: null, district_id: null,
    mahalla_id: null, min_price: '', max_price: '', only_active: false, search: ''
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
            :class="{ 'filter-chip--active': modelValue.category_id === cat.id }"
            :style="modelValue.category_id === cat.id ? {} : { background: getCategoryStyle(idx).bg, color: getCategoryStyle(idx).color }"
            @click="toggleCategory(cat.id)"
          >
            <span>{{ CATEGORY_ICONS[cat.name] || '🔨' }}</span>
            {{ cat.name }}
          </button>
        </div>
      </div>

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
            <input
              type="number"
              class="filter-sidebar__price-input"
              placeholder="Min"
              :value="modelValue.min_price"
              @input="update('min_price', $event.target.value)"
            />
            <span class="filter-sidebar__price-sep">—</span>
            <input
              type="number"
              class="filter-sidebar__price-input"
              placeholder="Max"
              :value="modelValue.max_price"
              @input="update('max_price', $event.target.value)"
            />
          </div>
          <button class="filter-sidebar__price-apply" @click="applyPrice">Qo'llash</button>
        </div>
      </div>

      <div class="filter-sidebar__section filter-sidebar__section--toggle">
        <div>
          <div class="filter-sidebar__toggle-label">Faqat faollar</div>
          <div class="filter-sidebar__toggle-sub">Bo'sh ustalar</div>
        </div>
        <label class="toggle-switch">
          <input
            type="checkbox"
            :checked="modelValue.only_active"
            @change="update('only_active', $event.target.checked); $emit('apply')"
          />
          <span class="toggle-switch__slider"></span>
        </label>
      </div>

      <button class="filter-sidebar__apply-btn" @click="applyAndClose">
        Filtrlash
      </button>
    </aside>
  </div>
</template>

<style scoped>
.filter-backdrop {
  display: none;
}

.filter-sidebar {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: fit-content;
  position: sticky;
  top: 80px;
}

.filter-sidebar__head {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 700;
  color: #111827;
}

.filter-sidebar__icon { width: 18px; height: 18px; color: #374151; }

.filter-sidebar__reset {
  margin-left: auto;
  background: none;
  border: none;
  font-size: 12px;
  color: #2563eb;
  cursor: pointer;
  font-weight: 600;
  font-family: inherit;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background 0.2s;
}

.filter-sidebar__reset:hover { background: #eff6ff; }

.filter-sidebar__close-btn {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: #6b7280;
  border-radius: 8px;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.filter-sidebar__close-btn svg { width: 20px; height: 20px; }
.filter-sidebar__close-btn:hover { background: #f3f4f6; }

.filter-sidebar__section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-top: 16px;
  border-top: 1px solid #f3f4f6;
}

.filter-sidebar__section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
  font-weight: 700;
  color: #374151;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.filter-sidebar__section-title svg { width: 16px; height: 16px; }

.filter-sidebar__section-title--btn {
  background: none;
  border: none;
  cursor: pointer;
  font-family: inherit;
  width: 100%;
  text-align: left;
  padding: 0;
}

.filter-sidebar__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.filter-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  border: 1.5px solid transparent;
  transition: all 0.2s;
  font-family: inherit;
}

.filter-chip:hover { opacity: 0.85; transform: translateY(-1px); }

.filter-chip--active {
  background: linear-gradient(135deg, #1d4ed8, #2563eb) !important;
  color: #fff !important;
  border-color: transparent;
  box-shadow: 0 2px 8px rgba(37,99,235,0.3);
}

.filter-sidebar__selects {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-sidebar__select-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.filter-sidebar__select-wrap svg {
  position: absolute;
  left: 12px;
  width: 16px;
  height: 16px;
  color: #9ca3af;
  pointer-events: none;
}

.filter-sidebar__select {
  width: 100%;
  padding: 10px 12px 10px 36px;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  font-size: 13px;
  color: #374151;
  background: #f9fafb;
  outline: none;
  cursor: pointer;
  transition: border-color 0.2s;
  appearance: none;
  font-family: inherit;
}

.filter-sidebar__select:focus { border-color: #2563eb; background: #fff; }

.filter-sidebar__price {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-sidebar__price-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-sidebar__price-input {
  flex: 1;
  padding: 9px 12px;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  font-size: 13px;
  color: #374151;
  background: #f9fafb;
  outline: none;
  transition: border-color 0.2s;
  font-family: inherit;
  min-width: 0;
}

.filter-sidebar__price-input:focus { border-color: #2563eb; background: #fff; }
.filter-sidebar__price-input::placeholder { color: #9ca3af; }

.filter-sidebar__price-sep { color: #9ca3af; font-size: 14px; flex-shrink: 0; }

.filter-sidebar__price-apply {
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  color: #fff;
  border: none;
  border-radius: 10px;
  padding: 9px 16px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: opacity 0.2s;
}

.filter-sidebar__price-apply:hover { opacity: 0.9; }

.filter-sidebar__section--toggle {
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

.filter-sidebar__toggle-label {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
}

.filter-sidebar__toggle-sub {
  font-size: 11px;
  color: #9ca3af;
  margin-top: 2px;
}

.filter-sidebar__apply-btn {
  display: none;
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  color: #fff;
  border: none;
  border-radius: 12px;
  padding: 14px;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  width: 100%;
  text-align: center;
  margin-top: 4px;
}

.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
  flex-shrink: 0;
}

.toggle-switch input { opacity: 0; width: 0; height: 0; }

.toggle-switch__slider {
  position: absolute;
  inset: 0;
  background: #e5e7eb;
  border-radius: 24px;
  cursor: pointer;
  transition: all 0.3s;
}

.toggle-switch__slider::before {
  content: '';
  position: absolute;
  width: 18px;
  height: 18px;
  left: 3px;
  bottom: 3px;
  background: #fff;
  border-radius: 50%;
  transition: all 0.3s;
  box-shadow: 0 1px 4px rgba(0,0,0,0.2);
}

.toggle-switch input:checked + .toggle-switch__slider { background: #2563eb; }
.toggle-switch input:checked + .toggle-switch__slider::before { transform: translateX(20px); }

/* Mobile */
@media (max-width: 900px) {
  .filter-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.45);
    z-index: 499;
    backdrop-filter: blur(2px);
  }

  .filter-sidebar {
    display: none;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    top: auto;
    max-height: 88vh;
    overflow-y: auto;
    border-radius: 20px 20px 0 0;
    z-index: 500;
    box-shadow: 0 -8px 32px rgba(0, 0, 0, 0.18);
    padding: 20px 20px 32px;
    border: none;
  }

  .filter-sidebar--mobile-open {
    display: flex;
  }

  .filter-sidebar__close-btn {
    display: flex;
  }

  .filter-sidebar__apply-btn {
    display: block;
  }
}
</style>
