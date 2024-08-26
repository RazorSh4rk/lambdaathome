<script lang="ts">
	import '../app.css';
	import { AppRail, AppRailTile, AppRailAnchor } from '@skeletonlabs/skeleton';
	import Icon from '@iconify/svelte';

	let currentTile: number = 0;
</script>

<div class="w-screen flex">
	<div class="h-screen inline-block">
		<AppRail>
			<svelte:fragment slot="lead">
				<AppRailAnchor href="/">
					<Icon icon="material-symbols:cloud" width="48" class="m-auto" />
					<span
						><Icon icon="mdi:lambda" class="inline" width="24" />
						@
						<Icon icon="mdi:house" class="inline" width="24" />
					</span>
				</AppRailAnchor>
				<div class="pb-6"></div>
			</svelte:fragment>

			<AppRailTile bind:group={currentTile} name="Key" value={0} title="Key">
				<svelte:fragment slot="lead">
					<a href="/">
						<Icon icon="mdi:key" width="36" class="m-auto" />
						<span>Key</span>
					</a>
				</svelte:fragment>
			</AppRailTile>

			<AppRailTile bind:group={currentTile} name="Installed" value={1} title="Installed">
				<svelte:fragment slot="lead">
					<a href="/installed">
						<Icon icon="mdi:function" width="36" class="m-auto" />
						<span>Installed</span>
					</a>
				</svelte:fragment>
			</AppRailTile>

			<AppRailTile bind:group={currentTile} name="Deploy" value={2} title="Deploy">
				<svelte:fragment slot="lead">
					<Icon icon="mdi:content-save-plus" width="36" class="m-auto" />
				</svelte:fragment>
				<span>Deploy</span>
			</AppRailTile>

			<AppRailTile bind:group={currentTile} name="Settings" value={3} title="Settings">
				<svelte:fragment slot="lead">
					<Icon icon="mdi:settings" width="36" class="m-auto" />
				</svelte:fragment>
				<span>Settings</span>
			</AppRailTile>

			<AppRailTile bind:group={currentTile} name="logs" value={4} title="Logs">
				<svelte:fragment slot="lead">
					<Icon icon="mdi:blog" width="36" class="m-auto" />
				</svelte:fragment>
				<span>Logs</span>
			</AppRailTile>

			<svelte:fragment slot="trail">
				<AppRailAnchor href="/" target="_blank" title="Account">
					<Icon icon="mdi:github" width="36" class="m-auto" />
				</AppRailAnchor>
			</svelte:fragment>
		</AppRail>
	</div>

	<div class="flex-1 inline-block"><slot></slot></div>
</div>

<!-- <style>
	div{
		border: 1px solid blue;
	}
</style> -->