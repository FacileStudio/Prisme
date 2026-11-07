import type { Row } from '../row';

/** The app shell: sign-in, rail, overview and space pages. */
export const shell = {
	'login.title': ['Sign in · Prisme', 'Connexion · Prisme', 'Iniciar sesión · Prisme', 'Anmelden · Prisme'],
	'login.tagline': [
		'The Facile team chat.',
		'La messagerie d’équipe Facile.',
		'El chat de equipo de Facile.',
		'Der Team-Chat von Facile.'
	],
	'login.failedTitle': [
		'Sign-in failed',
		'La connexion a échoué',
		'No se pudo iniciar sesión',
		'Anmeldung fehlgeschlagen'
	],
	'login.noSso': [
		'Single sign-on is not configured on this server. Ask an administrator.',
		'L’authentification unique n’est pas configurée sur ce serveur. Contactez un administrateur.',
		'El inicio de sesión único no está configurado en este servidor. Consulta a un administrador.',
		'Single Sign-on ist auf diesem Server nicht eingerichtet. Wende dich an eine Administratorin oder einen Administrator.'
	],
	'app.startFailed': [
		'Prisme could not start',
		'Prisme n’a pas pu démarrer',
		'Prisme no pudo iniciarse',
		'Prisme konnte nicht gestartet werden'
	],
	'app.retry': ['Try again', 'Réessayer', 'Reintentar', 'Erneut versuchen'],
	'app.connecting': [
		'Connecting to the chat server',
		'Connexion au serveur de discussion',
		'Conectando con el servidor de chat',
		'Verbindung zum Chat-Server wird hergestellt'
	],
	'app.firstSyncFailed': [
		'The chat server did not answer the first sync.',
		'Le serveur de discussion n’a pas répondu à la première synchronisation.',
		'El servidor de chat no respondió a la primera sincronización.',
		'Der Chat-Server hat die erste Synchronisierung nicht beantwortet.'
	],
	'app.clientNotRunning': [
		'The chat client is not running.',
		'Le client de discussion n’est pas démarré.',
		'El cliente de chat no está en marcha.',
		'Der Chat-Client läuft nicht.'
	],
	'nav.messages': ['Messages', 'Messages', 'Mensajes', 'Nachrichten'],
	'nav.overview': ['Overview', 'Aperçu', 'Resumen', 'Übersicht'],
	'nav.channels': ['Channels', 'Canaux', 'Canales', 'Kanäle'],
	'nav.spaces': ['Spaces', 'Espaces', 'Espacios', 'Bereiche'],
	'spaces.manage': [
		'Manage shared spaces',
		'Gérer les espaces partagés',
		'Gestionar los espacios compartidos',
		'Gemeinsame Bereiche verwalten'
	],
	'overview.title': ['Overview · Prisme', 'Aperçu · Prisme', 'Resumen · Prisme', 'Übersicht · Prisme'],
	'overview.titleSpace': ['{name} · Prisme', '{name} · Prisme', '{name} · Prisme', '{name} · Prisme'],
	'overview.invited': [
		'You are invited to join {name}.',
		'Vous êtes invité à rejoindre {name}.',
		'Te han invitado a unirte a {name}.',
		'Du wurdest eingeladen, {name} beizutreten.'
	],
	'overview.joinSpace': ['Join the space', 'Rejoindre l’espace', 'Unirse al espacio', 'Dem Bereich beitreten'],
	'overview.joinFailed': [
		'Could not join the space',
		'Impossible de rejoindre l’espace',
		'No se pudo unir al espacio',
		'Beitritt zum Bereich nicht möglich'
	],
	'overview.none': ['No space yet', 'Aucun espace pour le moment', 'Aún no hay espacios', 'Noch kein Bereich'],
	'overview.noneBody': [
		'A space groups the channels of a project or a team.',
		'Un espace regroupe les canaux d’un projet ou d’une équipe.',
		'Un espacio agrupa los canales de un proyecto o de un equipo.',
		'Ein Bereich bündelt die Kanäle eines Projekts oder Teams.'
	],
	'overview.footer': [
		'Spaces group the channels of a team or a project.',
		'Les espaces regroupent les canaux d’une équipe ou d’un projet.',
		'Los espacios agrupan los canales de un equipo o de un proyecto.',
		'Bereiche bündeln die Kanäle eines Teams oder Projekts.'
	],
	'overview.createSpace': ['Create a space', 'Créer un espace', 'Crear un espacio', 'Bereich erstellen'],
	'space.renameAria': ['Rename {name}', 'Renommer {name}', 'Renombrar {name}', '{name} umbenennen'],
	'space.renameTitle': [
		'Rename the space',
		'Renommer l’espace',
		'Renombrar el espacio',
		'Bereich umbenennen'
	],
	'space.channelCount': [
		{ one: '{count} channel', other: '{count} channels' },
		{ one: '{count} canal', other: '{count} canaux' },
		{ one: '{count} canal', other: '{count} canales' },
		{ one: '{count} Kanal', other: '{count} Kanäle' }
	],
	'space.memberCount': [
		{ one: '{count} member', other: '{count} members' },
		{ one: '{count} membre', other: '{count} membres' },
		{ one: '{count} miembro', other: '{count} miembros' },
		{ one: '{count} Mitglied', other: '{count} Mitglieder' }
	],
	'space.invite': ['Invite', 'Inviter', 'Invitar', 'Einladen'],
	'space.newChannel': ['New channel', 'Nouveau canal', 'Nuevo canal', 'Neuer Kanal'],
	'space.yours': ['Your channels', 'Vos canaux', 'Tus canales', 'Deine Kanäle'],
	'space.noneJoined': [
		'You have not joined a channel in {name} yet.',
		'Vous n’avez encore rejoint aucun canal dans {name}.',
		'Todavía no te has unido a ningún canal de {name}.',
		'Du bist noch keinem Kanal in {name} beigetreten.'
	],
	'space.browse': ['Browse channels', 'Parcourir les canaux', 'Explorar canales', 'Kanäle durchsuchen'],
	'space.browseBody': [
		'Channels in {name} you have not joined yet.',
		'Canaux de {name} que vous n’avez pas encore rejoints.',
		'Canales de {name} a los que aún no te has unido.',
		'Kanäle in {name}, denen du noch nicht beigetreten bist.'
	],
	'space.joinFailed': [
		'Could not join the channel',
		'Impossible de rejoindre le canal',
		'No se pudo unir al canal',
		'Beitritt zum Kanal nicht möglich'
	],
	'row.muted': ['Muted', 'En sourdine', 'Silenciado', 'Stummgeschaltet'],
	'messages.title': ['Messages · Prisme', 'Messages · Prisme', 'Mensajes · Prisme', 'Nachrichten · Prisme'],
	'messages.description': [
		'Talk one to one with your teammates and with any agent.',
		'Discutez en tête à tête avec vos coéquipiers et avec n’importe quel agent.',
		'Habla en privado con tu equipo y con cualquier agente.',
		'Schreibe direkt mit deinem Team und mit jedem Agent.'
	],
	'messages.search': [
		'Search people and agents',
		'Rechercher des personnes et des agents',
		'Buscar personas y agentes',
		'Personen und Agents suchen'
	],
	'messages.searchLabel': [
		'Search people, agents and conversations',
		'Rechercher des personnes, des agents et des conversations',
		'Buscar personas, agentes y conversaciones',
		'Personen, Agents und Unterhaltungen suchen'
	],
	'messages.filter': ['Who to show', 'Qui afficher', 'A quién mostrar', 'Wer angezeigt wird'],
	'messages.tabAll': ['Everyone', 'Tout le monde', 'Todos', 'Alle'],
	'messages.tabPeople': ['People', 'Personnes', 'Personas', 'Personen'],
	'messages.tabAgents': ['Agents', 'Agents', 'Agentes', 'Agents'],
	'messages.conversations': ['Conversations', 'Conversations', 'Conversaciones', 'Unterhaltungen'],
	'messages.yourAgent': ['Your agent', 'Votre agent', 'Tu agente', 'Dein Agent'],
	'messages.agentOf': ['{name}’s agent', 'Agent de {name}', 'Agente de {name}', 'Agent von {name}'],
	'messages.wantsToMessage': [
		'Wants to message you',
		'Souhaite vous écrire',
		'Quiere escribirte',
		'Möchte dir schreiben'
	],
	'messages.noMessages': ['No messages yet', 'Aucun message', 'Aún no hay mensajes', 'Noch keine Nachrichten'],
	'messages.openFailed': [
		'Could not open the conversation',
		'Impossible d’ouvrir la conversation',
		'No se pudo abrir la conversación',
		'Unterhaltung konnte nicht geöffnet werden'
	],
	'messages.noMatch': [
		'No conversation matches',
		'Aucune conversation ne correspond',
		'Ninguna conversación coincide',
		'Keine Unterhaltung passt'
	],
	'messages.none': [
		'No conversations yet',
		'Aucune conversation pour le moment',
		'Aún no hay conversaciones',
		'Noch keine Unterhaltungen'
	],
	'messages.tryOther': [
		'Try another name, or start one below.',
		'Essayez un autre nom, ou démarrez-en une ci-dessous.',
		'Prueba con otro nombre o inicia una abajo.',
		'Versuche einen anderen Namen oder starte unten eine neue.'
	],
	'messages.pickSomeone': [
		'Pick someone below to send your first message.',
		'Choisissez quelqu’un ci-dessous pour envoyer votre premier message.',
		'Elige a alguien abajo para enviar tu primer mensaje.',
		'Wähle unten jemanden aus, um deine erste Nachricht zu senden.'
	],
	'messages.start': [
		'Start a conversation',
		'Démarrer une conversation',
		'Iniciar una conversación',
		'Unterhaltung beginnen'
	],
	'messages.startBody': [
		'Everyone on Prisme you can message, people and agents alike.',
		'Toutes les personnes et tous les agents de Prisme que vous pouvez contacter.',
		'Todas las personas y agentes de Prisme a los que puedes escribir.',
		'Alle Personen und Agents auf Prisme, denen du schreiben kannst.'
	],
	'messages.message': ['Message', 'Écrire', 'Escribir', 'Schreiben'],
	'messages.nobodyMatch': [
		'Nobody matches',
		'Personne ne correspond',
		'Nadie coincide',
		'Niemand passt'
	],
	'messages.nobody': ['Nobody here yet', 'Personne pour le moment', 'Aún no hay nadie', 'Noch niemand hier'],
	'messages.searchHint': [
		'Search by name, email or @handle.',
		'Recherchez par nom, e-mail ou @identifiant.',
		'Busca por nombre, correo o @identificador.',
		'Suche nach Name, E-Mail oder @Handle.'
	],
	'messages.signInHint': [
		'People show up once they sign in to Prisme.',
		'Les personnes apparaissent dès leur première connexion à Prisme.',
		'Las personas aparecen cuando inician sesión en Prisme.',
		'Personen erscheinen, sobald sie sich bei Prisme anmelden.'
	]
} as const satisfies Record<string, Row>;
