<script lang="ts">
	import { onMount } from 'svelte';
	import type { RuntimeConfig } from '$lib/types';
	import { store } from '$lib/store';
	import Icon from '@iconify/svelte';

	let functions: RuntimeConfig[] = [];
	let selected: RuntimeConfig | null = null;

	const opts = {
		headers: {
			Authorization: $store || ''
		}
	};

	const loadFunctions = () => {
		fetch('http://localhost:8080/function/list', opts)
			.then((res) => res.json())
			.then((res) => {
				functions = res.functions || [];
			})
			.catch((err) => console.error(err));
	};

	const selectFunction = (name: string) => {
		fetch(`http://localhost:8080/function/get/${name}`, opts)
			.then((res) => res.json())
			.then((res) => {
				selected = res;
			})
			.catch((err) => console.error(err));
	};

	const deleteFunction = (name: string) => {
		fetch(`http://localhost:8080/function/delete/${name}`, {
			method: 'DELETE',
			headers: { Authorization: $store || '' }
		})
			.then((res) => res.json())
			.then(() => {
				selected = null;
				loadFunctions();
			})
			.catch((err) => console.error(err));
	};

	onMount(loadFunctions);
</script>

<div class="grid grid-cols-5 pl-2">
	<div class="flex flex-col justify-center">Deployed Functions</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each functions as fn}
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<!-- svelte-ignore a11y-no-static-element-interactions -->
			<div
				class="card variant-filled m-1 p-2 cursor-pointer"
				class:variant-ghost={selected?.name === fn.name}
				on:click={() => selectFunction(fn.name)}
			>
				<p class="h5"><Icon icon="mdi:lambda" class="inline-block" /> {fn.name}</p>
				<p><Icon icon="mdi:docker" class="inline-block" /> {fn.tag}</p>
				<p><Icon icon="mdi:cog" class="inline-block" /> {fn.runtime}</p>
				<p><Icon icon="mdi:server-network" class="inline-block" /> :{fn.port}</p>
			</div>
		{:else}
			<p class="p-4">No functions deployed</p>
		{/each}
	</div>

	{#if selected}
		<div class="flex flex-col justify-center">Details</div>
		<div class="col-span-4 col-start-2 p-4 border-b">
			<div class="card variant-ghost p-4">
				<p class="h4"><Icon icon="mdi:lambda" class="inline-block" /> {selected.name}</p>
				<div class="grid grid-cols-2 gap-2 pt-2">
					<p><Icon icon="mdi:docker" class="inline-block" /> Tag: {selected.tag}</p>
					<p><Icon icon="mdi:cog" class="inline-block" /> Runtime: {selected.runtime}</p>
					<p><Icon icon="mdi:server-network" class="inline-block" /> Port: {selected.port}</p>
					{#if selected.volume}
						<p><Icon icon="mdi:folder" class="inline-block" /> Volume: {selected.volume}</p>
					{/if}
				</div>
				<p class="pt-2 opacity-50 text-sm">ID: {selected.id}</p>
				<div class="pt-4">
					<button
						class="btn variant-filled-error"
						on:click={() => selected && deleteFunction(selected.name)}
					>
						<Icon icon="mdi:trash-can" class="inline-block" />
						<span>Delete</span>
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>
