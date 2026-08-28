import {definePageMessages} from './_define'

const zh = {
  keywordPlaceholder: '流水ID / 主播ID / 昵称',
  fetchFailed: '获取每日流水失败',
  dailyLiveIncome: '当天直播收益',
  dailyGiftIncome: '当天礼物收益',
  dailyPaidDanmakuIncome: '当天付费弹幕收益',
  dailyPrivateRoomTicketIncome: '当天私密房门票收益',
  dailyPrivateRoomWatchIncome: '当天私密房观看收益',
  dailyVideoCallIncome: '当天视频通话收益',
  dailyVideoTicketIncome: '当天视频门票收益',
  dailyVideoBillingIncome: '当天视频通话计费收益',
  dailyShortVideoIncome: '当天短视频收益',
  dailyGameIncome: '当天游戏收益',
  unsettledTotalIncome: '未结算总收益',
}

const en = {
  keywordPlaceholder: 'Flow ID / anchor ID / nickname',
  fetchFailed: 'Failed to load daily live flow',
  dailyLiveIncome: 'Daily Live Income',
  dailyGiftIncome: 'Daily Gift Income',
  dailyPaidDanmakuIncome: 'Daily Paid Danmaku Income',
  dailyPrivateRoomTicketIncome: 'Daily Private Room Ticket Income',
  dailyPrivateRoomWatchIncome: 'Daily Private Room Watch Income',
  dailyVideoCallIncome: 'Daily Video Call Income',
  dailyVideoTicketIncome: 'Daily Video Ticket Income',
  dailyVideoBillingIncome: 'Daily Video Call Billing Income',
  dailyShortVideoIncome: 'Daily Short Video Income',
  dailyGameIncome: 'Daily Game Income',
  unsettledTotalIncome: 'Unsettled Total Income',
}

const es = {
  keywordPlaceholder: 'ID flujo / ID anfitrión / apodo',
  fetchFailed: 'Error al cargar flujo diario en vivo',
  dailyLiveIncome: 'Ingresos en vivo del día',
  dailyGiftIncome: 'Ingresos regalos del día',
  dailyPaidDanmakuIncome: 'Ingresos danmaku de pago del día',
  dailyPrivateRoomTicketIncome: 'Ingresos entrada sala privada del día',
  dailyPrivateRoomWatchIncome: 'Ingresos visualización sala privada del día',
  dailyVideoCallIncome: 'Ingresos videollamada del día',
  dailyVideoTicketIncome: 'Ingresos entrada video del día',
  dailyVideoBillingIncome: 'Ingresos facturación videollamada del día',
  dailyShortVideoIncome: 'Ingresos videos cortos del día',
  dailyGameIncome: 'Ingresos de juego del día',
  unsettledTotalIncome: 'Ingresos totales no liquidados',
}

const pt = {
  keywordPlaceholder: 'ID fluxo / ID apresentador / apelido',
  fetchFailed: 'Falha ao carregar fluxo diário ao vivo',
  dailyLiveIncome: 'Receita ao vivo do dia',
  dailyGiftIncome: 'Receita de presentes do dia',
  dailyPaidDanmakuIncome: 'Receita danmaku pago do dia',
  dailyPrivateRoomTicketIncome: 'Receita ingresso sala privada do dia',
  dailyPrivateRoomWatchIncome: 'Receita visualização sala privada do dia',
  dailyVideoCallIncome: 'Receita videochamada do dia',
  dailyVideoTicketIncome: 'Receita ingresso vídeo do dia',
  dailyVideoBillingIncome: 'Receita cobrança videochamada do dia',
  dailyShortVideoIncome: 'Receita vídeos curtos do dia',
  dailyGameIncome: 'Receita de jogo do dia',
  unsettledTotalIncome: 'Receita total não liquidada',
}

const hi = {
  keywordPlaceholder: 'फ्लो ID / एंकर ID / उपनाम',
  fetchFailed: 'दैनिक लाइव फ्लो लोड विफल',
  dailyLiveIncome: 'दिन का लाइव आय',
  dailyGiftIncome: 'दिन का उपहार आय',
  dailyPaidDanmakuIncome: 'दिन का सशुल्क डैनमाकू आय',
  dailyPrivateRoomTicketIncome: 'दिन का निजी कमरा टिकट आय',
  dailyPrivateRoomWatchIncome: 'दिन का निजी कमरा देखने की आय',
  dailyVideoCallIncome: 'दिन का वीडियो कॉल आय',
  dailyVideoTicketIncome: 'दिन का वीडियो टिकट आय',
  dailyVideoBillingIncome: 'दिन का वीडियो कॉल बिलिंग आय',
  dailyShortVideoIncome: 'दिन का शॉर्ट वीडियो आय',
  dailyGameIncome: 'दिन का गेम आय',
  unsettledTotalIncome: 'अनिपटारा कुल आय',
}

const id = {
  keywordPlaceholder: 'ID aliran / ID host / nama',
  fetchFailed: 'Gagal memuat aliran live harian',
  dailyLiveIncome: 'Pendapatan live hari ini',
  dailyGiftIncome: 'Pendapatan hadiah hari ini',
  dailyPaidDanmakuIncome: 'Pendapatan danmaku berbayar hari ini',
  dailyPrivateRoomTicketIncome: 'Pendapatan tiket ruang pribadi hari ini',
  dailyPrivateRoomWatchIncome: 'Pendapatan tontonan ruang pribadi hari ini',
  dailyVideoCallIncome: 'Pendapatan panggilan video hari ini',
  dailyVideoTicketIncome: 'Pendapatan tiket video hari ini',
  dailyVideoBillingIncome: 'Pendapatan tagihan panggilan video hari ini',
  dailyShortVideoIncome: 'Pendapatan video pendek hari ini',
  dailyGameIncome: 'Pendapatan game hari ini',
  unsettledTotalIncome: 'Total pendapatan belum diselesaikan',
}

export const liveDailyEffectiveLiveListMessages = definePageMessages(zh, en, es, pt, hi, id)
