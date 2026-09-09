import {h} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {dataSyncApi} from '@/api/modules/data-sync'
import {i18n} from '@/i18n'

export type ConfirmDataSyncOptions = {
  /** 业务说明文案(同步条数等) */
  detail: string
  /** 弹窗标题,默认「同步数据」 */
  title?: string
}

/**
 * 跨环境同步前确认:先读取数据同步配置中的目标 API,弹窗展示后点确定才继续。
 * 取消或配置不可用时 reject('cancel'),与 ElMessageBox 行为一致。
 */
export async function confirmDataSync(options: ConfirmDataSyncOptions): Promise<void> {
  const t = i18n.global.t
  let target = ''
  try {
    const cfg = await dataSyncApi.getDataSyncCfg()
    target = String(cfg?.targetApiBase || '').trim()
  } catch (error) {
    console.error('fetch data sync cfg failed:', error)
    ElMessage.error(t('common.syncFetchTargetFailed'))
    return Promise.reject('cancel')
  }
  if (!target) {
    ElMessage.error(t('common.syncTargetEmpty'))
    return Promise.reject('cancel')
  }

  await ElMessageBox.confirm(
      h('div', {class: 'data-sync-confirm'}, [
        h('p', {style: 'margin: 0 0 12px; line-height: 1.5; white-space: pre-wrap'}, options.detail),
        h('p', {
          style: 'margin: 0 0 6px; color: var(--el-text-color-secondary); font-size: 13px',
        }, t('common.syncTargetApiLabel')),
        h('p', {
          style: 'margin: 0 0 12px; font-weight: 600; word-break: break-all; color: var(--el-color-warning); line-height: 1.5',
        }, target),
        h('p', {
          style: 'margin: 0; color: var(--el-text-color-regular); line-height: 1.5',
        }, t('common.syncTargetConfirmHint')),
      ]),
      options.title || t('common.syncData'),
      {
        confirmButtonText: t('common.confirmSync'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      },
  )
}
