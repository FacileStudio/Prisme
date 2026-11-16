import { request, send } from './http';

export type SessionUser = {
	id: number;
	email: string;
	name: string;
	avatar_url: string;
	matrix_user_id: string;
	locale: string;
};

export type MatrixLogin = {
	homeserver: string;
	user_id: string;
	access_token: string;
	device_id: string;
};

export type Session = { user: SessionUser; matrix: MatrixLogin; calls_enabled: boolean };
export type AuthConfig = { sso_only: boolean; oidc_enabled: boolean };
export type Channel = { room_id: string; name: string; alias: string };
export type Space = { id: number; room_id: string; name: string; slug: string; channels: Channel[] };
export type DirectoryUser = { user_id: string; name: string; avatar_url?: string; agent?: boolean };
export type AgentUser = { user_id: string; name: string; owner_user_id: string };
export type Agent = { user_id: string; name: string; room_ids: string[] };
export type CallGrant = { call_id: string; livekit_url: string; token: string; join_url: string; e2ee_key: string };

export type SealedCapsule = { id: string; url: string; expires_at?: string };
export type CapsuleMeta = { id: string; exists: boolean; expires_at?: string };

export type LinkPreview = {
	url: string;
	app: string;
	title?: string;
	description?: string;
	image?: string;
	embed_url?: string;
	duration_ms?: number;
};

export type PushConfig = { native: boolean };
export type Invitations = { count: number; href?: string };

export const LOGIN_PATH = '/api/auth/oidc';

/** Every Prisme API route the client calls. */
export const prisme = {
	authConfig: () => request<AuthConfig>('/auth/config'),
	session: () => request<Session>('/session'),
	logout: () => send<void>('/auth/logout', 'POST'),
	spaces: () => request<Space[]>('/spaces'),
	createSpace: (name: string) => send<Space>('/spaces', 'POST', { name }),
	createChannel: (spaceId: number, name: string, topic = '') =>
		send<Channel>(`/spaces/${spaceId}/channels`, 'POST', { name, topic }),
	users: (q: string) => request<DirectoryUser[]>(`/users?q=${encodeURIComponent(q)}`),
	usersAgents: () => request<AgentUser[]>('/users/agents'),
	agents: () => request<Agent[]>('/agents'),
	setAgentRooms: (userId: string, roomIds: string[]) =>
		send<{ room_ids: string[] }>(`/agents/${encodeURIComponent(userId)}/rooms`, 'PUT', { room_ids: roomIds }),
	call: (roomId: string, join = false) => send<CallGrant>('/calls', 'POST', { room_id: roomId, join }),
	sealCapsule: (ciphertext: string, expiresIn: string) =>
		send<SealedCapsule>('/capsules', 'POST', { ciphertext, expires_in: expiresIn }),
	capsuleMeta: (id: string) => request<CapsuleMeta>(`/capsules/${encodeURIComponent(id)}`),
	openCapsule: (id: string) => send<{ ciphertext: string }>(`/capsules/${encodeURIComponent(id)}/open`, 'POST'),
	unfurl: (url: string) => request<LinkPreview>(`/unfurl?url=${encodeURIComponent(url)}`),
	pushConfig: () => request<PushConfig>('/push'),
	invitations: () => request<Invitations>('/me/invitations'),
	registerPushDevice: (token: string) => send<void>('/push/devices', 'POST', { token })
};
