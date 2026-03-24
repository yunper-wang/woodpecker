<template>
  <div class="text-wp-text-100 space-y-4">
    <ListItem
      v-for="secret in secrets"
      :key="secret.id"
      class="bg-wp-background-200! dark:bg-wp-background-200! items-center"
    >
      <span :title="secret.note">{{ secret.name }}</span>
      <Badge
        v-if="secret.edit === false"
        class="ml-2"
        :value="secret.org_id === 0 ? $t('global_level_secret') : $t('org_level_secret')"
      />
      <div class="md:display-unset ml-auto hidden space-x-2">
        <Badge v-for="event in secret.events" :key="event" :value="event" />
      </div>
      <div v-if="secret.edit !== false" class="flex items-center gap-2">
        <IconButton icon="edit" class="h-8 w-8 md:ml-2" :title="$t('secrets.edit')" @click="editSecret(secret)" />
        <IconButton
          icon="trash"
          class="hover:text-wp-error-100 h-8 w-8"
          :is-loading="isDeleting"
          :title="$t('secrets.delete')"
          @click="deleteSecret(secret)"
        />
      </div>
    </ListItem>

    <div v-if="loading" class="flex justify-center">
      <Icon name="spinner" class="animate-spin" />
    </div>
    <div v-else-if="secrets?.length === 0" class="ml-2">{{ $t('secrets.none') }}</div>
  </div>
</template>

<script lang="ts" setup>
import { toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import Badge from '~/components/atomic/Badge.vue';
import Icon from '~/components/atomic/Icon.vue';
import IconButton from '~/components/atomic/IconButton.vue';
import ListItem from '~/components/atomic/ListItem.vue';
import { useConfirm } from '~/compositions/useConfirm';
import type { Secret } from '~/lib/api/types';

const props = defineProps<{
  modelValue: (Secret & { edit?: boolean })[];
  isDeleting: boolean;
  loading: boolean;
}>();

const emit = defineEmits<{
  (event: 'edit', secret: Secret): void;
  (event: 'delete', secret: Secret): void;
}>();

const i18n = useI18n();
const { confirm } = useConfirm();

const secrets = toRef(props, 'modelValue');

function editSecret(secret: Secret) {
  emit('edit', secret);
}

async function deleteSecret(secret: Secret) {
  const ok = await confirm(i18n.t('secrets.delete_confirm'), i18n.t('secrets.delete'));
  if (!ok) {
    return;
  }
  emit('delete', secret);
}
</script>
