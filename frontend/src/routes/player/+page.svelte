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
	let showEpisodes = $state(false); // Start closed on mobile
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
		// Close panel on mobile after selection
		if (window.innerWidth < 1024) {
			showEpisodes = false;
		}
	}

	function handleBack() {
		history.back();
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

<div class="h-screen bg-cyber-bg flex flex-col overflow-hidden">
	<!-- Header -->
	<header class="flex items-center gap-2 px-3 py-2 border-b border-cyber-border shrink-0">
		<!-- Back button -->
		<button
			type="button"
			onclick={handleBack}
			class="shrink-0 p-1.5 text-gray-400 hover:text-white transition-colors"
			aria-label={m.back()}
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
			</svg>
		</button>

		<!-- Title (scrollable) -->
		<div class="flex-1 min-w-0 overflow-hidden">
			<h1 class="text-sm font-medium text-white whitespace-nowrap overflow-x-auto scrollbar-hide">
				{title || m.loading()}{#if !loading && episodes.length > 1} · {m.episode({ n: currentIndex + 1 })}{/if}
			</h1>
		</div>

		<!-- Control buttons -->
		<div class="shrink-0 flex items-center gap-1.5">
			<!-- Episodes button -->
			{#if !loading && episodes.length > 1}
				<button
					type="button"
					onclick={toggleEpisodes}
					class="p-1.5 rounded transition-colors {showEpisodes ? 'text-cyber-blue bg-cyber-blue/10' : 'text-gray-400 hover:text-white'}"
					aria-label={m.episodes()}
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"></path>
					</svg>
				</button>
			{/if}

			<!-- Source selector -->
			{#if sourceCode}
				<div class="relative">
					<button
						type="button"
						onclick={toggleSourceSelector}
						disabled={!hasMultipleSources}
						class="flex items-center gap-0.5 px-2 py-1 text-xs text-cyber-blue bg-cyber-surface border border-cyber-border rounded transition-colors hover:bg-cyber-border disabled:opacity-60 disabled:cursor-default"
					>
						<span class="max-w-16 truncate">{currentSource?.source_name || sourceCode}</span>
						{#if hasMultipleSources}
							<svg class="w-3 h-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
						<div class="absolute right-0 top-full mt-1 z-20 bg-cyber-surface border border-cyber-border rounded shadow-xl py-1 min-w-32">
							{#each availableSources as source}
								<button
									type="button"
									onclick={() => handleSourceSelect(source)}
									class="w-full px-3 py-1.5 text-left text-xs hover:bg-cyber-border transition-colors flex items-center justify-between gap-2"
									class:text-cyber-blue={source.source_code === sourceCode && source.vod_id === vodId}
									class:text-gray-300={source.source_code !== sourceCode || source.vod_id !== vodId}
								>
									<span class="truncate">{source.source_name}</span>
									{#if source.source_code === sourceCode && source.vod_id === vodId}
										<svg class="w-3 h-3 shrink-0" fill="currentColor" viewBox="0 0 20 20">
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
	</header>

	<!-- Main content -->
	<div class="flex-1 flex flex-row min-h-0">
		<!-- Video player -->
		<div class="flex-1 relative bg-black min-h-0 min-w-0">
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
				<div class="absolute inset-0">
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

		<!-- Episode list sidebar (desktop only) -->
		{#if !loading && episodes.length > 1 && showEpisodes}
			<aside class="hidden lg:flex w-64 xl:w-72 border-l border-cyber-border bg-cyber-bg/50 backdrop-blur-sm flex-col min-h-0">
				<div class="p-3 flex-1 min-h-0 flex flex-col overflow-hidden">
					<EpisodeList {episodes} {currentIndex} onselect={handleEpisodeSelect} />
				</div>
			</aside>
		{/if}
	</div>

	<!-- Mobile bottom bar -->
	{#if !loading && episodes.length > 1}
		<div class="lg:hidden flex items-center justify-between px-4 py-3 border-t border-cyber-border shrink-0">
			<button
				type="button"
				onclick={handlePrev}
				disabled={!hasPrev}
				class="flex items-center gap-1 px-3 py-1.5 text-sm bg-cyber-surface border border-cyber-border rounded transition-colors disabled:opacity-40 disabled:cursor-not-allowed hover:bg-cyber-border"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
				</svg>
				<span>{m.prev()}</span>
			</button>
			<span class="text-gray-500 text-sm">
				{currentIndex + 1} / {episodes.length}
			</span>
			<button
				type="button"
				onclick={handleNext}
				disabled={!hasNext}
				class="flex items-center gap-1 px-3 py-1.5 text-sm bg-cyber-surface border border-cyber-border rounded transition-colors disabled:opacity-40 disabled:cursor-not-allowed hover:bg-cyber-border"
			>
				<span>{m.next()}</span>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
				</svg>
			</button>
		</div>
	{/if}
</div>

<!-- Mobile episode panel (bottom sheet) -->
{#if !loading && episodes.length > 1 && showEpisodes}
	<div class="lg:hidden fixed inset-0 z-30">
		<!-- Backdrop -->
		<button
			type="button"
			class="absolute inset-0 bg-black/60 backdrop-blur-sm"
			onclick={() => (showEpisodes = false)}
			aria-label="Close"
		></button>
		<!-- Panel -->
		<div class="absolute bottom-0 left-0 right-0 bg-cyber-bg border-t border-cyber-border rounded-t-2xl max-h-[70vh] flex flex-col animate-slide-up">
			<!-- Handle -->
			<div class="flex justify-center py-2 shrink-0">
				<div class="w-10 h-1 bg-gray-600 rounded-full"></div>
			</div>
			<!-- Content -->
			<div class="px-4 pb-4 flex-1 min-h-0 flex flex-col overflow-hidden">
				<EpisodeList {episodes} {currentIndex} onselect={handleEpisodeSelect} />
			</div>
		</div>
	</div>
{/if}

<style>
	@keyframes slide-up {
		from {
			transform: translateY(100%);
		}
		to {
			transform: translateY(0);
		}
	}

	.animate-slide-up {
		animation: slide-up 0.25s ease-out;
	}

	/* Hide scrollbar but keep functionality */
	.scrollbar-hide {
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
	.scrollbar-hide::-webkit-scrollbar {
		display: none;
	}
</style>
