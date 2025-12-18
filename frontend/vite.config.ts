import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		port: 9595,
		proxy: {
			'/api': {
				target: 'http://localhost:9898',
				changeOrigin: true
			}
		}
	}
});
