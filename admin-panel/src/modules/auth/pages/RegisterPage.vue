<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../store/authStore'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const form = ref({
  full_name: '',
  phone: '',
  password: '',
  role: 'user',
})
const loading = ref(false)
const error = ref('')

const submit = async () => {
  error.value = ''
  loading.value = true
  try {
    await auth.registerAction(form.value)
    router.push({ name: 'profile' })
  } catch (e) {
    error.value = e.response?.data?.message || e.response?.data?.error || "Ro'yxatdan o'tish xatoligi"
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <RouterLink :to="`/${route.params.lang || 'uz'}`" class="auth-logo">
        <div class="auth-logo__box">UT</div>
        <span>UstaTop</span>
      </RouterLink>

      <h1 class="auth-card__title">Ro'yxatdan o'tish</h1>
      <p class="auth-card__sub">Yangi hisob yarating</p>

      <div v-if="error" class="auth-error">{{ error }}</div>

      <form @submit.prevent="submit" class="auth-form">
        <div class="auth-field">
          <label>To'liq ism</label>
          <div class="auth-input-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
            <input v-model="form.full_name" type="text" placeholder="Abdullayev Alisher" required minlength="2" />
          </div>
        </div>

        <div class="auth-field">
          <label>Telefon raqam</label>
          <div class="auth-input-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.69 12 19.79 19.79 0 0 1 1.61 3.45 2 2 0 0 1 3.6 1.27h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L7.91 8.37a16 16 0 0 0 5.72 5.72l.87-.87a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 21.04 16z"/>
            </svg>
            <input v-model="form.phone" type="tel" placeholder="+998 90 123 45 67" required />
          </div>
        </div>

        <div class="auth-field">
          <label>Parol</label>
          <div class="auth-input-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
            <input v-model="form.password" type="password" placeholder="Kamida 6 ta belgi" required minlength="6" />
          </div>
        </div>

        <div class="auth-field">
          <label>Rol tanlang</label>
          <div class="role-group">
            <label class="role-option" :class="{ 'role-option--active': form.role === 'user' }">
              <input type="radio" v-model="form.role" value="user" hidden />
              <div class="role-option__icon">👷</div>
              <div>
                <div class="role-option__title">Usta / Ishchi</div>
                <div class="role-option__sub">Resume joylash</div>
              </div>
            </label>
            <label class="role-option" :class="{ 'role-option--active': form.role === 'employer' }">
              <input type="radio" v-model="form.role" value="employer" hidden />
              <div class="role-option__icon">🏢</div>
              <div>
                <div class="role-option__title">Ish beruvchi</div>
                <div class="role-option__sub">Vakansiya e'lon qilish</div>
              </div>
            </label>
          </div>
        </div>

        <button type="submit" class="auth-btn" :disabled="loading">
          <span v-if="loading" class="auth-btn__spinner"></span>
          {{ loading ? "Ro'yxatdan o'tilmoqda..." : "Ro'yxatdan o'tish" }}
        </button>
      </form>

      <p class="auth-switch">
        Hisobingiz bormi?
        <RouterLink :to="`/${route.params.lang || 'uz'}/login`">Kirish</RouterLink>
      </p>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a3fbd 0%, #1d56db 40%, #2563eb 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.auth-card {
  background: #fff;
  border-radius: 20px;
  padding: 40px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.2);
}

.auth-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  margin-bottom: 28px;
}

.auth-logo__box {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 800;
  font-size: 14px;
}

.auth-logo span {
  font-size: 20px;
  font-weight: 800;
  color: #111827;
}

.auth-card__title { font-size: 26px; font-weight: 800; color: #111827; margin-bottom: 6px; }
.auth-card__sub { font-size: 14px; color: #6b7280; margin-bottom: 24px; }

.auth-error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 13px;
  margin-bottom: 20px;
}

.auth-form { display: flex; flex-direction: column; gap: 16px; }
.auth-field { display: flex; flex-direction: column; gap: 6px; }

.auth-field label {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
}

.auth-input-wrap { position: relative; display: flex; align-items: center; }

.auth-input-wrap svg {
  position: absolute;
  left: 12px;
  width: 16px;
  height: 16px;
  color: #9ca3af;
  pointer-events: none;
}

.auth-input-wrap input {
  width: 100%;
  padding: 12px 12px 12px 38px;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  font-size: 14px;
  color: #111827;
  background: #f9fafb;
  outline: none;
  font-family: inherit;
  transition: border-color 0.2s, background 0.2s;
  box-sizing: border-box;
}

.auth-input-wrap input:focus { border-color: #2563eb; background: #fff; }
.auth-input-wrap input::placeholder { color: #9ca3af; }

.role-group { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }

.role-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  border: 1.5px solid #e5e7eb;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  background: #f9fafb;
}

.role-option:hover { border-color: #93c5fd; background: #eff6ff; }

.role-option--active {
  border-color: #2563eb;
  background: #eff6ff;
  box-shadow: 0 0 0 3px rgba(37,99,235,0.1);
}

.role-option__icon { font-size: 22px; }

.role-option__title {
  font-size: 13px;
  font-weight: 700;
  color: #111827;
}

.role-option__sub { font-size: 11px; color: #6b7280; margin-top: 1px; }

.auth-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  color: #fff;
  border: none;
  border-radius: 12px;
  padding: 14px;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  margin-top: 4px;
  transition: opacity 0.2s;
}

.auth-btn:hover:not(:disabled) { opacity: 0.9; }
.auth-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.auth-btn__spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.4);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.auth-switch {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #6b7280;
}

.auth-switch a { color: #2563eb; font-weight: 600; text-decoration: none; }
.auth-switch a:hover { text-decoration: underline; }
</style>
