<template>
  <footer class="site-footer">
    <div class="site-footer__inner">

      <div class="site-footer__top">
        <!-- Brand -->
        <div class="site-footer__brand">
          <RouterLink to="/" class="site-footer__logo-wrap">
            <div class="site-footer__logo">UT</div>
            <span class="site-footer__name">Usta<strong>Top</strong></span>
          </RouterLink>
          <p class="site-footer__tagline">
            {{ t('footer_tagline') }}
          </p>
          <div class="site-footer__socials">
            <a href="#" class="site-footer__social" aria-label="Telegram">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M11.944 0A12 12 0 1 0 24 12 12 12 0 0 0 11.944 0zm5.69 8.27-2.037 9.6c-.153.67-.553.832-1.12.517l-3.1-2.284-1.495 1.438c-.166.165-.305.305-.625.305l.222-3.18 5.77-5.21c.25-.222-.055-.345-.39-.123L7.29 14.44l-3.06-.952c-.666-.207-.68-.666.14-.987l11.95-4.607c.554-.2 1.038.135.857.987z"/>
              </svg>
            </a>
            <a href="#" class="site-footer__social" aria-label="Instagram">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zM12 0C8.741 0 8.333.014 7.053.072 2.695.272.273 2.69.073 7.052.014 8.333 0 8.741 0 12c0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98C8.333 23.986 8.741 24 12 24c3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98C15.668.014 15.259 0 12 0zm0 5.838a6.162 6.162 0 1 0 0 12.324 6.162 6.162 0 0 0 0-12.324zM12 16a4 4 0 1 1 0-8 4 4 0 0 1 0 8zm6.406-11.845a1.44 1.44 0 1 0 0 2.881 1.44 1.44 0 0 0 0-2.881z"/>
              </svg>
            </a>
            <a href="#" class="site-footer__social" aria-label="Facebook">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M24 12.073C24 5.405 18.627 0 12 0S0 5.405 0 12.073C0 18.1 4.388 23.094 10.125 24v-8.437H7.078v-3.49h3.047V9.41c0-3.025 1.79-4.697 4.533-4.697 1.312 0 2.686.236 2.686.236v2.97h-1.513c-1.491 0-1.956.93-1.956 1.886v2.267h3.328l-.532 3.49h-2.796V24C19.612 23.094 24 18.1 24 12.073z"/>
              </svg>
            </a>
          </div>
        </div>

        <!-- Navigation -->
        <div class="site-footer__col">
          <h4 class="site-footer__col-title">{{ t('footer_nav_title') }}</h4>
          <ul class="site-footer__links">
            <li><RouterLink to="/" class="site-footer__link">{{ t('footer_masters_link') }}</RouterLink></li>
            <li><RouterLink to="/vacancies" class="site-footer__link">{{ t('footer_vacancies_link') }}</RouterLink></li>
            <li><RouterLink to="/login" class="site-footer__link">{{ t('footer_login_link') }}</RouterLink></li>
            <li><RouterLink to="/register" class="site-footer__link">{{ t('footer_register_link') }}</RouterLink></li>
          </ul>
        </div>

        <!-- For users -->
        <div class="site-footer__col">
          <h4 class="site-footer__col-title">{{ t('footer_for_users_title') }}</h4>
          <ul class="site-footer__links">
            <li><RouterLink to="/profile/resumes/create" class="site-footer__link">{{ t('footer_add_resume') }}</RouterLink></li>
            <li><RouterLink to="/profile/vacancies/create" class="site-footer__link">{{ t('footer_add_vacancy') }}</RouterLink></li>
            <li><RouterLink to="/profile" class="site-footer__link">{{ t('footer_cabinet') }}</RouterLink></li>
          </ul>
        </div>

        <!-- Contact -->
        <div class="site-footer__col">
          <h4 class="site-footer__col-title">{{ t('footer_contact_title') }}</h4>
          <ul class="site-footer__links">
            <li>
              <a href="tel:+998901234567" class="site-footer__link site-footer__link--icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.69 12a19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 3.6 1.24h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L7.91 8.8a16 16 0 0 0 6.29 6.29l.91-.91a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
                </svg>
                +998 90 123 45 67
              </a>
            </li>
            <li>
              <a href="mailto:info@ustatop.uz" class="site-footer__link site-footer__link--icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                  <polyline points="22,6 12,13 2,6"/>
                </svg>
                info@ustatop.uz
              </a>
            </li>
            <li>
              <a href="#" class="site-footer__link site-footer__link--icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 1 1 18 0z"/>
                  <circle cx="12" cy="10" r="3"/>
                </svg>
                {{ t('footer_location') }}
              </a>
            </li>
          </ul>
        </div>
      </div>

      <!-- Bottom bar -->
      <div class="site-footer__bottom">
        <p class="site-footer__copy">© {{ year }} UstaTop. {{ t('footer_rights') }}</p>
        <div class="site-footer__bottom-links">
          <a href="#" class="site-footer__bottom-link">{{ t('footer_privacy') }}</a>
          <span class="site-footer__dot"></span>
          <a href="#" class="site-footer__bottom-link">{{ t('footer_terms') }}</a>
        </div>
      </div>

    </div>
  </footer>
