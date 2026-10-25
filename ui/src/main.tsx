import '@mantine/core/styles.css';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { MantineProvider, createTheme, type MantineColorsTuple } from '@mantine/core';
import { QueryClientProvider } from '@tanstack/react-query';
import { App } from './App';
import { queryClient } from './app/query-client';

const myColor: MantineColorsTuple = [
  '#f3edff',
  '#e0d7fa',
  '#beabf0',
  '#9a7de6',
  '#7c55de',
  '#693cd9',
  '#5f30d8',
  '#4f23c0',
  '#461eac',
  '#3b1898',
];

const theme = createTheme({
  colors: {
    myColor,
  },
  primaryColor: 'myColor',
  defaultRadius: 'sm',
  fontFamily:
    'Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  headings: {
    fontWeight: '750',
  },
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={theme}>
        <App />
      </MantineProvider>
    </QueryClientProvider>
  </StrictMode>,
);
