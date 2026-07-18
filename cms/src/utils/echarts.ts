import {BarChart, LineChart} from 'echarts/charts'
import {
    GridComponent,
    LegendComponent,
    MarkPointComponent,
    TitleComponent,
    TooltipComponent,
} from 'echarts/components'
import type {ComposeOption} from 'echarts/core'
import * as echarts from 'echarts/core'
import type {BarSeriesOption, LineSeriesOption} from 'echarts/charts'
import type {
    GridComponentOption,
    LegendComponentOption,
    MarkPointComponentOption,
    TitleComponentOption,
    TooltipComponentOption,
} from 'echarts/components'
import {CanvasRenderer} from 'echarts/renderers'

export type EChartsOption = ComposeOption<
    | BarSeriesOption
    | LineSeriesOption
    | TitleComponentOption
    | TooltipComponentOption
    | GridComponentOption
    | LegendComponentOption
    | MarkPointComponentOption
>

echarts.use([
    TitleComponent,
    TooltipComponent,
    GridComponent,
    LegendComponent,
    MarkPointComponent,
    BarChart,
    LineChart,
    CanvasRenderer,
])

export default echarts
