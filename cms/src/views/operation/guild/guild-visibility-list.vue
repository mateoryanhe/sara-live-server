<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.GuildVisibilityManagement') }}</span>
        </div>
      </template>

      <div class="content">
        <el-alert
            :closable="false"
            class="tip-alert"
            show-icon
            :title="t('pages.guildVisibilityList.tip')"
            type="info"
        />

        <el-form class="search-form" inline>
          <el-form-item :label="t('pages.guildVisibilityList.selectCmsUser')">
            <el-select
                v-model="selectedCmsUserId"
                clearable
                filterable
                remote
                reserve-keyword
                style="width: 320px"
                :loading="cmsUserLoading"
                :placeholder="t('pages.guildVisibilityList.selectCmsUserPlaceholder')"
                :remote-method="searchCmsUsers"
                @change="handleUserChange"
                @focus="searchCmsUsers('')"
            >
              <el-option
                  v-for="u in cmsUserOptions"
                  :key="u.id"
                  :label="`${u.name} (${u.id})`"
                  :value="u.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :disabled="!selectedCmsUserId" :loading="loading" @click="reloadForUser">
              {{ t('common.search') }}
            </el-button>
          </el-form-item>
        </el-form>

        <el-empty v-if="!selectedCmsUserId" :description="t('pages.guildVisibilityList.emptyHint')"/>

        <div v-else v-loading="loading" class="transfer-wrap">
          <el-transfer
              v-model="grantedIds"
              filterable
              :data="transferData"
              :titles="[
                t('pages.guildVisibilityList.leftTitle'),
                t('pages.guildVisibilityList.rightTitle'),
              ]"
              :button-texts="[
                t('pages.guildVisibilityList.toLeft'),
                t('pages.guildVisibilityList.toRight'),
              ]"
              :props="{key: 'key', label: 'label', disabled: 'disabled'}"
              :filter-placeholder="t('pages.guildVisibilityList.filterPlaceholder')"
              @change="handleTransferChange"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {useI18n} from 'vue-i18n'
import {guildApi} from '@/api/modules/guild'
import {cmsUserApi, type CMSUser} from '@/api/modules/cmsuser'
import type {Guild} from '@/types/api'
import {usePagePermission} from '@/composables/usePagePermission'

const {t} = useI18n()
const {can} = usePagePermission('GuildVisibilityManagement')

const loading = ref(false)
const syncing = ref(false)
const guildOptions = ref<Guild[]>([])
const grantedIds = ref<string[]>([])
const selectedCmsUserId = ref('')
const cmsUserLoading = ref(false)
const cmsUserOptions = ref<CMSUser[]>([])

const transferData = computed(() =>
    guildOptions.value.map((g) => ({
      key: g.id,
      label: `${g.name} (${g.id})`,
      disabled: false,
    })),
)

const fetchGuildOptions = async () => {
  try {
    const res = await guildApi.getGuildListForVisibility({pageIndex: 1, pageSize: 500, name: ''})
    guildOptions.value = res.data || []
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.guildVisibilityList.fetchGuildFailed'))
  }
}

const fetchGrantedForUser = async () => {
  if (!selectedCmsUserId.value) {
    grantedIds.value = []
    return
  }
  loading.value = true
  try {
    const res = await guildApi.getGuildVisibilityByUserList(selectedCmsUserId.value)
    grantedIds.value = (res.list || []).map((item) => item.guildId)
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.guildVisibilityList.fetchListFailed'))
  } finally {
    loading.value = false
  }
}

const reloadForUser = async () => {
  if (!selectedCmsUserId.value) {
    ElMessage.warning(t('pages.guildVisibilityList.cmsUserRequired'))
    return
  }
  await fetchGrantedForUser()
}

const handleUserChange = async () => {
  await fetchGrantedForUser()
}

const searchCmsUsers = async (keyword: string) => {
  cmsUserLoading.value = true
  try {
    const res = await cmsUserApi.getCMSUserList({
      pageIndex: 1,
      pageSize: 50,
      name: keyword || undefined,
    })
    cmsUserOptions.value = res.data || []
  } catch (e) {
    console.error(e)
    ElMessage.error(t('pages.guildVisibilityList.fetchCmsUserFailed'))
  } finally {
    cmsUserLoading.value = false
  }
}

const handleTransferChange = async (
    _value: string[],
    direction: 'left' | 'right',
    movedKeys: string[],
) => {
  if (syncing.value || !movedKeys.length || !selectedCmsUserId.value) {
    return
  }

  if (direction === 'right') {
    if (!can('grant')) {
      grantedIds.value = grantedIds.value.filter((id) => !movedKeys.includes(id))
      ElMessage.warning(t('pages.guildVisibilityList.noGrantPermission'))
      return
    }
    syncing.value = true
    try {
      const res = await guildApi.batchGrantGuildVisibility({
        cmsUserId: selectedCmsUserId.value,
        guildIds: movedKeys,
      })
      ElMessage.success(`${t('pages.guildVisibilityList.grantSuccess')} (${res.grantedCount ?? movedKeys.length})`)
    } catch (e) {
      console.error(e)
      grantedIds.value = grantedIds.value.filter((id) => !movedKeys.includes(id))
      ElMessage.error(t('pages.guildVisibilityList.grantFailed'))
    } finally {
      syncing.value = false
    }
    return
  }

  if (!can('revoke')) {
    grantedIds.value = [...grantedIds.value, ...movedKeys]
    ElMessage.warning(t('pages.guildVisibilityList.noRevokePermission'))
    return
  }

  try {
    await ElMessageBox.confirm(
        t('pages.guildVisibilityList.revokeConfirm'),
        t('common.confirm'),
        {type: 'warning'},
    )
  } catch {
    grantedIds.value = [...grantedIds.value, ...movedKeys]
    return
  }

  syncing.value = true
  try {
    const res = await guildApi.batchRevokeGuildVisibility({
      cmsUserId: selectedCmsUserId.value,
      guildIds: movedKeys,
    })
    ElMessage.success(`${t('pages.guildVisibilityList.revokeSuccess')} (${res.revokedCount ?? movedKeys.length})`)
  } catch (e) {
    console.error(e)
    grantedIds.value = [...grantedIds.value, ...movedKeys]
    ElMessage.error(t('pages.guildVisibilityList.revokeFailed'))
  } finally {
    syncing.value = false
  }
}

onMounted(async () => {
  await fetchGuildOptions()
})
</script>

<style scoped>
.tip-alert {
  margin-bottom: 12px;
}

.search-form {
  margin-bottom: 12px;
}

.transfer-wrap {
  width: 100%;
  min-height: calc(100vh - 280px);
}

.transfer-wrap :deep(.el-transfer) {
  display: flex;
  align-items: stretch;
  width: 100%;
}

.transfer-wrap :deep(.el-transfer-panel) {
  flex: 1;
  width: auto;
  max-width: none;
  height: calc(100vh - 280px);
  min-height: 480px;
}

.transfer-wrap :deep(.el-transfer__buttons) {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  width: auto;
}

.transfer-wrap :deep(.el-transfer__button) {
  margin: 8px 0 !important;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 88px;
}

.transfer-wrap :deep(.el-transfer-panel__body) {
  height: calc(100% - 55px);
}

.transfer-wrap :deep(.el-transfer-panel__list) {
  height: calc(100% - 50px);
}

.transfer-wrap :deep(.el-transfer-panel__item) {
  margin-right: 0;
}
</style>
