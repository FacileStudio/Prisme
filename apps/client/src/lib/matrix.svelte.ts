import type { ChatMessage, MatrixCredentials, RoomMember } from './types';
import {
	parseChatMessage,
	parseMember,
	type JoinedRoomData,
	type JoinedRoomsData,
	type RawMatrixEvent,
	type SyncResponse
} from './matrix-events';

export class MatrixClientState {
	credentials = $state<MatrixCredentials | null>(null);
	connected = $state(false);
	syncError = $state<string | null>(null);

	messagesByRoom = $state<Record<string, ChatMessage[]>>({});
	membersByRoom = $state<Record<string, RoomMember[]>>({});
	typingByRoom = $state<Record<string, string[]>>({});
	agentUserIds = $state<Set<string>>(new Set());

	private abortController: AbortController | null = null;
	private nextBatch = '';
	private running = false;
	private txnCounter = 0;

	setKnownAgents(ids: string[]) {
		this.agentUserIds = new Set(ids);
	}

	async start(credentials: MatrixCredentials) {
		this.stop();
		this.credentials = credentials;
		this.running = true;
		this.abortController = new AbortController();
		this.runSyncLoop();
	}

	stop() {
		this.running = false;
		if (this.abortController) {
			this.abortController.abort();
			this.abortController = null;
		}
		this.connected = false;
	}

	private authHeaders(): HeadersInit {
		if (!this.credentials?.access_token) return {};
		return {
			Authorization: `Bearer ${this.credentials.access_token}`,
			'Content-Type': 'application/json'
		};
	}

	private async runSyncLoop() {
		while (this.running) {
			await this.syncOnce();
		}
	}

	private async syncOnce() {
		try {
			const params = new URLSearchParams({ timeout: '30000' });
			if (this.nextBatch) params.set('since', this.nextBatch);
			const res = await fetch(`/_matrix/client/v3/sync?${params.toString()}`, {
				headers: this.authHeaders(),
				signal: this.abortController?.signal
			});
			if (!res.ok) throw new Error(`Status ${res.status}`);
			const data = (await res.json()) as SyncResponse;
			this.nextBatch = data.next_batch;
			this.connected = true;
			this.syncError = null;
			this.processRooms(data.rooms?.join);
		} catch (err) {
			if (!this.running) return;
			this.syncError = err instanceof Error ? err.message : 'Sync disconnected';
			await new Promise((resolve) => setTimeout(resolve, 3000));
		}
	}

	private processRooms(rooms?: JoinedRoomsData) {
		if (!rooms) return;
		for (const [roomId, room] of Object.entries(rooms)) {
			this.consumeRoom(roomId, room);
		}
	}

	private consumeRoom(roomId: string, room: JoinedRoomData) {
		this.consumeState(roomId, room.state?.events);
		this.consumeTimeline(roomId, room.timeline?.events);
		this.consumeTyping(roomId, room.ephemeral?.events);
	}

	private consumeState(roomId: string, events?: RawMatrixEvent[]) {
		if (!events) return;
		for (const event of events) {
			const member = parseMember(event, this.isUserAgent(event.sender));
			if (!member) continue;
			const current = (this.membersByRoom[roomId] || []).filter((m) => m.userId !== member.userId);
			this.membersByRoom[roomId] = [...current, member];
		}
	}

	private consumeTimeline(roomId: string, events?: RawMatrixEvent[]) {
		if (!events) return;
		for (const event of events) {
			if (event.type === 'm.reaction') {
				this.applyReaction(roomId, event);
			} else if (event.type === 'm.room.message') {
				this.applyMessage(roomId, event);
			}
		}
	}

	private consumeTyping(roomId: string, events?: RawMatrixEvent[]) {
		if (!events) return;
		for (const event of events) {
			if (event.type === 'm.typing') {
				const typers = (event.content?.user_ids as string[]) || [];
				this.typingByRoom[roomId] = typers.filter((id) => id !== this.credentials?.user_id);
			}
		}
	}

	private isUserAgent(userId: string): boolean {
		return this.agentUserIds.has(userId) || userId.startsWith('@agent-');
	}