</template>

<script setup>
import { useI18n } from '@/shared/composables/useI18n'
const { t } = useI18n()
const year = new Date().getFullYear()
</script>

<style scoped>
.site-footer {
  background: #fff;
  border-top: 1px solid #e5e7eb;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.site-footer__inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 20px;
}

/* ─── Top section ──────────────────────────────────── */
.site-footer__top {
  display: grid;
  grid-template-columns: 1.6fr 1fr 1fr 1fr;
  gap: 48px;
  padding: 52px 0 40px;
  border-bottom: 1px solid #f3f4f6;
}

/* Brand */
.site-footer__logo-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  margin-bottom: 14px;
}

.site-footer__logo {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: linear-gradient(135deg, #1d4ed8, #2563eb);
  color: #fff;
  font-size: 13px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  letter-spacing: 0.5px;
}

.site-footer__name {
  font-size: 18px;
  font-weight: 500;
  color: #111827;
}

.site-footer__name strong {
  font-weight: 800;
}

.site-footer__tagline {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.6;
  max-width: 260px;
  margin-bottom: 22px;
}

/* Socials */
.site-footer__socials {
  display: flex;
  gap: 8px;
}

.site-footer__social {
  width: 36px;
  height: 36px;
  border-radius: 9px;
  border: 1.5px solid #e5e7eb;
  background: #f9fafb;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  transition: all 0.2s ease;
}

.site-footer__social svg {
  width: 16px;
  height: 16px;
}

.site-footer__social:hover {
  background: #2563eb;
  border-color: #2563eb;
  color: #fff;
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.3);
}

/* Columns */
.site-footer__col-title {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #111827;
  margin-bottom: 16px;
}

.site-footer__links {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.site-footer__link {
  font-size: 14px;
  color: #6b7280;
  text-decoration: none;
  transition: color 0.2s ease;
  display: inline-flex;
  align-items: center;
}

.site-footer__link:hover {
  color: #2563eb;
}

.site-footer__link--icon {
  gap: 8px;
}

.site-footer__link--icon svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  color: #9ca3af;
}

.site-footer__link--icon:hover svg {
  color: #2563eb;
}

/* ─── Bottom bar ───────────────────────────────────── */
.site-footer__bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 0;
  gap: 12px;
  flex-wrap: wrap;
}

.site-footer__copy {
  font-size: 13px;
  color: #9ca3af;
}

.site-footer__bottom-links {
  display: flex;
  align-items: center;
  gap: 12px;
}

.site-footer__bottom-link {
  font-size: 13px;
  color: #9ca3af;
  text-decoration: none;
  transition: color 0.2s ease;
}

.site-footer__bottom-link:hover {
  color: #2563eb;
}

.site-footer__dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: #d1d5db;
}

/* ─── Responsive ───────────────────────────────────── */
@media (max-width: 1024px) {
  .site-footer__top {
    grid-template-columns: 1fr 1fr;
    gap: 32px;
  }

  .site-footer__brand {
    grid-column: span 2;
  }
}

@media (max-width: 640px) {
  .site-footer__top {
    grid-template-columns: 1fr 1fr;
    gap: 24px;
    padding: 36px 0 28px;
  }

  .site-footer__brand {
    grid-column: span 2;
  }

  .site-footer__tagline {
    max-width: none;
  }

  .site-footer__bottom {
    flex-direction: column;
    align-items: flex-start;
    gap: 6px;
    padding: 14px 0;
  }
}

@media (max-width: 400px) {
  .site-footer__top {
    grid-template-columns: 1fr;
  }

  .site-footer__brand {
    grid-column: span 1;
  }
}
</style>
