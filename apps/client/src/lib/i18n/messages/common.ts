import type { Row } from '../row';

/** Words used across the whole app. */
export const common = {
	'common.cancel': ['Cancel', 'Annuler', 'Cancelar', 'Abbrechen'],
	'common.working': ['Working…', 'En cours…', 'Procesando…', 'Wird ausgeführt…'],
	'common.agent': ['Agent', 'Agent', 'Agente', 'Agent'],
	'common.member': ['Member', 'Membre', 'Miembro', 'Mitglied'],
	'common.members': ['Members', 'Membres', 'Miembros', 'Mitglieder'],
	'common.name': ['Name', 'Nom', 'Nombre', 'Name'],
	'common.create': ['Create', 'Créer', 'Crear', 'Erstellen'],
	'common.search': ['Search', 'Rechercher', 'Buscar', 'Suchen'],
	'common.open': ['Open', 'Ouvrir', 'Abrir', 'Öffnen'],
	'common.join': ['Join', 'Rejoindre', 'Unirse', 'Beitreten'],
	'common.joining': ['Joining…', 'Connexion…', 'Uniéndose…', 'Beitritt…'],
	'common.delete': ['Delete', 'Supprimer', 'Eliminar', 'Löschen'],
	'avatar.online': ['Online', 'En ligne', 'En línea', 'Online'],
	'avatar.offline': ['Offline', 'Hors ligne', 'Desconectado', 'Offline'],
	'badge.mentions': [
		{ one: '{count} unread mention', other: '{count} unread mentions' },
		{ one: '{count} mention non lue', other: '{count} mentions non lues' },
		{ one: '{count} mención sin leer', other: '{count} menciones sin leer' },
		{ one: '{count} ungelesene Erwähnung', other: '{count} ungelesene Erwähnungen' }
	],
	'badge.messages': [
		{ one: '{count} unread message', other: '{count} unread messages' },
		{ one: '{count} message non lu', other: '{count} messages non lus' },
		{ one: '{count} mensaje sin leer', other: '{count} mensajes sin leer' },
		{ one: '{count} ungelesene Nachricht', other: '{count} ungelesene Nachrichten' }
	]
} as const satisfies Record<string, Row>;
