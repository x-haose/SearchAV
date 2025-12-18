<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { search, AuthError, type VideoResult } from '$lib';
	import * as m from '$lib/paraglide/messages.js';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import ResultGrid from '$lib/components/ResultGrid.svelte';
	import PasswordModal from '$lib/components/PasswordModal.svelte';
	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
	import AdultToggle from '$lib/components/AdultToggle.svelte';
	import Footer from '$lib/components/Footer.svelte';

	let results = $state<VideoResult[]>([]);
	let searching = $state(false);
	let searchError = $state<string | null>(null);
	let showPasswordModal = $state(false);
	let currentQuery = $state('');

	// Get query from URL
	let query = $derived($page.url.searchParams.get('wd') || '');

	// Perform search when query changes
	$effect(() => {
		if (query) {
			currentQuery = query;
			performSearch(query);
		}
	});

	async function performSearch(q: string) {
		if (!q.trim()) return;

		searching = true;
		searchError = null;
		results = [];

		try {
			results = await search(q);
		} catch (e) {
			if (e instanceof AuthError) {
				showPasswordModal = true;
			} else {
				searchError = e instanceof Error ? e.message : m.search_failed();
			}
		} finally {
			searching = false;
		}
	}

	function handleSearch(q: string) {
		if (!q.trim()) return;
		const trimmed = q.trim();
		// If same query, just re-search without navigation
		if (trimmed === query) {
			performSearch(trimmed);
		} else {
			goto(`/search?wd=${encodeURIComponent(trimmed)}`);
		}
	}

	function handlePasswordSubmit() {
		showPasswordModal = false;
		if (query) {
			performSearch(query);
		}
	}

	function handleSelectVideo(video: VideoResult) {
		const source = video.sources[0];
		if (!source) return;

		// Store video info for source switching
		sessionStorage.setItem(
			'searchav_current_video',
			JSON.stringify({
				vod_name: video.vod_name,
				vod_pic: video.vod_pic,
				sources: video.sources
			})
		);

		goto(`/player?source=${source.source_code}&id=${source.vod_id}`);
	}

	function goHome() {
		goto('/');
	}
</script>

<svelte:head>
	<title>{query ? `${query} - ` : ''}{m.app_name()}</title>
</svelte:head>

<PasswordModal open={showPasswordModal} onsubmit={handlePasswordSubmit} />

<!-- Top right controls (fixed position across all pages) -->
<div class="fixed top-4 right-4 z-20 flex items-center gap-2">
	<AdultToggle />
	<LanguageSwitcher />
</div>

<div class="min-h-screen flex flex-col">
	<!-- Top bar -->
	<div class="sticky top-0 z-10 bg-cyber-bg/95 backdrop-blur border-b border-gray-800">
		<div class="container mx-auto px-4 py-3 flex items-center gap-4 pr-28">
			<!-- Home button -->
			<button
				type="button"
				onclick={goHome}
				class="flex items-center gap-2 text-gray-400 hover:text-white transition-colors shrink-0"
			>
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
					></path>
				</svg>
				<span class="text-lg font-bold gradient-text hidden sm:inline">{m.app_name()}</span>
			</button>

			<!-- Search bar -->
			<div class="flex-1 max-w-2xl">
				<SearchBar bind:value={currentQuery} loading={searching} onsubmit={handleSearch} compact />
			</div>
		</div>
	</div>

	<!-- Results -->
	<div class="flex-1 container mx-auto px-4 py-6">
		{#if searching}
			<div class="flex items-center justify-center py-12">
				<div class="w-8 h-8 border-4 border-white border-t-transparent rounded-full animate-spin"></div>
				<span class="ml-3 text-gray-400">{m.searching()}</span>
			</div>
		{:else if searchError}
			<div class="text-center text-red-500 py-12">
				<p>{searchError}</p>
			</div>
		{:else if results.length > 0}
			<div class="mb-4 text-sm text-gray-400">
				{m.found_results({ count: results.length })}
			</div>
			<ResultGrid {results} onselect={handleSelectVideo} />
		{:else if query}
			<div class="text-center text-gray-500 py-12">
				<svg class="w-16 h-16 mx-auto mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>
				<p>{m.no_results()}</p>
			</div>
		{/if}
	</div>

	<Footer />
</div>
