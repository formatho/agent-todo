<template>
  <div
    class="agent-avatar"
    :class="[`size-${size}`, { 'clickable': clickable }]"
    :style="avatarStyle"
    :title="agent.name"
    @click="handleClick"
  >
    {{ initials }}
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { generateAgentColor, getAgentInitials } from '../utils/agentColors'

const props = defineProps({
  agent: {
    type: Object,
    required: true
  },
  size: {
    type: String,
    default: 'medium',
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  },
  clickable: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

const initials = computed(() => getAgentInitials(props.agent.name))

const avatarStyle = computed(() => {
  // Generate consistent color based on agent ID or name
  const index = props.agent.id ? parseInt(props.agent.id.slice(-8), 16) % 10 : 0
  const colors = generateAgentColor(index)

  return {
    backgroundColor: colors.primary,
    color: 'white'
  }
})

const handleClick = () => {
  if (props.clickable) {
    emit('click', props.agent)
  }
}
</script>

<style scoped>
.agent-avatar {
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  text-transform: uppercase;
  flex-shrink: 0;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.agent-avatar.size-small {
  width: 24px;
  height: 24px;
  font-size: 10px;
}

.agent-avatar.size-medium {
  width: 32px;
  height: 32px;
  font-size: 12px;
}

.agent-avatar.size-large {
  width: 40px;
  height: 40px;
  font-size: 14px;
}

.agent-avatar.clickable {
  cursor: pointer;
}

.agent-avatar.clickable:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
