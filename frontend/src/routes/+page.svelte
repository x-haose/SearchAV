<script lang="ts">
	import { goto } from '$app/navigation';
	import * as m from '$lib/paraglide/messages.js';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
	import AdultToggle from '$lib/components/AdultToggle.svelte';
	import Footer from '$lib/components/Footer.svelte';

	let query = $state('');

	function handleSearch(q: string) {
		if (!q.trim()) return;
		goto(`/search?wd=${encodeURIComponent(q.trim())}`);
	}
</script>

<svelte:head>
	<title>{m.app_name()} - {m.app_desc()}</title>
</svelte:head>

<!-- Top right controls (fixed position across all pages) -->
<div class="fixed top-4 right-4 z-20 flex items-center gap-2">
	<AdultToggle />
	<LanguageSwitcher />
</div>

<div class="min-h-screen flex flex-col">
	<div class="flex-1 flex flex-col items-center justify-center px-4">
		<header class="text-center mb-8">
			<div class="flex justify-center items-center mb-4">
				<svg class="w-12 h-12 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
					></path>
				</svg>
				<h1 class="text-6xl font-bold gradient-text">{m.app_name()}</h1>
			</div>
			<p class="text-gray-400 text-lg">{m.app_desc()}</p>
		</header>

		<div class="w-full max-w-2xl">
			<SearchBar bind:value={query} onsubmit={handleSearch} />
		</div>
	</div>

	<Footer />
</div>
