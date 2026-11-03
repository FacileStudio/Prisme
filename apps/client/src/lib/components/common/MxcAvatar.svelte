<script lang="ts">
	import { Avatar } from '@facile/muse';
	import { loadMedia } from '$lib/matrix/media';

	type Props = { mxc?: string | null; src?: string; name: string; size?: 'sm' | 'md' | 'lg'; class?: string };
	let { mxc = null, src = '', name, size = 'sm', class: className = '' }: Props = $props();

	let blob = $state('');

	$effect(() => {
		blob = '';
		if (!mxc) return;
		let live = true;
		loadMedia(mxc, 96).then(
			(url) => live && (blob = url),
			() => undefined
		);
		return () => (live = false);
	});
</script>

<Avatar src={blob || src || undefined} {name} {size} class={className} />
