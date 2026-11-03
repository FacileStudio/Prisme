import {
	ClientEvent,
	RoomEvent,
	RoomMemberEvent,
	RoomStateEvent,
	ThreadEvent,
	UserEvent,
	type MatrixClient
} from 'matrix-js-sdk';

const WATCHED = [
	ClientEvent.Sync,
	ClientEvent.AccountData,
	ClientEvent.Room,
	RoomEvent.Timeline,
	RoomEvent.TimelineReset,
	RoomEvent.Redaction,
	RoomEvent.Receipt,
	RoomEvent.Name,
	RoomEvent.MyMembership,
	RoomEvent.LocalEchoUpdated,
	RoomEvent.UnreadNotifications,
	RoomEvent.AccountData,
	RoomMemberEvent.Typing,
	RoomMemberEvent.Membership,
	RoomMemberEvent.Name,
	RoomStateEvent.Events,
	ThreadEvent.New,
	ThreadEvent.Update,
	ThreadEvent.NewReply,
	UserEvent.Presence
] as const;

/**
 * Calls `onChange` at most once per animation frame whenever the client sees
 * something the UI renders. Returns the unsubscribe function.
 */
export function watchClient(client: MatrixClient, onChange: () => void): () => void {
	let queued = false;
	const bump = () => {
		if (queued) return;
		queued = true;
		requestAnimationFrame(() => {
			queued = false;
			onChange();
		});
	};
	const emitter = client as unknown as { on(e: string, f: () => void): void; off(e: string, f: () => void): void };
	WATCHED.forEach((event) => emitter.on(event, bump));
	return () => WATCHED.forEach((event) => emitter.off(event, bump));
}
