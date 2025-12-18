<script lang="ts">
	import { savePassword } from '$lib';
	import * as m from '$lib/paraglide/messages.js';

	interface Props {
		open: boolean;
		onsubmit?: () => void;
	}

	let { open, onsubmit }: Props = $props();
	let password = $state('');
	let error = $state('');

	function handleSubmit(e: Event) {
		e.preventDefault();
		if (!password.trim()) {
			error = m.password_required();
			return;
		}
		savePassword(password.trim());
		password = '';
		error = '';
		onsubmit?.();
	}
</script>

{#if open}
	<div class="fixed inset-0 bg-black/90 flex items-center justify-center z-50 p-4">
		<div class="bg-cyber-surface border border-cyber-border rounded-lg w-full max-w-md p-6">
			<h2 class="text-xl font-bold text-white mb-4">{m.auth_required()}</h2>
			<p class="text-gray-400 text-sm mb-6">
				{m.auth_desc()}
			</p>

			<form onsubmit={handleSubmit}>
				<input
					type="password"
					bind:value={password}
					placeholder={m.enter_password()}
					class="w-full px-4 py-3 bg-cyber-bg border border-cyber-border rounded-lg text-white placeholder:text-gray-500 focus:outline-none focus:border-cyber-blue mb-4"
				/>

				{#if error}
					<p class="text-red-500 text-sm mb-4">{error}</p>
				{/if}

				<button
					type="submit"
					class="w-full py-3 bg-cyber-blue text-black font-medium rounded-lg hover:bg-cyan-400 transition-colors"
				>
					{m.submit()}
				</button>
			</form>
		</div>
	</div>
{/if}
