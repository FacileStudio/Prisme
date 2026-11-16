import type { Row } from '../row';

/** Link cards for messages that carry a link to another Facile app. */
export const links = {
	'link.feuilles': ['Sheet in Feuilles', 'Feuille dans Feuilles', 'Hoja en Feuilles', 'Tabelle in Feuilles'],
	'link.docs': ['Document in Docs', 'Document dans Docs', 'Documento en Docs', 'Dokument in Docs'],
	'link.echo': ['Video call in Echo', 'Visio dans Echo', 'Videollamada en Echo', 'Videoanruf in Echo'],
	'link.clap': ['Video in Clap', 'Vidéo dans Clap', 'Vídeo en Clap', 'Video in Clap'],
	'link.diapo': ['Presentation in Diapo', 'Présentation dans Diapo', 'Presentación en Diapo', 'Präsentation in Diapo'],
	'link.nuage': ['File in Nuage', 'Fichier dans Nuage', 'Archivo en Nuage', 'Datei in Nuage'],
	'link.open': ['Open', 'Ouvrir', 'Abrir', 'Öffnen'],
	'link.join': ['Join the call', 'Rejoindre l’appel', 'Unirse a la llamada', 'Anruf beitreten'],
	'link.play': ['Play here', 'Lire ici', 'Reproducir aquí', 'Hier abspielen'],
	'link.close': ['Close the player', 'Fermer le lecteur', 'Cerrar el reproductor', 'Player schließen']
} as const satisfies Record<string, Row>;
