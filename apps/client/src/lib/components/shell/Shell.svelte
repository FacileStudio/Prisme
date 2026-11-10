<script lang="ts">
	import type { Snippet } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { MobileNav, PageTransition, SideBar, SpaceSwitcher, Topbar, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { notify } from '$lib/state/notify.svelte';
	import { isApp } from '$lib/native/platform';
	import { roomHref } from '$lib/state/actions';
	import { BROWSE_PATH, MESSAGES_PATH, buildNav, railPages, type NavKind } from '$lib/state/nav';
	import { HOME, MANAGE_SPACES_HREF, ui } from '$lib/state/ui.svelte';
	import { t, rememberLocale } from '$lib/i18n/index.svelte';
	import { theme, type ThemePreference } from '$lib/state/theme.svelte';
	import { call } from '$lib/call/call.svelte';
	import { prisme } from '$lib/api/prisme';
	import Dialogs from '$lib/components/dialogs/Dialogs.svelte';
	import SidePanel from '$lib/components/panels/SidePanel.svelte';
	import CallPanel from '$lib/components/call/CallPanel.svelte';

	let { children }: { children: Snippet } = $props();

	const kindIcons: Record<NavKind, string> = {
		messages: glyphs.chat,
		browse: icons.home,
		channel: glyphs.hashtag
	};

	const joinedSpaces = $derived.by(() => {
		void app.version;
		return (app.client?.getRooms() ?? [])
			.filter((r) => r.isSpaceRoom() && r.getMyMembership() === 'join')
			.map((r) => r.roomId)
			.sort()
			.join(',');
	});

	const nav = $derived.by(() => {
		void app.version;
		return app.client ? buildNav(app.client, app.spaces, app.agents) : null;
	});

	const current = $derived(nav?.spaces.find((s) => s.space.room_id === ui.spaceId) ?? null);
	const path = $derived(page.url.pathname);
	const roomId = $derived(page.params.id ? decodeURIComponent(page.params.id) : null);
	const inRoom = $derived(roomId !== null);
	const onSettings = $derived(path.startsWith('/settings'));

	const pages = $derived(
		nav ? railPages(nav, current, path, roomId).map(({ kind, ...entry }) => ({ ...entry, icon: kindIcons[kind] })) : []
	);

	const spaces = $derived((nav?.spaces ?? []).map((s) => ({ id: s.space.room_id, name: s.name })));

	const user = $derived.by(() => {
		const me = app.session?.user;
		if (!me) return undefined;
		return { name: me.name || me.email.split('@')[0], email: me.email, avatar: me.avatar_url || undefined };
	});

	let invitations = $state<{ count: number; href: string } | null>(null);

	$effect(() => {
		if (!app.session) return;
		prisme
			.invitations()
			.then((answer) => (invitations = answer.href ? { count: answer.count, href: answer.href } : null))
			.catch(() => (invitations = null));
	});

	const userMenu = $derived({
		invitations,
		items: isApp() ? [] : [{ label: t('downloads.menu'), icon: icons.download, href: '/downloads' }],
		subscription: { href: 'https://sso.facile.studio/spaces' },
		helpHref: 'https://facile.studio/docs/prisme',
		onSignOut: () => void app.logout(),
		onLocaleChange: rememberLocale,
		theme: theme.preference,
		onThemeChange: (next: string) => theme.set(next as ThemePreference)
	});

	const mobileItems = $derived([
		{ href: MESSAGES_PATH, label: t('nav.messages'), icon: glyphs.chat, active: path === MESSAGES_PATH },
		{ href: BROWSE_PATH, label: t('nav.spaces'), icon: icons.home, active: path === BROWSE_PATH }
	]);

	function selectSpace(id: string | null) {
		if (id) ui.spaceId = id;
	}

	$effect(() => {
		if (!nav || current) return;
		if (ui.spaceId === HOME && nav.spaces.length > 0) ui.spaceId = nav.spaces[0].space.room_id;
	});

	$effect(() => {
		if (!ui.drawer) return;
		ui.drawer = false;
		const inDm = !!roomId && !!nav?.dms.some((row) => row.roomId === roomId);
		void goto(inDm ? MESSAGES_PATH : BROWSE_PATH);
	});

	$effect(() => {
		void joinedSpaces;
		void app.refreshSpaces();
	});

	$effect(() => {
		const client = app.client;
		if (!client) return;
		void client.setPresence({ presence: 'online' }).catch(() => undefined);
		return notify.watch(client, (id) => void goto(roomHref(id)));
	});
</script>

<div class="flex h-dvh w-full overflow-hidden bg-fc-page text-fc-fg">
	<div class="hidden h-full shrink-0 p-3 md:block">
		<SideBar
			class="h-full"
			icon={glyphs.logo}
			title="Prisme"
			bind:collapsed={ui.railCollapsed}
			{pages}
			{user}
			userHref="/settings"
			userActive={onSettings}
			{userMenu}
			{spaces}
			activeSpaceId={current?.space.room_id ?? null}
			onSpaceSelect={selectSpace}
			personalSpaceLabel={null}
			manageSpacesHref={MANAGE_SPACES_HREF}
			manageSpacesLabel={t('spaces.manage')}
		/>
	</div>

	<div class="flex min-w-0 flex-1 flex-col">
		{#if !inRoom}
			<Topbar class="shrink-0 md:hidden">
				<span class="text-fc-md font-semibold text-fc-fg">Prisme</span>
				{#if spaces.length > 0}
					<div class="min-w-0 max-w-56 flex-1">
						<SpaceSwitcher
							{spaces}
							activeId={current?.space.room_id ?? null}
							onSelect={selectSpace}
							personalLabel={null}
							manageHref={MANAGE_SPACES_HREF}
							manageLabel={t('spaces.manage')}
						/>
					</div>
				{/if}
			</Topbar>
		{/if}

		<main class="flex min-h-0 min-w-0 flex-1 flex-col {inRoom ? 'overflow-hidden' : 'overflow-auto overscroll-contain'}">
			{#if call.active}
				<CallPanel />
			{/if}
			<PageTransition key={path} class="flex flex-col {inRoom ? 'min-h-0 flex-1' : 'flex-[1_0_auto] pb-28 md:pb-0'}">
				{@render children()}
			</PageTransition>
		</main>
	</div>

	<SidePanel />

	{#if !inRoom}
		<MobileNav items={mobileItems} {user} profileHref="/settings" profileActive={onSettings} {userMenu} />
	{/if}
</div>

<Dialogs />
