<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store/authStore'
import axios from '@/app/providers/axios'
import { getResume, updateResume } from '@/modules/client/api/resumeApi'
import { getCategories } from '@/modules/category/api/categoryApi'
import { getCountries } from '@/modules/country/api/countryApi'
import ClientHeader from '@/modules/client/components/ClientHeader.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const resumeId = ref(null)
const form = ref({
  name: '', title: '', adress: '', text: '', contact: '',
  photo: '', skills: '', price: '', experience_year: '',
  region_id: null, district_id: null, mahalla_id: null,
  is_active: true, category_ids: [],
})

const categories        = ref([])
const regions           = ref([])
const districts         = ref([])
const mahallas          = ref([])
const loading           = ref(false)
const pageLoading       = ref(true)
const error             = ref('')
const fieldErrors       = ref({})
const photoPreview      = ref('')
const photoLoading      = ref(false)
const photoUploadError  = ref('')

const apiHost = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1').replace(/\/api\/v1\/?$/, '')

const uploadPhoto = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  photoLoading.value = true
  photoUploadError.value = ''
  const fd = new FormData()
  fd.append('photo', file)
  try {
    const res = await axios.post('/upload', fd, { headers: { 'Content-Type': null } })
    form.value.photo = apiHost + res.data.url
    photoPreview.value = form.value.photo
  } catch (err) {
    photoUploadError.value = err.response?.data?.error || 'Rasm yuklanmadi'
  } finally {
    photoLoading.value = false
  }
}

const selectRegion = async (val) => {
  form.value.region_id = val ? Number(val) : null
  form.value.district_id = null
  form.value.mahalla_id = null
  districts.value = []
  mahallas.value = []
  if (val) {
    const res = await getCountries({ parent_id: val, limit: 100 })
    districts.value = res.data?.data || []
  }
}

const selectDistrict = async (val) => {
  form.value.district_id = val ? Number(val) : null
  form.value.mahalla_id = null
  mahallas.value = []
  if (val) {
    const res = await getCountries({ parent_id: val, limit: 100 })
    mahallas.value = res.data?.data || []
  }
}

const toggleCategory = (id) => {
  const idx = form.value.category_ids.indexOf(id)
  if (idx === -1) form.value.category_ids.push(id)
  else form.value.category_ids.splice(idx, 1)
}

const submit = async () => {
  error.value = ''
  fieldErrors.value = {}
  loading.value = true
  try {
    const payload = {
      name: form.value.name,
      title: form.value.title,
      adress: form.value.adress,
      text: form.value.text,
      contact: form.value.contact,
      photo: form.value.photo,
      skills: form.value.skills,
      is_active: form.value.is_active,
    }
    if (form.value.price) payload.price = Number(form.value.price)
    if (form.value.experience_year) payload.experience_year = Number(form.value.experience_year)
    if (form.value.region_id) payload.region_id = form.value.region_id
    if (form.value.district_id) payload.district_id = form.value.district_id
    if (form.value.mahalla_id) payload.mahalla_id = form.value.mahalla_id

    await updateResume(resumeId.value, payload)
    router.push({ name: 'profile' })
  } catch (e) {
    const data = e.response?.data
    if (data?.errors) fieldErrors.value = data.errors
    else error.value = data?.message || data?.error || 'Xatolik yuz berdi'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (!auth.isLoggedIn) { router.push({ name: 'login' }); return }
  try {
    const slug = route.params.id
    const [resumeRes, catRes, regRes] = await Promise.all([
      getResume(slug),
      getCategories({ limit: 50 }),
      getCountries({ parent_id: 196, limit: 100 }),
    ])
    const r = resumeRes.data
    resumeId.value = r.id
    form.value = {
      name: r.name || '',
      title: r.title || '',
      adress: r.adress || '',
      text: r.text || '',
      contact: r.contact || '',
      photo: r.photo || '',
      skills: r.skills || '',
      price: r.price || '',
      experience_year: r.experience_year || '',
      region_id: r.region_id || null,
      district_id: r.district_id || null,
      mahalla_id: r.mahalla_id || null,
      is_active: r.is_active ?? true,
      category_ids: (r.categories || []).map(c => c.id),
    }
    categories.value = catRes.data?.data || []
    regions.value = regRes.data?.data || []
    if (r.photo) photoPreview.value = r.photo

    if (form.value.region_id) {
      const dr = await getCountries({ parent_id: form.value.region_id, limit: 100 })
      districts.value = dr.data?.data || []
    }
    if (form.value.district_id) {
      const mr = await getCountries({ parent_id: form.value.district_id, limit: 100 })
      mahallas.value = mr.data?.data || []
    }
  } catch {
    router.push({ name: 'profile' })
  } finally {
    pageLoading.value = false
  }
})
</script>

