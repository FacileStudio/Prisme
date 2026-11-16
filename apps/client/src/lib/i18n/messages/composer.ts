import type { Row } from '../row';

/** The composer's "+" menu and the emoji picker. */
export const composer = {
	'composerMenu.open': ['Add to the message', 'Ajouter au message', 'Añadir al mensaje', 'Zur Nachricht hinzufügen'],
	'composerMenu.label': ['Add', 'Ajouter', 'Añadir', 'Hinzufügen'],
	'composerMenu.file': ['File', 'Fichier', 'Archivo', 'Datei'],
	'composerMenu.capsule': ['Capsule', 'Capsule', 'Cápsula', 'Kapsel'],
	'composerMenu.capsuleOff': ['Stop sealing', 'Ne plus sceller', 'Dejar de sellar', 'Nicht mehr versiegeln'],
	'emoji.open': ['Emoji', 'Emoji', 'Emoji', 'Emoji'],
	'emoji.search': ['Search emoji', 'Chercher un emoji', 'Buscar emoji', 'Emoji suchen'],
	'emoji.none': ['No emoji found', 'Aucun emoji trouvé', 'No se encontró ningún emoji', 'Kein Emoji gefunden'],
	'emoji.more': ['More emoji', 'Plus d’emojis', 'Más emojis', 'Mehr Emojis'],
	'emoji.tone': ['Skin tone', 'Teint', 'Tono de piel', 'Hautton'],
	'emoji.loading': ['Loading emoji', 'Chargement des emojis', 'Cargando emojis', 'Emojis werden geladen'],
	'emoji.group.recent': ['Recent', 'Récents', 'Recientes', 'Zuletzt'],
	'emoji.group.0': ['Smileys and emotions', 'Smileys et émotions', 'Caras y emociones', 'Smileys und Gefühle'],
	'emoji.group.1': ['People and body', 'Personnes et corps', 'Personas y cuerpo', 'Menschen und Körper'],
	'emoji.group.3': ['Animals and nature', 'Animaux et nature', 'Animales y naturaleza', 'Tiere und Natur'],
	'emoji.group.4': ['Food and drink', 'Nourriture et boissons', 'Comida y bebida', 'Essen und Trinken'],
	'emoji.group.5': ['Travel and places', 'Voyages et lieux', 'Viajes y lugares', 'Reisen und Orte'],
	'emoji.group.6': ['Activities', 'Activités', 'Actividades', 'Aktivitäten'],
	'emoji.group.7': ['Objects', 'Objets', 'Objetos', 'Objekte'],
	'emoji.group.8': ['Symbols', 'Symboles', 'Símbolos', 'Symbole'],
	'emoji.group.9': ['Flags', 'Drapeaux', 'Banderas', 'Flaggen'],
	'panel.close': ['Close', 'Fermer', 'Cerrar', 'Schließen']
} as const satisfies Record<string, Row>;
