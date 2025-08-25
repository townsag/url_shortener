import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	// https://vite.dev/config/server-options.html#server-proxy
	server: {
		proxy: {
			"/api": {
				"target": "http://localhost:8000"
			}
		}
	}
});