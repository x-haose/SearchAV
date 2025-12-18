<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';

	interface Props {
		value: string;
		loading?: boolean;
		compact?: boolean;
		onsubmit?: (query: string) => void;
	}

	let { value = $bindable(''), loading = false, compact = false, onsubmit }: Props = $props();

	function handleSubmit(e: Event) {
		e.preventDefault();
		if (value.trim() && onsubmit) {
			onsubmit(value.trim());
		}
	}

	function handleClear() {
		value = '';
	}
</script>

<form onsubmit={handleSubmit} class="w-full {compact ? '' : 'max-w-2xl mx-auto'}">
	<div class="flex items-stretch {compact ? 'h-10' : 'h-14'} shadow-lg rounded-lg overflow-hidden">
		<input
			type="text"
			bind:value
			placeholder={m.search_placeholder()}
			disabled={loading}
			class="flex-1 bg-cyber-surface border-y border-cyber-border text-white {compact ? 'px-4 text-sm' : 'px-6'} focus:outline-none transition-colors min-w-0 placeholder:text-gray-500"
		/>
		{#if value}
			<button
				type="button"
				onclick={handleClear}
				aria-label="Clear search"
				class="flex px-2 bg-cyber-surface border-y border-cyber-border items-center justify-center text-gray-400 hover:text-white"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					></path>
				</svg>
			</button>
		{/if}
		<button
			type="submit"
			disabled={loading || !value.trim()}
			class="{compact ? 'w-16 text-sm' : 'w-24'} flex items-center justify-center bg-white text-black font-medium hover:bg-gray-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
		>
			{#if loading}
				<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
					<circle
						class="opacity-25"
						cx="12"
						cy="12"
						r="10"
						stroke="currentColor"
						stroke-width="4"
					></circle>
					<path
						class="opacity-75"
						fill="currentColor"
						d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
					></path>
				</svg>
			{:else}
				{m.search()}
			{/if}
		</button>
	</div>
</form>
