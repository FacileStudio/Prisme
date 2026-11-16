import type { Row } from '../row';

export const channels = {
	'channel.menu': ['Channel actions', 'Actions du canal', 'Acciones del canal', 'Kanalaktionen'],
	'channel.open': ['Open', 'Ouvrir', 'Abrir', 'Öffnen'],
	'channel.join': ['Join', 'Rejoindre', 'Unirse', 'Beitreten'],
	'channel.markRead': ['Mark as read', 'Marquer comme lu', 'Marcar como leído', 'Als gelesen markieren'],
	'channel.markedRead': ['Marked as read', 'Marqué comme lu', 'Marcado como leído', 'Als gelesen markiert'],
	'channel.copyLink': ['Copy link', 'Copier le lien', 'Copiar enlace', 'Link kopieren'],
	'channel.copied': ['Link copied', 'Lien copié', 'Enlace copiado', 'Link kopiert'],
	'channel.copyFailed': [
		'Could not copy the link',
		'Impossible de copier le lien',
		'No se pudo copiar el enlace',
		'Link konnte nicht kopiert werden'
	],
	'channel.leave': ['Leave channel', 'Quitter le canal', 'Salir del canal', 'Kanal verlassen'],
	'channel.leaveTitle': ['Leave #{name}?', 'Quitter #{name} ?', '¿Salir de #{name}?', '#{name} verlassen?'],
	'channel.leaveBody': [
		'You can join it again from the space overview.',
		'Vous pourrez le rejoindre depuis la vue d’ensemble de l’espace.',
		'Puedes volver a unirte desde la vista general del espacio.',
		'Du kannst ihm in der Übersicht des Bereichs wieder beitreten.'
	],
	'channel.left': ['You left #{name}', 'Vous avez quitté #{name}', 'Saliste de #{name}', 'Du hast #{name} verlassen'],
	'channel.leaveFailed': [
		'Could not leave the channel',
		'Impossible de quitter le canal',
		'No se pudo salir del canal',
		'Kanal konnte nicht verlassen werden'
	],
	'channel.new': ['New channel', 'Nouveau canal', 'Nuevo canal', 'Neuer Kanal'],
	'channel.startTitle': ['This is the start of #{name}', 'Début de #{name}', 'Este es el inicio de #{name}', 'Hier beginnt #{name}'],
	'channel.startDirect': [
		'This is the start of your conversation with {name}',
		'Début de votre conversation avec {name}',
		'Este es el inicio de tu conversación con {name}',
		'Hier beginnt dein Gespräch mit {name}'
	],
	'channel.startBody': [
		'Nothing here yet. Say hello, share a file or start a call.',
		'Rien pour l’instant. Dites bonjour, partagez un fichier ou lancez un appel.',
		'Aún no hay nada. Saluda, comparte un archivo o inicia una llamada.',
		'Noch nichts hier. Sag Hallo, teile eine Datei oder starte einen Anruf.'
	],
	'channel.startDirectBody': [
		'Messages you send here are only between the two of you.',
		'Les messages envoyés ici restent entre vous deux.',
		'Los mensajes que envíes aquí quedan entre vosotros dos.',
		'Nachrichten hier bleiben zwischen euch beiden.'
	],
	'channel.inviteHere': ['Invite people', 'Inviter des personnes', 'Invitar personas', 'Personen einladen'],
	'channel.sayHi': ['Say hello', 'Dire bonjour', 'Saludar', 'Hallo sagen'],
	'channel.hello': ['Hello 👋', 'Bonjour 👋', 'Hola 👋', 'Hallo 👋'],
	'channel.noneJoinedTitle': ['No channels yet', 'Aucun canal pour l’instant', 'Aún no hay canales', 'Noch keine Kanäle'],
	'messages.emptyTitle': ['No conversations yet', 'Aucune conversation', 'Aún no hay conversaciones', 'Noch keine Unterhaltungen'],
	'messages.emptyBody': [
		'Pick someone below to send your first message.',
		'Choisissez quelqu’un ci-dessous pour envoyer votre premier message.',
		'Elige a alguien abajo para enviar tu primer mensaje.',
		'Wähle unten jemanden für deine erste Nachricht.'
	]
} satisfies Record<string, Row>;
