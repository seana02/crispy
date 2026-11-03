/**
 * Shared chart configuration applying the "Neon Night" palette.
 * Chart.js defaults are overridden here to match the dark theme.
 */

// Chart colour palette matches the spec's chart palette order:
// primary, neon-green, neon-cyan, neon-yellow, neon-magenta, neon-violet, neon-blue, neon-orange
export const CHART_COLORS = [
  '#4d8eff', // --color-primary
  '#39ff85', // --color-neon-green
  '#22d3ee', // --color-neon-cyan
  '#ffe03d', // --color-neon-yellow
  '#f472b6', // --color-neon-magenta
  '#a78bfa', // --color-neon-violet
  '#3b82f6', // --color-neon-blue
  '#ff8c42', // --color-neon-orange
];

export const CHART_GRID_COLOR = 'rgba(51, 51, 102, 0.5)';  // --color-border at 50% opacity
export const CHART_TEXT_COLOR = '#8888aa';                  // --color-text-secondary
export const CHART_BG_COLOR   = '#22223a';                  // --color-bg-surface

/** Returns a color from the palette, cycling if index exceeds the palette length. */
export function paletteColor(index: number): string {
  return CHART_COLORS[index % CHART_COLORS.length];
}

/** Shared Chart.js plugin/scale defaults for the dark theme. */
export const darkDefaults = {
  color: CHART_TEXT_COLOR,
  borderColor: CHART_GRID_COLOR,
  scales: {
    x: {
      ticks: { color: CHART_TEXT_COLOR },
      grid:  { color: CHART_GRID_COLOR },
    },
    y: {
      ticks: { color: CHART_TEXT_COLOR },
      grid:  { color: CHART_GRID_COLOR },
    },
  },
  plugins: {
    legend: {
      labels: { color: CHART_TEXT_COLOR },
    },
    tooltip: {
      backgroundColor: '#2a2a4a', // --color-bg-elevated
      titleColor: '#e0e0f0',       // --color-text-primary
      bodyColor:  '#8888aa',       // --color-text-secondary
      borderColor: '#333366',      // --color-border
      borderWidth: 1,
    },
  },
};