<template>
  <div class="create-page">
    <ClientHeader />

    <div v-if="pageLoading" class="page-loading">
      <div class="spinner"></div>
    </div>

    <template v-else>
      <div class="create-hero">
        <div class="create-hero__inner">
          <RouterLink :to="`/${route.params.lang}/profile`" class="create-back">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m15 18-6-6 6-6"/></svg>
            Profilga qaytish
          </RouterLink>
          <h1>Resumeni tahrirlash</h1>
          <p>Ma'lumotlarni yangilang</p>
        </div>
      </div>

      <div class="create-main">
        <div class="create-main__inner">
          <div v-if="error" class="form-error-banner">{{ error }}</div>

          <form @submit.prevent="submit" class="create-form">
            <div class="form-section">
              <h3 class="form-section__title">
                <span class="form-section__num">1</span>
                Asosiy ma'lumotlar
              </h3>
              <div class="form-grid">
                <div class="form-field" :class="{ 'form-field--error': fieldErrors.name }">
                  <label>To'liq ism <span class="req">*</span></label>
                  <input v-model="form.name" type="text" placeholder="Abdullayev Alisher" required />
                  <span v-if="fieldErrors.name" class="form-field__err">{{ fieldErrors.name }}</span>
                </div>
                <div class="form-field" :class="{ 'form-field--error': fieldErrors.title }">
                  <label>Kasb / Mutaxassislik <span class="req">*</span></label>
                  <input v-model="form.title" type="text" placeholder="Elektrik usta, 10 yil tajriba" required />
                  <span v-if="fieldErrors.title" class="form-field__err">{{ fieldErrors.title }}</span>
                </div>
                <div class="form-field" :class="{ 'form-field--error': fieldErrors.contact }">
                  <label>Aloqa (telefon) <span class="req">*</span></label>
                  <input v-model="form.contact" type="text" placeholder="+998 90 123 45 67" required />
                  <span v-if="fieldErrors.contact" class="form-field__err">{{ fieldErrors.contact }}</span>
                </div>
              </div>

              <!-- Rasm yuklash -->
              <div class="form-field form-field--full" :class="{ 'form-field--error': fieldErrors.photo || photoUploadError }">
                <label>Rasm <span class="req">*</span></label>
                <label class="photo-upload-area" :class="{ 'photo-upload-area--has': photoPreview, 'photo-upload-area--loading': photoLoading }">
                  <input
                    type="file"
                    accept="image/jpeg,image/png,image/gif,image/webp"
                    class="photo-file-input"
                    :disabled="photoLoading"
                    @change="uploadPhoto"
                  />
                  <template v-if="photoLoading">
                    <div class="photo-state">
                      <div class="photo-spinner"></div>
                      <span>Yuklanmoqda...</span>
                    </div>
                  </template>
                  <template v-else-if="photoPreview">
                    <img :src="photoPreview" class="photo-preview-img" alt="preview" />
                    <div class="photo-change-overlay">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                        <polyline points="17 8 12 3 7 8"/>
                        <line x1="12" y1="3" x2="12" y2="15"/>
                      </svg>
                      <span>O'zgartirish</span>
                    </div>
                  </template>
                  <template v-else>
                    <div class="photo-state">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <rect x="3" y="3" width="18" height="18" rx="2"/>
                        <circle cx="8.5" cy="8.5" r="1.5"/>
                        <polyline points="21 15 16 10 5 21"/>
                      </svg>
                      <span class="photo-state__title">Rasm yuklash</span>
                      <span class="photo-state__hint">JPEG, PNG, GIF, WEBP · max 3 MB</span>
                    </div>
                  </template>
                </label>
                <span v-if="photoUploadError" class="form-field__err">{{ photoUploadError }}</span>
                <span v-else-if="fieldErrors.photo" class="form-field__err">{{ fieldErrors.photo }}</span>
              </div>

              <div class="form-grid">
                <div class="form-field" :class="{ 'form-field--error': fieldErrors.price }">
                  <label>Narx (so'm / oyiga)</label>
                  <input v-model="form.price" type="number" placeholder="3000000" min="0" />
                </div>
                <div class="form-field">
                  <label>Tajriba (yil)</label>
                  <input v-model="form.experience_year" type="number" placeholder="5" min="0" max="60" />
                </div>
              </div>

              <div class="form-field form-field--full" :class="{ 'form-field--error': fieldErrors.skills }">
                <label>Ko'nikmalar <span class="req">*</span></label>
                <input v-model="form.skills" type="text" placeholder="Excel, Word, Elektr o'rnatish" required />
                <span class="form-field__hint">Vergul bilan ajrating</span>
                <span v-if="fieldErrors.skills" class="form-field__err">{{ fieldErrors.skills }}</span>
              </div>

              <div class="form-field form-field--full" :class="{ 'form-field--error': fieldErrors.text }">
                <label>O'zim haqimda <span class="req">*</span></label>
                <textarea v-model="form.text" rows="5" placeholder="Tajribangiz haqida..." required></textarea>
                <span v-if="fieldErrors.text" class="form-field__err">{{ fieldErrors.text }}</span>
              </div>
            </div>

            <div class="form-section">
              <h3 class="form-section__title">
                <span class="form-section__num">2</span>
                Joylashuv
              </h3>
              <div class="form-grid">
                <div class="form-field">
                  <label>Viloyat</label>
                  <select :value="form.region_id || ''" @change="selectRegion($event.target.value)">
                    <option value="">Tanlang</option>
                    <option v-for="r in regions" :key="r.id" :value="r.id">{{ r.name }}</option>
                  </select>
                </div>
                <div class="form-field" v-if="districts.length">
                  <label>Tuman</label>
                  <select :value="form.district_id || ''" @change="selectDistrict($event.target.value)">
                    <option value="">Tanlang</option>
                    <option v-for="d in districts" :key="d.id" :value="d.id">{{ d.name }}</option>
                  </select>
                </div>
                <div class="form-field" v-if="mahallas.length">
                  <label>Mahalla</label>
                  <select :value="form.mahalla_id || ''" @change="form.mahalla_id = $event.target.value ? Number($event.target.value) : null">
                    <option value="">Tanlang</option>
                    <option v-for="m in mahallas" :key="m.id" :value="m.id">{{ m.name }}</option>
                  </select>
                </div>
                <div class="form-field form-field--full" :class="{ 'form-field--error': fieldErrors.adress }">
                  <label>To'liq manzil <span class="req">*</span></label>
                  <input v-model="form.adress" type="text" placeholder="Toshkent, Chilonzor tumani..." required />
                  <span v-if="fieldErrors.adress" class="form-field__err">{{ fieldErrors.adress }}</span>
                </div>
              </div>
            </div>

            <div class="form-section" v-if="categories.length">
              <h3 class="form-section__title">
                <span class="form-section__num">3</span>
                Kategoriyalar
              </h3>
              <div class="cat-chips">
                <button v-for="cat in categories" :key="cat.id" type="button" class="cat-chip"
                  :class="{ 'cat-chip--active': form.category_ids.includes(cat.id) }"
                  @click="toggleCategory(cat.id)">
                  {{ cat.name }}
                </button>
              </div>
            </div>

            <div class="form-section">
              <h3 class="form-section__title">
                <span class="form-section__num">{{ categories.length ? '4' : '3' }}</span>
                Holat
              </h3>
              <label class="toggle-row">
                <div>
                  <div class="toggle-row__label">Resumeni faol qilish</div>
                  <div class="toggle-row__sub">Faol bo'lsa, boshqalar ko'rishi mumkin</div>
                </div>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="form.is_active" />
                  <span class="toggle-switch__slider"></span>
                </label>
              </label>
            </div>

            <div class="form-actions">
              <RouterLink :to="`/${route.params.lang}/profile`" class="form-cancel">Bekor qilish</RouterLink>
              <button type="submit" class="form-submit" :disabled="loading || photoLoading">
                <span v-if="loading" class="btn-spinner"></span>
                {{ loading ? 'Saqlanmoqda...' : 'O\'zgarishlarni saqlash' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.create-page { min-height: 100vh; background: #f3f4f6; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; }

.page-loading { display: flex; justify-content: center; align-items: center; height: 60vh; }
.spinner { width: 44px; height: 44px; border: 3px solid #e5e7eb; border-top-color: #2563eb; border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.create-hero { background: linear-gradient(135deg, #1a3fbd 0%, #1d56db 40%, #2563eb 100%); padding: 32px 20px 40px; }
.create-hero__inner { max-width: 800px; margin: 0 auto; }
.create-back { display: inline-flex; align-items: center; gap: 6px; color: rgba(255,255,255,0.8); text-decoration: none; font-size: 14px; font-weight: 600; margin-bottom: 16px; transition: color 0.2s; }
.create-back svg { width: 16px; height: 16px; }
.create-back:hover { color: #fff; }
.create-hero__inner h1 { font-size: 28px; font-weight: 800; color: #fff; margin-bottom: 6px; }
.create-hero__inner p { font-size: 14px; color: rgba(255,255,255,0.75); }

.create-main { padding: 28px 20px 48px; }
.create-main__inner { max-width: 800px; margin: 0 auto; }

.form-error-banner { background: #fef2f2; border: 1px solid #fecaca; color: #dc2626; padding: 14px 18px; border-radius: 12px; font-size: 14px; margin-bottom: 20px; }

.create-form { display: flex; flex-direction: column; gap: 20px; }

.form-section { background: #fff; border: 1px solid #e5e7eb; border-radius: 16px; padding: 24px; display: flex; flex-direction: column; gap: 16px; }

.form-section__title { display: flex; align-items: center; gap: 10px; font-size: 16px; font-weight: 700; color: #111827; margin-bottom: 4px; }
.form-section__num { width: 28px; height: 28px; border-radius: 50%; background: linear-gradient(135deg, #1d4ed8, #2563eb); color: #fff; font-size: 13px; font-weight: 700; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }

.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.form-field { display: flex; flex-direction: column; gap: 6px; }
.form-field--full { grid-column: 1 / -1; }
.form-field label { font-size: 13px; font-weight: 600; color: #374151; }
.req { color: #ef4444; }

.form-field input, .form-field select, .form-field textarea {
  padding: 11px 14px; border: 1.5px solid #e5e7eb; border-radius: 10px; font-size: 14px; color: #111827;
  background: #f9fafb; outline: none; font-family: inherit; transition: border-color 0.2s, background 0.2s; resize: vertical; width: 100%; box-sizing: border-box;
}
.form-field input:focus, .form-field select:focus, .form-field textarea:focus { border-color: #2563eb; background: #fff; }
.form-field input::placeholder, .form-field textarea::placeholder { color: #9ca3af; }
.form-field--error input, .form-field--error select, .form-field--error textarea { border-color: #ef4444; }
.form-field__hint { font-size: 11px; color: #9ca3af; }
.form-field__err { font-size: 12px; color: #ef4444; font-weight: 500; }

.cat-chips { display: flex; flex-wrap: wrap; gap: 8px; }
.cat-chip { padding: 8px 16px; border: 1.5px solid #e5e7eb; border-radius: 20px; font-size: 13px; font-weight: 600; color: #374151; background: #f9fafb; cursor: pointer; font-family: inherit; transition: all 0.2s; }
.cat-chip:hover { border-color: #93c5fd; background: #eff6ff; color: #1d4ed8; }
.cat-chip--active { background: linear-gradient(135deg, #1d4ed8, #2563eb); color: #fff; border-color: transparent; box-shadow: 0 2px 8px rgba(37,99,235,0.3); }

.toggle-row { display: flex; align-items: center; justify-content: space-between; cursor: pointer; }
.toggle-row__label { font-size: 14px; font-weight: 600; color: #374151; }
.toggle-row__sub { font-size: 12px; color: #9ca3af; margin-top: 2px; }
.toggle-switch { position: relative; display: inline-block; width: 44px; height: 24px; flex-shrink: 0; }
.toggle-switch input { opacity: 0; width: 0; height: 0; }
.toggle-switch__slider { position: absolute; inset: 0; background: #e5e7eb; border-radius: 24px; cursor: pointer; transition: all 0.3s; }
.toggle-switch__slider::before { content: ''; position: absolute; width: 18px; height: 18px; left: 3px; bottom: 3px; background: #fff; border-radius: 50%; transition: all 0.3s; box-shadow: 0 1px 4px rgba(0,0,0,0.2); }
.toggle-switch input:checked + .toggle-switch__slider { background: #2563eb; }
.toggle-switch input:checked + .toggle-switch__slider::before { transform: translateX(20px); }

.form-actions { display: flex; align-items: center; justify-content: flex-end; gap: 12px; }
.form-cancel { padding: 12px 24px; border: 1.5px solid #e5e7eb; border-radius: 10px; font-size: 14px; font-weight: 600; color: #374151; text-decoration: none; transition: all 0.2s; }
.form-cancel:hover { border-color: #93c5fd; color: #1d4ed8; }
.form-submit { display: flex; align-items: center; gap: 8px; background: linear-gradient(135deg, #1d4ed8, #2563eb); color: #fff; border: none; border-radius: 10px; padding: 12px 28px; font-size: 15px; font-weight: 700; cursor: pointer; font-family: inherit; transition: opacity 0.2s; }
.form-submit:hover:not(:disabled) { opacity: 0.9; }
.form-submit:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-spinner { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.4); border-top-color: #fff; border-radius: 50%; animation: spin 0.7s linear infinite; }

.photo-upload-area {
  position: relative; display: flex; align-items: center; justify-content: center;
  min-height: 160px; border: 2px dashed #d1d5db; border-radius: 12px;
  background: #f9fafb; cursor: pointer; overflow: hidden; transition: border-color 0.2s, background 0.2s;
}
.photo-upload-area:hover { border-color: #2563eb; background: #eff6ff; }
.photo-upload-area--loading { cursor: not-allowed; opacity: 0.7; }
.photo-upload-area--has { border-style: solid; border-color: #2563eb; }
.form-field--error .photo-upload-area { border-color: #ef4444; }
.photo-file-input { position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; height: 100%; padding: 0 !important; border: none !important; background: transparent !important; }
.photo-file-input:disabled { cursor: not-allowed; }
.photo-state { display: flex; flex-direction: column; align-items: center; gap: 6px; padding: 20px; pointer-events: none; }
.photo-state svg { width: 40px; height: 40px; color: #9ca3af; }
.photo-state__title { font-size: 14px; font-weight: 600; color: #374151; }
.photo-state__hint  { font-size: 12px; color: #9ca3af; }
.photo-spinner { width: 32px; height: 32px; border: 3px solid #dbeafe; border-top-color: #2563eb; border-radius: 50%; animation: spin 0.7s linear infinite; }
.photo-preview-img { width: 100%; height: 160px; object-fit: cover; display: block; }
.photo-change-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.45); display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 6px; opacity: 0; transition: opacity 0.2s; pointer-events: none; }
.photo-upload-area:hover .photo-change-overlay { opacity: 1; }
.photo-change-overlay svg  { width: 28px; height: 28px; color: #fff; }
.photo-change-overlay span { font-size: 13px; font-weight: 600; color: #fff; }

@media (max-width: 640px) {
  .form-grid { grid-template-columns: 1fr; }
  .form-actions { flex-direction: column-reverse; }
  .form-cancel, .form-submit { width: 100%; justify-content: center; text-align: center; }
}
</style>
