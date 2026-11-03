import { onMount, onCleanup, createEffect, Show } from 'solid-js';
import type { JSX } from 'solid-js';
import Chart from 'chart.js/auto';
import { CHART_COLORS, CHART_TEXT_COLOR, CHART_BG_COLOR, darkDefaults } from './chartConfig';
import './charts.css';

export interface DonutSlice {
  label: string;
  value: number;
}

interface DonutChartProps {
  slices: DonutSlice[];
  height?: number;
  /** Cut-out percentage. Defaults to 60 (donut). Set to 0 for a full pie. */
  cutout?: number;
}

export function DonutChart(props: DonutChartProps): JSX.Element {
  let canvasRef: HTMLCanvasElement | undefined;
  let chart: Chart | undefined;

  const buildData = () => ({
    labels: props.slices.map(s => s.label),
    datasets: [{
      data: props.slices.map(s => s.value),
      backgroundColor: CHART_COLORS.slice(0, props.slices.length).map(c => c + 'cc'),
      borderColor: CHART_COLORS.slice(0, props.slices.length),
      borderWidth: 1,
    }],
  });

  onMount(() => {
    if (!canvasRef) return;
    chart = new Chart(canvasRef, {
      type: 'doughnut',
      data: buildData(),
      options: {
        responsive: true,
        maintainAspectRatio: true,
        animation: false,
        cutout: `${props.cutout ?? 60}%`,
        plugins: {
          ...darkDefaults.plugins,
          legend: {
            position: 'right',
            labels: { color: CHART_TEXT_COLOR },
          },
        },
      },
    });
  });

  createEffect(() => {
    if (!chart) return;
    chart.data = buildData();
    chart.update();
  });

  onCleanup(() => chart?.destroy());

  return (
    <Show
      when={props.slices.length > 0}
      fallback={<div class="chart-empty">No data to display</div>}
    >
      <div class="chart-container" style={props.height ? { height: `${props.height}px` } : undefined}>
        <canvas ref={canvasRef} class="chart-canvas" />
      </div>
    </Show>
  );
}
