import { M_POLL_END, M_POLL_RESPONSE, M_POLL_START, type MatrixClient, type MatrixEvent, type Poll, type Room } from 'matrix-js-sdk';
import { isPollStartType, responseAnswers, type PollVote } from '$lib/chat/poll';
import { pollContent, type PollDraft } from '$lib/chat/poll-draft';
import { relation, type Target } from './messages';

/** Sends a new poll, in a thread when the target names one. */
export async function sendPoll(client: MatrixClient, target: Target, draft: PollDraft, lastInThread?: string) {
	const content = pollContent(draft, () => crypto.randomUUID());
	const rel = relation({ roomId: target.roomId, threadId: target.threadId }, lastInThread);
	if (rel) content['m.relates_to'] = rel;
	await client.sendEvent(target.roomId, target.threadId ?? null, M_POLL_START.name as never, content as never);
}

/** Records the user's answers to a poll. An empty list withdraws the vote. */
export async function vote(client: MatrixClient, poll: MatrixEvent, answers: string[]) {
	const content = {
		[M_POLL_RESPONSE.name]: { answers },
		'm.relates_to': { rel_type: 'm.reference', event_id: poll.getId() }
	};
	await client.sendEvent(poll.getRoomId()!, null, M_POLL_RESPONSE.name as never, content as never);
}

/** Closes a poll, so no later vote counts. */
export async function endPoll(client: MatrixClient, poll: MatrixEvent, text: string) {
	const content = {
		[M_POLL_END.name]: {},
		'm.text': text,
		body: text,
		'm.relates_to': { rel_type: 'm.reference', event_id: poll.getId() }
	};
	await client.sendEvent(poll.getRoomId()!, null, M_POLL_END.name as never, content as never);
}

/** The SDK's poll model for a start event, which fetches every response and tracks the end. */
export async function pollModel(room: Room, event: MatrixEvent): Promise<Poll | undefined> {
	const id = event.getId();
	if (!id || !isPollStartType(event.getType())) return undefined;
	if (!room.polls.has(id)) await room.processPollEvents([event]);
	return room.polls.get(id);
}

/** The votes a poll's responses hold, ready for `tally`. */
export async function pollVotes(poll: Poll): Promise<PollVote[]> {
	const responses = await poll.getResponses();
	return responses.getRelations().filter((e) => !e.isRedacted()).map((e) => ({
		sender: e.getSender() ?? '',
		ts: e.getTs(),
		answers: responseAnswers(e.getContent())
	}));
}

/** When a poll closed, if it did. */
export function endedAt(room: Room, poll: Poll): number | undefined {
	const id = poll.endEventId;
	return id ? (room.findEventById(id)?.getTs() ?? Date.now()) : undefined;
}
