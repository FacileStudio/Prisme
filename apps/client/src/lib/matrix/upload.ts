import { MsgType, THREAD_RELATION_TYPE, type MatrixClient } from 'matrix-js-sdk';

type Size = { w: number; h: number };

function imageSize(file: File): Promise<Size | null> {
	if (!file.type.startsWith('image/')) return Promise.resolve(null);
	return new Promise((resolve) => {
		const url = URL.createObjectURL(file);
		const img = new Image();
		img.onload = () => {
			resolve({ w: img.naturalWidth, h: img.naturalHeight });
			URL.revokeObjectURL(url);
		};
		img.onerror = () => {
			resolve(null);
			URL.revokeObjectURL(url);
		};
		img.src = url;
	});
}

/**
 * Uploads a file to the media repository and posts it as `m.image` or `m.file`,
 * with the mimetype, size and, for images, the dimensions.
 */
export async function sendFile(client: MatrixClient, roomId: string, file: File, threadId?: string | null) {
	const [size, upload] = await Promise.all([imageSize(file), client.uploadContent(file, { type: file.type })]);
	const info: Record<string, unknown> = { mimetype: file.type || 'application/octet-stream', size: file.size };
	if (size) Object.assign(info, size);
	const content: Record<string, unknown> = {
		msgtype: size ? MsgType.Image : MsgType.File,
		body: file.name,
		filename: file.name,
		url: upload.content_uri,
		info
	};
	if (threadId) {
		content['m.relates_to'] = {
			rel_type: THREAD_RELATION_TYPE.name,
			event_id: threadId,
			is_falling_back: true,
			'm.in_reply_to': { event_id: threadId }
		};
	}
	await client.sendMessage(roomId, threadId ?? null, content as never);
}
