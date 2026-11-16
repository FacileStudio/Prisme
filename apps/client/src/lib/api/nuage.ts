import { request, send } from './http';
import type { NuageShared } from '$lib/chat/nuage';

export type NuageSpace = { id: number; name: string; role: string };
export type NuageFolder = { id: number; name: string; parent_id: number | null };
export type NuageFile = { id: number; name: string; mime_type: string; size: number; folder_id: number | null; updated_at: string };
export type NuageListing = { folders: NuageFolder[]; files: NuageFile[] };

/** Where a listing looks: a folder, a space's root, or the user's own root when both are empty. */
export type NuagePlace = { folderId?: number; spaceId?: number };

function placeQuery(place: NuagePlace): string {
	const query = new URLSearchParams();
	if (place.folderId) query.set('folder_id', String(place.folderId));
	if (place.spaceId) query.set('space_id', String(place.spaceId));
	const text = query.toString();
	return text ? `?${text}` : '';
}

/** Prisme's picker over the signed-in user's Nuage drive. Nuage decides what the user can see. */
export const nuage = {
	spaces: () => request<{ spaces: NuageSpace[] }>('/nuage/spaces'),
	list: (place: NuagePlace) => request<NuageListing>(`/nuage/list${placeQuery(place)}`),
	search: (q: string) => request<{ files: NuageFile[] }>(`/nuage/search?q=${encodeURIComponent(q)}`),
	share: (fileId: number) => send<NuageShared>('/nuage/share', 'POST', { file_id: fileId })
};
