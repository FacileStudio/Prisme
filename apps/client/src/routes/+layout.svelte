<script lang="ts">
	import '../app.css';
	import { Toaster } from '@facile/muse';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href="/favicon.ico" type="image/x-icon" />
</svelte:head>

<Toaster />
{@render children()}
