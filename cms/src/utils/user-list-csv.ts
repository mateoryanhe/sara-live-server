type TranslateFn = (key: string, params?: Record<string, unknown>) => string

export function buildUserListCsvHeaders(t: TranslateFn, ns = 'pages.userList'): string[] {
  return [
    'ID',
    t(`${ns}.inviterId`),
    t('common.createdAt'),
    t(`${ns}.lastLoginTime`),
    t('common.nickname'),
    t(`${ns}.goldBalance`),
    t(`${ns}.diamondBalance`),
    t(`${ns}.vipLevel`),
    t(`${ns}.userType`),
    t(`${ns}.anchor`),
    'IP',
    t(`${ns}.loginCountry`),
    t(`${ns}.registerIp`),
    t(`${ns}.registerCountry`),
    t(`${ns}.deviceType`),
    t(`${ns}.packageName`),
    t(`${ns}.appVersion`),
    t(`${ns}.banStatus`),
    t(`${ns}.cancelStatus`),
    t(`${ns}.rechargeWhitelist`),
    t(`${ns}.banTime`),
    t(`${ns}.openId`),
    t(`${ns}.channel`),
    t('common.remark'),
  ]
}

export function buildUserListUserTypeLabels(t: TranslateFn, ns = 'pages.userList') {
  return {
    '0': t(`${ns}.userTypeNormal`),
    '1': t(`${ns}.userTypeAnchor`),
    '2': t(`${ns}.userTypeBotAnchor`),
    '3': t(`${ns}.userTypeBotViewer`),
    '4': t(`${ns}.userTypeTester`),
    '5': t(`${ns}.userTypeCmsAuthor`),
    '6': t(`${ns}.userTypeCoinMerchant`),
    '7': t(`${ns}.userTypeSeniorAnchor`),
  }
}

export function buildUserListYesNoLabels(t: TranslateFn) {
  return {
    '1': t('common.yes'),
    '0': t('common.no'),
  }
}

export function buildUserListBanLabels(t: TranslateFn, ns = 'pages.userList') {
  return {
    '1': t(`${ns}.banned`),
    '0': t('common.normal'),
  }
}

export function buildUserListCancelLabels(t: TranslateFn, ns = 'pages.userList') {
  return {
    '1': t(`${ns}.canceled`),
    '0': t('common.normal'),
  }
}
