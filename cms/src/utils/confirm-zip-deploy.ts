import {ElMessageBox} from 'element-plus'
import type {ComposerTranslation} from 'vue-i18n'

const I18N_PREFIX = 'pages.deployZipConfirm'

export type ConfirmZipDeployInput = {
  t: ComposerTranslation
  siteLabel: string
  file: File
  domain?: string
  deployPath?: string
  urlPrefix?: string
  uploadMode?: 'chunked' | 'direct'
}

export function formatDeployFileSize(fileSize: number): string {
  if (fileSize >= 1024 * 1024 * 1024) {
    return `${(fileSize / (1024 * 1024 * 1024)).toFixed(2)} GB`
  }
  if (fileSize >= 1024 * 1024) {
    return `${(fileSize / (1024 * 1024)).toFixed(2)} MB`
  }
  return `${(fileSize / 1024).toFixed(2)} KB`
}

/** 上传 ZIP 部署前二次确认；取消返回 false */
export async function confirmZipDeploy(input: ConfirmZipDeployInput): Promise<boolean> {
  const {t, siteLabel, file, domain, deployPath, urlPrefix, uploadMode = 'direct'} = input
  const lines = [
    t(`${I18N_PREFIX}.intro`),
    t(`${I18N_PREFIX}.site`, {site: siteLabel}),
    t(`${I18N_PREFIX}.file`, {file: file.name}),
    t(`${I18N_PREFIX}.size`, {size: formatDeployFileSize(file.size)}),
  ]
  const domainText = domain?.trim()
  if (domainText) {
    lines.push(t(`${I18N_PREFIX}.domain`, {domain: domainText}))
  }
  const pathText = deployPath?.trim()
  if (pathText) {
    lines.push(t(`${I18N_PREFIX}.deployPath`, {path: pathText}))
  }
  const prefixText = urlPrefix?.trim()
  if (prefixText) {
    lines.push(t(`${I18N_PREFIX}.urlPrefix`, {prefix: prefixText}))
  }
  lines.push(t(uploadMode === 'chunked' ? `${I18N_PREFIX}.chunked` : `${I18N_PREFIX}.direct`))

  try {
    await ElMessageBox.confirm(lines.join('\n'), t(`${I18N_PREFIX}.title`), {
      type: 'warning',
      confirmButtonText: t(`${I18N_PREFIX}.confirmBtn`),
      cancelButtonText: t(`${I18N_PREFIX}.cancelBtn`),
      closeOnClickModal: false,
      customClass: 'zip-deploy-confirm-dialog',
    })
    return true
  } catch {
    return false
  }
}