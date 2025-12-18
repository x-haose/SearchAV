<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { getDetail, AuthError, type SourceInfo } from '$lib';
	import * as m from '$lib/paraglide/messages.js';
	import VideoPlayer from '$lib/components/VideoPlayer.svelte';
	import EpisodeList from '$lib/components/EpisodeList.svelte';
	import PasswordModal from '$lib/components/PasswordModal.svelte';
	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
	import AdultToggle from '$lib/components/AdultToggle.svelte';

	let sourceCode = $state('');
	let vodId = $state(0);
	let title = $state('');
	let currentIndex = $state(0);
	let episodes = $state<string[]>([]);
	let showEpisodes = $state(true);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let showPasswordModal = $state(false);

	// Available sources from sessionStorage
	let availableSources = $state<SourceInfo[]>([]);
	let showSourceSelector = $state(false);

	// Current playing URL
	const currentUrl = $derived(episodes[currentIndex] || '');

	// Current source info
	const currentSource = $derived(
		availableSources.find((s) => s.source_code === sourceCode && s.vod_id === vodId)
	);

	// Load available sources from sessionStorage
	function loadAvailableSources() {
		if (!browser) return;
		try {
			const stored = sessionStorage.getItem('searchav_current_video');
			if (stored) {
				const data = JSON.parse(stored);
				availableSources = data.sources || [];
			}
		} catch {
			availableSources = [];
		}
	}

	// Parse URL params and fetch detail
	$effect(() => {
		if (browser) {
			const params = $page.url.searchParams;
			const newSource = params.get('source') || '';
			const newId = parseInt(params.get('id') || '0', 10);
			const newIndex = parseInt(params.get('index') || '0', 10);

			if (!newSource || !newId) {
				error = m.missing_params();
				loading = false;
				return;
			}

			// Load available sources on first render
			if (availableSources.length === 0) {
				loadAvailableSources();
			}

			// If source or ID changed, re-fetch detail
			if (newSource !== sourceCode || newId !== vodId) {
				sourceCode = newSource;
				vodId = newId;
				currentIndex = newIndex;
				fetchDetail();
			} else if (newIndex !== currentIndex) {
				// Just switching episode
				currentIndex = newIndex;
			}
		}
	});

	async function fetchDetail() {
		loading = true;
		error = null;

		try {
			const detail = await getDetail(sourceCode, vodId);
			title = detail.vod_name;
			episodes = detail.episodes || [];

			if (episodes.length === 0) {
				error = m.no_playable_episodes();
			}
		} catch (e) {
			if (e instanceof AuthError) {
				showPasswordModal = true;
			} else {
				error = e instanceof Error ? e.message : m.detail_failed();
			}
		} finally {
			loading = false;
		}
	}

	function handlePasswordSubmit() {
		showPasswordModal = false;
		fetchDetail();
	}

	function updateIndex(newIndex: number) {
		if (newIndex < 0 || newIndex >= episodes.length) return;

		const params = new URLSearchParams({
			source: sourceCode,
			id: String(vodId),
			index: String(newIndex)
		});

		goto(`/player?${params.toString()}`, { replaceState: true });
	}

	function handlePrev() {
		if (currentIndex > 0) {
			updateIndex(currentIndex - 1);
		}
	}

	function handleNext() {
		if (currentIndex < episodes.length - 1) {
			updateIndex(currentIndex + 1);
		}
	}

	function handleEpisodeSelect(index: number) {
		updateIndex(index);
	}

	function handleBack() {
		goto('/');
	}

	function toggleEpisodes() {
		showEpisodes = !showEpisodes;
	}

	function toggleSourceSelector() {
		showSourceSelector = !showSourceSelector;
	}

	function handleSourceSelect(source: SourceInfo) {
		showSourceSelector = false;
		// Navigate to new source, reset to episode 0
		goto(`/player?source=${source.source_code}&id=${source.vod_id}&index=0`);
	}

	const hasPrev = $derived(currentIndex > 0);
	const hasNext = $derived(currentIndex < episodes.length - 1);
	const hasMultipleSources = $derived(availableSources.length > 1);
</script>

<svelte:head>
	<title>{title} - {m.app_name()}</title>
</svelte:head>

<PasswordModal open={showPasswordModal} onsubmit={handlePasswordSubmit} />

