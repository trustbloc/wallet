import { reactive } from 'vue';
import resolveConfig from 'tailwindcss/resolveConfig';
import tailwindConfig from '../../tailwind.js';

const { theme } = resolveConfig(tailwindConfig);

// Declare screens object and define a property for each screen size from the theme
const screens = {};
Object.keys(theme.screens).map((key) =>
  Object.defineProperty(screens, key, { value: parseFloat(theme.screens[key]) })
);

const breakpoints = reactive({
  xs: true,
  sm: false,
  md: false,
  lg: false,
  xl: false,
  '2xl': false,
});

function setBreakpointTo(breakpoint) {
  Object.keys(breakpoints).map((key) => (breakpoints[key] = false));
  breakpoints[breakpoint] = true;
}

function updateBreakpoints() {
  const width = window.innerWidth;
  if (width >= screens['2xl']) {
    setBreakpointTo('2xl');
    return;
  }
  if (width >= screens.xl) {
    setBreakpointTo('xl');
    return;
  }
  if (width >= screens.lg) {
    setBreakpointTo('lg');
    return;
  }
  if (width >= screens.md) {
    setBreakpointTo('md');
    return;
  }
  if (width >= screens.sm) {
    setBreakpointTo('sm');
    return;
  }
  if (width >= screens.xs) {
    setBreakpointTo('xs');
    return;
  }
  return;
}

export default function useBreakpoints() {
  updateBreakpoints();
  window.addEventListener('resize', () => {
    updateBreakpoints();
  });

  return breakpoints;
}
