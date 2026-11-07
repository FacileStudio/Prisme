import type { Row } from '../row';

/** The side panels: thread, members and search. */
export const panels = {
	'panel.thread': ['Thread', 'Fil', 'Hilo', 'Thread'],
	'panel.members': ['Members', 'Membres', 'Miembros', 'Mitglieder'],
	'panel.search': ['Search', 'Recherche', 'Búsqueda', 'Suche'],
	'panel.closeThread': [
		'Close the thread panel',
		'Fermer le panneau du fil',
		'Cerrar el panel del hilo',
		'Thread-Bereich schließen'
	],
	'panel.closeMembers': [
		'Close the members panel',
		'Fermer le panneau des membres',
		'Cerrar el panel de miembros',
		'Mitgliederbereich schließen'
	],
	'panel.closeSearch': [
		'Close the search panel',
		'Fermer le panneau de recherche',
		'Cerrar el panel de búsqueda',
		'Suchbereich schließen'
	],
	'members.invited': ['{name} (invited)', '{name} (invité)', '{name} (invitado)', '{name} (eingeladen)'],
	'search.placeholderSpace': [
		'Search in {name}',
		'Rechercher dans {name}',
		'Buscar en {name}',
		'In {name} suchen'
	],
	'search.placeholderRoom': [
		'Search in this conversation',
		'Rechercher dans cette conversation',
		'Buscar en esta conversación',
		'In dieser Unterhaltung suchen'
	],
	'search.hintSpace': [
		'Search the messages in {name}.',
		'Recherchez dans les messages de {name}.',
		'Busca en los mensajes de {name}.',
		'Durchsuche die Nachrichten in {name}.'
	],
	'search.hintRoom': [
		'Search the messages in this conversation.',
		'Recherchez dans les messages de cette conversation.',
		'Busca en los mensajes de esta conversación.',
		'Durchsuche die Nachrichten in dieser Unterhaltung.'
	],
	'search.busy': ['Searching', 'Recherche en cours', 'Buscando', 'Suche läuft'],
	'search.failed': ['Search failed', 'La recherche a échoué', 'La búsqueda falló', 'Suche fehlgeschlagen'],
	'search.noMatch': [
		'No messages match “{term}”.',
		'Aucun message ne correspond à « {term} ».',
		'Ningún mensaje coincide con «{term}».',
		'Keine Nachricht passt zu „{term}“.'
	],
	'search.inRoom': ['in #{room}', 'dans #{room}', 'en #{room}', 'in #{room}'],
	'thread.loading': [
		'Loading the thread',
		'Chargement du fil',
		'Cargando el hilo',
		'Thread wird geladen'
	],
	'thread.placeholder': [
		'Reply in thread',
		'Répondre dans le fil',
		'Responder en el hilo',
		'Im Thread antworten'
	]
} as const satisfies Record<string, Row>;
