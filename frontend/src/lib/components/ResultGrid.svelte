<script lang="ts">
	import type { VideoResult } from '$lib/types';
	import * as m from '$lib/paraglide/messages.js';
	import ResultCard from './ResultCard.svelte';

	interface Props {
		results: VideoResult[];
		onselect?: (video: VideoResult) => void;
		loadingVideo?: string | null;
	}

	let { results, onselect, loadingVideo = null }: Props = $props();
</script>

{#if results.length === 0}
	<div class="text-center text-gray-500 py-12">
		<svg class="w-16 h-16 mx-auto mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
		</svg>
		<p>{m.no_results()}</p>
	</div>
{:else}
	<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
		{#each results as video (video.vod_name)}
			<ResultCard {video} onclick={onselect} loading={loadingVideo === video.vod_name} />
		{/each}
	</div>
{/if}
