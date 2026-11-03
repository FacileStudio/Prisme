import { EventType, type MatrixClient, type Room } from 'matrix-js-sdk';

/** Reports whether the signed-in user may rename the room. */
export function canRename(client: MatrixClient, room: Room): boolean {
	const me = client.getUserId();
	return !!me && room.currentState.maySendStateEvent(EventType.RoomName, me);
}
