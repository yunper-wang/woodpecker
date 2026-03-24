<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="$emit('update:modelValue', false)"
      >
        <div
          class="bg-wp-background-100 dark:bg-wp-background-200 border-wp-background-400 w-full max-w-md rounded-xl border shadow-2xl"
        >
          <div class="flex items-start gap-3 p-6">
            <div class="bg-wp-error-100/10 mt-0.5 flex h-10 w-10 shrink-0 items-center justify-center rounded-full">
              <Icon name="alert" class="text-wp-error-100 h-5 w-5" />
            </div>
            <div class="flex-1">
              <h3 class="text-wp-text-100 mb-1 text-base font-semibold">{{ title }}</h3>
              <p class="text-wp-text-alt-100 text-sm">{{ text }}</p>
            </div>
          </div>
          <div class="border-wp-background-400 flex justify-end gap-3 border-t px-6 py-4">
            <Button color="gray" :text="cancelText || $t('cancel')" @click="$emit('update:modelValue', false)" />
            <Button color="red" :text="confirmText || $t('confirm_modal.confirm')" @click="onConfirm" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script lang="ts" setup>
import Button from '~/components/atomic/Button.vue';
import Icon from '~/components/atomic/Icon.vue';

const props = defineProps<{
  modelValue: boolean;
  title?: string;
  text: string;
  confirmText?: string;
  cancelText?: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', val: boolean): void;
  (e: 'confirm'): void;
}>();

function onConfirm() {
  emit('confirm');
  emit('update:modelValue', false);
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
