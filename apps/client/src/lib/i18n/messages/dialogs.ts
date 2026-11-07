import type { Row } from '../row';

/** The modal dialogs: create, invite, rename and add an agent. */
export const dialogs = {
	'addAgent.title': ['Add an agent', 'Ajouter un agent', 'Añadir un agente', 'Agent hinzufügen'],
	'addAgent.loading': [
		'Loading your agents',
		'Chargement de vos agents',
		'Cargando tus agentes',
		'Deine Agents werden geladen'
	],
	'addAgent.none': ['No agents yet', 'Aucun agent pour le moment', 'Aún no hay agentes', 'Noch keine Agents'],
	'addAgent.noneBody': [
		'Your Bulle assistants appear here once Bulle connects them to Prisme.',
		'Vos assistants Bulle apparaissent ici dès que Bulle les connecte à Prisme.',
		'Tus asistentes de Bulle aparecen aquí cuando Bulle los conecta a Prisme.',
		'Deine Bulle-Assistenten erscheinen hier, sobald Bulle sie mit Prisme verbindet.'
	],
	'addAgent.here': ['In channel', 'Dans le canal', 'En el canal', 'Im Kanal'],
	'addAgent.add': ['Add', 'Ajouter', 'Añadir', 'Hinzufügen'],
	'addAgent.joined': [
		'{name} joined the channel',
		'{name} a rejoint le canal',
		'{name} se ha unido al canal',
		'{name} ist dem Kanal beigetreten'
	],
	'addAgent.failed': [
		'Could not add the agent',
		'Impossible d’ajouter l’agent',
		'No se pudo añadir el agente',
		'Agent konnte nicht hinzugefügt werden'
	],
	'addAgent.listFailed': [
		'Could not list your agents',
		'Impossible de lister vos agents',
		'No se pudieron listar tus agentes',
		'Deine Agents konnten nicht geladen werden'
	],
	'createChannel.title': [
		'Create a channel in {name}',
		'Créer un canal dans {name}',
		'Crear un canal en {name}',
		'Kanal in {name} erstellen'
	],
	'createChannel.topic': ['Topic', 'Sujet', 'Tema', 'Thema'],
	'createChannel.topicHelp': [
		'Optional. Shown under the channel name.',
		'Facultatif. Affiché sous le nom du canal.',
		'Opcional. Se muestra bajo el nombre del canal.',
		'Optional. Wird unter dem Kanalnamen angezeigt.'
	],
	'createChannel.topicPlaceholder': [
		'What this channel is for',
		'À quoi sert ce canal',
		'Para qué sirve este canal',
		'Wofür dieser Kanal gedacht ist'
	],
	'createChannel.failed': [
		'Could not create the channel',
		'Impossible de créer le canal',
		'No se pudo crear el canal',
		'Kanal konnte nicht erstellt werden'
	],
	'createSpace.title': ['Create a space', 'Créer un espace', 'Crear un espacio', 'Bereich erstellen'],
	'createSpace.body': [
		'A space groups the channels of a project or a team. It starts with a #general channel.',
		'Un espace regroupe les canaux d’un projet ou d’une équipe. Il démarre avec un canal #general.',
		'Un espacio agrupa los canales de un proyecto o de un equipo. Empieza con un canal #general.',
		'Ein Bereich bündelt die Kanäle eines Projekts oder Teams. Er beginnt mit einem Kanal #general.'
	],
	'createSpace.failed': [
		'Could not create the space',
		'Impossible de créer l’espace',
		'No se pudo crear el espacio',
		'Bereich konnte nicht erstellt werden'
	],
	'invite.title': [
		'Invite people to {name}',
		'Inviter des personnes dans {name}',
		'Invitar a personas a {name}',
		'Personen zu {name} einladen'
	],
	'invite.body': [
		'Invited people join the space and its #general channel when they accept.',
		'Les personnes invitées rejoignent l’espace et son canal #general quand elles acceptent.',
		'Las personas invitadas se unen al espacio y a su canal #general cuando aceptan.',
		'Eingeladene Personen treten dem Bereich und seinem Kanal #general bei, sobald sie annehmen.'
	],
	'invite.done': [
		'{name} is invited to {space}',
		'{name} est invité(e) dans {space}',
		'{name} ha sido invitado a {space}',
		'{name} wurde zu {space} eingeladen'
	],
	'invite.failed': [
		'Could not send the invite',
		'Impossible d’envoyer l’invitation',
		'No se pudo enviar la invitación',
		'Einladung konnte nicht gesendet werden'
	],
	'rename.titleChannel': [
		'Rename the channel',
		'Renommer le canal',
		'Renombrar el canal',
		'Kanal umbenennen'
	],
	'rename.titleSpace': ['Rename the space', 'Renommer l’espace', 'Renombrar el espacio', 'Bereich umbenennen'],
	'rename.submit': ['Rename', 'Renommer', 'Renombrar', 'Umbenennen'],
	'rename.doneChannel': [
		'The channel is renamed',
		'Le canal est renommé',
		'El canal se ha renombrado',
		'Der Kanal wurde umbenannt'
	],
	'rename.doneSpace': [
		'The space is renamed',
		'L’espace est renommé',
		'El espacio se ha renombrado',
		'Der Bereich wurde umbenannt'
	],
	'rename.failedChannel': [
		'Could not rename the channel',
		'Impossible de renommer le canal',
		'No se pudo renombrar el canal',
		'Kanal konnte nicht umbenannt werden'
	],
	'rename.failedSpace': [
		'Could not rename the space',
		'Impossible de renommer l’espace',
		'No se pudo renombrar el espacio',
		'Bereich konnte nicht umbenannt werden'
	],
	'picker.placeholder': [
		'Search people by name or email',
		'Rechercher des personnes par nom ou e-mail',
		'Buscar personas por nombre o correo',
		'Personen nach Name oder E-Mail suchen'
	],
	'picker.label': ['Search people', 'Rechercher des personnes', 'Buscar personas', 'Personen suchen'],
	'picker.searching': ['Searching…', 'Recherche…', 'Buscando…', 'Suche läuft…'],
	'picker.nobody': ['Nobody matches.', 'Personne ne correspond.', 'Nadie coincide.', 'Niemand passt.']
} as const satisfies Record<string, Row>;
