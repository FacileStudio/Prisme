import { Track, type RemoteTrack } from 'livekit-client';

const SINK_ID = 'prisme-call-audio';

function sink(): HTMLElement {
	const existing = document.getElementById(SINK_ID);
	if (existing) return existing;
	const element = document.createElement('div');
	element.id = SINK_ID;
	element.hidden = true;
	document.body.append(element);
	return element;
}

/**
 * Plays a remote audio track. livekit-client never plays remote audio by
 * itself: each track needs an attached `<audio>` element in the DOM.
 */
export function attachAudio(track: RemoteTrack): void {
	if (track.kind !== Track.Kind.Audio) return;
	const element = track.attach();
	element.dataset.track = track.sid ?? '';
	sink().append(element);
}

/** Stops playing a remote audio track and removes its elements. */
export function detachAudio(track: RemoteTrack): void {
	if (track.kind !== Track.Kind.Audio) return;
	track.detach().forEach((element) => element.remove());
}

/** Removes every attached audio element, after leaving a call. */
export function clearAudio(): void {
	document.getElementById(SINK_ID)?.replaceChildren();
}
