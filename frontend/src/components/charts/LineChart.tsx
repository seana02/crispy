import { onMount, onCleanup, createEffect, Show } from 'solid-js';
import type { JSX } from 'solid-js';
import Chart from 'chart.js/auto';
import { paletteColor, darkDefaults } from './chartConfig';
import './charts.css';

export interface LineChartDataset {
  label: string;
  data: number[];
  colorIndex?: number;
}

interface LineChartProps {
  labels: string[];
  datasets: LineChartDataset[];
  height?: number;
}

export function LineChart(props: LineChartProps): JSX.Element {
  let canvasRef: HTMLCanvasElement | undefined;
  let chart: Chart | undefined;

  onMount(() => {
    if (!canvasRef) return;
    chart = new Chart(canvasRef, {
      type: 'line',
      data: {
        labels: props.labels,
        datasets: props.datasets.map((ds, i) => {
          const color = paletteColor(ds.colorIndex ?? i);
          return {
            label: ds.label,
            data: ds.data,
            borderColor: color,
            backgroundColor: color + '22',
            pointBackgroundColor: color,
            tension: 0.3,
            fill: false,
          };
        }),
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        animation: false,
        color: darkDefaults.color,
        scales: darkDefaults.scales,
        plugins: darkDefaults.plugins,
      },
    });
  });

  createEffect(() => {
    if (!chart) return;
    chart.data.labels = props.labels;
    chart.data.datasets = props.datasets.map((ds, i) => {
      const color = paletteColor(ds.colorIndex ?? i);
      return {
        label: ds.label,
        data: ds.data,
        borderColor: color,
        backgroundColor: color + '22',
        pointBackgroundColor: color,
        tension: 0.3,
        fill: false,
      };
    });
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
