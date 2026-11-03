import { icons } from '@facile/muse';
import type { AgentRecord, ChannelRecord, SpaceRecord } from './types';

export type NavPage = {
	label: string;
	href: string;
	active?: boolean;
	icon?: string;
	group?: string;
};

export function buildNavPages(
	activeSpace: SpaceRecord | null,
	dmRooms: Array<{ roomId: string; name: string }>,
	agents: AgentRecord[],
	activeRoomId: string | null
): NavPage[] {
	const list: NavPage[] = [];

	if (activeSpace) {
		for (const ch of activeSpace.channels) {
			list.push({
				label: '#' + ch.name,
				href: '?room=' + ch.room_id,
				active: activeRoomId === ch.room_id,
				icon: 'solar:hashtag-linear',
				group: 'Channels'
			});
		}
	}
	list.push({ label: 'Add Channel', href: '#create-channel', icon: icons.plus, group: 'Channels' });

	for (const dm of dmRooms) {
		list.push({
			label: dm.name,
			href: '?room=' + dm.roomId,
			active: activeRoomId === dm.roomId,
			icon: 'solar:user-linear',
			group: 'Direct Messages'
		});
	}
	list.push({ label: 'New Direct Message', href: '#new-dm', icon: icons.plus, group: 'Direct Messages' });

	for (const ag of agents) {
		list.push({
			label: ag.name,
			href: '?agent=' + ag.user_id,
			icon: 'solar:bot-linear',
			group: 'Agents'
		});
	}
	list.push({ label: 'Add Agent', href: '#add-agent', icon: icons.plus, group: 'Agents' });

	return list;
}

export function findChannelRoom(
	spaces: SpaceRecord[],
	roomId: string
): { spaceId: number; channel: ChannelRecord } | null {
	for (const sp of spaces) {
		const channel = sp.channels.find((c) => c.room_id === roomId);
		if (channel) return { spaceId: sp.id, channel };
	}
	return null;
}
