import type {
	AgentAvailabilityResponse,
	AgentRecord,
	CallResponse,
	ChannelRecord,
	MatrixCredentials,
	SessionResponse,
	SpaceRecord,
	UserRecord
} from './types';

const backendBaseUrl = '/api';

export const ssoLoginUrl = `${backendBaseUrl}/auth/oidc`;

type ApiErrorPayload = {
	error?: {
		code?: string;
		message?: string;
	};
};

export class ApiError extends Error {
	status: number;
	code?: string;

	constructor(message: string, status: number, code?: string) {
		super(message);
		this.status = status;
		this.code = code;
	}
}

async function authedInit(options: RequestInit = {}): Promise<RequestInit> {
	const headers = new Headers(options.headers);
	if (options.body && !headers.has('Content-Type')) {
		headers.set('Content-Type', 'application/json');
	}
	if (options.method && options.method !== 'GET') {
		headers.set('X-Facile-CSRF', '1');
	}
	return { ...options, headers, credentials: 'same-origin' };
}

export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
	const response = await fetch(`${backendBaseUrl}${path}`, await authedInit(options));
	if (!response.ok) {
		const payload = (await response.json().catch(() => undefined)) as ApiErrorPayload | undefined;
		throw new ApiError(
			payload?.error?.message || `Request failed with status ${response.status}`,
			response.status,
			payload?.error?.code
		);
	}
	if (response.status === 204) {
		return undefined as T;
	}
	return (await response.json()) as T;
}

export const backend = {
	baseUrl: backendBaseUrl,

	getSession(): Promise<SessionResponse> {
		return apiFetch<SessionResponse>('/session');
	},

	logout(): Promise<void> {
		return apiFetch<void>('/auth/logout', { method: 'POST' });
	},

	listSpaces(): Promise<SpaceRecord[]> {
		return apiFetch<SpaceRecord[]>('/spaces');
	},

	createSpace(name: string): Promise<SpaceRecord> {
		return apiFetch<SpaceRecord>('/spaces', {
			method: 'POST',
			body: JSON.stringify({ name })
		});
	},

	createChannel(spaceId: number, name: string, topic?: string): Promise<ChannelRecord> {
		return apiFetch<ChannelRecord>(`/spaces/${spaceId}/channels`, {
			method: 'POST',
			body: JSON.stringify({ name, topic: topic || '' })
		});
	},

	listAgents(): Promise<AgentRecord[]> {
		return apiFetch<AgentRecord[]>('/agents');
	},

	checkAgentName(name: string): Promise<AgentAvailabilityResponse> {
		return apiFetch<AgentAvailabilityResponse>(`/agents/check?name=${encodeURIComponent(name)}`);
	},

	createAgent(name: string, avatarUrl?: string): Promise<MatrixCredentials> {
		return apiFetch<MatrixCredentials>('/agents', {
			method: 'POST',
			body: JSON.stringify({ name, avatar_url: avatarUrl || '' })
		});
	},

	setAgentRooms(userId: string, roomIds: string[]): Promise<{ room_ids: string[] }> {
		return apiFetch<{ room_ids: string[] }>(`/agents/${encodeURIComponent(userId)}/rooms`, {
			method: 'PUT',
			body: JSON.stringify({ room_ids: roomIds })
		});
	},

	deactivateAgent(userId: string): Promise<void> {
		return apiFetch<void>(`/agents/${encodeURIComponent(userId)}`, {
			method: 'DELETE'
		});
	},

	searchUsers(query: string): Promise<UserRecord[]> {
		return apiFetch<UserRecord[]>(`/users?q=${encodeURIComponent(query)}`);
	},

	startCall(roomId: string, join = false): Promise<CallResponse> {
		return apiFetch<CallResponse>('/calls', {
			method: 'POST',
			body: JSON.stringify({ room_id: roomId, join })
		});
	}
};
