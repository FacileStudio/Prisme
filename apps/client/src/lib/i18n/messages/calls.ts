import type { Row } from '../row';

/** Calls: the call card, the call panel and the media errors. */
export const calls = {
	'call.start': ['Start a call', 'Démarrer un appel', 'Iniciar una llamada', 'Anruf starten'],
	'call.unavailable': [
		'Calls are not available on this server',
		'Les appels ne sont pas disponibles sur ce serveur',
		'Las llamadas no están disponibles en este servidor',
		'Anrufe sind auf diesem Server nicht verfügbar'
	],
	'call.startFailed': [
		'Could not start the call',
		'Impossible de démarrer l’appel',
		'No se pudo iniciar la llamada',
		'Anruf konnte nicht gestartet werden'
	],
	'call.joinFailed': [
		'Could not join the call',
		'Impossible de rejoindre l’appel',
		'No se pudo unir a la llamada',
		'Anruf konnte nicht beigetreten werden'
	],
	'call.started': ['Call started', 'Appel démarré', 'Llamada iniciada', 'Anruf gestartet'],
	'call.ended': ['Call ended', 'Appel terminé', 'Llamada finalizada', 'Anruf beendet'],
	'call.startedBy': ['by {name}', 'par {name}', 'de {name}', 'von {name}'],
	'call.endedMeta': [
		'{length} · started by {name}',
		'{length} · lancé par {name}',
		'{length} · iniciada por {name}',
		'{length} · gestartet von {name}'
	],
	'call.join': ['Join', 'Rejoindre', 'Unirse', 'Beitreten'],
	'call.joined': ['Joined', 'Rejoint', 'Te has unido', 'Beigetreten'],
	'call.underMinute': ['Under a minute', 'Moins d’une minute', 'Menos de un minuto', 'Unter einer Minute'],
	'call.minutes': ['{count} min', '{count} min', '{count} min', '{count} Min.'],
	'call.hours': ['{count} h', '{count} h', '{count} h', '{count} Std.'],
	'call.hoursMinutes': [
		'{hours} h {minutes} min',
		'{hours} h {minutes} min',
		'{hours} h {minutes} min',
		'{hours} Std. {minutes} Min.'
	],
	'call.panelTitle': ['Call in {room}', 'Appel dans {room}', 'Llamada en {room}', 'Anruf in {room}'],
	'call.defaultRoom': ['Call', 'Appel', 'Llamada', 'Anruf'],
	'call.e2eeNote': [
		'Only the people in this room hold the key. Echo and LiveKit cannot see or hear the call.',
		'Seules les personnes de ce salon détiennent la clé. Echo et LiveKit ne peuvent ni voir ni entendre l’appel.',
		'Solo las personas de esta sala tienen la clave. Echo y LiveKit no pueden ver ni oír la llamada.',
		'Nur die Personen in diesem Raum besitzen den Schlüssel. Echo und LiveKit können den Anruf weder sehen noch hören.'
	],
	'call.e2ee': [
		'End-to-end encrypted',
		'Chiffré de bout en bout',
		'Cifrado de extremo a extremo',
		'Ende-zu-Ende-verschlüsselt'
	],
	'call.connecting': ['Connecting', 'Connexion', 'Conectando', 'Verbindung wird hergestellt'],
	'call.left': ['You left the call', 'Vous avez quitté l’appel', 'Has salido de la llamada', 'Du hast den Anruf verlassen'],
	'call.dropped': [
		'The connection dropped.',
		'La connexion a été interrompue.',
		'Se perdió la conexión.',
		'Die Verbindung wurde unterbrochen.'
	],
	'call.rejoin': ['Rejoin', 'Rejoindre à nouveau', 'Volver a unirse', 'Erneut beitreten'],
	'call.enableSound': ['Enable sound', 'Activer le son', 'Activar el sonido', 'Ton aktivieren'],
	'call.micMute': [
		'Mute the microphone',
		'Couper le microphone',
		'Silenciar el micrófono',
		'Mikrofon stummschalten'
	],
	'call.micUnmute': [
		'Unmute the microphone',
		'Réactiver le microphone',
		'Activar el micrófono',
		'Mikrofon einschalten'
	],
	'call.camOff': [
		'Turn the camera off',
		'Désactiver la caméra',
		'Apagar la cámara',
		'Kamera ausschalten'
	],
	'call.camOn': ['Turn the camera on', 'Activer la caméra', 'Encender la cámara', 'Kamera einschalten'],
	'call.shareStop': [
		'Stop sharing the screen',
		'Arrêter le partage d’écran',
		'Dejar de compartir la pantalla',
		'Bildschirmfreigabe beenden'
	],
	'call.share': [
		'Share the screen',
		'Partager l’écran',
		'Compartir la pantalla',
		'Bildschirm freigeben'
	],
	'call.leave': ['Leave', 'Quitter', 'Salir', 'Verlassen'],
	'call.yourScreen': ['Your screen', 'Votre écran', 'Tu pantalla', 'Dein Bildschirm'],
	'call.screenOf': ['{name}’s screen', 'Écran de {name}', 'Pantalla de {name}', 'Bildschirm von {name}'],
	'call.you': ['{name} (you)', '{name} (vous)', '{name} (tú)', '{name} (du)'],
	'call.micOnLabel': ['Microphone on', 'Microphone activé', 'Micrófono activado', 'Mikrofon an'],
	'call.mutedLabel': ['Muted', 'Coupé', 'Silenciado', 'Stummgeschaltet'],
	'call.unsupported': [
		'This browser cannot join end-to-end encrypted calls. Use a recent Chrome, Edge, Firefox or Safari.',
		'Ce navigateur ne peut pas rejoindre les appels chiffrés de bout en bout. Utilisez une version récente de Chrome, Edge, Firefox ou Safari.',
		'Este navegador no puede unirse a llamadas cifradas de extremo a extremo. Usa una versión reciente de Chrome, Edge, Firefox o Safari.',
		'Dieser Browser kann an Ende-zu-Ende-verschlüsselten Anrufen nicht teilnehmen. Verwende eine aktuelle Version von Chrome, Edge, Firefox oder Safari.'
	],
	'call.decryptFailed': [
		'Some media could not be decrypted. Someone in the call may be on an old or foreign key: ask them to rejoin from Prisme.',
		'Certains flux n’ont pas pu être déchiffrés. Une personne de l’appel utilise peut-être une clé ancienne ou étrangère : demandez-lui de rejoindre à nouveau depuis Prisme.',
		'No se pudo descifrar parte del contenido multimedia. Es posible que alguien de la llamada use una clave antigua o ajena: pídele que se vuelva a unir desde Prisme.',
		'Einige Medien konnten nicht entschlüsselt werden. Möglicherweise nutzt jemand im Anruf einen alten oder fremden Schlüssel: Bitte die Person, über Prisme erneut beizutreten.'
	],
	'call.shareBlocked': [
		'Your system blocks screen sharing for this browser. Allow screen recording in the system settings.',
		'Votre système bloque le partage d’écran pour ce navigateur. Autorisez l’enregistrement de l’écran dans les réglages système.',
		'Tu sistema bloquea el uso compartido de pantalla en este navegador. Permite la grabación de pantalla en los ajustes del sistema.',
		'Dein System blockiert die Bildschirmfreigabe für diesen Browser. Erlaube die Bildschirmaufnahme in den Systemeinstellungen.'
	],
	'call.shareFailed': [
		'Could not share the screen',
		'Impossible de partager l’écran',
		'No se pudo compartir la pantalla',
		'Bildschirm konnte nicht freigegeben werden'
	],
	'call.micUnavailable': [
		'Microphone unavailable. You joined muted.',
		'Microphone indisponible. Vous avez rejoint l’appel en sourdine.',
		'Micrófono no disponible. Te has unido silenciado.',
		'Mikrofon nicht verfügbar. Du bist stummgeschaltet beigetreten.'
	],
	'call.deviceUnavailable': [
		'Device unavailable',
		'Appareil indisponible',
		'Dispositivo no disponible',
		'Gerät nicht verfügbar'
	],
	'notify.title': ['{sender} in {room}', '{sender} dans {room}', '{sender} en {room}', '{sender} in {room}']
} as const satisfies Record<string, Row>;
