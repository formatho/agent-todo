<template>
  <div class="view-toggle">
    <button
      v-for="option in options"
      :key="option.value"
      :class="['toggle-button', { active: modelValue === option.value }]"
      :title="option.label"
      @click="$emit('update:modelValue', option.value)"
    >
      {{ option.icon }}
    </button>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: String,
    required: true
  },
  options: {
    type: Array,
    default: () => [
      { value: 'grid', label: 'Grid View', icon: '⊞' },
      { value: 'list', label: 'List View', icon: '☰' }
    ]
  }
})

defineEmits(['update:modelValue', 'change'])
</script>

<style scoped>
.view-toggle {
  display: flex;
  gap: 4px;
  background: #F3F4F6;
  padding: 4px;
  border-radius: 8px;
}

.toggle-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 18px;
}

.toggle-button:hover {
  background: rgba(0, 0, 0, 0.05);
}

.toggle-button.active {
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
</style>
