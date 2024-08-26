<script lang="ts">
	import Icon from '@iconify/svelte';
	import { store } from '$lib/store';

	$store;

	let loading: boolean = false;
	const showSpinner = () => {
        console.log(".")
		loading = true;
		setTimeout(() => (loading = false), 250);
	};
</script>

<div class="w-full pt-12">
	<div class="flex justify-center">
		<div>Set your generated API key here to be used for all the calls</div>
	</div>
	<div class="grid w-full grid-cols-7">
		<div class="col-span-3 col-start-3 flex p-1">
			<div class="input-group input-group-divider grid-cols-[auto_1fr_auto]">
				<div class="input-group-shim">
					{#if loading}
						<Icon icon="mdi:timer-sand" />
					{:else}
						<Icon icon="mdi:check-bold" />
					{/if}
				</div>
				<input
					type="text"
					placeholder="api key"
					bind:value={$store}
					on:keyup={showSpinner}
					class="input m-4 p-4"
				/>
			</div>
		</div>
	</div>
</div>

<!-- <style>
	div {
		border: 1px solid red;
	}
</style> -->
