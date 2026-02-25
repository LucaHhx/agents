import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: '{{APP_IDENTIFIER}}',
  appName: '{{PROJECT_TITLE}}',
  webDir: 'dist',
  server: {
    androidScheme: 'https',
  },
};

export default config;
