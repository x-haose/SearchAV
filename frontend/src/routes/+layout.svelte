<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>SearchAV</title>
	<meta name="description" content="Video Search & Playback" />
</svelte:head>

<div class="min-h-screen bg-cyber-bg">
	{@render children()}
</div>
