<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';

	interface Props {
		episodes: string[];
		currentIndex: number;
		onselect?: (index: number) => void;
	}

	let { episodes, currentIndex, onselect }: Props = $props();

	let jumpInput = $state('');
	let gridContainer: HTMLDivElement;

	function handleClick(index: number) {
		if (index !== currentIndex) {
			onselect?.(index);
		}
	}

	function handleJump() {
		const num = parseInt(jumpInput, 10);
		if (num >= 1 && num <= episodes.length && num - 1 !== currentIndex) {
			onselect?.(num - 1);
		}
		jumpInput = '';
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			handleJump();
		}
	}

	// Auto scroll to current episode
	$effect(() => {
		if (gridContainer) {
			const activeBtn = gridContainer.querySelector('.active') as HTMLElement;
			if (activeBtn) {
				activeBtn.scrollIntoView({ block: 'center', behavior: 'smooth' });
			}
		}
	});
</script>

<div class="episode-list flex flex-col h-full">
	<!-- Header -->
	<div class="flex items-center justify-between mb-3 shrink-0">
		<h3 class="text-base font-medium text-white">{m.episodes()}</h3>
		<span class="text-sm text-gray-500">{currentIndex + 1} / {episodes.length}</span>
	</div>

	<!-- Episodes grid (scrollable) -->
	<div
		bind:this={gridContainer}
		class="flex-1 overflow-y-auto overflow-x-hidden pr-1 min-h-0"
	>
		<div class="grid grid-cols-5 gap-1.5">
			{#each episodes as _, index}
				<button
					type="button"
					onclick={() => handleClick(index)}
					class="episode-btn {index === currentIndex ? 'active' : ''}"
				>
					{index + 1}
				</button>
			{/each}
		</div>
	</div>

	<!-- Quick jump (only show if many episodes) -->
	{#if episodes.length > 30}
		<div class="flex items-center gap-2 mt-3 pt-3 border-t border-cyber-border shrink-0">
			<input
				type="number"
				bind:value={jumpInput}
				onkeydown={handleKeydown}
				placeholder="1-{episodes.length}"
				min="1"
				max={episodes.length}
				class="flex-1 h-8 px-2 text-sm bg-cyber-surface border border-cyber-border rounded text-white placeholder-gray-500 focus:outline-none focus:border-cyber-blue [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
			/>
			<button
				type="button"
				onclick={handleJump}
				class="h-8 px-3 text-sm bg-cyber-blue/20 text-cyber-blue border border-cyber-blue/50 rounded hover:bg-cyber-blue/30 transition-colors"
			>
				Go
			</button>
		</div>
	{/if}
</div>

<style>
	.episode-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		aspect-ratio: 1;
		border-radius: 6px;
		font-size: 12px;
		font-weight: 500;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.08);
		color: #6b7280;
		transition: all 0.15s ease;
		cursor: pointer;
	}

	.episode-btn:hover:not(.active) {
		background: rgba(255, 255, 255, 0.08);
		border-color: rgba(255, 255, 255, 0.15);
		color: #fff;
	}

	.episode-btn.active {
		background: linear-gradient(135deg, #00ccff 0%, #0099cc 100%);
		border-color: #00ccff;
		color: #000;
		font-weight: 600;
		box-shadow: 0 0 12px rgba(0, 204, 255, 0.4);
	}

	/* Custom scrollbar */
	.overflow-y-auto::-webkit-scrollbar {
		width: 4px;
	}

	.overflow-y-auto::-webkit-scrollbar-track {
		background: transparent;
	}

	.overflow-y-auto::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.1);
		border-radius: 2px;
	}

	.overflow-y-auto::-webkit-scrollbar-thumb:hover {
		background: rgba(255, 255, 255, 0.2);
	}
</style>
