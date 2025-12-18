<script lang="ts">
	import { getAdultMode, setAdultMode } from '$lib';

	let enabled = $state(getAdultMode());

	function toggle() {
		enabled = !enabled;
		setAdultMode(enabled);
	}
</script>

<button
	type="button"
	onclick={toggle}
	class="w-8 h-8 flex items-center justify-center text-xs rounded transition-all duration-200 {enabled
		? 'text-pink-500 bg-pink-500/10 border border-pink-500'
		: 'opacity-0 hover:opacity-10'}"
>
	{#if enabled}
		<span>18+</span>
	{/if}
</button>
