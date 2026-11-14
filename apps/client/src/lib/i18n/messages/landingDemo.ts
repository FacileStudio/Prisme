import type { Row } from '../row';

/** The words drawn inside the landing's demo windows. */
export const landingDemo = {
	'landing.demo.label': [
		'Prisme demo: Sam mentions Pixel, a Bulle assistant, in the launch room and it answers with the open tasks.',
		'Démo de Prisme : Sam mentionne Pixel, un assistant Bulle, dans le salon launch et il répond avec les tâches ouvertes.',
		'Demo de Prisme: Sam menciona a Pixel, un asistente de Bulle, en la sala launch y responde con las tareas abiertas.',
		'Prisme-Demo: Sam erwähnt Pixel, einen Bulle-Assistenten, im Raum launch, und es antwortet mit den offenen Aufgaben.'
	],
	'landing.demo.space': ['Studio', 'Studio', 'Studio', 'Studio'],
	'landing.demo.first': [
		'The launch visuals are approved.',
		'Les visuels du lancement sont validés.',
		'Los visuales del lanzamiento están aprobados.',
		'Die Launch-Visuals sind freigegeben.'
	],
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
	'landing.demo.now': ['now', 'maintenant', 'ahora', 'jetzt'],
	'landing.demo.account': ['Signed in with Facile', 'Connecté avec Facile', 'Sesión iniciada con Facile', 'Mit Facile angemeldet'],
	'landing.demo.member': ['In the space', 'Dans l’espace', 'En el espacio', 'Im Bereich'],
	'landing.demo.ask': [
		'@pixel what is left for Friday?',
		'@pixel il reste quoi pour vendredi ?',
		'@pixel ¿qué queda para el viernes?',
		'@pixel was ist bis Freitag offen?'
	],
	'landing.demo.taskPricing': ['Pricing page', 'Page tarifs', 'Página de precios', 'Preisseite'],
	'landing.demo.taskNotes': ['Release notes', 'Notes de version', 'Notas de versión', 'Release Notes'],
	'landing.demo.friday': ['Friday', 'vendredi', 'viernes', 'Freitag'],
	'landing.demo.callAsk': ['Quick call?', 'On s’appelle ?', '¿Hacemos una llamada?', 'Kurzer Anruf?'],
	'landing.demo.encrypted': ['Encrypted', 'Chiffré', 'Cifrada', 'Verschlüsselt']
} as const satisfies Record<string, Row>;
