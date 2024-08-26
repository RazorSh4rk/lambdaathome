<script lang="ts">
	import { onMount } from 'svelte';
	import type { RuntimeConfig, Container } from '$lib/types';
	import { store } from '$lib/store';
	import Icon from '@iconify/svelte';

	let runtimes: string[] = [];
	let functions: string[] = [];
	let running: Container[] = [];

	let obj: RuntimeConfig = JSON.parse(`{
			"name": "test_runtime",
			"tag": "go-test",
			"runtime": "golang-8080",
			"port": "9002",
			"volume": "./docs:/docs,./:/codefolder",
			"source": "./buildcache/f0a7cd12-c506-4eb4-820e-fa992aa9b3a0/code/",
			"id": "5c91e57f167e0bd93a69376ad94e296c14e7769c3520ed52a3799478867aca93"
		}`);
	let lambdas: RuntimeConfig[] = [obj];

	const opts = {
		headers: {
			Authorization: $store
		}
	};

	onMount(() => {
		fetch('http://localhost:8080/runtime/list', opts)
			.then((res) => res.json())
			.then((res) => {
				runtimes = res;
			})
			.catch((err) => console.error(err));

		fetch('http://localhost:8080/function/listinstalled', opts)
			.then((res) => res.json())
			.then((res) => {
				functions = res;
			})
			.catch((err) => console.error(err));

		fetch('http://localhost:8080/function/listrunning', opts)
			.then((res) => res.text())
			.then((res) => {
				let r: Container[] = JSON.parse(res);
				running = r;
			})
			.catch((err) => console.error(err));
	});
</script>

<div class="grid grid-cols-5 pl-2">
	<div class="flex flex-col justify-center">Installed Runtimes</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each runtimes as runtime}
			<div class="variant-filled badge m-1">{runtime}</div>
		{/each}
	</div>

	<div class="flex flex-col justify-center">Installed Functions</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each functions as fn}
			<div class="variant-filled badge m-1">{fn}</div>
		{/each}
	</div>

	<div class="flex flex-col justify-center">Running Functions</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each running as fn}
			<div class="card variant-filled m-1 p-2">
				<p class="h5"><Icon icon="mdi:lambda" class="inline-block" />{fn.Names[0]}</p>
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
		{/each}
	</div>

	<div class="flex flex-col justify-center">Installed Runtimes</div>
	<div class="col-span-4 col-start-2 flex flex-wrap justify-center p-4 border-b">
		{#each runtimes as runtime}
			<div class="card-p4 variant-ghost m-1">{runtime}</div>
		{/each}
	</div>
</div>

<!-- <style>
	div,
	p {
		border: 1px solid red;
	}
</style> -->
