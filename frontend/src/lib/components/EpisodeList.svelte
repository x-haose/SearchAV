<script lang="ts">
	interface Props {
		episodes: string[];
		currentIndex: number;
		onselect?: (index: number) => void;
	}

	let { episodes, currentIndex, onselect }: Props = $props();

	function handleClick(index: number) {
		if (index !== currentIndex) {
			onselect?.(index);
		}
	}
</script>

<div class="episode-list">
	<div class="flex items-center justify-between mb-4">
		<h3 class="text-base font-medium text-white">Episodes</h3>
		<span class="text-sm text-gray-500">{currentIndex + 1} / {episodes.length}</span>
	</div>
	<div class="grid grid-cols-4 sm:grid-cols-5 md:grid-cols-6 gap-2">
		{#each episodes as _, index}
			<button
				type="button"
				onclick={() => handleClick(index)}
				class="episode-btn {index === currentIndex ? 'active' : ''}"
			>
				<span class="episode-num">{index + 1}</span>
				{#if index === currentIndex}
					<span class="playing-indicator"></span>
				{/if}
			</button>
		{/each}
	</div>
</div>

<style>
	.episode-btn {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		height: 40px;
		border-radius: 8px;
		font-size: 14px;
		font-weight: 500;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		color: #9ca3af;
		transition: all 0.2s ease;
		cursor: pointer;
	}

	.episode-btn:hover:not(.active) {
		background: rgba(255, 255, 255, 0.1);
		border-color: rgba(255, 255, 255, 0.2);
		color: #fff;
		transform: translateY(-1px);
	}

	.episode-btn.active {
		background: linear-gradient(135deg, #00ccff 0%, #0099cc 100%);
		border-color: #00ccff;
		color: #000;
		box-shadow: 0 0 12px rgba(0, 204, 255, 0.4);
	}

	.episode-num {
		z-index: 1;
	}

	.playing-indicator {
		position: absolute;
		bottom: 4px;
		left: 50%;
		transform: translateX(-50%);
		width: 16px;
		height: 2px;
		background: #000;
		border-radius: 1px;
	}
</style>
