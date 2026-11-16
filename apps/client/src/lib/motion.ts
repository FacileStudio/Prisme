import { cubicIn, cubicOut } from 'svelte/easing';
import type { TransitionConfig } from 'svelte/transition';
import { prefersReducedMotion } from '@facile/muse';

type Options = { direction?: 'in' | 'out' | 'both' };

/**
 * A floating surface growing out of the control that opened it: scale and fade from the
 * element's `transform-origin`, quicker on the way out. Transform and opacity only, so the
 * layout never moves, and instant under reduced motion.
 */
export function pop(_node: Element, { start = 0.92, duration = 180 } = {}, options: Options = {}): TransitionConfig {
	const out = options.direction === 'out';
	return {
		duration: prefersReducedMotion() ? 0 : out ? Math.round(duration * 0.65) : duration,
		easing: out ? cubicIn : cubicOut,
		css: (t, u) => `opacity: ${t}; transform: translateY(${u * 6}px) scale(${start + (1 - start) * t});`
	};
}

/** A tile or a row settling in place, for staggered lists. Instant under reduced motion. */
export function rise(_node: Element, { delay = 0, duration = 220, distance = 8 } = {}): TransitionConfig {
	const reduced = prefersReducedMotion();
	return {
		delay: reduced ? 0 : delay,
		duration: reduced ? 0 : duration,
		easing: cubicOut,
		css: (t, u) => `opacity: ${t}; transform: translateY(${u * distance}px);`
	};
}

/** The delay of the nth item of a staggered list, capped so a long list never lags. */
export function stagger(index: number, step = 22, cap = 12): number {
	return Math.min(index, cap) * step;
}
