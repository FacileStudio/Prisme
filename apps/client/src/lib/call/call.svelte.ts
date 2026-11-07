import { Room, RoomEvent, type Participant } from 'livekit-client';
import { prisme, type CallGrant } from '$lib/api/prisme';
import { attachAudio, clearAudio, detachAudio } from './audio';
import { t } from '$lib/i18n/index.svelte';
import { canEncrypt, canShareScreen, encryptedRoom, shareError } from './media';

export type CallStatus = 'idle' | 'connecting' | 'connected' | 'disconnected';

class CallState {
	active = $state(false);
	roomId = $state('');
	status = $state<CallStatus>('idle');
	error = $state('');
	participants = $state.raw<Participant[]>([]);
	micOn = $state(false);
	camOn = $state(false);
	screenOn = $state(false);
	encrypted = $state(false);
	canPlayAudio = $state(true);
	version = $state(0);
	readonly canShare = canShareScreen();
	private room: Room | null = null;
	private worker: Worker | null = null;
	private grant: CallGrant | null = null;

	/**
	 * Asks Prisme for a call token and key for a Matrix room and joins it,
	 * end-to-end encrypted. A browser that cannot encrypt is refused.
	 */
	async start(roomId: string, join = false) {
		if (this.active && this.roomId === roomId && this.status === 'connected') return;
		this.leave();
		if (!canEncrypt()) {
			const message = t('call.unsupported');
			this.error = message;
			throw new Error(message);
		}
		this.active = true;
		this.roomId = roomId;
		this.status = 'connecting';
		this.error = '';
		try {
			this.grant = await prisme.call(roomId, join);
		} catch (err) {
			this.active = false;
			this.status = 'idle';
			throw err;
		}
		try {
			await this.connect(this.grant);
		} catch (err) {
			this.status = 'disconnected';
			this.error = err instanceof Error ? err.message : t('call.startFailed');
			throw err;
		}
	}

	private async connect(grant: CallGrant) {
		const { room, worker } = await encryptedRoom(grant.e2ee_key);
		this.room = room;
		this.worker = worker;
		const refresh = () => this.refresh();
		room
			.on(RoomEvent.TrackSubscribed, (track) => (attachAudio(track), refresh()))
			.on(RoomEvent.TrackUnsubscribed, (track) => (detachAudio(track), refresh()))
			.on(RoomEvent.AudioPlaybackStatusChanged, () => (this.canPlayAudio = room.canPlaybackAudio))
			.on(RoomEvent.EncryptionError, () => this.room === room && (this.error = t('call.decryptFailed')))
			.on(RoomEvent.Disconnected, () => this.room === room && (this.status = 'disconnected'));
		[RoomEvent.ParticipantConnected, RoomEvent.ParticipantDisconnected, RoomEvent.LocalTrackPublished,
			RoomEvent.LocalTrackUnpublished, RoomEvent.TrackMuted, RoomEvent.TrackUnmuted, RoomEvent.ActiveSpeakersChanged,
			RoomEvent.ParticipantEncryptionStatusChanged]
			.forEach((event) => room.on(event, refresh));
		await room.connect(grant.livekit_url, grant.token);
		this.status = 'connected';
		this.canPlayAudio = room.canPlaybackAudio;
		await room.localParticipant.setMicrophoneEnabled(true).catch(() => (this.error = t('call.micUnavailable')));
		this.refresh();
	}

	private refresh() {
		if (!this.room) return;
		this.participants = [this.room.localParticipant, ...this.room.remoteParticipants.values()];
		this.micOn = this.room.localParticipant.isMicrophoneEnabled;
		this.camOn = this.room.localParticipant.isCameraEnabled;
		this.screenOn = this.room.localParticipant.isScreenShareEnabled;
		this.encrypted = this.room.isE2EEEnabled;
		this.version++;
	}

	/** Reconnects with a fresh token after a drop. */
	async rejoin() {
		const roomId = this.roomId;
		await this.start(roomId, true);
	}

	/**
	 * Toggles the microphone, the camera or the screen share. Closing the
	 * screen picker is not an error.
	 */
	async toggle(kind: 'mic' | 'cam' | 'screen') {
		const local = this.room?.localParticipant;
		if (!local) return;
		try {
			if (kind === 'mic') await local.setMicrophoneEnabled(!this.micOn);
			else if (kind === 'cam') await local.setCameraEnabled(!this.camOn);
			else await local.setScreenShareEnabled(!this.screenOn, { audio: true });
		} catch (err) {
			const message = kind === 'screen' ? shareError(err) : err instanceof Error ? err.message : t('call.deviceUnavailable');
			if (message) this.error = message;
		}
		this.refresh();
	}

	/** Unblocks audio playback after a user gesture, when autoplay was refused. */
	async enableAudio() {
		await this.room?.startAudio();
		this.canPlayAudio = this.room?.canPlaybackAudio ?? true;
	}

	/** Leaves the call, releases devices and stops the encryption worker. */
	leave() {
		const room = this.room;
		const worker = this.worker;
		this.room = null;
		this.worker = null;
		void room?.disconnect().finally(() => worker?.terminate());
		clearAudio();
		this.active = false;
		this.status = 'idle';
		this.encrypted = false;
		this.screenOn = false;
		this.participants = [];
	}
}

export const call = new CallState();
