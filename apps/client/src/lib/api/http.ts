import { apiBase } from '$lib/native/platform';
import { clearToken, getToken } from '$lib/native/token';

const SAFE_METHODS = new Set(['GET', 'HEAD', 'OPTIONS']);

/** An API failure that keeps the HTTP status and the server's error code. */
export class ApiError extends Error {
	status: number;
	code: string;

	constructor(message: string, status: number, code = '') {
		super(message);
		this.name = 'ApiError';
		this.status = status;
		this.code = code;
	}
}

type ErrorPayload = { error?: { code?: string; message?: string } };

async function readError(response: Response): Promise<ApiError> {
	const payload = (await response.json().catch(() => ({}))) as ErrorPayload;
	const message = payload.error?.message ?? `Request failed (${response.status})`;
	return new ApiError(message, response.status, payload.error?.code ?? '');
}

/**
 * Calls the Prisme API. The web app sends the porte session cookie, and its mutating
 * requests carry `X-Facile-CSRF`, which porte requires of every cookie-authenticated write.
 * The installed apps send their bearer token instead, and forget it when porte refuses it.
 */
export async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
	const method = (init.method ?? 'GET').toUpperCase();
	const headers = new Headers(init.headers);
	const token = await getToken();
	if (token) headers.set('Authorization', `Bearer ${token}`);
	else if (!SAFE_METHODS.has(method)) headers.set('X-Facile-CSRF', '1');
	if (init.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json');
	const credentials = token ? 'omit' : 'same-origin';
	const response = await fetch(`${apiBase()}/api${path}`, { ...init, method, headers, credentials });
	if (response.status === 401 && token) await clearToken();
	if (!response.ok) throw await readError(response);
	if (response.status === 204) return undefined as T;
	const text = await response.text();
	return (text ? JSON.parse(text) : undefined) as T;
}

/** Sends a JSON body with the given method. */
export function send<T>(path: string, method: string, body?: unknown): Promise<T> {
	return request<T>(path, { method, body: body === undefined ? undefined : JSON.stringify(body) });
}
