import type { Row } from '../row';

/** The Nuage picker in the composer and the card of a shared Nuage file. */
export const nuage = {
	'nuage.menu': ['File from Nuage', 'Fichier Nuage', 'Archivo de Nuage', 'Datei aus Nuage'],
	'nuage.title': ['Share from Nuage', 'Partager depuis Nuage', 'Compartir desde Nuage', 'Aus Nuage teilen'],
	'nuage.search': ['Search your Nuage', 'Chercher dans Nuage', 'Buscar en Nuage', 'In Nuage suchen'],
	'nuage.close': ['Close', 'Fermer', 'Cerrar', 'Schließen'],
	'nuage.back': ['Back', 'Retour', 'Atrás', 'Zurück'],
	'nuage.mine': ['My drive', 'Mon espace', 'Mi unidad', 'Meine Ablage'],
	'nuage.spaces': ['Shared spaces', 'Espaces partagés', 'Espacios compartidos', 'Geteilte Bereiche'],
	'nuage.empty': ['Nothing here yet', 'Rien ici pour l’instant', 'Aún no hay nada', 'Noch nichts hier'],
	'nuage.noMatch': ['No file matches', 'Aucun fichier ne correspond', 'Ningún archivo coincide', 'Keine Datei passt'],
	'nuage.loadFailed': ['Nuage did not answer', 'Nuage n’a pas répondu', 'Nuage no respondió', 'Nuage hat nicht geantwortet'],
	'nuage.reconnect': [
		'Sign in to Prisme again to reach your Nuage',
		'Reconnectez-vous à Prisme pour accéder à Nuage',
		'Vuelve a iniciar sesión en Prisme para acceder a Nuage',
		'Melden Sie sich erneut bei Prisme an, um Nuage zu erreichen'
	],
	'nuage.retry': ['Try again', 'Réessayer', 'Reintentar', 'Erneut versuchen'],
	'nuage.share': ['Share {name}', 'Partager {name}', 'Compartir {name}', '{name} teilen'],
	'nuage.shareFailed': ['Could not share this file', 'Impossible de partager ce fichier', 'No se pudo compartir el archivo', 'Datei konnte nicht geteilt werden'],
	'nuage.openFolder': ['Open {name}', 'Ouvrir {name}', 'Abrir {name}', '{name} öffnen'],
	'nuage.open': ['Open in Nuage', 'Ouvrir dans Nuage', 'Abrir en Nuage', 'In Nuage öffnen'],
	'nuage.card': ['File in Nuage', 'Fichier dans Nuage', 'Archivo en Nuage', 'Datei in Nuage'],
	'nuage.preview': ['📎 {name}', '📎 {name}', '📎 {name}', '📎 {name}']
} as const satisfies Record<string, Row>;
