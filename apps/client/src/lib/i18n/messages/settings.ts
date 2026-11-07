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
	'settings.profileBody': [
		'Your name and picture come from your Facile account.',
		'Votre nom et votre photo viennent de votre compte Facile.',
		'Tu nombre y tu foto proceden de tu cuenta de Facile.',
		'Dein Name und dein Bild stammen aus deinem Facile-Konto.'
	],
	'settings.matrixId': ['Matrix ID', 'Identifiant Matrix', 'ID de Matrix', 'Matrix-ID'],
	'settings.appearance': ['Appearance', 'Apparence', 'Apariencia', 'Darstellung'],
	'settings.appearanceBody': [
		'Stored in this browser. Every device keeps its own choice.',
		'Enregistré dans ce navigateur. Chaque appareil garde son propre choix.',
		'Se guarda en este navegador. Cada dispositivo conserva su propia elección.',
		'In diesem Browser gespeichert. Jedes Gerät behält seine eigene Wahl.'
	],
	'settings.scheme': ['Colour scheme', 'Thème de couleurs', 'Esquema de color', 'Farbschema'],
	'settings.schemeBody': [
		'System follows your operating system while Prisme is open.',
		'« Système » suit votre système d’exploitation tant que Prisme est ouvert.',
		'«Sistema» sigue tu sistema operativo mientras Prisme está abierto.',
		'„System“ folgt deinem Betriebssystem, solange Prisme geöffnet ist.'
	],
	'settings.modeSystem': ['System', 'Système', 'Sistema', 'System'],
	'settings.modeLight': ['Light', 'Clair', 'Claro', 'Hell'],
	'settings.modeDark': ['Dark', 'Sombre', 'Oscuro', 'Dunkel'],
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
	'settings.session': ['Session', 'Session', 'Sesión', 'Sitzung'],
	'settings.signOut': ['Sign out', 'Se déconnecter', 'Cerrar sesión', 'Abmelden']
} as const satisfies Record<string, Row>;
