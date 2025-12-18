<script lang="ts">
	import type { VideoResult } from '$lib/types';
	import * as m from '$lib/paraglide/messages.js';

	interface Props {
		video: VideoResult;
		onclick?: (video: VideoResult) => void;
		loading?: boolean;
	}

	let { video, onclick, loading = false }: Props = $props();

	function handleClick() {
		if (loading) return;
		onclick?.(video);
	}

	function handleImageError(e: Event) {
		const img = e.target as HTMLImageElement;
		img.src =
			'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSIjNjY2IiBzdHJva2Utd2lkdGg9IjIiPjxyZWN0IHg9IjMiIHk9IjMiIHdpZHRoPSIxOCIgaGVpZ2h0PSIxOCIgcng9IjIiPjwvcmVjdD48cGF0aCBkPSJtOSAxMGw2IDRsLTYgNFYxMHoiPjwvcGF0aD48L3N2Zz4=';
	}

	// Get source display text
	const sourceText = $derived(() => {
		if (video.sources.length === 1) {
			return video.sources[0].source_name;
		}
		return m.sources({ count: video.sources.length });
	});
</script>

<button
	type="button"
	onclick={handleClick}
	disabled={loading}
	class="card-hover w-full text-left bg-cyber-surface border border-cyber-border rounded-lg overflow-hidden cursor-pointer disabled:cursor-wait"
>
	<div class="aspect-[2/3] relative overflow-hidden">
		<img
			src={video.vod_pic}
			alt={video.vod_name}
			onerror={handleImageError}
			class="w-full h-full object-cover"
			loading="lazy"
		/>
		{#if loading}
			<div class="absolute inset-0 bg-black/70 flex items-center justify-center">
				<div class="w-8 h-8 border-3 border-cyber-blue border-t-transparent rounded-full animate-spin"></div>
			</div>
		{:else if video.vod_remarks}
			<div
				class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent px-2 py-1"
			>
				<span class="text-xs text-gray-300">{video.vod_remarks}</span>
			</div>
		{/if}
	</div>
	<div class="p-2">
		<h3 class="text-sm font-medium text-white line-clamp-2 break-words leading-tight">
			{video.vod_name}
		</h3>
		<div class="flex items-center justify-between mt-1">
			{#if video.type_name}
				<span class="text-xs text-gray-500">{video.type_name}</span>
			{/if}
			<span class="text-xs text-cyber-blue">{sourceText()}</span>
		</div>
	</div>
</button>
