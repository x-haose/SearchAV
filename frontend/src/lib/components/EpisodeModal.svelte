<script lang="ts">
	import type { VideoDetail, VideoResult, SourceInfo } from '$lib/types';

	interface Props {
		open: boolean;
		video: VideoResult | null;
		source: SourceInfo | null;
		detail: VideoDetail | null;
		loading?: boolean;
		error?: string | null;
		onclose?: () => void;
		onplay?: (url: string, index: number) => void;
		onsourcechange?: (source: SourceInfo) => void;
	}

	let {
		open,
		video,
		source,
		detail,
		loading = false,
		error = null,
		onclose,
		onplay,
		onsourcechange
	}: Props = $props();

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			onclose?.();
		}
	}

	function handleEpisodeClick(url: string, index: number) {
		onplay?.(url, index);
	}

	function handleSourceChange(src: SourceInfo) {
		onsourcechange?.(src);
	}
</script>

{#if open}
	<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
	<div
		class="fixed inset-0 bg-black/90 flex items-center justify-center z-50 p-4"
		onclick={handleBackdropClick}
		onkeydown={(e) => e.key === 'Escape' && onclose?.()}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div
			class="bg-cyber-surface border border-cyber-border rounded-lg w-full max-w-4xl max-h-[90vh] flex flex-col"
		>
			<div class="flex justify-between items-center p-4 border-b border-cyber-border">
				<h2 class="text-xl font-bold gradient-text truncate pr-4">
					{video?.vod_name || 'Video Detail'}
				</h2>
				<button
					type="button"
					onclick={onclose}
					class="text-gray-400 hover:text-white text-2xl leading-none"
				>
					&times;
				</button>
			</div>

			<div class="flex-1 overflow-auto p-4">
				{#if video && video.sources.length > 1}
					<div class="mb-4">
						<div class="text-sm text-gray-400 mb-2">Select source</div>
						<div class="flex flex-wrap gap-2">
							{#each video.sources as src}
								<button
									type="button"
									onclick={() => handleSourceChange(src)}
									class="px-3 py-1.5 rounded text-sm transition-colors {source?.source_code === src.source_code
										? 'bg-cyber-blue text-black font-medium'
										: 'bg-cyber-surface-100 hover:bg-cyber-border border border-cyber-border text-white'}"
								>
									{src.source_name}
								</button>
							{/each}
						</div>
					</div>
				{/if}

				{#if loading}
					<div class="flex items-center justify-center py-12">
						<div
							class="w-8 h-8 border-4 border-white border-t-transparent rounded-full animate-spin"
						></div>
						<span class="ml-3 text-gray-400">Loading...</span>
					</div>
				{:else if error}
					<div class="text-center text-red-500 py-12">
						<p>{error}</p>
					</div>
				{:else if detail && detail.episodes.length > 0}
					<div class="mb-4 text-sm text-gray-400">
						{detail.episodes.length} episodes
					</div>
					<div class="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 gap-2">
						{#each detail.episodes as url, index}
							<button
								type="button"
								onclick={() => handleEpisodeClick(url, index)}
								class="px-3 py-2 bg-cyber-surface-100 hover:bg-cyber-border border border-cyber-border rounded-lg transition-colors text-center text-sm"
							>
								{index + 1}
							</button>
						{/each}
					</div>
				{:else if detail}
					<div class="text-center text-gray-500 py-12">
						<p>No playable sources</p>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
