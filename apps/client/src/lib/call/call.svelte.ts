import { Room, RoomEvent, type Participant } from 'livekit-client';
import { prisme, type CallGrant } from '$lib/api/prisme';
import { attachAudio, clearAudio, detachAudio } from './audio';

export type CallStatus = 'idle' | 'connecting' | 'connected' | 'disconnected';

class CallState {
	active = $state(false);
	roomId = $state('');
	status = $state<CallStatus>('idle');
	error = $state('');
	participants = $state.raw<Participant[]>([]);
	micOn = $state(false);
	camOn = $state(false);
	canPlayAudio = $state(true);
	version = $state(0);
	private room: Room | null = null;
	private grant: CallGrant | null = null;

	/** Asks Prisme for a call token for a Matrix room and joins it. */
	async start(roomId: string, join = false) {
		if (this.active && this.roomId === roomId && this.status === 'connected') return;
		this.leave();
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
			this.error = err instanceof Error ? err.message : 'Could not start the call';
			throw err;
		}
	}

	private async connect(grant: CallGrant) {
		const room = new Room({ adaptiveStream: true, dynacast: true });
		this.room = room;
		const refresh = () => this.refresh();
		room
			.on(RoomEvent.TrackSubscribed, (track) => (attachAudio(track), refresh()))
			.on(RoomEvent.TrackUnsubscribed, (track) => (detachAudio(track), refresh()))
			.on(RoomEvent.AudioPlaybackStatusChanged, () => (this.canPlayAudio = room.canPlaybackAudio))
			.on(RoomEvent.Disconnected, () => this.room === room && (this.status = 'disconnected'));
		[RoomEvent.ParticipantConnected, RoomEvent.ParticipantDisconnected, RoomEvent.LocalTrackPublished,
			RoomEvent.LocalTrackUnpublished, RoomEvent.TrackMuted, RoomEvent.TrackUnmuted, RoomEvent.ActiveSpeakersChanged]
			.forEach((event) => room.on(event, refresh));
		await room.connect(grant.livekit_url, grant.token);
		this.status = 'connected';
		this.canPlayAudio = room.canPlaybackAudio;
		await room.localParticipant.setMicrophoneEnabled(true).catch(() => (this.error = 'Microphone unavailable. You joined muted.'));
		this.refresh();
	}

	private refresh() {
		if (!this.room) return;
		this.participants = [this.room.localParticipant, ...this.room.remoteParticipants.values()];
		this.micOn = this.room.localParticipant.isMicrophoneEnabled;
		this.camOn = this.room.localParticipant.isCameraEnabled;
		this.version++;
	}

	/** Reconnects with a fresh token after a drop. */
	async rejoin() {
		const roomId = this.roomId;
		await this.start(roomId, true);
	}

	/** Toggles the microphone or the camera. */
	async toggle(kind: 'mic' | 'cam') {
		const local = this.room?.localParticipant;
		if (!local) return;
		try {
			if (kind === 'mic') await local.setMicrophoneEnabled(!this.micOn);
			else await local.setCameraEnabled(!this.camOn);
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Device unavailable';
		}
		this.refresh();
	}

	/** Unblocks audio playback after a user gesture, when autoplay was refused. */
	async enableAudio() {
		await this.room?.startAudio();
		this.canPlayAudio = this.room?.canPlaybackAudio ?? true;
	}

	/** Leaves the call and releases devices. */
	leave() {
		const room = this.room;
		this.room = null;
		void room?.disconnect();
		clearAudio();
		this.active = false;
		this.status = 'idle';
		this.participants = [];
	}
}

export const call = new CallState();
