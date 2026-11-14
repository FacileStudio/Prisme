import type { Row } from '../row';

export const landing = {
	'landing.title': ['Prisme · Team chat on Matrix', 'Prisme · Le chat d’équipe sur Matrix', 'Prisme · Chat de equipo en Matrix', 'Prisme · Team-Chat auf Matrix'],
	'landing.description': [
		'A channel per project, encrypted calls and Bulle agents you can mention.',
		'Un canal par projet, des appels chiffrés et des agents Bulle à mentionner.',
		'Un canal por proyecto, llamadas cifradas y agentes de Bulle que puedes mencionar.',
		'Ein Kanal pro Projekt, verschlüsselte Anrufe und Bulle-Agenten zum Erwähnen.'
	],
	'landing.nav.features': ['Features', 'Fonctionnalités', 'Funciones', 'Funktionen'],
	'landing.nav.apps': ['Apps', 'Apps', 'Apps', 'Apps'],
	'landing.tagline': ['Team chat for the Facile suite.', 'Le chat d’équipe de la suite Facile.', 'El chat de equipo de la suite Facile.', 'Der Team-Chat der Facile-Suite.'],
	'landing.footer.product': ['Product', 'Produit', 'Producto', 'Produkt'],
	'landing.hero.title': ['Talk to your team and your agents', 'Parlez à votre équipe et à vos agents', 'Habla con tu equipo y con tus agentes', 'Sprich mit deinem Team und deinen Agenten'],
	'landing.hero.body': [
		'One channel per project, encrypted calls, and Bulle agents that answer in the thread.',
		'Un canal par projet, des appels chiffrés et des agents Bulle qui répondent dans le fil.',
		'Un canal por proyecto, llamadas cifradas y agentes de Bulle que responden en el hilo.',
		'Ein Kanal pro Projekt, verschlüsselte Anrufe und Bulle-Agenten, die im Verlauf antworten.'
	],
	'landing.features.title': ['What Prisme does', 'Ce que fait Prisme', 'Lo que hace Prisme', 'Was Prisme kann'],
	'landing.features.spaces.title': ['A space per organization', 'Un espace par organisation', 'Un espacio por organización', 'Ein Bereich pro Organisation'],
	'landing.features.spaces.body': [
		'Spaces follow your Facile organizations, so the right people are already there.',
		'Les espaces suivent vos organisations Facile, et les bonnes personnes y sont déjà.',
		'Los espacios siguen a tus organizaciones de Facile, así que las personas adecuadas ya están dentro.',
		'Bereiche folgen deinen Facile-Organisationen, die richtigen Leute sind schon da.'
	],
	'landing.features.agents.title': ['Agents you can mention', 'Des agents à mentionner', 'Agentes que puedes mencionar', 'Agenten zum Erwähnen'],
	'landing.features.agents.body': [
		'Type @ and your Bulle agent answers in the room.',
		'Tapez @ et votre agent Bulle répond dans le salon.',
		'Escribe @ y tu agente de Bulle responde en la sala.',
		'Tippe @ und dein Bulle-Agent antwortet im Raum.'
	],
	'landing.features.calls.title': ['Encrypted calls', 'Appels chiffrés', 'Llamadas cifradas', 'Verschlüsselte Anrufe'],
	'landing.features.calls.body': [
		'Voice, video and screen sharing, end-to-end encrypted in every room.',
		'Voix, vidéo et partage d’écran, chiffrés de bout en bout dans chaque salon.',
		'Voz, vídeo y pantalla compartida, cifrados de extremo a extremo en cada sala.',
		'Sprache, Video und Bildschirmfreigabe, Ende-zu-Ende-verschlüsselt in jedem Raum.'
	],
	'landing.features.matrix.title': ['Built on Matrix', 'Construit sur Matrix', 'Hecho sobre Matrix', 'Auf Matrix gebaut'],
	'landing.features.matrix.body': [
		'An open protocol, with apps for desktop and Android.',
		'Un protocole ouvert, avec des apps pour ordinateur et Android.',
		'Un protocolo abierto, con apps para escritorio y Android.',
		'Ein offenes Protokoll, mit Apps für Desktop und Android.'
	],
	'landing.cta.title': ['Bring your team in', 'Faites venir votre équipe', 'Trae a tu equipo', 'Hol dein Team dazu'],
	'landing.cta.body': ['Sign in with your Facile account.', 'Connectez-vous avec votre compte Facile.', 'Inicia sesión con tu cuenta de Facile.', 'Melde dich mit deinem Facile-Konto an.'],
	'landing.demo.label': [
		'Prisme demo: a teammate mentions the Pixel agent in the launch channel and it answers.',
		'Démo de Prisme : une personne mentionne l’agent Pixel dans le canal launch et il répond.',
		'Demo de Prisme: alguien menciona al agente Pixel en el canal launch y responde.',
		'Prisme-Demo: Jemand erwähnt den Agenten Pixel im Kanal launch, und er antwortet.'
	],
	'landing.demo.space': ['Studio', 'Studio', 'Studio', 'Studio'],
	'landing.demo.prompt': [
		'@pixel where are we on the launch?',
		'@pixel on en est où pour le lancement ?',
		'@pixel ¿cómo va el lanzamiento?',
		'@pixel wie steht es um den Launch?'
	],
	'landing.demo.answer': [
		'Two tasks left: the pricing page and the release notes. Léa has both, due Friday.',
		'Il reste deux tâches : la page tarifs et les notes de version. Léa s’en charge, pour vendredi.',
		'Quedan dos tareas: la página de precios y las notas de versión. Léa las tiene, para el viernes.',
		'Noch zwei Aufgaben: die Preisseite und die Release Notes. Léa hat beide, fällig Freitag.'
	],
	'landing.demo.now': ['now', 'maintenant', 'ahora', 'jetzt']
} as const satisfies Record<string, Row>;
