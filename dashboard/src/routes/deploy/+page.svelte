<script lang="ts">
	import { onMount } from 'svelte';
	import { store } from '$lib/store';
	import Icon from '@iconify/svelte';

	let runtimes: string[] = [];

	let name = '';
	let tag = '';
	let runtime = '';
	let port = '';
	let volume = '';
	let files: FileList;
	let message = '';
	let loading = false;

	const opts = {
		headers: {
			Authorization: $store || ''
		}
	};

	onMount(() => {
		fetch('http://localhost:8080/runtime/list', opts)
			.then((res) => res.json())
			.then((res) => {
				runtimes = res;
			})
			.catch((err) => console.error(err));
	});

	const deploy = () => {
		if (!name || !tag || !runtime || !port || !files?.length) {
			message = 'Please fill in all required fields';
			return;
		}

		loading = true;
		message = '';

		const formData = new FormData();
		formData.append('name', name);
		formData.append('tag', tag);
		formData.append('runtime', runtime);
		formData.append('port', port);
		formData.append('volume', volume);
		formData.append('file', files[0]);

		fetch('http://localhost:8080/function/upload', {
			method: 'POST',
			headers: { Authorization: $store || '' },
			body: formData
		})
			.then((res) => res.json())
			.then((res) => {
				message = res.message || JSON.stringify(res);
				loading = false;
			})
			.catch((err) => {
				message = `Error: ${err}`;
				loading = false;
			});
	};
</script>

<div class="grid grid-cols-5 pl-2">
	<div class="flex flex-col justify-center">Deploy Function</div>
	<div class="col-span-4 col-start-2 p-4 border-b">
		<p class="pb-4 opacity-75">
			Upload a function to deploy. If a function with the same name already exists, it will be
			replaced.
		</p>

		<div class="flex flex-col gap-4">
			<label class="label">
				<span>Name</span>
				<input class="input p-2" type="text" placeholder="my-function" bind:value={name} />
			</label>

			<label class="label">
				<span>Tag</span>
				<input
					class="input p-2"
					type="text"
					placeholder="my-function:latest"
					bind:value={tag}
				/>
			</label>

			<label class="label">
				<span>Runtime</span>
				<select class="select p-2" bind:value={runtime}>
					<option value="" disabled selected>Select a runtime</option>
					{#each runtimes as rt}
						<option value={rt}>{rt}</option>
					{/each}
				</select>
			</label>

			<label class="label">
				<span>Port</span>
				<input class="input p-2" type="text" placeholder="9002" bind:value={port} />
			</label>

			<label class="label">
				<span>Volume (optional)</span>
				<input
					class="input p-2"
					type="text"
					placeholder="./data:/data"
					bind:value={volume}
				/>
			</label>

			<label class="label">
				<span>Code (zip)</span>
				<input class="input p-2" type="file" accept=".zip" bind:files />
			</label>

			<button class="btn variant-filled w-fit" on:click={deploy} disabled={loading}>
				{#if loading}
					<Icon icon="mdi:timer-sand" class="inline-block" />
					<span>Deploying...</span>
				{:else}
					<Icon icon="mdi:content-save-plus" class="inline-block" />
					<span>Deploy</span>
				{/if}
			</button>

			{#if message}
				<div class="variant-ghost badge p-2">{message}</div>
			{/if}
		</div>
	</div>
</div>
