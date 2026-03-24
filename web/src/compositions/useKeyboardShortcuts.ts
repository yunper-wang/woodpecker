import { onMounted, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';

export function useKeyboardShortcuts() {
  const router = useRouter();

  function handleKeydown(e: KeyboardEvent) {
    // 忽略输入框内的快捷键
    const tag = (e.target as HTMLElement).tagName;
    if (tag === 'INPUT' || tag === 'TEXTAREA' || (e.target as HTMLElement).isContentEditable) {
      return;
    }

    // g h → 仓库列表
    if (e.key === 'h' && (e as any)._gPrefix) {
      router.push({ name: 'repos' });
      return;
    }

    // g s → 用户设置
    if (e.key === 's' && (e as any)._gPrefix) {
      router.push({ name: 'user' });
      return;
    }

    // g 前缀
    if (e.key === 'g' && !e.ctrlKey && !e.metaKey && !e.altKey) {
      const syntheticFlag = { _gPrefix: true };
      const handler = (e2: KeyboardEvent) => {
        Object.assign(e2, syntheticFlag);
        window.removeEventListener('keydown', handler);
      };
      window.addEventListener('keydown', handler, { once: true });
      return;
    }

    // Escape → 返回
    if (e.key === 'Escape') {
      router.back();
      return;
    }
  }

  onMounted(() => window.addEventListener('keydown', handleKeydown));
  onBeforeUnmount(() => window.removeEventListener('keydown', handleKeydown));
}
