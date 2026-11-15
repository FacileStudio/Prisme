import type { Row } from '../row';

/** Capsules: the composer's seal toggle and the sealed message card. */
export const capsule = {
	'capsule.toggle': ['Seal as a capsule', 'Sceller en capsule', 'Sellar como cápsula', 'Als Kapsel versiegeln'],
	'capsule.on': [
		'Sealed: burns once read',
		'Scellé : s’efface une fois lu',
		'Sellado: se borra al leerlo',
		'Versiegelt: löscht sich nach dem Lesen'
	],
	'capsule.expiry': ['Expires after', 'Expire après', 'Caduca tras', 'Läuft ab nach'],
	'capsule.1h': ['1 hour', '1 heure', '1 hora', '1 Stunde'],
	'capsule.24h': ['1 day', '1 jour', '1 día', '1 Tag'],
	'capsule.7d': ['7 days', '7 jours', '7 días', '7 Tage'],
	'capsule.30d': ['30 days', '30 jours', '30 días', '30 Tage'],
	'capsule.placeholder': ['Write the secret', 'Écrivez le secret', 'Escribe el secreto', 'Schreiben Sie das Geheimnis'],
	'capsule.sealed': ['Sealed capsule', 'Capsule scellée', 'Cápsula sellada', 'Versiegelte Kapsel'],
	'capsule.mine': [
		'Only the first person to open it can read it.',
		'Seule la première personne qui l’ouvre peut la lire.',
		'Solo la primera persona que la abra puede leerla.',
		'Nur die erste Person, die sie öffnet, kann sie lesen.'
	],
	'capsule.warn': [
		'Opening it destroys it, for everyone.',
		'L’ouvrir la détruit, pour tout le monde.',
		'Abrirla la destruye para todos.',
		'Beim Öffnen wird sie für alle gelöscht.'
	],
	'capsule.open': ['Open', 'Ouvrir', 'Abrir', 'Öffnen'],
	'capsule.opened': [
		'Opened. It is gone from Capsule; only you saw this.',
		'Ouverte. Elle n’existe plus dans Capsule ; vous seul l’avez vue.',
		'Abierta. Ya no existe en Capsule; solo tú la viste.',
		'Geöffnet. Sie ist aus Capsule gelöscht; nur Sie haben sie gesehen.'
	],
	'capsule.gone': ['Opened or expired', 'Ouverte ou expirée', 'Abierta o caducada', 'Geöffnet oder abgelaufen'],
	'capsule.failed': [
		'The capsule could not be opened.',
		'La capsule n’a pas pu être ouverte.',
		'No se pudo abrir la cápsula.',
		'Die Kapsel konnte nicht geöffnet werden.'
	],
	'capsule.sealFailed': [
		'The capsule could not be sealed.',
		'La capsule n’a pas pu être scellée.',
		'No se pudo sellar la cápsula.',
		'Die Kapsel konnte nicht versiegelt werden.'
	],
	'capsule.body': [
		'🔒 Sealed capsule, open it in Prisme',
		'🔒 Capsule scellée, ouvrez-la dans Prisme',
		'🔒 Cápsula sellada, ábrela en Prisme',
		'🔒 Versiegelte Kapsel, öffnen Sie sie in Prisme'
	]
} as const satisfies Record<string, Row>;
