<script lang="ts">
	import { onDestroy } from 'svelte';
	import { untrack } from 'svelte';
	import type Artplayer from 'artplayer';
	import type Hls from 'hls.js';
	import * as m from '$lib/paraglide/messages.js';

	interface Props {
		url: string;
		title?: string;
		onprev?: () => void;
		onnext?: () => void;
		hasPrev?: boolean;
		hasNext?: boolean;
	}

	let { url, title = '', onprev, onnext, hasPrev = false, hasNext = false }: Props = $props();

	let container: HTMLDivElement;
	let art: Artplayer | null = null;
	let hls: Hls | null = null;
	let loading = $state(true);
	let error = $state<string | null>(null);
	let currentUrl = '';  // Track currently loaded URL

	function isHlsUrl(url: string): boolean {
		return url.includes('.m3u8') || url.includes('m3u8');
	}

	/**
	 * Extract base path from URL (for ad detection)
	 * e.g.: https://example.com/video/20251218/abc/001.ts -> https://example.com/video/20251218/abc/
	 */
	function getBasePath(url: string): string {
		try {
			const urlObj = new URL(url);
			const pathParts = urlObj.pathname.split('/');
			pathParts.pop(); // Remove filename
			return urlObj.origin + pathParts.join('/') + '/';
		} catch {
			// Relative path
			const parts = url.split('/');
			parts.pop();
			return parts.join('/') + '/';
		}
	}

	/**
	 * Filter m3u8 content, remove ad segments
	 * Principle: Count ts segment base paths, keep segments from the most common path
	 */
	function filterAdsFromM3u8(m3u8Content: string, baseUrl: string): string {
		const lines = m3u8Content.split('\n');
		const segments: { line: string; tsLine: string; basePath: string }[] = [];
		const pathCount: Record<string, number> = {};

		// First pass: collect all ts segments and their paths
		for (let i = 0; i < lines.length; i++) {
			const line = lines[i].trim();
			if (line && !line.startsWith('#')) {
				// This is a ts file line
				let fullUrl = line;
				if (!line.startsWith('http')) {
					// Relative path, convert to absolute
					try {
						fullUrl = new URL(line, baseUrl).href;
					} catch {
						fullUrl = line;
					}
				}

				const basePath = getBasePath(fullUrl);
				pathCount[basePath] = (pathCount[basePath] || 0) + 1;

				// Find corresponding #EXTINF line
				let extinfLine = '';
				for (let j = i - 1; j >= 0; j--) {
					if (lines[j].trim().startsWith('#EXTINF')) {
						extinfLine = lines[j];
						break;
					}
				}

				segments.push({ line: extinfLine, tsLine: lines[i], basePath });
			}
		}

		// Find the most common path (main content path)
		let mainPath = '';
		let maxCount = 0;
		for (const [path, count] of Object.entries(pathCount)) {
			if (count > maxCount) {
				maxCount = count;
				mainPath = path;
			}
		}

		// If only one path or no clear difference, don't filter
		if (Object.keys(pathCount).length <= 1) {
			return m3u8Content;
		}

		// Second pass: rebuild m3u8, keep only main content path segments
		const resultLines: string[] = [];
		let skipNext = false;
		let adsRemoved = 0;

		for (const line of lines) {
			const trimmed = line.trim();

			if (trimmed.startsWith('#EXTINF')) {
				// Check if next line's ts is an ad
				const nextLineIndex = lines.indexOf(line) + 1;
				if (nextLineIndex < lines.length) {
					const nextLine = lines[nextLineIndex].trim();
					if (nextLine && !nextLine.startsWith('#')) {
						let fullUrl = nextLine;
						if (!nextLine.startsWith('http')) {
							try {
								fullUrl = new URL(nextLine, baseUrl).href;
							} catch {
								fullUrl = nextLine;
							}
						}
						const basePath = getBasePath(fullUrl);
						if (basePath !== mainPath) {
							skipNext = true;
							adsRemoved++;
							continue; // Skip this #EXTINF line
						}
					}
				}
				resultLines.push(line);
			} else if (trimmed && !trimmed.startsWith('#')) {
				// ts file line
				if (skipNext) {
					skipNext = false;
					continue; // Skip ad ts
				}
				resultLines.push(line);
			} else {
				// Other lines (header info, etc.)
				resultLines.push(line);
			}
		}

		if (adsRemoved > 0) {
			console.log(`[AdFilter] Removed ${adsRemoved} ad segments`);
		}

		return resultLines.join('\n');
	}

	function destroyPlayer() {
		if (hls) {
			hls.destroy();
			hls = null;
		}
		if (art) {
			art.destroy();
			art = null;
		}
		currentUrl = '';
	}

	async function initPlayer(videoUrl: string) {
		// Prevent re-initialization of same URL
		if (currentUrl === videoUrl && art) {
			return;
		}

		destroyPlayer();
		loading = true;
		error = null;
		currentUrl = videoUrl;

		if (!container || !videoUrl) {
			loading = false;
			return;
		}

		try {
			// Dynamic imports for code splitting
			const [ArtplayerModule, HlsModule] = await Promise.all([
				import('artplayer'),
				import('hls.js')
			]);

			const Artplayer = ArtplayerModule.default;
			const Hls = HlsModule.default;

			const options: ConstructorParameters<typeof Artplayer>[0] = {
				container,
				url: videoUrl,
				volume: 0.7,
				isLive: false,
				muted: false,
				autoplay: true,
				pip: true,
				autoSize: false,
				autoMini: true,
				screenshot: true,
				setting: true,
				loop: false,
				flip: true,
				playbackRate: true,
				aspectRatio: true,
				fullscreen: true,
				fullscreenWeb: true,
				subtitleOffset: true,
				miniProgressBar: true,
				mutex: true,
				backdrop: true,
				playsInline: true,
				autoPlayback: true,
				airplay: true,
				theme: '#00ccff',
				lang: navigator.language.toLowerCase().includes('zh') ? 'zh-cn' : 'en',
				hotkey: true,
				controls: [
					{
						name: 'prev',
						position: 'left',
						html: '<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M6 6h2v12H6V6zm3.5 6l8.5 6V6l-8.5 6z"/></svg>',
						tooltip: 'Previous',
						click: () => {
							if (hasPrev && onprev) onprev();
						}
					},
					{
						name: 'next',
						position: 'left',
						html: '<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M6 18l8.5-6L6 6v12zm2 0h2V6h-2v12z" transform="rotate(180 12 12)"/></svg>',
						tooltip: 'Next',
						click: () => {
							if (hasNext && onnext) onnext();
						}
					}
				]
			};

			// HLS support with ad filtering
			if (isHlsUrl(videoUrl) && Hls.isSupported()) {
				options.customType = {
					m3u8: (video: HTMLVideoElement, url: string) => {
						// Create custom playlist loader to filter ads
						class AdFilterLoader extends Hls.DefaultConfig.loader {
							load(context: any, config: any, callbacks: any) {
								const onSuccess = callbacks.onSuccess;
								callbacks.onSuccess = (response: any, stats: any, context: any) => {
									// Only process m3u8 playlist
									if (context.type === 'manifest' || context.type === 'level') {
										const originalData = response.data;
										if (typeof originalData === 'string' && originalData.includes('#EXTINF')) {
											response.data = filterAdsFromM3u8(originalData, context.url);
										}
									}
									onSuccess(response, stats, context);
								};
								super.load(context, config, callbacks);
							}
						}

						hls = new Hls({
							maxBufferLength: 30,
							maxMaxBufferLength: 60,
							loader: AdFilterLoader
						});
						hls.loadSource(url);
						hls.attachMedia(video);
					}
				};
				options.type = 'm3u8';
			}

			art = new Artplayer(options);

			// Keyboard shortcuts
			art.on('ready', () => {
				loading = false;
				document.addEventListener('keydown', handleKeydown);
			});

			art.on('error', (err) => {
				console.error('Player error:', err);
				error = m.play_error();
			});
		} catch (err) {
			console.error('Failed to load player:', err);
			error = m.load_error();
			loading = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (!art) return;

		// Ignore if typing in input
		if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
			return;
		}

		switch (e.key) {
			case 'ArrowLeft':
				art.currentTime = Math.max(0, art.currentTime - 5);
				break;
			case 'ArrowRight':
				art.currentTime = Math.min(art.duration, art.currentTime + 5);
				break;
			case 'ArrowUp':
				art.volume = Math.min(1, art.volume + 0.1);
				break;
			case 'ArrowDown':
				art.volume = Math.max(0, art.volume - 0.1);
				break;
			case ' ':
				e.preventDefault();
				art.playing ? art.pause() : art.play();
				break;
			case 'f':
			case 'F':
				art.fullscreen = !art.fullscreen;
				break;
			case 'm':
			case 'M':
				art.muted = !art.muted;
				break;
			case 'p':
			case 'P':
				if (hasPrev && onprev) onprev();
				break;
			case 'n':
			case 'N':
				if (hasNext && onnext) onnext();
				break;
		}
	}

	onDestroy(() => {
		document.removeEventListener('keydown', handleKeydown);
		destroyPlayer();
	});

	// Watch URL changes, initialize player
	$effect(() => {
		const videoUrl = url;
		const playerContainer = container;

		if (videoUrl && playerContainer) {
			// Use untrack to prevent initPlayer internal state changes from triggering re-execution
			untrack(() => {
				initPlayer(videoUrl);
			});
		}
	});
</script>

<div class="video-player-container">
	{#if loading}
		<div class="absolute inset-0 flex items-center justify-center bg-black">
			<div class="flex flex-col items-center gap-3">
				<div class="w-10 h-10 border-4 border-cyber-blue border-t-transparent rounded-full animate-spin"></div>
				<span class="text-gray-400 text-sm">{m.loading()}</span>
			</div>
		</div>
	{/if}
	{#if error}
		<div class="absolute inset-0 flex items-center justify-center bg-black">
			<div class="text-center text-red-500">
				<p>{error}</p>
			</div>
		</div>
	{/if}
	<div bind:this={container} class="w-full h-full"></div>
</div>

<style>
	.video-player-container {
		width: 100%;
		height: 100%;
		background: #000;
		position: relative;
	}

	.video-player-container :global(.art-video-player) {
		width: 100% !important;
		height: 100% !important;
	}
</style>
