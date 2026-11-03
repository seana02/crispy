import { onMount, onCleanup, createEffect, Show } from 'solid-js';
import type { JSX } from 'solid-js';
import Chart from 'chart.js/auto';
import { paletteColor, darkDefaults } from './chartConfig';
import './charts.css';

export interface BarChartDataset {
  label: string;
  data: number[];
  colorIndex?: number;
}

interface BarChartProps {
  labels: string[];
  datasets: BarChartDataset[];
  height?: number;
  horizontal?: boolean;
}

export function BarChart(props: BarChartProps): JSX.Element {
  let canvasRef: HTMLCanvasElement | undefined;
  let chart: Chart | undefined;

  const buildDatasets = () =>
    props.datasets.map((ds, i) => {
      const color = paletteColor(ds.colorIndex ?? i);
      return {
        label: ds.label,
        data: ds.data,
        backgroundColor: color + 'cc',
        borderColor: color,
        borderWidth: 1,
      };
    });

  onMount(() => {
    if (!canvasRef) return;
    chart = new Chart(canvasRef, {
      type: 'bar',
      data: {
        labels: props.labels,
        datasets: buildDatasets(),
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        animation: false,
        indexAxis: props.horizontal ? 'y' : 'x',
        color: darkDefaults.color,
        scales: darkDefaults.scales,
        plugins: darkDefaults.plugins,
      },
    });
  });

  createEffect(() => {
    if (!chart) return;
    chart.data.labels = props.labels;
    chart.data.datasets = buildDatasets();
    chart.update();
  });

  onCleanup(() => chart?.destroy());

  return (
    <Show
      when={props.labels.length > 0 && props.datasets.length > 0}
      fallback={<div class="chart-empty">No data to display</div>}
    >
      <div class="chart-container" style={props.height ? { height: `${props.height}px` } : undefined}>
        <canvas ref={canvasRef} class="chart-canvas" />
      </div>
    </Show>
  );
}
