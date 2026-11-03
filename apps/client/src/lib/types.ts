export type UserProfile = {
	id: number;
	email: string;
	name: string;
	avatar_url?: string;
	matrix_user_id: string;
};

export type MatrixCredentials = {
	homeserver: string;
	user_id: string;
	access_token: string;
	device_id: string;
};

export type SessionResponse = {
	user: UserProfile;
	matrix: MatrixCredentials;
	calls_enabled: boolean;
};

export type ChannelRecord = {
	room_id: string;
	name: string;
	alias: string;
};

export type SpaceRecord = {
	id: number;
	room_id: string;
	name: string;
	slug: string;
	channels: ChannelRecord[];
};

export type AgentRecord = {
	user_id: string;
	name: string;
	room_ids: string[];
};

export type AgentAvailabilityResponse = {
	name: string;
	slug: string;
	available: boolean;
	reason?: string;
};

export type UserRecord = {
	user_id: string;
	name: string;
};

export type CallResponse = {
	call_id: string;
	livekit_url: string;
	token: string;
	join_url: string;
};

export type CallPayload = {
	call_id: string;
	join_url: string;
	started_by: string;
};

export type ChatMessage = {
	id: string;
	roomId: string;
	sender: string;
	senderName: string;
	senderAvatar?: string;
	isAgent: boolean;
	body: string;
	timestamp: number;
	call?: CallPayload;
	replyToId?: string;
	threadRootId?: string;
	reactions: Record<string, string[]>;
};

export type RoomMember = {
	userId: string;
	displayName: string;
	avatarUrl?: string;
	isAgent: boolean;
};
