import type { Row } from '../row';

/** The installed apps: sign-in, first-run screens, updates and the downloads page. */
export const apps = {
	'app.signInMismatch': [
		'The sign-in did not match. Try again.',
		'La connexion ne correspond pas. Réessayez.',
		'El inicio de sesión no coincide. Inténtalo de nuevo.',
		'Die Anmeldung passt nicht. Versuch es noch einmal.'
	],
	'app.signInFailed': [
		'Sign-in could not be completed. Try again.',
		'La connexion n’a pas pu aboutir. Réessayez.',
		'No se pudo completar el inicio de sesión. Inténtalo de nuevo.',
		'Die Anmeldung konnte nicht abgeschlossen werden. Versuch es noch einmal.'
	],
	'app.signInOpenFailed': [
		'The sign-in page could not open.',
		'La page de connexion n’a pas pu s’ouvrir.',
		'No se pudo abrir la página de inicio de sesión.',
		'Die Anmeldeseite konnte nicht geöffnet werden.'
	],
	'app.continue': ['Continue with Facile', 'Continuer avec Facile', 'Continuar con Facile', 'Weiter mit Facile'],
	'app.browserHint': [
		'Sign-in opens in your browser, then brings you back here.',
		'La connexion s’ouvre dans votre navigateur, puis vous ramène ici.',
		'El inicio de sesión se abre en tu navegador y luego te trae de vuelta.',
		'Die Anmeldung öffnet sich im Browser und bringt dich dann hierher zurück.'
	],
	'app.skip': ['Skip', 'Passer', 'Saltar', 'Überspringen'],
	'app.next': ['Next', 'Suivant', 'Siguiente', 'Weiter'],
	'app.slides': ['Introduction', 'Présentation', 'Presentación', 'Einführung'],
	'app.goToSlide': ['Go to step {n}', 'Aller à l’étape {n}', 'Ir al paso {n}', 'Zu Schritt {n}'],
	'app.updateReady': [
		'Prisme {version} is ready to install.',
		'Prisme {version} est prêt à être installé.',
		'Prisme {version} está listo para instalarse.',
		'Prisme {version} ist bereit zur Installation.'
	],
	'app.updateInstall': ['Install and restart', 'Installer et redémarrer', 'Instalar y reiniciar', 'Installieren und neu starten'],
	'app.updateFailed': [
		'The update could not be installed.',
		'La mise à jour n’a pas pu être installée.',
		'No se pudo instalar la actualización.',
		'Das Update konnte nicht installiert werden.'
	],
	'app.windowClose': ['Close', 'Fermer', 'Cerrar', 'Schließen'],
	'app.windowMinimize': ['Minimize', 'Réduire', 'Minimizar', 'Minimieren'],
	'app.windowMaximize': ['Maximize', 'Agrandir', 'Maximizar', 'Maximieren'],
	'app.window': ['Window', 'Fenêtre', 'Ventana', 'Fenster'],
	'intro.spaces.title': [
		'Your team, in one place.',
		'Votre équipe, au même endroit.',
		'Tu equipo, en un solo lugar.',
		'Dein Team, an einem Ort.'
	],
	'intro.spaces.body': [
		'Every space and channel you have on the web, ready the moment you open Prisme.',
		'Tous vos espaces et salons du web, prêts dès l’ouverture de Prisme.',
		'Todos tus espacios y canales de la web, listos en cuanto abres Prisme.',
		'Alle Spaces und Kanäle aus dem Web, sofort da, wenn du Prisme öffnest.'
	],
	'intro.people.title': [
		'People and agents, side by side.',
		'Des humains et des agents, côte à côte.',
		'Personas y agentes, codo con codo.',
		'Menschen und Agents, Seite an Seite.'
	],
	'intro.people.body': [
		'Write to a teammate or to your Bulle agent in the same conversation list.',
		'Écrivez à un collègue ou à votre agent Bulle depuis la même liste de conversations.',
		'Escribe a un compañero o a tu agente de Bulle desde la misma lista de conversaciones.',
		'Schreib Kolleginnen oder deinem Bulle-Agent aus derselben Unterhaltungsliste.'
	],
	'intro.calls.title': [
		'Only what needs you.',
		'Seulement ce qui vous concerne.',
		'Solo lo que te necesita.',
		'Nur, was dich betrifft.'
	],
	'intro.calls.body': [
		'Direct messages and mentions reach you. The rest of the chatter waits until you look.',
		'Les messages directs et les mentions vous parviennent. Le reste attend que vous regardiez.',
		'Te llegan los mensajes directos y las menciones. El resto espera a que mires.',
		'Direktnachrichten und Erwähnungen erreichen dich. Der Rest wartet, bis du nachsiehst.'
	],
	'intro.signIn.title': ['Sign in to Prisme', 'Connectez-vous à Prisme', 'Inicia sesión en Prisme', 'Bei Prisme anmelden'],
	'intro.signIn.body': [
		'Use your Facile account. Your spaces and conversations follow you.',
		'Utilisez votre compte Facile. Vos espaces et conversations vous suivent.',
		'Usa tu cuenta de Facile. Tus espacios y conversaciones te siguen.',
		'Nutze dein Facile-Konto. Deine Spaces und Unterhaltungen kommen mit.'
	],
	'welcome.title': [
		'Never miss a message',
		'Ne manquez aucun message',
		'No te pierdas ningún mensaje',
		'Verpasse keine Nachricht'
	],
	'welcome.body': [
		'Prisme tells you when someone writes to you directly or mentions you. Nothing else.',
		'Prisme vous prévient quand quelqu’un vous écrit directement ou vous mentionne. Rien d’autre.',
		'Prisme te avisa cuando alguien te escribe directamente o te menciona. Nada más.',
		'Prisme meldet sich, wenn dir jemand direkt schreibt oder dich erwähnt. Sonst nichts.'
	],
	'welcome.enable': ['Turn on notifications', 'Activer les notifications', 'Activar notificaciones', 'Benachrichtigungen einschalten'],
	'welcome.later': ['Not now', 'Plus tard', 'Ahora no', 'Später'],
	'welcome.open': ['Open Prisme', 'Ouvrir Prisme', 'Abrir Prisme', 'Prisme öffnen'],
	'welcome.on': [
		'Notifications are on.',
		'Les notifications sont activées.',
		'Las notificaciones están activadas.',
		'Benachrichtigungen sind eingeschaltet.'
	],
	'welcome.blocked': [
		'Notifications are blocked. Turn them on in your system settings, then in Prisme’s settings.',
		'Les notifications sont bloquées. Activez-les dans les réglages du système, puis dans ceux de Prisme.',
		'Las notificaciones están bloqueadas. Actívalas en los ajustes del sistema y luego en los de Prisme.',
		'Benachrichtigungen sind blockiert. Schalte sie in den Systemeinstellungen und dann in Prisme ein.'
	],
	'settings.deviceNotifications': [
		'Notifications on this device',
		'Notifications sur cet appareil',
		'Notificaciones en este dispositivo',
		'Benachrichtigungen auf diesem Gerät'
	],
	'settings.notifyDesktopHint': [
		'For mentions and direct messages while the Prisme window is in the background.',
		'Pour les mentions et les messages directs lorsque la fenêtre de Prisme est en arrière-plan.',
		'Para menciones y mensajes directos mientras la ventana de Prisme está en segundo plano.',
		'Für Erwähnungen und Direktnachrichten, während das Prisme-Fenster im Hintergrund ist.'
	],
	'settings.notifyAndroidHint': [
		'For mentions and direct messages, even when Prisme is closed.',
		'Pour les mentions et les messages directs, même quand Prisme est fermé.',
		'Para menciones y mensajes directos, incluso con Prisme cerrado.',
		'Für Erwähnungen und Direktnachrichten, auch wenn Prisme geschlossen ist.'
	],
	'settings.notifySystemBlocked': [
		'Blocked in your system settings.',
		'Bloqué dans les réglages du système.',
		'Bloqueado en los ajustes del sistema.',
		'In den Systemeinstellungen blockiert.'
	],
	'settings.notifyServerOff': [
		'This server does not send push notifications yet.',
		'Ce serveur n’envoie pas encore de notifications push.',
		'Este servidor aún no envía notificaciones push.',
		'Dieser Server sendet noch keine Push-Benachrichtigungen.'
	],
	'downloads.menu': ['Get the apps', 'Télécharger les apps', 'Descargar las apps', 'Apps herunterladen'],
	'downloads.title': ['Get the apps · Prisme', 'Télécharger les apps · Prisme', 'Descargar las apps · Prisme', 'Apps herunterladen · Prisme'],
	'downloads.description': [
		'Prisme for Android, macOS, Windows and Linux.',
		'Prisme pour Android, macOS, Windows et Linux.',
		'Prisme para Android, macOS, Windows y Linux.',
		'Prisme für Android, macOS, Windows und Linux.'
	],
	'downloads.heading': ['Prisme, everywhere', 'Prisme, partout', 'Prisme, en todas partes', 'Prisme, überall'],
	'downloads.intro': [
		'The same chat on your phone and your computer, with notifications for what needs you.',
		'La même messagerie sur votre téléphone et votre ordinateur, avec des notifications pour ce qui vous concerne.',
		'El mismo chat en tu teléfono y tu ordenador, con avisos de lo que te necesita.',
		'Derselbe Chat auf Handy und Computer, mit Benachrichtigungen für das, was dich betrifft.'
	]
} as const satisfies Record<string, Row>;