	private applyMessage(roomId: string, event: RawMatrixEvent) {
		const current = this.messagesByRoom[roomId] || [];
		if (current.some((m) => m.id === event.event_id)) return;
		const msg = parseChatMessage(roomId, event, this.isUserAgent(event.sender));
		this.messagesByRoom[roomId] = [...current, msg].sort((a, b) => a.timestamp - b.timestamp);
	}

	private applyReaction(roomId: string, event: RawMatrixEvent) {
		const relatesTo = event.content?.['m.relates_to'];
		if (relatesTo?.rel_type !== 'm.annotation') return;
		const targetId = relatesTo.event_id;
		const key = relatesTo.key;
		if (!targetId || !key) return;
		const current = this.messagesByRoom[roomId] || [];
		const target = current.find((m) => m.id === targetId);
		if (!target) return;
		const users = target.reactions[key] || [];
		if (!users.includes(event.sender)) {
			target.reactions[key] = [...users, event.sender];
			this.messagesByRoom[roomId] = [...current];
		}
	}

	async sendMessage(roomId: string, body: string, threadRootId?: string): Promise<string> {
		this.txnCounter++;
		const txnId = `m${Date.now()}_${this.txnCounter}`;
		const content: Record<string, any> = { msgtype: 'm.text', body };
		if (threadRootId) {
			content['m.relates_to'] = { rel_type: 'm.thread', event_id: threadRootId, is_falling_back: true };
		}
		const res = await fetch(`/_matrix/client/v3/rooms/${encodeURIComponent(roomId)}/send/m.room.message/${txnId}`, {
			method: 'PUT',
			headers: this.authHeaders(),
			body: JSON.stringify(content)
		});
		if (!res.ok) throw new Error(`Send failed: ${res.status}`);
		const data = (await res.json()) as { event_id: string };
		return data.event_id;
	}

	async sendReaction(roomId: string, eventId: string, key: string): Promise<void> {
		this.txnCounter++;
		const txnId = `m${Date.now()}_${this.txnCounter}`;
		const content = { 'm.relates_to': { rel_type: 'm.annotation', event_id: eventId, key } };
		await fetch(`/_matrix/client/v3/rooms/${encodeURIComponent(roomId)}/send/m.reaction/${txnId}`, {
			method: 'PUT',
			headers: this.authHeaders(),
			body: JSON.stringify(content)
		});
	}

	async sendTyping(roomId: string, typing: boolean): Promise<void> {
		if (!this.credentials) return;
		await fetch(`/_matrix/client/v3/rooms/${encodeURIComponent(roomId)}/typing/${encodeURIComponent(this.credentials.user_id)}`, {
			method: 'PUT',
			headers: this.authHeaders(),
			body: JSON.stringify({ typing, timeout: 4000 })
		});
	}

	async createDirectMessage(targetUserId: string): Promise<string> {
		const res = await fetch('/_matrix/client/v3/createRoom', {
			method: 'POST',
			headers: this.authHeaders(),
			body: JSON.stringify({ preset: 'trusted_private_chat', invite: [targetUserId], is_direct: true })
		});
		if (!res.ok) throw new Error(`Create DM failed: ${res.status}`);
		const data = (await res.json()) as { room_id: string };
		return data.room_id;
	}

	async loadRoomMembers(roomId: string): Promise<RoomMember[]> {
		const res = await fetch(`/_matrix/client/v3/rooms/${encodeURIComponent(roomId)}/joined_members`, {
			headers: this.authHeaders()
		});
		if (!res.ok) return [];
		const data = (await res.json()) as { joined: Record<string, { display_name?: string; avatar_url?: string }> };
		const members: RoomMember[] = Object.entries(data.joined).map(([userId, d]) => ({
			userId,
			displayName: d.display_name || userId.split(':')[0].slice(1),
			avatarUrl: d.avatar_url,
			isAgent: this.isUserAgent(userId)
		}));
		this.membersByRoom[roomId] = members;
		return members;
	}
}

export const matrix = new MatrixClientState();
