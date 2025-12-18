<script lang="ts">
	import { getLocale, setLocale, locales, type Locale } from '$lib/paraglide/runtime.js';

	let currentLocale = $state(getLocale());

	function toggle() {
		const idx = locales.indexOf(currentLocale as Locale);
		const nextIdx = (idx + 1) % locales.length;
		const nextLocale = locales[nextIdx];
		setLocale(nextLocale);
		currentLocale = nextLocale;
	}
</script>

<button
	type="button"
	onclick={toggle}
	class="flex items-center gap-1 px-2 py-1 text-sm text-gray-400 hover:text-white bg-cyber-surface border border-cyber-border rounded transition-colors"
	title="Switch language"
>
	<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129"></path>
	</svg>
	<span>{currentLocale === 'en' ? 'EN' : '中'}</span>
</button>
