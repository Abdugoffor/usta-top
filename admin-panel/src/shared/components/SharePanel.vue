<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  url:   { type: String, default: '' },
  title: { type: String, default: '' },
  phone: { type: String, default: '' },
  dark:  { type: Boolean, default: false },
})

const open   = ref(false)
const copied = ref(false)
const root   = ref(null)

const shareUrl = computed(() => props.url || window.location.href)

const onTelegram = () => {
  window.open(
    `https://t.me/share/url?url=${encodeURIComponent(shareUrl.value)}&text=${encodeURIComponent(props.title)}`,
    '_blank'
  )
  open.value = false
}

const onInstagram = async () => {
  if (navigator.share) {
    try { await navigator.share({ title: props.title, url: shareUrl.value }) } catch {}
  } else {
    await copyFn()
  }
  open.value = false
}

const copyFn = async () => {
  try {
    await navigator.clipboard.writeText(shareUrl.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {}
}

const onCopy = async () => {
  await copyFn()
  open.value = false
}

const onOutside = (e) => {
  if (root.value && !root.value.contains(e.target)) open.value = false
}
onMounted(() => document.addEventListener('mousedown', onOutside))
onUnmounted(() => document.removeEventListener('mousedown', onOutside))
</script>

<template>
  <div class="sp" ref="root">
    <button
      class="sp__btn"
      :class="dark ? 'sp__btn--dark' : 'sp__btn--light'"
      @click.stop="open = !open"
      title="Ulashish"
    >
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
        <circle cx="18" cy="5" r="3"/>
        <circle cx="6" cy="12" r="3"/>
        <circle cx="18" cy="19" r="3"/>
        <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/>
        <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
      </svg>
    </button>

    <transition name="sp-anim">
      <div v-if="open" class="sp__panel" @click.stop>

        <a v-if="phone" :href="'tel:' + phone" class="sp__item" @click="open = false">
          <span class="sp__icon sp__icon--green">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.69 12 19.79 19.79 0 0 1 1.61 3.4 2 2 0 0 1 3.6 1.22h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L7.91 8.78A16 16 0 0 0 15 15.87l.85-.85a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/>
            </svg>
          </span>
          <div class="sp__meta">
            <span class="sp__name">Qo'ng'iroq</span>
            <span class="sp__sub">{{ phone }}</span>
          </div>
        </a>

        <button class="sp__item" @click="onTelegram">
          <span class="sp__icon sp__icon--tg">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.562 8.248-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.447 1.394c-.16.16-.295.295-.605.295l.213-3.053 5.56-5.023c.242-.213-.054-.333-.373-.12l-6.871 4.326-2.962-.924c-.643-.204-.657-.643.136-.953l11.57-4.461c.537-.194 1.006.131.833.941z"/>
            </svg>
          </span>
          <div class="sp__meta">
            <span class="sp__name">Telegram</span>
            <span class="sp__sub">Ulashish</span>
          </div>
        </button>

        <button class="sp__item" @click="onInstagram">
          <span class="sp__icon sp__icon--ig">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"/>
            </svg>
          </span>
          <div class="sp__meta">
            <span class="sp__name">Instagram</span>
            <span class="sp__sub">Ulashish</span>
          </div>
        </button>

        <button class="sp__item" @click="onCopy">
          <span class="sp__icon sp__icon--copy" :class="{ 'sp__icon--done': copied }">
            <svg v-if="!copied" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2"/>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
          </span>
          <div class="sp__meta">
            <span class="sp__name" :class="{ 'sp__name--done': copied }">
              {{ copied ? 'Nusxalandi!' : 'Havolani nusxalash' }}
            </span>
            <span class="sp__sub">URL</span>
          </div>
        </button>

      </div>
    </transition>
  </div>
</template>

<style scoped>
.sp { position: relative; display: inline-flex; }

.sp__btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}
.sp__btn svg { width: 18px; height: 18px; }

.sp__btn--dark {
  background: rgba(255,255,255,0.13);
  border: 1.5px solid rgba(255,255,255,0.25);
  color: #fff;
}
.sp__btn--dark:hover { background: rgba(255,255,255,0.24); }

.sp__btn--light {
  background: #f3f4f6;
  border: 1.5px solid #e5e7eb;
  color: #6b7280;
}
.sp__btn--light:hover { background: #eff6ff; border-color: #93c5fd; color: #2563eb; }

.sp__panel {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  box-shadow: 0 12px 40px rgba(0,0,0,0.14), 0 2px 8px rgba(0,0,0,0.05);
  min-width: 230px;
  z-index: 200;
  padding: 6px;
}

.sp__item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 9px 10px;
  border: none;
  border-radius: 10px;
  background: none;
  cursor: pointer;
  font-family: inherit;
  text-decoration: none;
  color: #111827;
  transition: background 0.15s;
  text-align: left;
}
.sp__item:hover { background: #f3f4f6; }
.sp__item + .sp__item { margin-top: 2px; }

.sp__icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
}
.sp__icon svg { width: 18px; height: 18px; }

.sp__icon--green { background: #dcfce7; color: #16a34a; }
.sp__icon--tg    { background: #e0f2fe; color: #0088cc; }
.sp__icon--ig    { background: linear-gradient(135deg, #f9ce34, #ee2a7b, #6228d7); color: #fff; }
.sp__icon--copy  { background: #f3f4f6; color: #374151; }
.sp__icon--done  { background: #dcfce7; color: #16a34a; }

.sp__meta { display: flex; flex-direction: column; gap: 1px; }
.sp__name { font-size: 13px; font-weight: 600; line-height: 1.3; }
.sp__name--done { color: #16a34a; }
.sp__sub  { font-size: 11px; color: #9ca3af; line-height: 1.2; }

.sp-anim-enter-active, .sp-anim-leave-active {
  transition: all 0.18s cubic-bezier(0.4, 0, 0.2, 1);
}
.sp-anim-enter-from { opacity: 0; transform: translateY(-6px) scale(0.97); }
.sp-anim-leave-to   { opacity: 0; transform: translateY(-4px) scale(0.97); }
</style>
