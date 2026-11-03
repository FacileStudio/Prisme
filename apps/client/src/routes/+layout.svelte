<script lang="ts">
	import '../app.css';
	import { Toaster } from '@facile/muse';
	import { registerGlyphs } from '$lib/glyphs';
	import { theme } from '$lib/state/theme.svelte';
	import { notify } from '$lib/state/notify.svelte';

	let { children } = $props();

	registerGlyphs();
	theme.restore();
	notify.restore();
</script>

{@render children()}

<Toaster />
