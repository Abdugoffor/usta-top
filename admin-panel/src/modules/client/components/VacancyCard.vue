<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  item: { type: Object, required: true }
})

const router = useRouter()

const initials = computed(() => {
  const name = props.item.name || props.item.title || ''
  return name.split(' ').slice(0, 2).map(w => w[0] || '').join('').toUpperCase() || 'IS'
})

const iconColor = computed(() => {
  const colors = ['#f59e0b','#10b981','#8b5cf6','#ef4444','#06b6d4','#f97316','#ec4899']
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

const timeAgo = computed(() => {
  if (!props.item.created_at) return ''
  const date = new Date(props.item.created_at)
  const now = new Date()
  const diff = Math.floor((now - date) / 1000)
  if (diff < 3600) return `${Math.floor(diff/60)} daqiqa oldin`
  if (diff < 86400) return `${Math.floor(diff/3600)} soat oldin`
  return `${Math.floor(diff/86400)} kun oldin`
})

const goDetail = () => {
  if (props.item.slug) router.push({ name: 'vacancy-detail', params: { slug: props.item.slug } })
}
</script>

<template>
  <div class="vacancy-card" @click="goDetail">
    <div class="vacancy-card__top">
      <div class="vacancy-card__icon" :style="{ background: iconColor + '22', color: iconColor }">
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path d="M20 6h-2.18c.07-.44.18-.88.18-1.36C18 2.96 16.04 1 13.64 1c-1.3 0-2.55.57-3.42 1.56L9 4 7.78 2.56C6.91 1.57 5.66 1 4.36 1 1.96 1 0 2.96 0 5.36c0 .48.11.92.18 1.36H0v2h20V6zm-7.37-2.21c.43-.51 1.06-.79 1.73-.79.97 0 1.64.67 1.64 1.64 0 .96-.77 1.72-1.64 1.72h-2.55l.82-.97.82-.97.18-.63zm-9.26.64C3.37 4.43 4.04 3.76 5.01 3.76c.67 0 1.3.28 1.73.79l.82.97.82.97-.18.63H5.64C4.77 7.12 4 6.36 4 5.4c0-.34.1-.65.24-.94-.14-.02-.24.01-.37.03zM0 19c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V9H0v10z"/>
        </svg>
      </div>
      <div class="vacancy-card__info">
        <h3 class="vacancy-card__title">{{ item.name || item.title }}</h3>
        <p class="vacancy-card__subtitle">{{ item.title !== item.name ? item.title : '' }}</p>
      </div>
      <span class="vacancy-card__badge" :class="item.is_active ? 'badge--active' : 'badge--inactive'">
        {{ item.is_active ? 'Faol' : 'Yopiq' }}
      </span>
    </div>

    <p v-if="item.text" class="vacancy-card__desc">{{ item.text }}</p>

    <div class="vacancy-card__footer">
      <span v-if="location" class="vacancy-card__location">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
          <circle cx="12" cy="9" r="2.5"/>
        </svg>
        {{ location }}
      </span>
      <div class="vacancy-card__bottom">
        <span v-if="price" class="vacancy-card__price">{{ price }} so'm</span>
        <span v-if="timeAgo" class="vacancy-card__time">{{ timeAgo }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vacancy-card {
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

.vacancy-card:hover {
  border-color: #93c5fd;
  box-shadow: 0 4px 20px rgba(37,99,235,0.1);
  transform: translateY(-2px);
}

.vacancy-card__top {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.vacancy-card__icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.vacancy-card__icon svg {
  width: 22px;
  height: 22px;
}

.vacancy-card__info {
  flex: 1;
  min-width: 0;
}

.vacancy-card__title {
  font-size: 15px;
  font-weight: 700;
  color: #111827;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.vacancy-card__subtitle {
  font-size: 13px;
  color: #6b7280;
  margin-top: 2px;
}

.vacancy-card__badge {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 20px;
  flex-shrink: 0;
}

.badge--active { background: #d1fae5; color: #065f46; }
.badge--inactive { background: #f3f4f6; color: #6b7280; }

.vacancy-card__desc {
  font-size: 13px;
  color: #4b5563;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.vacancy-card__footer {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: auto;
}

.vacancy-card__location {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #9ca3af;
}

.vacancy-card__location svg {
  width: 13px;
  height: 13px;
  flex-shrink: 0;
}

.vacancy-card__bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.vacancy-card__price {
  font-size: 17px;
  font-weight: 800;
  color: #1d4ed8;
}

.vacancy-card__time {
  font-size: 12px;
  color: #9ca3af;
}
</style>
