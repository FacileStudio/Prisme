import type { Row } from '../row';

/** Polls and voice messages: the composer entries, the dialogs and the cards in the timeline. */
export const polls = {
	'poll.menu': ['Poll', 'Sondage', 'Encuesta', 'Umfrage'],
	'poll.title': ['Create a poll', 'Créer un sondage', 'Crear una encuesta', 'Umfrage erstellen'],
	'poll.question': ['Question', 'Question', 'Pregunta', 'Frage'],
	'poll.questionPlaceholder': ['What do you want to ask?', 'Que voulez-vous demander ?', '¿Qué quieres preguntar?', 'Was möchten Sie fragen?'],
	'poll.options': ['Options', 'Options', 'Opciones', 'Optionen'],
	'poll.option': ['Option {n}', 'Option {n}', 'Opción {n}', 'Option {n}'],
	'poll.addOption': ['Add an option', 'Ajouter une option', 'Añadir una opción', 'Option hinzufügen'],
	'poll.removeOption': ['Remove option {n}', 'Retirer l’option {n}', 'Quitar la opción {n}', 'Option {n} entfernen'],
	'poll.multiple': ['Allow several answers', 'Autoriser plusieurs réponses', 'Permitir varias respuestas', 'Mehrere Antworten erlauben'],
	'poll.hidden': [
		'Hide results until the poll ends',
		'Masquer les résultats jusqu’à la fin',
		'Ocultar los resultados hasta el final',
		'Ergebnisse bis zum Ende verbergen'
	],
	'poll.create': ['Create poll', 'Créer le sondage', 'Crear encuesta', 'Umfrage erstellen'],
	'poll.duplicate': [
		'Two options say the same thing.',
		'Deux options disent la même chose.',
		'Dos opciones dicen lo mismo.',
		'Zwei Optionen sagen dasselbe.'
	],
	'poll.sendFailed': ['The poll could not be sent.', 'Le sondage n’a pas pu être envoyé.', 'No se pudo enviar la encuesta.', 'Die Umfrage konnte nicht gesendet werden.'],
	'poll.label': ['Poll', 'Sondage', 'Encuesta', 'Umfrage'],
	'poll.ended': ['Poll ended', 'Sondage terminé', 'Encuesta terminada', 'Umfrage beendet'],
	'poll.votes': [
		{ one: '{count} vote', other: '{count} votes' },
		{ one: '{count} vote', other: '{count} votes' },
		{ one: '{count} voto', other: '{count} votos' },
		{ one: '{count} Stimme', other: '{count} Stimmen' }
	],
	'poll.single': ['One answer', 'Une réponse', 'Una respuesta', 'Eine Antwort'],
	'poll.several': ['Up to {count} answers', 'Jusqu’à {count} réponses', 'Hasta {count} respuestas', 'Bis zu {count} Antworten'],
	'poll.hiddenUntilEnd': [
		'Results show when the poll ends.',
		'Les résultats s’affichent à la fin du sondage.',
		'Los resultados se muestran al terminar.',
		'Ergebnisse erscheinen nach dem Ende.'
	],
	'poll.yourVote': ['Your vote', 'Votre vote', 'Tu voto', 'Ihre Stimme'],
	'poll.winner': ['Most votes', 'Le plus de votes', 'Más votos', 'Meiste Stimmen'],
	'poll.end': ['End poll', 'Terminer le sondage', 'Terminar encuesta', 'Umfrage beenden'],
	'poll.endConfirm': [
		'End this poll? Nobody can vote after that.',
		'Terminer ce sondage ? Plus personne ne pourra voter.',
		'¿Terminar esta encuesta? Nadie podrá votar después.',
		'Diese Umfrage beenden? Danach kann niemand mehr abstimmen.'
	],
	'poll.endText': ['The poll has ended.', 'Le sondage est terminé.', 'La encuesta ha terminado.', 'Die Umfrage ist beendet.'],
	'poll.voteFailed': ['Your vote was not saved.', 'Votre vote n’a pas été enregistré.', 'Tu voto no se guardó.', 'Ihre Stimme wurde nicht gespeichert.'],
	'poll.preview': ['Poll: {question}', 'Sondage : {question}', 'Encuesta: {question}', 'Umfrage: {question}'],
	'voice.record': ['Record a voice message', 'Enregistrer un message vocal', 'Grabar un mensaje de voz', 'Sprachnachricht aufnehmen'],
	'voice.recording': ['Recording', 'Enregistrement', 'Grabando', 'Aufnahme'],
	'voice.cancel': ['Discard the recording', 'Supprimer l’enregistrement', 'Descartar la grabación', 'Aufnahme verwerfen'],
	'voice.send': ['Send the voice message', 'Envoyer le message vocal', 'Enviar el mensaje de voz', 'Sprachnachricht senden'],
	'voice.body': ['Voice message', 'Message vocal', 'Mensaje de voz', 'Sprachnachricht'],
	'voice.preview': ['Voice message', 'Message vocal', 'Mensaje de voz', 'Sprachnachricht'],
	'voice.denied': [
		'Prisme cannot use the microphone. Allow it in your browser or system settings.',
		'Prisme ne peut pas utiliser le micro. Autorisez-le dans les réglages du navigateur ou du système.',
		'Prisme no puede usar el micrófono. Permítelo en los ajustes del navegador o del sistema.',
		'Prisme kann das Mikrofon nicht nutzen. Erlauben Sie es in den Browser- oder Systemeinstellungen.'
	],
	'voice.tooShort': ['Hold a little longer to record.', 'Parlez un peu plus longtemps.', 'Graba un poco más.', 'Nehmen Sie etwas länger auf.'],
	'voice.sendFailed': ['The voice message could not be sent.', 'Le message vocal n’a pas pu être envoyé.', 'No se pudo enviar el mensaje de voz.', 'Die Sprachnachricht konnte nicht gesendet werden.'],
	'voice.play': ['Play the voice message', 'Écouter le message vocal', 'Reproducir el mensaje de voz', 'Sprachnachricht abspielen'],
	'voice.pause': ['Pause', 'Pause', 'Pausa', 'Pause'],
	'voice.seek': ['Position in the voice message', 'Position dans le message vocal', 'Posición en el mensaje de voz', 'Position in der Sprachnachricht'],
	'voice.failed': ['This voice message cannot be played.', 'Ce message vocal ne peut pas être lu.', 'Este mensaje de voz no se puede reproducir.', 'Diese Sprachnachricht kann nicht abgespielt werden.'],
	'voice.speed': ['Playback speed', 'Vitesse de lecture', 'Velocidad de reproducción', 'Wiedergabegeschwindigkeit']
} as const satisfies Record<string, Row>;
