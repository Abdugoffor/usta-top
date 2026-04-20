<script setup>
defineProps({
  variant: {
    type: String,
    default: 'primary'
  },
  type: {
    type: String,
    default: 'button'
  }
})
</script>

<template>
  <button :type="type" :class="['btn', `btn--${variant}`]">
    <slot />
  </button>
</template>