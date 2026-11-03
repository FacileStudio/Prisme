import { createClient, ClientEvent, HttpApiEvent, SyncState, type MatrixClient } from 'matrix-js-sdk';
import { ApiError } from '$lib/api/http';
import { prisme, type Session, type Space } from '$lib/api/prisme';
import { watchClient } from '$lib/matrix/events';
import { setMediaToken } from '$lib/matrix/media';

export type Status = 'loading' | 'ready' | 'signed-out' | 'error';

function waitForSync(client: MatrixClient): Promise<void> {
	return new Promise((resolve, reject) => {
		const onSync = (state: SyncState) => {
			if (state === SyncState.Prepared || state === SyncState.Syncing) {
				client.off(ClientEvent.Sync, onSync);
				resolve();
			} else if (state === SyncState.Error && !client.isInitialSyncComplete()) {
				client.off(ClientEvent.Sync, onSync);
				reject(new Error('The chat server did not answer the first sync.'));
			}
		};
		client.on(ClientEvent.Sync, onSync);
	});
}

class AppState {
	status = $state<Status>('loading');
	error = $state('');
	session = $state<Session | null>(null);
	spaces = $state<Space[]>([]);
	client = $state.raw<MatrixClient | null>(null);
	version = $state(0);
	private unwatch: (() => void) | null = null;
	private retried = false;

	/** Fetches the session and starts the Matrix client. */
	async boot() {
		try {
			this.session = await prisme.session();
		} catch (err) {
			this.fail(err);
			return;
		}
		await this.start(this.session);
		void this.refreshSpaces();
	}

	private async start(session: Session) {
		const { homeserver, access_token, user_id, device_id } = session.matrix;
		const client = createClient({ baseUrl: homeserver, accessToken: access_token, userId: user_id, deviceId: device_id });
		setMediaToken(homeserver, access_token);
		client.on(HttpApiEvent.SessionLoggedOut, () => void this.recover());
		this.unwatch = watchClient(client, () => this.version++);
		this.client = client;
		await client.startClient({ initialSyncLimit: 30, threadSupport: true, lazyLoadMembers: true });
		try {
			await waitForSync(client);
			this.status = 'ready';
		} catch (err) {
			this.fail(err);
		}
	}

	private async recover() {
		this.stop();
		if (this.retried) {
			this.status = 'signed-out';
			return;
		}
		this.retried = true;
		await this.boot();
	}

	private fail(err: unknown) {
		if (err instanceof ApiError && err.status === 401) {
			this.status = 'signed-out';
			return;
		}
		this.error = err instanceof Error ? err.message : String(err);
		this.status = 'error';
	}

	/** Reloads the spaces the user belongs to from the Prisme API. */
	async refreshSpaces() {
		this.spaces = await prisme.spaces().catch(() => this.spaces);
	}

	/** Stops the Matrix client and drops its listeners. */
	stop() {
		this.unwatch?.();
		this.unwatch = null;
		this.client?.stopClient();
		this.client = null;
	}

	/** Revokes the session server-side, then leaves for the login page. */
	async logout() {
		await prisme.logout().catch(() => undefined);
		this.stop();
		this.status = 'signed-out';
	}
}

export const app = new AppState();
