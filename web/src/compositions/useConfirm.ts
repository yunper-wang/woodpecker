import { ref } from 'vue';

const visible = ref(false);
const text = ref('');
const title = ref('');
let resolveFunc: ((val: boolean) => void) | null = null;

export function useConfirm() {
  function confirm(message: string, titleStr = ''): Promise<boolean> {
    text.value = message;
    title.value = titleStr;
    visible.value = true;
    return new Promise((resolve) => {
      resolveFunc = resolve;
    });
  }

  function onConfirm() {
    visible.value = false;
    resolveFunc?.(true);
    resolveFunc = null;
  }

  function onCancel() {
    visible.value = false;
    resolveFunc?.(false);
    resolveFunc = null;
  }

  return { confirm, onConfirm, onCancel, visible, text, title };
}
