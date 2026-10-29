import type {EChartsOption} from '@/utils/echarts'

/** 至少 2 个点即显示底部滑块（便于浏览 7/30 天等较长范围） */
export const TREND_CHART_ZOOM_MIN_POINTS = 2

const SLIDER_HEIGHT = 22
const SLIDER_BOTTOM = 4
/** 双行日期+时间（不旋转）占用高度 */
const LABEL_BAND_TWO_LINE = 52
const LABEL_BAND_FLAT = 28

export function trendChartHasZoom(pointCount: number): boolean {
  return pointCount > TREND_CHART_ZOOM_MIN_POINTS
}

/** 为 X 轴刻度（及可选底部滑块）留出足够高度 */
export function trendChartGridBottom(pointCount: number, withSlider = true): number {
  const labelBand = pointCount > 12 ? LABEL_BAND_TWO_LINE : LABEL_BAND_FLAT
  if (!withSlider) {
    return labelBand + 20
  }
  if (!trendChartHasZoom(pointCount)) {
    return 36
  }
  return SLIDER_BOTTOM + SLIDER_HEIGHT + labelBand + 10
}

export function buildTrendChartXAxisLabel(pointCount: number) {
  const dense = pointCount > 12
  const veryDense = pointCount > 36
  let interval: number | 'auto' = 'auto'
  if (veryDense) {
    interval = Math.max(1, Math.floor(pointCount / 8))
  } else if (pointCount > 20) {
    interval = Math.max(0, Math.floor(pointCount / 12))
  }
  return {
    rotate: 0,
    align: 'center',
    hideOverlap: true,
    showMinLabel: true,
    showMaxLabel: true,
    margin: dense ? 10 : 8,
    lineHeight: dense ? 16 : 14,
    overflow: 'none',
    fontSize: 12,
    fontWeight: 500,
    color: '#334155',
    interval,
    formatter: (value: string) => formatTrendChartAxisTime(value, dense, veryDense),
  }
}

/**
 * 配合 dataZoom：刻度按当前可视窗口自适应，勿按全量点数算 interval，
 * 否则密采样全量时标签会乱、曲线也像「龟缩」在左侧。
 */
export function buildTrendChartXAxisLabelZoomFriendly() {
  return {
    rotate: 0,
    align: 'center' as const,
    hideOverlap: true,
    showMinLabel: true,
    showMaxLabel: true,
    margin: 10,
    lineHeight: 16,
    overflow: 'none' as const,
    fontSize: 12,
    fontWeight: 500,
    color: '#334155',
    interval: 'auto' as const,
    formatter: (value: string) => formatTrendChartAxisTime(value, true, true),
  }
}

/** 密刻度：日期与时间分两行，避免旋转长串被裁切；极密时省略年份 */
function formatTrendChartAxisTime(value: string, twoLine: boolean, dropYear: boolean): string {
  if (!value) {
    return ''
  }
  const dateTime = value.match(/^(\d{4}-\d{2}-\d{2})\s+(\d{2}:\d{2})/)
  if (dateTime) {
    if (twoLine) {
      const datePart = dropYear ? dateTime[1].slice(5).replace(/^-/, '') : dateTime[1]
      return `${datePart}\n${dateTime[2]}`
    }
    return value
  }
  if (!twoLine) {
    return value
  }
  const timeOnly = value.match(/(\d{2}:\d{2}:\d{2}|\d{2}:\d{2})/)
  if (timeOnly) {
    return timeOnly[1].length > 5 ? timeOnly[1].slice(0, 5) : timeOnly[1]
  }
  return value.length > 12 ? `${value.slice(0, 12)}…` : value
}

export type TrendChartDataZoomOptions = {
  /** 初始窗口为全量（滑块 0–100%） */
  defaultFullRange?: boolean
  /** 初始窗口对齐到序列末尾（时间序列看最近一段） */
  alignEnd?: boolean
  /**
   * 非全量时，初始可视窗口目标点数（点少则铺满全量，点多则压缩到该数量撑满 X 轴）。
   * 10 秒采样建议 120–360（约 20 分钟～1 小时）。
   */
  initialWindowPoints?: number
}

/** 密采样折线默认可视点数：约撑满常见 CMS 图表宽度，避免全量 3 天挤成左侧细条 */
export const TREND_CHART_DENSE_VISIBLE_POINTS = 180

export function buildTrendChartDataZoom(
    pointCount: number,
    options?: TrendChartDataZoomOptions,
): EChartsOption['dataZoom'] {
  if (!trendChartHasZoom(pointCount)) {
    return undefined
  }
  const windowPoints = Math.max(2, options?.initialWindowPoints ?? TREND_CHART_DENSE_VISIBLE_POINTS)
  const windowPct = options?.defaultFullRange
      ? 100
      : pointCount > windowPoints
          ? Math.min(100, (windowPoints / pointCount) * 100)
          : 100
  const alignEnd = options?.alignEnd === true && windowPct < 100
  const initialStart = alignEnd ? Math.max(0, 100 - windowPct) : 0
  const initialEnd = alignEnd ? 100 : windowPct
  return [
    {
      type: 'inside',
      xAxisIndex: 0,
      start: initialStart,
      end: initialEnd,
      zoomOnMouseWheel: true,
      moveOnMouseMove: true,
      filterMode: 'none',
    },
    {
      type: 'slider',
      xAxisIndex: 0,
      height: SLIDER_HEIGHT,
      bottom: SLIDER_BOTTOM,
      start: initialStart,
      end: initialEnd,
      show: true,
      showDetail: false,
      brushSelect: false,
      filterMode: 'none',
      borderColor: '#dcdfe6',
      fillerColor: 'rgba(64, 158, 255, 0.15)',
      handleSize: '110%',
      handleStyle: {color: '#409EFF', borderColor: '#409EFF'},
      moveHandleSize: 6,
      dataBackground: {
        lineStyle: {color: '#c0c4cc', width: 1},
        areaStyle: {color: 'rgba(144, 147, 153, 0.12)'},
      },
      selectedDataBackground: {
        lineStyle: {color: '#409EFF', width: 1},
        areaStyle: {color: 'rgba(64, 158, 255, 0.2)'},
      },
    },
  ]
}
