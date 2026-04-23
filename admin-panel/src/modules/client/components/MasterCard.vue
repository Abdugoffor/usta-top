<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SharePanel from '@/shared/components/SharePanel.vue'
import { useI18n } from '@/shared/composables/useI18n'

const props = defineProps({
  item: { type: Object, required: true }
})

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const initials = computed(() => {
  const name = props.item.name || ''
  return name.split(' ').slice(0, 2).map(w => w[0] || '').join('').toUpperCase() || 'U'
})

const avatarColor = computed(() => {
  const colors = ['#2563eb','#7c3aed','#059669','#d97706','#dc2626','#0891b2','#9333ea']
  const idx = (props.item.id || 0) % colors.length
  return colors[idx]
})

const price = computed(() => {
  const p = props.item.price
  if (!p) return null
  return Number(p).toLocaleString('uz-UZ')
})

const location = computed(() => {
  const parts = [props.item.region_name, props.item.district_name].filter(Boolean)
  return parts.join(', ') || props.item.adress || ''
})

const itemUrl = computed(() => `${window.location.origin}/${route.params.lang}/masters/${props.item.slug}`)

const goDetail = () => {
  if (props.item.slug) {
    router.push({
      name: 'master-detail',
      params: { lang: route.params.lang, slug: props.item.slug },
      query: { ...route.query },
    })
  }
}
</script>

<template>
  <div class="master-card" @click="goDetail">
    <div class="master-card__top">
      <div class="master-card__avatar" :style="{ background: avatarColor }">
        {{ initials }}
      </div>
      <div class="master-card__info">
        <h3 class="master-card__name">{{ item.name }}</h3>
        <p class="master-card__title">{{ item.title }}</p>
      </div>
      <span class="master-card__badge" :class="item.is_active ? 'badge--active' : 'badge--inactive'">
        {{ item.is_active ? t('status_active') : t('status_inactive') }}
      </span>
    </div>

    <p v-if="item.text" class="master-card__desc">{{ item.text }}</p>

    <div class="master-card__tags">
      <span v-if="item.category_name" class="master-card__tag">{{ item.category_name }}</span>
      <span v-if="item.skills" class="master-card__skill">{{ item.skills }}</span>
    </div>

    <div class="master-card__footer">
      <div class="master-card__meta">
        <span v-if="location" class="master-card__location">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
            <circle cx="12" cy="9" r="2.5"/>
          </svg>
          {{ location }}
        </span>
      </div>
      <div class="master-card__bottom">
        <span v-if="price" class="master-card__price">{{ price }} {{ t('price_sum') }}</span>
        <div class="master-card__actions">
          <span v-if="item.experience_year" class="master-card__exp">{{ item.experience_year }} {{ t('label_years_experience') }}</span>
          <span v-if="item.views_count" class="master-card__views">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
              <circle cx="12" cy="12" r="3"/>
            </svg>
            {{ item.views_count }}
          </span>
          <SharePanel :url="itemUrl" :title="item.name" :phone="item.contact" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.master-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.25s ease;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}

.master-card:hover {
  border-color: #93c5fd;
  box-shadow: 0 4px 20px rgba(37,99,235,0.1);
  transform: translateY(-2px);
}

.master-card__top {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.master-card__avatar {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
  font-weight: 700;
  flex-shrink: 0;
}

.master-card__info {
  flex: 1;
  min-width: 0;
}

.master-card__name {
  font-size: 15px;
  font-weight: 700;
  color: #111827;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.master-card__title {
  font-size: 13px;
  color: #6b7280;
  margin-top: 2px;
}

.master-card__badge {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 20px;
  flex-shrink: 0;
}

.badge--active {
  background: #d1fae5;
  color: #065f46;
}

.badge--inactive {
  background: #f3f4f6;
  color: #6b7280;
}

.master-card__desc {
  font-size: 13px;
  color: #4b5563;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.master-card__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.master-card__tag {
  background: #eff6ff;
  color: #1d4ed8;
  font-size: 12px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 8px;
}

.master-card__skill {
  background: #f0fdf4;
  color: #166534;
  font-size: 12px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 8px;
}

.master-card__footer {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: auto;
}

.master-card__meta {
  display: flex;
  align-items: center;
}

.master-card__location {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #9ca3af;
}

.master-card__location svg {
  width: 13px;
  height: 13px;
  flex-shrink: 0;
}

.master-card__bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.master-card__actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.master-card__price {
  font-size: 17px;
  font-weight: 800;
  color: #1d4ed8;
}

.master-card__exp {
  font-size: 12px;
  color: #9ca3af;
}

.master-card__views {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #9ca3af;
}

.master-card__views svg {
  width: 13px;
  height: 13px;
}
</style>
