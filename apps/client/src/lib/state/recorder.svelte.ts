import { downsample, MAX_VOICE_MS, MIN_VOICE_MS, recordingType } from '$lib/chat/voice';

export type Recording = { blob: Blob; duration: number; waveform: number[] };
type Phase = 'idle' | 'starting' | 'recording';

const LIVE_BARS = 40;

/**
 * A microphone recording with a live level meter. `start` asks for the mic,
 * `stop` resolves the recording (or null when too short), `cancel` drops it.
 */
export class Recorder {
	phase = $state<Phase>('idle');
	elapsed = $state(0);
	levels = $state<number[]>(Array.from({ length: LIVE_BARS }, () => 0));

	private stream: MediaStream | null = null;
	private recorder: MediaRecorder | null = null;
	private context: AudioContext | null = null;
	private frame = 0;
	private chunks: Blob[] = [];
	private samples: number[] = [];
	private startedAt = 0;
	private onLimit: (() => void) | null = null;

	static supported(): boolean {
		return typeof window !== 'undefined' && !!navigator.mediaDevices?.getUserMedia && typeof MediaRecorder !== 'undefined';
	}

	async start(onLimit: () => void): Promise<void> {
		if (this.phase !== 'idle') return;
		this.phase = 'starting';
		this.onLimit = onLimit;
		try {
			this.stream = await navigator.mediaDevices.getUserMedia({ audio: { echoCancellation: true, noiseSuppression: true } });
			const mimeType = recordingType((type) => MediaRecorder.isTypeSupported(type));
			this.recorder = new MediaRecorder(this.stream, mimeType ? { mimeType } : undefined);
			this.chunks = [];
			this.samples = [];
			this.recorder.ondataavailable = (e) => {
				if (e.data.size > 0) this.chunks.push(e.data);
			};
			const level = this.meter(this.stream);
			this.recorder.start(250);
			this.startedAt = performance.now();
			this.tick(level);
			this.elapsed = 0;
			this.phase = 'recording';
		} catch (err) {
			this.release();
			throw err;
		}
	}

	stop(): Promise<Recording | null> {
		const recorder = this.recorder;
		if (!recorder || this.phase !== 'recording') {
			this.release();
			return Promise.resolve(null);
		}
		const duration = Math.round(performance.now() - this.startedAt);
		return new Promise((resolve) => {
			recorder.onstop = () => {
				const blob = new Blob(this.chunks, { type: recorder.mimeType || 'audio/webm' });
				const waveform = downsample(this.samples);
				this.release();
				resolve(duration < MIN_VOICE_MS || blob.size === 0 ? null : { blob, duration, waveform });
			};
			recorder.stop();
		});
	}

	cancel(): void {
		if (this.recorder && this.recorder.state !== 'inactive') {
			this.recorder.onstop = null;
			this.recorder.stop();
		}
		this.release();
	}

	private meter(stream: MediaStream): () => number {
		const AudioCtx = window.AudioContext ?? (window as unknown as { webkitAudioContext?: typeof AudioContext }).webkitAudioContext;
		if (!AudioCtx) return () => 0;
		this.context = new AudioCtx();
		const analyser = this.context.createAnalyser();
		analyser.fftSize = 512;
		this.context.createMediaStreamSource(stream).connect(analyser);
		const data = new Uint8Array(analyser.fftSize);
		return () => {
			analyser.getByteTimeDomainData(data);
			let sum = 0;
			for (const v of data) sum += ((v - 128) / 128) ** 2;
			return Math.min(1, Math.sqrt(sum / data.length) * 3);
		};
	}

	private tick(level: () => number) {
		const step = () => {
			const value = level();
			this.samples.push(value);
			this.levels = [...this.levels.slice(1), value];
			this.elapsed = Math.round(performance.now() - this.startedAt);
			if (this.elapsed >= MAX_VOICE_MS) this.onLimit?.();
			this.frame = requestAnimationFrame(step);
		};
		this.frame = requestAnimationFrame(step);
	}

	private release() {
		cancelAnimationFrame(this.frame);
		this.frame = 0;
		this.stream?.getTracks().forEach((track) => track.stop());
		this.stream = null;
		void this.context?.close().catch(() => undefined);
		this.context = null;
		this.recorder = null;
		this.phase = 'idle';
		this.levels = Array.from({ length: LIVE_BARS }, () => 0);
	}
}
