import type { Row } from '../row';

/** The settings page. */
export const settings = {
	'settings.title': ['Settings · Prisme', 'Paramètres · Prisme', 'Ajustes · Prisme', 'Einstellungen · Prisme'],
	'settings.heading': ['Settings', 'Paramètres', 'Ajustes', 'Einstellungen'],
	'settings.subtitle': [
		'You, this browser, and the agents that chat on your behalf.',
		'Vous, ce navigateur et les agents qui discutent en votre nom.',
		'Tú, este navegador y los agentes que chatean en tu nombre.',
		'Du, dieser Browser und die Agents, die in deinem Namen schreiben.'
	],
	'settings.profile': ['Profile', 'Profil', 'Perfil', 'Profil'],
	'settings.preferences': ['Preferences', 'Préférences', 'Preferencias', 'Einstellungen'],
	'settings.matrixId': ['Matrix ID', 'Identifiant Matrix', 'ID de Matrix', 'Matrix-ID'],
	'settings.notifications': ['Notifications', 'Notifications', 'Notificaciones', 'Benachrichtigungen'],
	'settings.browserNotifications': [
		'Browser notifications',
		'Notifications du navigateur',
		'Notificaciones del navegador',
		'Browser-Benachrichtigungen'
	],
	'settings.notifyUnsupported': [
		'Not supported by this browser.',
		'Non pris en charge par ce navigateur.',
		'Este navegador no lo admite.',
		'Von diesem Browser nicht unterstützt.'
	],
	'settings.notifyBlocked': [
		'Blocked in your browser settings.',
		'Bloqué dans les réglages de votre navigateur.',
		'Bloqueado en los ajustes de tu navegador.',
		'In den Browser-Einstellungen blockiert.'
	],
	'settings.notifyHint': [
		'For mentions and direct messages while Prisme is in the background.',
		'Pour les mentions et les messages directs lorsque Prisme est en arrière-plan.',
		'Para menciones y mensajes directos mientras Prisme está en segundo plano.',
		'Für Erwähnungen und Direktnachrichten, während Prisme im Hintergrund läuft.'
	],
	'settings.agents': ['Agents', 'Agents', 'Agentes', 'Agents'],
	'settings.agentsBody': [
		'Bots you own. Bulle creates them and adds them to your spaces.',
		'Les bots qui vous appartiennent. Bulle les crée et les ajoute à vos espaces.',
		'Los bots que te pertenecen. Bulle los crea y los añade a tus espacios.',
		'Bots, die dir gehören. Bulle erstellt sie und fügt sie deinen Bereichen hinzu.'
	],
	'settings.agentsNone': [
		'You have no agents yet.',
		'Vous n’avez pas encore d’agent.',
		'Todavía no tienes agentes.',
		'Du hast noch keine Agents.'
	],
	'settings.agentsFailed': [
		'Could not load your agents.',
		'Impossible de charger vos agents.',
		'No se pudieron cargar tus agentes.',
		'Deine Agents konnten nicht geladen werden.'
	],
	'settings.agentRooms': [
		{ one: 'In {count} room', other: 'In {count} rooms' },
		{ one: 'Dans {count} salon', other: 'Dans {count} salons' },
		{ one: 'En {count} sala', other: 'En {count} salas' },
		{ one: 'In {count} Raum', other: 'In {count} Räumen' }
	],
} as const satisfies Record<string, Row>;
