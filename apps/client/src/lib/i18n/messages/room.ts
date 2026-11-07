import type { Row } from '../row';

/** A room: header, timeline, messages, composer and invitations. */
export const room = {
	'room.unavailable': [
		'This conversation is not available',
		'Cette conversation n’est pas disponible',
		'Esta conversación no está disponible',
		'Diese Unterhaltung ist nicht verfügbar'
	],
	'room.unavailableBody': [
		'You are not a member of it, or it does not exist.',
		'Vous n’en faites pas partie, ou elle n’existe pas.',
		'No formas parte de ella o no existe.',
		'Du bist kein Mitglied oder sie existiert nicht.'
	],
	'header.openMenu': ['Open the menu', 'Ouvrir le menu', 'Abrir el menú', 'Menü öffnen'],
	'header.searchAria': [
		'Search this space',
		'Rechercher dans cet espace',
		'Buscar en este espacio',
		'Diesen Bereich durchsuchen'
	],
	'header.membersAria': [
		'Members ({count})',
		'Membres ({count})',
		'Miembros ({count})',
		'Mitglieder ({count})'
	],
	'menu.more': ['More actions', 'Plus d’actions', 'Más acciones', 'Weitere Aktionen'],
	'menu.close': ['Close the menu', 'Fermer le menu', 'Cerrar el menú', 'Menü schließen'],
	'menu.muteChannel': ['Mute channel', 'Mettre le canal en sourdine', 'Silenciar el canal', 'Kanal stummschalten'],
	'menu.unmuteChannel': [
		'Unmute channel',
		'Réactiver le son du canal',
		'Activar el sonido del canal',
		'Kanal laut schalten'
	],
	'menu.muteDirect': [
		'Mute conversation',
		'Mettre la conversation en sourdine',
		'Silenciar la conversación',
		'Unterhaltung stummschalten'
	],
	'menu.unmuteDirect': [
		'Unmute conversation',
		'Réactiver le son de la conversation',
		'Activar el sonido de la conversación',
		'Unterhaltung laut schalten'
	],
	'menu.addAgent': ['Add an agent', 'Ajouter un agent', 'Añadir un agente', 'Agent hinzufügen'],
	'menu.renameChannel': ['Rename channel', 'Renommer le canal', 'Renombrar el canal', 'Kanal umbenennen'],
	'menu.muted': ['Muted', 'Mis en sourdine', 'Silenciado', 'Stummgeschaltet'],
	'menu.unmuted': [
		'Notifications are back on',
		'Les notifications sont réactivées',
		'Las notificaciones vuelven a estar activas',
		'Benachrichtigungen sind wieder aktiv'
	],
	'menu.muteFailed': [
		'Could not change notifications',
		'Impossible de modifier les notifications',
		'No se pudieron cambiar las notificaciones',
		'Benachrichtigungen konnten nicht geändert werden'
	],
	'composer.placeholder': ['Message {name}', 'Écrire à {name}', 'Escribir a {name}', 'Nachricht an {name}'],
	'composer.sendFailed': [
		'Could not send the message',
		'Impossible d’envoyer le message',
		'No se pudo enviar el mensaje',
		'Nachricht konnte nicht gesendet werden'
	],
	'composer.attach': ['Attach a file', 'Joindre un fichier', 'Adjuntar un archivo', 'Datei anhängen'],
	'composer.saveEdit': ['Save the edit', 'Enregistrer la modification', 'Guardar el cambio', 'Änderung speichern'],
	'composer.send': ['Send', 'Envoyer', 'Enviar', 'Senden'],
	'timeline.label': ['Messages', 'Messages', 'Mensajes', 'Nachrichten'],
	'timeline.loadingOlder': [
		'Loading older messages',
		'Chargement des anciens messages',
		'Cargando mensajes anteriores',
		'Ältere Nachrichten werden geladen'
	],
	'timeline.loadOlder': [
		'Load older messages',
		'Charger les anciens messages',
		'Cargar mensajes anteriores',
		'Ältere Nachrichten laden'
	],
	'timeline.start': [
		'This is the start of {name}.',
		'Voici le début de {name}.',
		'Este es el principio de {name}.',
		'Das ist der Anfang von {name}.'
	],
	'message.edited': ['(edited)', '(modifié)', '(editado)', '(bearbeitet)'],
	'message.deleted': [
		'This message was deleted.',
		'Ce message a été supprimé.',
		'Este mensaje se eliminó.',
		'Diese Nachricht wurde gelöscht.'
	],
	'actions.toolbar': ['Message actions', 'Actions du message', 'Acciones del mensaje', 'Nachrichtenaktionen'],
	'actions.reactAria': ['Add a reaction', 'Ajouter une réaction', 'Añadir una reacción', 'Reaktion hinzufügen'],
	'actions.react': ['React', 'Réagir', 'Reaccionar', 'Reagieren'],
	'actions.reply': ['Reply', 'Répondre', 'Responder', 'Antworten'],
	'actions.replyThread': [
		'Reply in thread',
		'Répondre dans le fil',
		'Responder en el hilo',
		'Im Thread antworten'
	],
	'actions.edit': ['Edit', 'Modifier', 'Editar', 'Bearbeiten'],
	'actions.reactFailed': [
		'Could not react',
		'Impossible de réagir',
		'No se pudo reaccionar',
		'Reaktion nicht möglich'
	],
	'actions.deleteFailed': [
		'Could not delete the message',
		'Impossible de supprimer le message',
		'No se pudo eliminar el mensaje',
		'Nachricht konnte nicht gelöscht werden'
	],
	'actions.confirmTitle': [
		'Delete this message?',
		'Supprimer ce message ?',
		'¿Eliminar este mensaje?',
		'Diese Nachricht löschen?'
	],
	'actions.confirmBody': [
		'Everyone in the conversation will see that it was deleted.',
		'Toute la conversation verra qu’il a été supprimé.',
		'Todos en la conversación verán que se eliminó.',
		'Alle in der Unterhaltung sehen, dass sie gelöscht wurde.'
	],
	'reactions.label': [
		'{key} {count}, by {who}',
		'{key} {count}, par {who}',
		'{key} {count}, de {who}',
		'{key} {count}, von {who}'
	],
	'emoji.pick': ['Pick a reaction', 'Choisir une réaction', 'Elegir una reacción', 'Reaktion auswählen'],
	'emoji.react': ['React with {emoji}', 'Réagir avec {emoji}', 'Reaccionar con {emoji}', 'Mit {emoji} reagieren'],
	'quote.original': ['Original message', 'Message d’origine', 'Mensaje original', 'Ursprüngliche Nachricht'],
	'preview.deleted': ['Deleted message', 'Message supprimé', 'Mensaje eliminado', 'Gelöschte Nachricht'],
	'preview.image': ['Sent an image', 'A envoyé une image', 'Envió una imagen', 'Hat ein Bild gesendet'],
	'preview.file': ['Sent a file', 'A envoyé un fichier', 'Envió un archivo', 'Hat eine Datei gesendet'],
	'preview.video': ['Sent a video', 'A envoyé une vidéo', 'Envió un vídeo', 'Hat ein Video gesendet'],
	'preview.audio': ['Sent audio', 'A envoyé un audio', 'Envió un audio', 'Hat eine Audiodatei gesendet'],
	'preview.callStarted': ['Started a call', 'A lancé un appel', 'Inició una llamada', 'Hat einen Anruf gestartet'],
	'preview.callEnded': ['Call ended', 'Appel terminé', 'Llamada finalizada', 'Anruf beendet'],
	'preview.mine': ['You: {text}', 'Vous : {text}', 'Tú: {text}', 'Du: {text}'],
	'draft.editing': [
		'Editing your message',
		'Modification de votre message',
		'Editando tu mensaje',
		'Deine Nachricht wird bearbeitet'
	],
	'draft.replying': [
		'Replying to {name}: {text}',
		'Réponse à {name} : {text}',
		'Respondiendo a {name}: {text}',
		'Antwort an {name}: {text}'
	],
	'invitePrompt.wants': [
		'{name} wants to talk with you.',
		'{name} souhaite discuter avec vous.',
		'{name} quiere hablar contigo.',
		'{name} möchte mit dir sprechen.'
	],
	'invitePrompt.generic': [
		'You are invited to this conversation.',
		'Vous êtes invité à cette conversation.',
		'Te han invitado a esta conversación.',
		'Du wurdest zu dieser Unterhaltung eingeladen.'
	],
	'invitePrompt.accept': ['Accept', 'Accepter', 'Aceptar', 'Annehmen'],
	'invitePrompt.decline': ['Decline', 'Refuser', 'Rechazar', 'Ablehnen'],
	'invitePrompt.failed': [
		'Could not answer the invite',
		'Impossible de répondre à l’invitation',
		'No se pudo responder a la invitación',
		'Einladung konnte nicht beantwortet werden'
	],
	'seen.by': ['Seen by', 'Vu par', 'Visto por', 'Gesehen von'],
	'seen.byNames': ['Seen by {names}', 'Vu par {names}', 'Visto por {names}', 'Gesehen von {names}'],
	'mention.label': ['Mention someone', 'Mentionner quelqu’un', 'Mencionar a alguien', 'Jemanden erwähnen'],
	'files.attachments': [
		'Attachments to send',
		'Pièces jointes à envoyer',
		'Archivos adjuntos por enviar',
		'Anhänge zum Senden'
	],
	'files.remove': ['Remove {name}', 'Retirer {name}', 'Quitar {name}', '{name} entfernen'],
	'thread.replies': [
		{ one: '{count} reply', other: '{count} replies' },
		{ one: '{count} réponse', other: '{count} réponses' },
		{ one: '{count} respuesta', other: '{count} respuestas' },
		{ one: '{count} Antwort', other: '{count} Antworten' }
	],
	'thread.lastFrom': ['last from {name}', 'dernière de {name}', 'última de {name}', 'zuletzt von {name}'],
	'typing.one': ['{a} is typing…', '{a} écrit…', '{a} está escribiendo…', '{a} schreibt…'],
	'typing.two': [
		'{a} and {b} are typing…',
		'{a} et {b} écrivent…',
		'{a} y {b} están escribiendo…',
		'{a} und {b} schreiben…'
	],
	'typing.many': [
		'{a}, {b} and {count} more are typing…',
		'{a}, {b} et {count} autres écrivent…',
		'{a}, {b} y {count} más están escribiendo…',
		'{a}, {b} und {count} weitere schreiben…'
	],
	'image.fallback': ['image', 'image', 'imagen', 'Bild'],
	'image.open': ['Open {name}', 'Ouvrir {name}', 'Abrir {name}', '{name} öffnen'],
	'image.unavailable': [
		'Image unavailable',
		'Image indisponible',
		'Imagen no disponible',
		'Bild nicht verfügbar'
	],
	'file.fallback': ['file', 'fichier', 'archivo', 'Datei'],
	'file.kind': ['File', 'Fichier', 'Archivo', 'Datei']
} as const satisfies Record<string, Row>;
