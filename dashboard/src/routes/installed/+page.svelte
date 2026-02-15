<script lang="ts">
	import { onMount } from 'svelte';
	import type { RuntimeConfig, Container } from '$lib/types';
	import { store } from '$lib/store';
	import Icon from '@iconify/svelte';

	let runtimes: string[] = [];
	let allImages: string[] = [];
	let deployedFunctions: RuntimeConfig[] = [];
	let running: Container[] = [];

	$: functionTags = deployedFunctions.map((fn) => fn.tag);
	$: baseImages = allImages.filter((img) => !functionTags.includes(img));
	$: functionImages = allImages.filter((img) => functionTags.includes(img));

	const opts = {
		headers: {
			Authorization: $store
		}
	};

	const reload = () => {
		fetch('http://localhost:8080/runtime/list', opts)
			.then((res) => res.json())
			.then((res) => {
				runtimes = res;
			})
			.catch((err) => console.error(err));

		fetch('http://localhost:8080/function/listinstalled', opts)
			.then((res) => res.json())
			.then((res) => {
				allImages = res;
			})
			.catch((err) => console.error(err));

		fetch('http://localhost:8080/function/list', opts)
			.then((res) => res.json())
			.then((res) => {
				deployedFunctions = res.functions || [];
			})
			.catch((err) => console.error(err));

		fetch('http://localhost:8080/function/listrunning', opts)
			.then((res) => res.text())
			.then((res) => {
				let r: Container[] = JSON.parse(res);
				running = r;
			})
			.catch((err) => console.error(err));
	};

	const deleteRuntime = (name: string) => {
		fetch(`http://localhost:8080/runtime/delete/${name}`, {
			method: 'DELETE',
			headers: { Authorization: $store || '' }
		})
			.then((res) => res.json())
			.then(() => reload())
			.catch((err) => console.error(err));
	};

	const deleteFunctionByTag = (tag: string) => {
		const fn = deployedFunctions.find((f) => f.tag === tag);
		if (!fn) return;
		fetch(`http://localhost:8080/function/delete/${fn.name}`, {
			method: 'DELETE',
			headers: { Authorization: $store || '' }
		})
			.then((res) => res.json())
			.then(() => reload())
			.catch((err) => console.error(err));
	};

	const deleteRunningFunction = (container: Container) => {
		const fn = deployedFunctions.find((f) => f.tag === container.Image);
		if (!fn) return;
		fetch(`http://localhost:8080/function/delete/${fn.name}`, {
			method: 'DELETE',
			headers: { Authorization: $store || '' }
		})
			.then((res) => res.json())
			.then(() => reload())
			.catch((err) => console.error(err));
	};

	onMount(reload);
</script>

<div class="grid grid-cols-5 pl-2">
	<div class="flex flex-col justify-center">Installed Runtimes</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each runtimes as runtime}
			<div class="variant-filled badge m-1">
				{runtime}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<!-- svelte-ignore a11y-no-static-element-interactions -->
				<span class="ml-1 cursor-pointer opacity-50 hover:opacity-100" on:click={() => deleteRuntime(runtime)}>
					<Icon icon="mdi:close" width="14" />
				</span>
			</div>
		{:else}
			<p class="opacity-50">No runtimes installed</p>
		{/each}
	</div>

	<div class="flex flex-col justify-center">Function Images</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each functionImages as fn}
			<div class="variant-filled badge m-1">
				{fn}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<!-- svelte-ignore a11y-no-static-element-interactions -->
				<span class="ml-1 cursor-pointer opacity-50 hover:opacity-100" on:click={() => deleteFunctionByTag(fn)}>
					<Icon icon="mdi:close" width="14" />
				</span>
			</div>
		{:else}
			<p class="opacity-50">No function images</p>
		{/each}
	</div>

	<div class="flex flex-col justify-center">Base Images</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each baseImages as img}
			<div class="variant-ghost badge m-1">{img}</div>
		{:else}
			<p class="opacity-50">No base images</p>
		{/each}
	</div>

	<div class="flex flex-col justify-center">Running Functions</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each running as fn}
			<div class="card variant-filled m-1 p-2">
				<div class="flex justify-between">
					<p class="h5"><Icon icon="mdi:lambda" class="inline-block" />{fn.Names[0]}</p>
					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<!-- svelte-ignore a11y-no-static-element-interactions -->
					<span class="cursor-pointer opacity-50 hover:opacity-100" on:click={() => deleteRunningFunction(fn)}>
						<Icon icon="mdi:close" width="18" />
					</span>
				</div>
				<p><Icon icon="mdi:docker" class="inline-block" /> {fn.Image}</p>
				<p><Icon icon="mdi:terminal" class="inline-block" />{fn.Command}</p>
				<p><Icon icon="mdi:server-network" class="inline-block" />{fn.Ports[0].PublicPort || fn.Ports[1].PublicPort}</p>
				<div>
					<p class="inline-block">
						{#if fn.State == 'running'}
							<Icon icon="mdi:play" />
						{:else}
							<Icon icon="mdi:stop" />
						{/if}
					</p>
					<p class="inline-block">{fn.Status}</p>
				</div>
				<p><Icon icon="mdi:folder" class="inline-block" />{fn.Mounts.map((el) => `${el.Source} -> ${el.Destination}`)}</p>
			</div>
		{:else}
			<p class="opacity-50">No functions running</p>
		{/each}
	</div>
</div>
