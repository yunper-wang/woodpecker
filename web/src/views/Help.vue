<template>
  <div class="max-w-4xl mx-auto px-4 py-8 text-wp-text-100">
    <h1 class="text-3xl font-bold mb-2">{{ $t('help.title') }}</h1>
    <p class="text-wp-text-alt-100 mb-8">{{ $t('help.subtitle') }}</p>

    <!-- 快速导航 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-10">
      <a
        v-for="section in sections"
        :key="section.id"
        :href="'#' + section.id"
        class="bg-wp-background-200 border border-wp-background-400 rounded-lg p-4 hover:border-wp-primary-200 transition-colors cursor-pointer text-center"
      >
        <Icon :name="section.icon" class="mx-auto mb-2 h-6 w-6 text-wp-primary-200" />
        <div class="text-sm font-medium">{{ section.title }}</div>
      </a>
    </div>

    <!-- 快速开始 -->
    <section :id="sections[0].id" class="mb-10">
      <h2 class="text-xl font-semibold mb-4 pb-2 border-b border-wp-background-400">🚀 {{ $t('help.quickstart.title') }}</h2>
      <div class="space-y-4">
        <div class="bg-wp-background-200 rounded-lg p-4">
          <h3 class="font-semibold mb-2">{{ $t('help.quickstart.step1.title') }}</h3>
          <p class="text-wp-text-alt-100 text-sm">{{ $t('help.quickstart.step1.desc') }}</p>
        </div>
        <div class="bg-wp-background-200 rounded-lg p-4">
          <h3 class="font-semibold mb-2">{{ $t('help.quickstart.step2.title') }}</h3>
          <p class="text-wp-text-alt-100 text-sm">{{ $t('help.quickstart.step2.desc') }}</p>
          <pre class="mt-2 bg-wp-background-300 rounded p-3 text-xs overflow-x-auto">{{ pipelineExample }}</pre>
        </div>
        <div class="bg-wp-background-200 rounded-lg p-4">
          <h3 class="font-semibold mb-2">{{ $t('help.quickstart.step3.title') }}</h3>
          <p class="text-wp-text-alt-100 text-sm">{{ $t('help.quickstart.step3.desc') }}</p>
        </div>
      </div>
    </section>

    <!-- 流水线配置 -->
    <section :id="sections[1].id" class="mb-10">
      <h2 class="text-xl font-semibold mb-4 pb-2 border-b border-wp-background-400">⚙️ {{ $t('help.pipeline.title') }}</h2>
      <div class="space-y-3">
        <div v-for="item in pipelineTopics" :key="item.key" class="bg-wp-background-200 rounded-lg p-4">
          <h3 class="font-semibold mb-1">{{ item.title }}</h3>
          <p class="text-wp-text-alt-100 text-sm">{{ item.desc }}</p>
        </div>
      </div>
    </section>

    <!-- 密钥管理 -->
    <section :id="sections[2].id" class="mb-10">
      <h2 class="text-xl font-semibold mb-4 pb-2 border-b border-wp-background-400">🔒 {{ $t('help.secrets.title') }}</h2>
      <div class="space-y-3">
        <div v-for="item in secretTopics" :key="item.key" class="bg-wp-background-200 rounded-lg p-4">
          <h3 class="font-semibold mb-1">{{ item.title }}</h3>
          <p class="text-wp-text-alt-100 text-sm">{{ item.desc }}</p>
        </div>
      </div>
    </section>

    <!-- 常见问题 -->
    <section :id="sections[3].id" class="mb-10">
      <h2 class="text-xl font-semibold mb-4 pb-2 border-b border-wp-background-400">❓ {{ $t('help.faq.title') }}</h2>
      <div class="space-y-2">
        <div
          v-for="faq in faqs"
          :key="faq.key"
          class="bg-wp-background-200 rounded-lg overflow-hidden"
        >
          <button
            type="button"
            class="w-full flex items-center justify-between p-4 text-left hover:bg-wp-background-300 transition-colors"
            @click="faq.open = !faq.open"
          >
            <span class="font-medium">{{ faq.q }}</span>
            <Icon name="chevron-right" class="h-4 w-4 transition-transform" :class="{ 'rotate-90': faq.open }" />
          </button>
          <div v-if="faq.open" class="px-4 pb-4 text-wp-text-alt-100 text-sm">
            {{ faq.a }}
          </div>
        </div>
      </div>
    </section>

    <!-- 键盘快捷键 -->
    <section :id="sections[4].id" class="mb-10">
      <h2 class="text-xl font-semibold mb-4 pb-2 border-b border-wp-background-400">⌨️ {{ $t('help.shortcuts.title') }}</h2>
      <div class="bg-wp-background-200 rounded-lg overflow-hidden">
        <table class="w-full text-sm">
          <thead class="bg-wp-background-300">
            <tr>
              <th class="text-left p-3">{{ $t('help.shortcuts.key') }}</th>
              <th class="text-left p-3">{{ $t('help.shortcuts.action') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="shortcut in shortcuts" :key="shortcut.key" class="border-t border-wp-background-300">
              <td class="p-3">
                <kbd class="bg-wp-background-300 border border-wp-background-400 rounded px-2 py-0.5 font-mono text-xs">{{ shortcut.key }}</kbd>
              </td>
              <td class="p-3 text-wp-text-alt-100">{{ shortcut.action }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
import { useI18n } from 'vue-i18n';

import Icon from '~/components/atomic/Icon.vue';
import { useWPTitle } from '~/compositions/useWPTitle';

const { t } = useI18n();
useWPTitle([t('help.title')]);

const pipelineExample = `steps:
  - name: 构建
    image: node:18
    commands:
      - npm install
      - npm run build

  - name: 测试
    image: node:18
    commands:
      - npm test`;

const sections = [
  { id: 'quickstart', title: t('help.quickstart.title'), icon: 'play' },
  { id: 'pipeline', title: t('help.pipeline.title'), icon: 'settings' },
  { id: 'secrets', title: t('help.secrets.nav'), icon: 'secret' },
  { id: 'faq', title: t('help.faq.title'), icon: 'question' },
  { id: 'shortcuts', title: t('help.shortcuts.title'), icon: 'duration' },
];

const pipelineTopics = [
  { key: 'steps', title: t('help.pipeline.steps.title'), desc: t('help.pipeline.steps.desc') },
  { key: 'events', title: t('help.pipeline.events.title'), desc: t('help.pipeline.events.desc') },
  { key: 'matrix', title: t('help.pipeline.matrix.title'), desc: t('help.pipeline.matrix.desc') },
  { key: 'services', title: t('help.pipeline.services.title'), desc: t('help.pipeline.services.desc') },
  { key: 'cron', title: t('help.pipeline.cron.title'), desc: t('help.pipeline.cron.desc') },
];

const secretTopics = [
  { key: 'repo', title: t('help.secrets.repo.title'), desc: t('help.secrets.repo.desc') },
  { key: 'org', title: t('help.secrets.org.title'), desc: t('help.secrets.org.desc') },
  { key: 'global', title: t('help.secrets.global.title'), desc: t('help.secrets.global.desc') },
  { key: 'usage', title: t('help.secrets.usage.title'), desc: t('help.secrets.usage.desc') },
];

const faqs = reactive([
  { key: 'agent', q: t('help.faq.agent.q'), a: t('help.faq.agent.a'), open: false },
  { key: 'fail', q: t('help.faq.fail.q'), a: t('help.faq.fail.a'), open: false },
  { key: 'env', q: t('help.faq.env.q'), a: t('help.faq.env.a'), open: false },
  { key: 'skip', q: t('help.faq.skip.q'), a: t('help.faq.skip.a'), open: false },
  { key: 'badge', q: t('help.faq.badge.q'), a: t('help.faq.badge.a'), open: false },
]);

const shortcuts = [
  { key: '?', action: t('help.shortcuts.show_help') },
  { key: 'g h', action: t('help.shortcuts.go_home') },
  { key: 'g s', action: t('help.shortcuts.go_settings') },
  { key: 'Esc', action: t('help.shortcuts.close') },
];
</script>