<!-- Top right controls (fixed position across all pages) -->
<div class="fixed top-4 right-4 z-20 flex items-center gap-2">
	<!-- Source selector -->
	{#if sourceCode}
		<div class="relative">
			<button
				type="button"
				onclick={toggleSourceSelector}
				disabled={!hasMultipleSources}
				class="flex items-center gap-1 px-2 py-1 text-sm text-cyber-blue bg-cyber-surface border border-cyber-border rounded transition-colors hover:bg-cyber-border disabled:cursor-default disabled:hover:bg-cyber-surface"
			>
				<span>{currentSource?.source_name || sourceCode}</span>
				{#if hasMultipleSources}
					<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
					</svg>
				{/if}
			</button>

			{#if showSourceSelector}
				<button
					type="button"
					class="fixed inset-0 z-10"
					onclick={() => (showSourceSelector = false)}
					aria-label="Close"
				></button>
				<div class="absolute left-0 top-full mt-1 z-20 bg-cyber-surface border border-cyber-border rounded shadow-xl py-1 w-full">
					{#each availableSources as source}
						<button
							type="button"
							onclick={() => handleSourceSelect(source)}
							class="w-full px-2 py-1 text-left text-sm hover:bg-cyber-border transition-colors flex items-center justify-between"
							class:text-cyber-blue={source.source_code === sourceCode && source.vod_id === vodId}
							class:text-gray-300={source.source_code !== sourceCode || source.vod_id !== vodId}
						>
							<span class="truncate">{source.source_name}</span>
							{#if source.source_code === sourceCode && source.vod_id === vodId}
								<svg class="w-3 h-3 shrink-0 ml-1" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
								</svg>
							{/if}
						</button>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
	<AdultToggle />
	<LanguageSwitcher />
</div>

<div class="min-h-screen bg-cyber-bg flex flex-col">
	<!-- Header -->
	<header class="flex items-center justify-between px-4 py-3 pr-28 border-b border-cyber-border">
		<div class="flex items-center gap-4">
			<button
				type="button"
				onclick={handleBack}
				class="flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					></path>
				</svg>
				<span>{m.back()}</span>
			</button>
			<div class="h-6 w-px bg-cyber-border"></div>
			<h1 class="text-lg font-medium text-white truncate max-w-md">
				{title || m.loading()}
				{#if !loading && episodes.length > 1}
					<span class="text-gray-500 text-sm ml-2">{m.episode({ n: currentIndex + 1 })}</span>
				{/if}
			</h1>
		</div>

		<div class="flex items-center gap-2">
			{#if !loading && episodes.length > 1}
				<button
					type="button"
					onclick={toggleEpisodes}
					class="flex items-center gap-1 px-3 py-1.5 text-sm text-gray-400 hover:text-white bg-cyber-surface border border-cyber-border rounded transition-colors"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 6h16M4 12h16M4 18h7"
						></path>
					</svg>
					<span>{m.episodes()}</span>
				</button>
			{/if}
		</div>
	</header>

	<!-- Main content -->
	<div class="flex-1 flex flex-col lg:flex-row">
		<!-- Video player -->
		<div class="flex-1 relative bg-black">
			{#if loading}
				<div class="absolute inset-0 flex items-center justify-center">
					<div class="flex flex-col items-center gap-3">
						<div class="w-10 h-10 border-4 border-cyber-blue border-t-transparent rounded-full animate-spin"></div>
						<span class="text-gray-400 text-sm">{m.loading()}</span>
					</div>
				</div>
			{:else if error}
				<div class="absolute inset-0 flex items-center justify-center">
					<div class="text-center text-red-500">
						<svg class="w-16 h-16 mx-auto mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							></path>
						</svg>
						<p>{error}</p>
						<button
							type="button"
							onclick={handleBack}
							class="mt-4 px-4 py-2 bg-cyber-surface border border-cyber-border rounded hover:bg-cyber-border transition-colors"
						>
							{m.back_to_search()}
						</button>
					</div>
				</div>
			{:else if currentUrl}
				<div class="aspect-video lg:aspect-auto lg:absolute lg:inset-0">
					<VideoPlayer
						url={currentUrl}
						{title}
						{hasPrev}
						{hasNext}
						onprev={handlePrev}
						onnext={handleNext}
					/>
				</div>
			{/if}
		</div>

		<!-- Episode list sidebar (desktop) -->
		{#if !loading && episodes.length > 1 && showEpisodes}
			<aside class="w-full lg:w-80 border-t lg:border-t-0 lg:border-l border-cyber-border bg-cyber-bg overflow-auto">
				<div class="p-4">
					<EpisodeList {episodes} {currentIndex} onselect={handleEpisodeSelect} />
				</div>
			</aside>
		{/if}
	</div>

	<!-- Navigation controls (mobile) -->
	{#if !loading && episodes.length > 1}
		<div class="lg:hidden flex items-center justify-center gap-4 p-4 border-t border-cyber-border">
			<button
				type="button"
				onclick={handlePrev}
				disabled={!hasPrev}
				class="flex items-center gap-2 px-4 py-2 bg-cyber-surface border border-cyber-border rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed hover:bg-cyber-border"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"
					></path>
				</svg>
				<span>{m.prev()}</span>
			</button>
			<span class="text-gray-500">
				{currentIndex + 1} / {episodes.length}
			</span>
			<button
				type="button"
				onclick={handleNext}
				disabled={!hasNext}
				class="flex items-center gap-2 px-4 py-2 bg-cyber-surface border border-cyber-border rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed hover:bg-cyber-border"
			>
				<span>{m.next()}</span>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"
					></path>
				</svg>
			</button>
		</div>
	{/if}
</div>

<style>
	/* Ensure video fills container on desktop */
	@media (min-width: 1024px) {
		.aspect-video {
			aspect-ratio: auto;
		}
	}
</style>
